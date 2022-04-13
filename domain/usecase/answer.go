package usecase

import (
	"time"

	"github.com/miekg/dns"
)

// Answer represents ALL answers returned to a query message.
// Such a query may contain any valid number of questions
// (i.e. different domain names or record types).
type Answer struct {
	RRs []dns.RR
	RTT time.Duration
}

type CachedAnswer []dns.RR
