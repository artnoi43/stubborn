package dnsutils

import (
	"fmt"
)

type Key struct {
	Dom  string
	Type uint16
}

func NewKey(dom string, dnsType uint16) Key {
	return Key{
		Dom:  dom,
		Type: dnsType,
	}
}

func (k Key) String() string {
	typeStr, valid := DnsTypes[k.Type]
	if valid {
		return fmt.Sprintf("%s:%s", typeStr, k.Dom)
	}
	return fmt.Sprintf("%d:%s", k.Type, k.Dom)
}
