package dotclient

import (
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

type dotClient struct {
	conf   *Config
	client *dns.Client
}

func New(conf *Config) *dotClient {
	c := new(dns.Client)
	c.Net = "tcp-tls"
	c.Dialer = &net.Dialer{
		Timeout: time.Duration(conf.UpStreamTimeout) * time.Second,
	}
	log.Printf("DoT upstream: %s:%s", conf.UpStreamIp, conf.UpStreamPort)
	return &dotClient{
		conf:   conf,
		client: c,
	}
}

// Implements handler.dnsClient
func (c *dotClient) Shutdown() {}
