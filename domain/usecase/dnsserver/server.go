package dnsserver

import (
	"github.com/miekg/dns"
)

// dnsServer implements interface handler.dnsServer
type dnsServer struct {
	server *dns.Server
}

func New(conf *Config) *dnsServer {
	return &dnsServer{
		server: &dns.Server{
			Addr: conf.Address,
			Net:  "udp",
		},
	}
}

func (s *dnsServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *dnsServer) Addr() string {
	return s.server.Addr
}

func (s *dnsServer) Shutdown() {
	s.server.Shutdown()
}
