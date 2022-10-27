package redis

import (
	"github.com/go-redis/redis/v8"
)

const (
	maxRetries = 3
	idleConns  = 4
)

func Client(addr string) (*redis.Client, error) {
	store := redis.NewClient(&redis.Options{
		Addr:         addr,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
	})
	return store, nil
}

func FailoverClient(masterName string, sentinelAddr ...string) (*redis.Client, error) {
	store := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddr,
		MaxRetries:    maxRetries,
		MinIdleConns:  idleConns,
	})
	return store, nil
}
