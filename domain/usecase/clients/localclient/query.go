package localclient

import (
	"fmt"
	"reflect"
	"time"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/usecase"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

func (j *jsonClient) Query(v interface{}) (*usecase.Answer, error) {
	if name, ok := v.(string); ok {
		return j.QueryUsecase(name)
	}
	return nil, fmt.Errorf("invalid input type %s - expecting string", reflect.TypeOf(v))
}

func (j *jsonClient) QueryUsecase(name string) (*usecase.Answer, error) {
	start := time.Now()
	var rrs []dns.RR
	for addr, hostnames := range localNetworkTable {
		for _, hostname := range hostnames {
			if name == hostname || name == hostname+"." || name == hostname+".local." {
				rr, err := dnsutils.NewRR(
					name,
					dnsutils.DnsTypes[dns.TypeA],
					addr,
				)
				if err != nil {
					return nil, err
				}
				rrs = append(rrs, rr)
			}
		}
	}
	return &usecase.Answer{
		RRs: rrs,
		RTT: time.Since(start),
	}, nil
}
