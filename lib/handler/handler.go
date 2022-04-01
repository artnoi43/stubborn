package handler

import (
	"context"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/lib/enums"
)

type handleFunc func(w dns.ResponseWriter, r *dns.Msg)

type networkHandleFuncMap map[enums.Network]handleFunc

type Handler interface {
	HandlerFunc(enums.Network) handleFunc
	Start(string) error
	Shutdown(context.Context)
}

func (h *handler) mapFuncs() *handler {
	h.fMap = networkHandleFuncMap{
		enums.Internet:     h.HandleDnsReq,
		enums.LocalNetwork: h.HandleLocalDnsReq,
	}
	return h
}

func (h *handler) HandlerFunc(n enums.Network) handleFunc {
	return h.fMap[n]
}
