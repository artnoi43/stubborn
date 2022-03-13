package rediswrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisCli struct {
	Ctx  context.Context
	Conf *Config
	Cli  *redis.Client
}

func New(ctx context.Context, conf *Config) *RedisCli {
	j, _ := json.Marshal(conf)
	log.Printf("Redis client configuration:\n%s\n", j)
	cli := redis.NewClient(&redis.Options{DB: conf.DB})
	return &RedisCli{
		Ctx:  ctx,
		Conf: conf,
		Cli:  cli,
	}
}

func (r *RedisCli) HSet(key, field string, v []string) error {
	log.Println("HSET Redis", key, field, v)
	valJson, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HSET: failed to marshal JSON for %s[%s]:%s", key, field, v))
	}
	vStr := string(valJson)
	if _, err := r.Cli.HSet(r.Ctx, key, field, vStr).Result(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("HSET failed for %s[%s]:%s", key, field, v))
	}
	return nil
}

func (r *RedisCli) HGet(key, field string) ([]string, error) {
	v, err := r.Cli.HGet(r.Ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.Wrap(err, fmt.Sprintf("Redis cache missed for key %s:%s", key, field))
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("HGET failed for %s:%s", key, field))
	}
	log.Printf("Redis cache hit: %s:%s", key, field)
	var answers []string
	if err := json.Unmarshal([]byte(v), &answers); err != nil {
		return nil, fmt.Errorf("interface{} conversion to []string failed %v", v)
	}
	return answers, nil
}
