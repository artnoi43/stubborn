package entity

// Network is an enum used during path registration.
type Network int

const (
	InvalidNetwork Network = iota
	LocalNetwork
	Internet
)

var (
	validNetworks = []Network{
		LocalNetwork,
		Internet,
	}
)

func (t Network) IsValid() bool {
	switch t {
	case LocalNetwork, Internet:
		return true
	}
	return false
}
