package cacher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/artnoi43/stubborn/lib/rediswrapper"
)

type Cacher struct {
	MemCache *cache.Cache
	Redis    *rediswrapper.RedisCli
}

func New(conf *Config, redisCli *rediswrapper.RedisCli) *Cacher {
	j, _ := json.Marshal(conf)
	log.Printf("Cacher configuration:\n%s\n", j)
	c := cache.New(
		time.Duration(conf.Expiration)*time.Second,
		time.Duration(conf.CleanUp)*time.Second,
	)
	return &Cacher{
		MemCache: c,
		Redis:    redisCli,
	}
}
