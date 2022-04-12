package dohclient

import (
	"github.com/artnoi43/mgl/str"
	"github.com/likexian/doh-go"
)

type Config struct {
	Upstream []provider `mapstructure:"upstream" json:"upstream"`
}

var (
	providerMap = map[provider]int{
		cloudflare: doh.CloudflareProvider,
		quad9:      doh.Quad9Provider,
		google:     doh.GoogleProvider,
		dnspod:     doh.DNSPodProvider,
	}
)

func (p provider) IsValid() bool {
	_, isValid := providerMap[str.ToUpper(p)]
	return isValid
}
