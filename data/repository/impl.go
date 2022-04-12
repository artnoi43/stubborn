package repository

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"

	"github.com/artnoi43/stubborn/domain/usecase"
)

func (c *cacher) GetRR(k usecase.CacheKey) (usecase.CachedAnswer, error) {
	v, found := c.memCache.Get(k.String())
	if !found {
		return nil, ErrCacheMissed
	}
	log.Println("memCache hit:", k.String())
	answers, ok := v.(usecase.CachedAnswer)
	if !ok {
		return nil, fmt.Errorf("type assertion failed to string failed for %v", v)
	}
	return answers, nil
}

func (c *cacher) SetMsg(k usecase.CacheKey, msg dns.Msg) {
	panic("not implemented")
}

func (c *cacher) SetRR(k usecase.CacheKey, v usecase.CachedAnswer) {
	ttl := time.Duration(c.conf.Expiration) * time.Second
	log.Println("memCache SET for key", k.String(), "ttl", ttl)
	c.memCache.Set(k.String(), v, ttl)
}

func (c *cacher) SetMapRR(m map[usecase.CacheKey]usecase.CachedAnswer) {
	var wg sync.WaitGroup
	for k, v := range m {
		wg.Add(1)
		go func(key usecase.CacheKey, value usecase.CachedAnswer) {
			defer wg.Done()
			c.SetRR(key, value)
		}(k, v)
	}
	wg.Wait()
}
