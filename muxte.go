package redlock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

const (
	DefaultSinglePoolTimeout = 50 * time.Millisecond
)

func (r *RedLock) NewMutex(key string, ttl time.Duration) *Mutex {
	return &Mutex{
		key:   key,
		Pools: r.Pools,
		ttl:   ttl,
	}
}

func (m *Mutex) Lock() error {

	ctx := context.Background()

	m.value = m.getRand()

	start := time.Now()
	count := 0
	for _, pool := range m.Pools {
		ctx, cancel := context.WithTimeout(ctx, DefaultSinglePoolTimeout)
		res := pool.SetNX(ctx, m.key, m.value, m.ttl)
		if res.Val() {
			count++
		}
		cancel()
	}

	end := time.Now()

	if end.Sub(start) < m.ttl {
		if count > len(m.Pools)/2 {
			return nil
		}
	}

	return errors.New("lock failed")
}

func (m *Mutex) getRand() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

var unlockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
	return redis.call("del",KEYS[1])
elseif val == false then
	return -1 
else
	return 0
end`

func (m *Mutex) Unlock() error {
	ctx := context.Background()
	for _, pool := range m.Pools {
		res := pool.Eval(ctx, unlockScript, []string{m.key}, m.value)
		if res.Val() != int64(0) {
			return errors.New("lock expired")
		}
	}
	return nil
}
