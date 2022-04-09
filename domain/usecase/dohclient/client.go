package dohclient

import (
	"context"
	"fmt"

	"github.com/likexian/doh-go"
	"github.com/likexian/doh-go/dns"
	"github.com/pkg/errors"
)

type queryFunc func(context.Context, *doh.DoH, dns.Domain, dns.Type) ([]dns.Answer, error)

// FuncMap maps server.all_types configuration to queryFunc.
var FuncMap = map[bool]queryFunc{
	false: queryType,
	true:  queryAll,
}

func NewDoH() *doh.DoH {
	return doh.Use().EnableCache(true)
}

func query(ctx context.Context, client *doh.DoH, d dns.Domain, t dns.Type) ([]dns.Answer, error) {
	resp, err := client.Query(ctx, dns.Domain(d), t)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("DoH query failed for %s", d))
	}
	return resp.Answer, nil
}

func queryType(ctx context.Context, client *doh.DoH, d dns.Domain, t dns.Type) ([]dns.Answer, error) {
	answers, err := query(ctx, client, d, t)
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func queryAll(ctx context.Context, client *doh.DoH, d dns.Domain, t dns.Type) ([]dns.Answer, error) {
	answers, err := query(ctx, client, d, dns.TypeANY)
	if err != nil {
		return nil, err
	}
	return answers, nil
}
