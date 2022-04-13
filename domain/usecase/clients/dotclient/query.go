package dotclient

import (
	"fmt"
	"reflect"

	"github.com/artnoi43/stubborn/domain/usecase"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func (c *dotClient) Query(v interface{}) (*usecase.Answer, error) {
	if msg, ok := v.(*dns.Msg); ok {
		return c.QueryUsecase(msg)
	}
	return nil, fmt.Errorf("invalid input type %s - expecting *dns.Msg", reflect.TypeOf(v))
}

func (c *dotClient) QueryUsecase(
	msg *dns.Msg,
) (
	*usecase.Answer,
	error,
) {
	s := fmt.Sprintf("%v:%v", c.conf.UpStreamIp, c.conf.UpStreamPort)
	ans, rtt, err := c.client.Exchange(msg, s)
	if err != nil {
		return nil, errors.Wrapf(err, "DoT outbound failed for %s", msg.String())
	}
	return &usecase.Answer{
		RRs: ans.Answer,
		RTT: rtt,
	}, nil
}
