package handler

import (
	"time"

	"github.com/miekg/dns"
)

type dnsServer interface {
	ListenAndServe() error
	Shutdown()
	Addr() string
}

type dotClient interface {
	Exchange(*dns.Msg) (*dns.Msg, time.Duration, error)

	Addr() string
	Port() string
}
