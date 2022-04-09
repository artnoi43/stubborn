package repository

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

func (c *cacher) GetRR(k dnsutils.Key) (entity.CachedAnswer, error) {
	v, found := c.MemCache.Get(k.String())
	if !found {
		return nil, ErrCacheMissed
	}
	log.Println("MemCache hit:", k.String())
	answers, ok := v.(entity.CachedAnswer)
	if !ok {
		return nil, fmt.Errorf("type assertion failed to string failed for %v", v)
	}
	return answers, nil
}

func (c *cacher) SetMsg(k dnsutils.Key, msg dns.Msg) {
	panic("not implemented")
}

func (c *cacher) SetRR(k dnsutils.Key, v entity.CachedAnswer) {
	ttl := time.Duration(c.conf.Expiration) * time.Second
	log.Println("SET MemCache for key", k.String(), "ttl", ttl)
	c.MemCache.Set(k.String(), v, ttl)
}

func (c *cacher) SetMapRR(m map[dnsutils.Key]entity.CachedAnswer) {
	var wg sync.WaitGroup
	for k, v := range m {
		wg.Add(1)
		go func(key dnsutils.Key, value entity.CachedAnswer) {
			defer wg.Done()
			c.SetRR(key, value)
		}(k, v)
	}
	wg.Wait()
}
