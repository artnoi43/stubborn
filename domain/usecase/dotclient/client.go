package dotclient

import (
	"fmt"
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
	return &dotClient{
		conf:   conf,
		client: c,
	}
}

func (c *dotClient) Exchange(
	msg *dns.Msg,
) (
	*dns.Msg,
	time.Duration,
	error,
) {
	s := fmt.Sprintf("%v:%v", c.Addr(), c.Port())
	return c.client.Exchange(msg, s)
}

func (c *dotClient) Addr() string {
	return c.conf.UpStreamIp
}

func (c *dotClient) Port() string {
	return c.conf.UpStreamPort
}
