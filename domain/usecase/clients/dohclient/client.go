package dohclient

import (
	"context"
	"log"

	"github.com/artnoi43/mgl/str"
	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
	"github.com/pkg/errors"
)

type queryFunc func(context.Context, *doh.DoH, dohdns.Domain, dohdns.Type) ([]dohdns.Answer, error)
type dohClient struct {
	client *doh.DoH
}

func New(dohProviders []provider) *dohClient {
	var providers []int
	var providersStr []provider
	for _, p := range dohProviders {
		// Using provider.IsValid will also run similar code
		// so I dont call IsValid() here
		provider, valid := providerMap[str.ToUpper(p)]
		if valid {
			providersStr = append(providersStr, p)
			providers = append(providers, int(provider))
		}
	}
	log.Printf("DoH providers: %v", providersStr)
	return &dohClient{
		client: doh.Use(providers...),
	}
}

func query(ctx context.Context, client *doh.DoH, d dohdns.Domain, t dohdns.Type) ([]dohdns.Answer, error) {
	resp, err := client.Query(ctx, d, t)
	if err != nil {
		return nil, errors.Wrapf(err, "DoH query failed for %s", d)
	}
	return resp.Answer, nil
}

func (c *dohClient) Shutdown() {
	c.client.Close()
}
