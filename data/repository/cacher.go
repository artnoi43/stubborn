package repository

import (
	"encoding/json"
	"log"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/artnoi43/stubborn/domain/usecase"
)

type Cacher interface {
	GetRR(k usecase.CacheKey) (usecase.CachedAnswer, error)
	SetRR(k usecase.CacheKey, v usecase.CachedAnswer)
	SetMapRR(map[usecase.CacheKey]usecase.CachedAnswer)
	Flush()
}

type memCache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
	Flush()
}

type cacher struct {
	conf     *Config
	memCache memCache
}

func New(conf *Config) *cacher {
	j, _ := json.Marshal(conf)
	log.Printf("cacher configuration:\n%s\n", j)
	c := cache.New(
		time.Duration(conf.Expiration)*time.Second,
		time.Duration(conf.CleanUp)*time.Second,
	)
	return &cacher{
		memCache: c,
		conf:     conf,
	}
}

func (c *cacher) Flush() {
	c.memCache.Flush()
}
