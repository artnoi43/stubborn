package handler

import (
	"log"
	"sync"

	"github.com/miekg/dns"
	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/domain/usecase"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dohclient"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

// answerMap maps cache CacheKey to answers (CacheKey: "domain:type", val: "val")
type answerMap map[usecase.CacheKey]usecase.CachedAnswer

// clientInputs stores inputs for both
// dotClient.QueryUsecase() and dohClient.QueryUsecase
type clientInputs struct {
	dotInput  *dns.Msg
	dohInput  *dohclient.Input
	jsonInput string
}

// handle handles incoming DNS queries. It takes in isLocal (bool)
// to change behavior (local network vs internet)
func (h *handler) handle(isLocal bool) handleFunc {
	var outboundClient dnsClient
	if isLocal {
		outboundClient = h.jsonClient
	} else {
		outboundClient = h.dnsClient
	}

	return func(w dns.ResponseWriter, r *dns.Msg) {
		var m answerMap // nil, but will be 'made' when cache missed
		errChan := make(chan error)
		var wg sync.WaitGroup
		for _, question := range r.Question {
			wg.Add(1)
			go func(q dns.Question) {
				// msg is created in each goroutine to avoid race condition
				msg := new(dns.Msg).SetReply(r)
				msg.Compress = true
				msg.RecursionDesired = true
				msg.RecursionAvailable = true
				defer wg.Done()
				tString, supported := dnsutils.DnsTypes[q.Qtype]
				if !supported {
					// Without string types, dnsutils.NewRR will fail
					return
				}
				k := usecase.NewCacheKey(q.Name, q.Qtype)
				cached, found := h.repository.Get(k)
				if found {
					log.Printf("cache hit for \"%s\" [%s]", k.String(), tString)
					msg.Answer = append(msg.Answer, cached...)
				} else {
					if m == nil {
						m = make(answerMap)
					}
					log.Printf("cache missed for \"%s\" [%s]", k.String(), tString)
					// Outgoing query starts here
					var properInput interface{}
					if isLocal {
						properInput = q.Name
					} else {
						inputs := h.inputsFromQuestion(&q)
						properInput = h.selectInput(inputs, isLocal)
					}
					// To use your own outbound client,
					// implement the dnsClient interface.
					ans, err := outboundClient.Query(properInput)
					if err != nil {
						errChan <- errors.Wrap(err, "to connect upstream")
						return
					}
					// Query done
					log.Printf("question: %s rtt: %v", q.String(), ans.RTT)
					if ans != nil {
						msg.Answer = append(msg.Answer, ans.RRs...)
						m[k] = append(m[k], ans.RRs...)
					}
				}
				if err := w.WriteMsg(msg); err != nil {
					log.Println("reply error", err.Error())
				}
			}(question)
			// Collects and prints errors
			go func() {
				if err := <-errChan; err != nil {
					log.Println("handle error:", err.Error())
					return
				}
			}()
		}
		wg.Wait()
		close(errChan)
		if m != nil {
			h.repository.SetMap(m)
		}
	}
}
