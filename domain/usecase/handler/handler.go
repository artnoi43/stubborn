package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/datagateway"
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dohclient"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dotclient"
	"github.com/artnoi43/stubborn/domain/usecase/clients/localclient"
)

type handleFunc func(w dns.ResponseWriter, r *dns.Msg)

type networkHandleFuncMap map[entity.Outbound]map[entity.Network]handleFunc

type Handler interface {
	// HandlerFunc determines which handlerFunc to use
	// based on outbound and network config
	// HandlerFunc(entity.Outbound, entity.Network) handleFunc
	Start() error
	Shutdown()
}

// handler implements Handler, and is created once.
type handler struct {
	ctx        context.Context
	conf       *Config
	repository datagateway.AnswerDataGateway
	dnsServer  dnsServer
	dnsClient  dnsClient
	jsonClient dnsClient
}

func New(ctx context.Context, conf *Config, s dnsServer, c datagateway.AnswerDataGateway) *handler {
	j, _ := json.MarshalIndent(conf, "", "  ")
	log.Printf("stubborn configuration:\n%s\n", j)
	// Base handler (w/o outbound)
	localNetworkClient := localclient.New(conf.HostsFile)
	h := &handler{
		ctx:        ctx,
		conf:       conf,
		repository: c,
		dnsServer:  s,
		jsonClient: localNetworkClient,
	}
	// Configure outbound
	switch conf.EntityOutbound {
	case entity.OutboundDoT:
		log.Println("setting up DoT client")
		h.dnsClient = dotclient.New(&conf.DoT)
	case entity.OutboundDoH:
		log.Println("setting up DoH client")
		h.dnsClient = dohclient.New(conf.DoH.Upstream)
	}
	return h
}

func (h *handler) Start() error {
	dns.HandleFunc(".", h.handle(false))
	dns.HandleFunc("local.", h.handle(true))
	log.Printf("starting stubborn DNS resolver on: %s\n", h.dnsServer.Addr())
	return h.dnsServer.ListenAndServe()
}

func (h *handler) Shutdown() {
	log.Println("shutting down handlers")
	h.dnsServer.Shutdown()
	if h.dnsClient != nil {
		h.dnsClient.Shutdown()
	}
	h.repository.Shutdown()
	log.Println("handlers gracefully shutdowm")
}
