package dohclient

import (
	"context"
	"fmt"
	"reflect"
	"time"

	dohdns "github.com/likexian/doh-go/dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"

	"github.com/artnoi43/stubborn/domain/usecase"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

type Input struct {
	Ctx      context.Context
	Domain   string
	DnsType  uint16
	AllTypes bool
}

func (c *dohClient) Query(v interface{}) (*usecase.Answer, error) {
	if input, ok := v.(*Input); ok {
		return c.QueryUsecase(input)
	}
	return nil, fmt.Errorf("invalid input type %s - expecting *Input", reflect.TypeOf(v))
}

func (c *dohClient) QueryUsecase(
	input *Input,
) (
	*usecase.Answer,
	error,
) {
	start := time.Now()
	domName := dohdns.Domain(input.Domain)
	tStr := dohdns.Type(dnsutils.DnsTypes[input.DnsType])
	if input.AllTypes {
		tStr = dohdns.TypeANY
	}
	var dohAnswers []dohdns.Answer
	var err error
	dohAnswers, err = query(input.Ctx, c.client, domName, tStr)
	if err != nil {
		return nil, err
	}
	var rrs []dns.RR
	for _, answer := range dohAnswers {
		t, supported := dnsutils.DnsTypes[uint16(answer.Type)]
		if !supported {
			return nil, fmt.Errorf("unsupported DNS record type: %d", answer.Type)
		}
		rr, err := dnsutils.NewRR(answer.Name, t, answer.Data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create new RR")
		}
		rrs = append(rrs, rr)
	}
	return &usecase.Answer{
		RRs: rrs,
		RTT: time.Since(start),
	}, nil
}
