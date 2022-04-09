package dnsutils

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// NewRRA returns new RR for supported DNS record types (in dnstypes.go)
func NewRR(dom string, typeStr DnsType, v string) (dns.RR, error) {
	rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", dom, typeStr, v))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("new RR failed for %s:%s", dom, v))
	}
	return rr, nil
}

func AppendRootServer(s string) string {
	if s[len(s)-1] != '.' {
		return s + "."
	}
	return s
}

func AppendRootServerKey(k Key) Key {
	if k.Dom[len(k.Dom)-1] != '.' {
		k.Dom += "."
	}
	return k
}
