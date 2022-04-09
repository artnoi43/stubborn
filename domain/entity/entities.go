package entity

import (
	"github.com/artnoi43/mgl/str"
	"github.com/miekg/dns"
)

type Network string
type Outbound string
type CachedAnswer []dns.RR

const (
	LocalNetwork Network = "LOCAL"
	Internet     Network = "INTERNET"

	OutboundDoT Outbound = "DOT"
	OutboundDoH Outbound = "DOH"
)

func (s Network) IsValid() bool {
	switch str.ToUpper(s) {
	case LocalNetwork, Internet:
		return true
	}
	return false
}

func (s Outbound) IsValid() bool {
	switch str.ToUpper(s) {
	case OutboundDoT, OutboundDoH:
		return true
	}
	return false
}
