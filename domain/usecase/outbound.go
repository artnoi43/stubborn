package usecase

import "github.com/artnoi43/mgl/str"

// Outbound is an enum used to determine outbound DNS protocol
// Currently, 2 protocols (DoT, DoH) are supported.
type Outbound int

const (
	InvalidOutbound Outbound = iota
	OutboundDoT
	OutboundDoH
)

// validOutbounds is not used for IsValid(), but instead to
// convert a string to Outbound based on the String() method.
var (
	validOutbounds = []Outbound{
		OutboundDoT,
		OutboundDoH,
	}
)

func (o Outbound) String() string {
	switch o {
	case OutboundDoT:
		return "DOT"
	case OutboundDoH:
		return "DOH"
	}
	return ""
}

func (o Outbound) IsValid() bool {
	switch o {
	case OutboundDoT, OutboundDoH:
		return true
	}
	return false
}

// OutboundFromString returns Outbound (int) based on
// its String() method and the validOutbounds slice.
func OutboundFromString(s string) Outbound {
	for _, v := range validOutbounds {
		if str.ToUpper(s) == v.String() {
			return v
		}
	}
	return InvalidOutbound
}
