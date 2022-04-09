package repository

import (
	"encoding/json"
	"log"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/lib/dnsutils"
)

type Cacher interface {
	GetRR(k dnsutils.Key) (entity.CachedAnswer, error)
	SetRR(k dnsutils.Key, v entity.CachedAnswer)
	SetMapRR(map[dnsutils.Key]entity.CachedAnswer)
	Flush()
}

type memCache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
	Flush()
}

type cacher struct {
	conf     *Config
	MemCache memCache
}

func New(conf *Config) *cacher {
	j, _ := json.Marshal(conf)
	log.Printf("cacher configuration:\n%s\n", j)
	c := cache.New(
		time.Duration(conf.Expiration)*time.Second,
		time.Duration(conf.CleanUp)*time.Second,
	)
	return &cacher{
		MemCache: c,
		conf:     conf,
	}
}

func (c *cacher) Flush() {
	c.MemCache.Flush()
}
