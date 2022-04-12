package handler

import (
	"context"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dohclient"
)

func (h *handler) inputsFromQuestion(q *dns.Question) *clientInputs {
	switch h.conf.EntityOutbound {
	case entity.OutboundDoT:
		m := new(dns.Msg).SetQuestion(q.Name, q.Qtype)
		return &clientInputs{
			dotInput:  m,
			jsonInput: q.Name,
		}
	case entity.OutboundDoH:
		return &clientInputs{
			dohInput: &dohclient.Input{
				Ctx:      context.Background(),
				Domain:   q.Name,
				DnsType:  q.Qtype,
				AllTypes: h.conf.QueryAllTypes,
			},
			jsonInput: q.Name,
		}
	}
	return nil
}

func (h *handler) selectInput(input *clientInputs, isLocal bool) interface{} {
	var thisInput interface{}
	switch h.conf.EntityOutbound {
	case entity.OutboundDoT:
		thisInput = input.dotInput
	case entity.OutboundDoH:
		thisInput = input.dohInput
	}
	return thisInput
}
