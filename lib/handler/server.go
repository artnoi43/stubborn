package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/lib/cacher"
	"github.com/artnoi43/stubborn/lib/dohclient"
)

type Handler struct {
	Ctx       context.Context
	Config    *Config
	DnsServer *dns.Server
	DohClient *doh.DoH
	Cacher    *cacher.Cacher
}

// answerMap maps cache key to answers (domain-type:val)
type answerMap map[cacher.Key][]string

func NewDNSServer(conf *Config) *dns.Server {
	return &dns.Server{
		Addr: conf.Address,
		Net:  conf.Protocol,
	}
}

func New(ctx context.Context, conf *Config, s *dns.Server, h *doh.DoH, c *cacher.Cacher) *Handler {
	j, _ := json.Marshal(conf)
	log.Printf("DNS server configuration:\n%s\n", j)
	return &Handler{
		Ctx:       ctx,
		Config:    conf,
		DnsServer: s,
		DohClient: h,
		Cacher:    c,
	}
}

// NewRRA returns new RR for supported DNS record types (in dnstypes.go)
func NewRR(dom string, t string, v string) (dns.RR, error) {
	rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", dom, t, v))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("new RR failed for %s:%s", dom, v))
	}
	return rr, nil
}

func (h *Handler) Action(m *dns.Msg) error {
	dohFunc := dohclient.FuncMap[h.Config.AllTypes]
	for _, q := range m.Question {
		// First we look in cache
		t, supported := dnsTypes[q.Qtype]
		if !supported {
			return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
		}
		k := cacher.NewKey(q.Name, t, -1)
		answers, err := h.Cacher.HGet(k)
		if answers != nil && err == nil {
			for _, answer := range answers {
				rr, err := NewRR(q.Name, t, answer)
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				return err
			}
		} else {
			// Then we use DoH to query uncached domains
			log.Println("All cache missed:", k.MemCacheKey())
			dom := dohdns.Domain(q.Name)
			t, supported := dnsTypes[q.Qtype]
			if !supported {
				return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
			}
			dohAnswers, err := dohFunc(h.Ctx, h.DohClient, dom, dohdns.Type(t))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to get DoH answer for %s", dom))
			}
			if l := len(dohAnswers); l == 0 {
				return fmt.Errorf("record not found: %s %s", q.Name, t)
			} else {
				log.Printf("Got %d answer(s)\n", l)
			}
			answerMap := make(answerMap)
			for _, dohAnswer := range dohAnswers {
				t, supported := dnsTypes[uint16(dohAnswer.Type)]
				if !supported {
					return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
				}
				rr, err := NewRR(q.Name, t, dohAnswer.Data)
				if err != nil {
					return err
				}
				m.Answer = append(m.Answer, rr)
				k := cacher.NewKey(q.Name, t, dohAnswer.TTL)
				answerMap[k] = append(answerMap[k], dohAnswer.Data)
			}
			for k, a := range answerMap {
				if err := h.Cacher.HSet(k, a); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (h *Handler) HandleDnsReq(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg).SetReply(r)
	m.Compress = false
	m.RecursionDesired = true
	m.RecursionAvailable = true

	switch r.Opcode {
	case dns.OpcodeQuery:
		if err := h.Action(m); err != nil {
			log.Println("Action error:", err.Error())
		}
	}
	w.WriteMsg(m)
}

func (h *Handler) Start() error {
	log.Println("Starting stubborn DNS resolver on", h.Config.Address)
	return h.DnsServer.ListenAndServe()
}
