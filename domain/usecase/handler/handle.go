package handler

import (
	"encoding/json"
	"log"
	"os"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

// answerMap maps cache key to answers (key: "domain:type", val: "val")
type answerMap map[dnsutils.Key]entity.CachedAnswer

// Handle local A record DNS queries
func (h *handler) HandleLocal(m *dns.Msg, z map[string]string) error {
	for _, q := range m.Question {
		if q.Qtype != dns.TypeA {
			continue
		}
		for k, v := range z {
			if q.Name == k || q.Name == k+"." {
				rr, err := dnsutils.NewRR(q.Name, dnsutils.DnsTypes[dns.TypeA], v)
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

func (h *handler) HandleLocalDnsReq(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg).SetReply(r)
	m.Compress = false

	fp, err := os.Open(h.conf.HostsFile)
	if err != nil {
		log.Printf("failed to open hosts file %s: %v", h.conf.HostsFile, err.Error())
	}
	var table = make(map[string]string)
	if err := json.NewDecoder(fp).Decode(&table); err != nil {
		log.Printf("failed to read hosts file %s: %v", h.conf.HostsFile, err.Error())
	}
	if err := h.HandleLocal(m, table); err != nil {
		log.Println("HandleLocal error:", err.Error())
	}
	if err := w.WriteMsg(m); err != nil {
		log.Println("Error writing reply", err.Error())
	}
}
