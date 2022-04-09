package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/likexian/doh-go"
	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/datagateway"
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/domain/usecase/dohclient"
	"github.com/artnoi43/stubborn/domain/usecase/dotclient"
)

type handleFunc func(w dns.ResponseWriter, r *dns.Msg)

type networkHandleFuncMap map[entity.Outbound]map[entity.Network]handleFunc

type Handler interface {
	// HandlerFunc determines which handlerFunc to use
	// based on outbound and network config
	HandlerFunc(entity.Outbound, entity.Network) handleFunc
	Start() error
	Shutdown()
}

// handler implements Handler
// TODO: replace remaining server/client structs with interfaces
type handler struct {
	ctx        context.Context
	conf       *Config
	repository datagateway.AnswerDataGateway
	dnsServer  dnsServer
	dotClient  dotClient
	dohClient  *doh.DoH
	fMap       networkHandleFuncMap
}

func New(ctx context.Context, conf *Config, s dnsServer, c datagateway.AnswerDataGateway) *handler {
	j, _ := json.Marshal(conf)
	log.Printf("DNS server configuration:\n%s\n", j)
	// Base handler (w/o outbound)
	h := &handler{
		ctx:        ctx,
		conf:       conf,
		repository: c,
		dnsServer:  s,
	}
	// Configure outbound
	switch conf.Outbound {
	case entity.OutboundDoT:
		client := dotclient.New(&conf.DoT)
		h.dotClient = client
	case entity.OutboundDoH:
		log.Printf("setting up DoH client")
		dohClient := dohclient.NewDoH()
		h.dohClient = dohClient
	}

	return h.mapFuncs()
}

func (h *handler) Start() error {
	o := h.conf.Outbound
	dns.HandleFunc(".", h.HandlerFunc(o, entity.Internet))
	dns.HandleFunc("local.", h.HandlerFunc(o, entity.LocalNetwork))
	log.Printf("starting stubborn DNS resolver on: %s\n", h.dnsServer.Addr())
	return h.dnsServer.ListenAndServe()
}

func (h *handler) Shutdown() {
	log.Println("shutting down handlers")
	h.dnsServer.Shutdown()
	h.repository.Shutdown()
	log.Println("handlers gracefully shutdowm")
}

func (h *handler) mapFuncs() *handler {
	h.fMap = networkHandleFuncMap{
		entity.OutboundDoH: {
			entity.Internet:     h.HandleWithDoH,
			entity.LocalNetwork: h.HandleLocalDnsReq,
		},
		entity.OutboundDoT: {
			entity.Internet:     h.HandleWithDoT,
			entity.LocalNetwork: h.HandleLocalDnsReq,
		},
	}
	return h
}

func (h *handler) HandlerFunc(o entity.Outbound, n entity.Network) handleFunc {
	return h.fMap[o][n]
}
