package redis

import (
	"context"
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	unlockScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`
	expireScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("EXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type redisLock struct {
	client *redis.Client
	expire time.Duration
	key    string
	value  string
}

type RedisLock interface {
	Lock() (success bool, err error)
	LockCtx(ctx context.Context) (success bool, err error)
	Unlock() (success bool, err error)
	UnlockCtx(ctx context.Context) (success bool, err error)
	//Expire is in order to defer the deadline of the lock when the process is not finished in lock time
	Expire(seconds int) (success bool, err error)
	ExpireCtx(ctx context.Context, seconds int) (success bool, err error)
}

/*
Parameter:
 client : redis client instance
 key    : lock name
 expire : the lock will exist in the expire
*/
func NewRedisLock(client *redis.Client, key string, expire time.Duration) RedisLock {
	return &redisLock{
		client: client,
		key:    key,
		expire: expire,
		value:  genValue(),
	}
}

func (rl *redisLock) Lock() (success bool, err error) {
	return rl.LockCtx(context.Background())
}

func (rl *redisLock) LockCtx(ctx context.Context) (success bool, err error) {
	b, e := rl.client.SetNX(ctx, rl.key, rl.value, rl.expire).Result()
	if e != nil {
		return false, e
	}
	return b, nil
}

func (rl *redisLock) Unlock() (success bool, err error) {
	return rl.UnlockCtx(context.Background())
}

func (rl *redisLock) UnlockCtx(ctx context.Context) (success bool, err error) {
	resp, e := rl.client.Eval(ctx, unlockScript, []string{rl.key}, rl.value).Result()
	if e != nil {
		return false, e
	}
	val, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return 1 == val, nil
}

func (rl *redisLock) Expire(seconds int) (success bool, err error) {
	return rl.ExpireCtx(context.Background(), seconds)
}

func (rl *redisLock) ExpireCtx(ctx context.Context, seconds int) (success bool, err error) {
	resp, e := rl.client.Eval(ctx, expireScript, []string{rl.key}, rl.value, seconds).Result()
	if e != nil {
		return false, e
	}
	val, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return 1 == val, nil
}

func genValue() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
