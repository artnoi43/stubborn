package handler

import (
	"fmt"
	"log"

	dohdns "github.com/likexian/doh-go/dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/domain/usecase/dohclient"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

func (h *handler) HandleDoH(m *dns.Msg) error {
	dohFunc := dohclient.FuncMap[h.conf.QueryAllTypes]
	// var wg sync.WaitGroup
	for _, q := range m.Question {
		// First we look in cache
		_, supported := dnsutils.DnsTypes[q.Qtype]
		if !supported {
			return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
		}
		k := dnsutils.Key{
			Dom:  q.Name,
			Type: q.Qtype,
		}
		cached, found := h.repository.Get(k)
		if len(cached) >= 1 && found {
			log.Printf("cache hit for: %s\n", k.String())
			for _, answer := range cached {
				m.Answer = append(m.Answer, answer)
			}
		} else {
			// Then we use DoH to query uncached domains
			log.Printf("cache missed for: %s\n", k.String())
			t, supported := dnsutils.DnsTypes[q.Qtype]
			if !supported {
				return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
			}
			dohAnswers, err := dohFunc(
				h.ctx,
				h.dohClient,
				dohdns.Domain(q.Name),
				dohdns.Type(t),
			)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to get DoH answer for %s", q.Name))
			}
			if l := len(dohAnswers); l == 0 {
				return fmt.Errorf("record not found: %s %s", q.Name, t)
			} else {
				log.Printf("Got %d answer(s)\n", l)
			}
			answerMap := make(answerMap)
			for _, dohAnswer := range dohAnswers {
				t, supported := dnsutils.DnsTypes[uint16(dohAnswer.Type)]
				if !supported {
					return fmt.Errorf("unsupported DNS record type: %d", q.Qtype)
				}
				rr, err := dnsutils.NewRR(q.Name, t, dohAnswer.Data)
				if err != nil {
					return err
				}
				m.Answer = append(m.Answer, rr)
				k := dnsutils.NewKey(q.Name, q.Qtype)
				answerMap[k] = append(answerMap[k], rr)
			}
			h.repository.SetMap(answerMap)
		}
	}
	// wg.Wait()
	return nil
}

func (h *handler) HandleWithDoH(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg).SetReply(r)
	m.Compress = false
	m.RecursionDesired = true
	m.RecursionAvailable = true

	switch r.Opcode {
	case dns.OpcodeQuery:
		if err := h.HandleDoH(m); err != nil {
			log.Println("handler error:", err.Error())
		}
	}
	w.WriteMsg(m)
}
