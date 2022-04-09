package dnsserver

import (
	"github.com/miekg/dns"
)

type DnsServer struct {
	server *dns.Server
}

func NewDNSServer(conf *Config) *DnsServer {
	return &DnsServer{
		server: &dns.Server{
			Addr: conf.Address,
			Net:  "udp",
		},
	}
}

func (s *DnsServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *DnsServer) Shutdown() {
	s.server.Shutdown()
}

func (s *DnsServer) Addr() string {
	return s.server.Addr
}
