package redlock

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	Addr string
}

type RedLock struct {
	Pools []*redis.Client
}

type Mutex struct {
	key   string
	Pools []*redis.Client
	ttl   time.Duration
	value string
}

func New(options ...*Options) *RedLock {

	pools := make([]*redis.Client, 0, len(options))
	for _, option := range options {
		pools = append(pools, NewPool(option))
	}

	return &RedLock{
		Pools: pools,
	}
}

func NewPool(options *Options) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: options.Addr,
	})
}
