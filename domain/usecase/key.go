package usecase

import (
	"fmt"

	"github.com/artnoi43/stubborn/lib/dnsutils"
)

type CacheKey struct {
	Dom  string
	Type uint16
}

func NewCacheKey(dom string, dnsType uint16) CacheKey {
	return CacheKey{
		Dom:  dom,
		Type: dnsType,
	}
}

func (k CacheKey) String() string {
	typeStr, valid := dnsutils.DnsTypes[k.Type]
	if valid {
		return fmt.Sprintf("%s:%s", typeStr, k.Dom)
	}
	return fmt.Sprintf("%d:%s", k.Type, k.Dom)
}
