package dnsserver

import "github.com/miekg/dns"

func NewDNSServer(conf *Config) *dns.Server {
	return &dns.Server{
		Addr: conf.Address,
		Net:  conf.Protocol,
	}
}
