package cacher

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Key struct {
	Dom  string
	Type string
	TTL  int
}

func NewKey(d, t string, ttl int) Key {
	if ttl < 0 {
		ttl = 60
	}
	return AppendRootServerKey(Key{
		Dom:  d,
		Type: t,
		TTL:  ttl,
	})
}

func (k Key) MemCacheKey() string {
	return fmt.Sprintf("%s:%s", k.Dom, k.Type)
}

func (k Key) RedisKey() string {
	return k.Dom
}

func (c *Cacher) HGet(k Key) ([]string, error) {
	var memCacheMissed bool
	v, found := c.MemCache.Get(k.MemCacheKey())
	if !found {
		var err error
		memCacheMissed = true
		log.Println("MemCache missed", k.MemCacheKey())
		v, err = c.Redis.HGet(k.RedisKey(), k.Type)
		if errors.Is(err, redis.Nil) {
			log.Println("Redis cache missed", k.RedisKey())
		}
		if err != nil {
			return nil, err
		}
	}
	log.Println("MemCache hit:", k.MemCacheKey())
	answers, ok := v.([]string)
	if !ok {
		return nil, fmt.Errorf("interface{} conversion to string failed for %v", v)
	}
	// Write to cache again
	if memCacheMissed {
		if err := c.HSet(k, answers); err != nil {
			log.Println(err.Error())
		}
	}
	return answers, nil
}

func (c *Cacher) HSet(k Key, v interface{}) error {
	dur := time.Duration(k.TTL) * time.Second
	log.Println("HSET MemCache", k.Dom, k.Type, v, "ttl", k.TTL)
	c.MemCache.Set(k.MemCacheKey(), v, dur)
	vStr, ok := v.([]string)
	if !ok {
		return fmt.Errorf("interface{} conversion to []string failed %v", v)
	}
	return c.Redis.HSet(k.RedisKey(), k.Type, vStr)
}
