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
	addr   string
	port   string
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
	s := fmt.Sprintf("%v:%v", c.addr, c.port)
	return c.client.Exchange(msg, s)
}

func (c *dotClient) Addr() string {
	return c.addr
}

func (c *dotClient) Port() string {
	return c.port
}
