package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/lib/cacher"
	"github.com/artnoi43/stubborn/lib/dohclient"
	"github.com/artnoi43/stubborn/lib/enums"
)

type handler struct {
	ctx       context.Context
	config    *Config
	dnsServer *dns.Server
	dohClient *doh.DoH
	cacher    *cacher.Cacher
	fMap      networkHandleFuncMap
}

// answerMap maps cache key to answers (domain-type:val)
type answerMap map[cacher.Key][]string

func New(ctx context.Context, conf *Config, s *dns.Server, h *doh.DoH, c *cacher.Cacher) *handler {
	j, _ := json.Marshal(conf)
	log.Printf("DNS server configuration:\n%s\n", j)
	_handler := &handler{
		ctx:       ctx,
		config:    conf,
		dnsServer: s,
		dohClient: h,
		cacher:    c,
	}
	return _handler.mapFuncs()
}

// NewRRA returns new RR for supported DNS record types (in dnstypes.go)
func NewRR(dom string, t string, v string) (dns.RR, error) {
	rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", dom, t, v))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("new RR failed for %s:%s", dom, v))
	}
	return rr, nil
}

func (h *handler) Handle(m *dns.Msg) error {
	dohFunc := dohclient.FuncMap[h.config.AllTypes]
	for _, q := range m.Question {
		// First we look in cache
		t, supported := dnsTypes[q.Qtype]
		if !supported {
			return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
		}
		k := cacher.NewKey(q.Name, t, -1)
		answers, err := h.cacher.HGet(k)
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
			dohAnswers, err := dohFunc(h.ctx, h.dohClient, dom, dohdns.Type(t))
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
				if err := h.cacher.HSet(k, a); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Handle local A record DNS queries
func (h *handler) HandleLocal(m *dns.Msg, t map[string]string) error {
	for _, q := range m.Question {
		if q.Qtype != dns.TypeA {
			continue
		}
		for k, v := range t {
			if q.Name == k || q.Name == k+"." {
				rr, err := NewRR(q.Name, enums.DnsTypes[dns.TypeA], v)
				if err != nil {
					return err
				}
				m.Answer = append(m.Answer, rr)
				break
			}
		}
	}
	return nil
}

func (h *handler) HandleDnsReq(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg).SetReply(r)
	m.Compress = false
	m.RecursionDesired = true
	m.RecursionAvailable = true

	switch r.Opcode {
	case dns.OpcodeQuery:
		if err := h.Handle(m); err != nil {
			log.Println("Handle error:", err.Error())
		}
	}
	w.WriteMsg(m)
}

func (h *handler) HandleLocalDnsReq(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg).SetReply(r)
	m.Compress = false

	b, err := os.ReadFile(h.config.HostsFile)
	if err != nil {
		log.Printf("failed to read hosts file %s: %v", h.config.HostsFile, err.Error())
	}
	var table = make(map[string]string)
	if err := json.Unmarshal(b, &table); err != nil {
		log.Printf("failed to parse hosts file %s: %v", h.config.HostsFile, err.Error())
	}
	if err := h.HandleLocal(m, table); err != nil {
		log.Println("HandleLocal error:", err.Error())
	}
	if err := w.WriteMsg(m); err != nil {
		log.Println("Error writing reply", err.Error())
	}
}

func (h *handler) Start(addr string) error {
	log.Println("Starting stubborn DNS resolver on", addr)
	return h.dnsServer.ListenAndServe()
}

func (h *handler) Shutdown(ctx context.Context) {
	log.Println("shutting down handlers")
	h.dohClient.Close()
	h.cacher.Redis.Cli.FlushDB(ctx)
	h.cacher.Redis.Cli.Close()
	log.Println("handlers gracefully shutdowm")
}
