package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisLock(t *testing.T) {
	c, err := Client("127.0.0.1:6379")
	assert.Nil(t, err)

	key := "mutex"
	r1 := NewRedisLock(c, key, 5*time.Second)
	r2 := NewRedisLock(c, key, 5*time.Second)

	succ, err := r1.Lock() //intance 1 lock
	assert.Nil(t, err)
	assert.True(t, succ)

	succ, err = r2.Unlock() //can't unlock same key by other instance
	assert.Nil(t, err)
	assert.False(t, succ)

	succ, err = r2.Lock() // can't lock if other has locked
	assert.Nil(t, err)
	assert.False(t, succ)

	succ, err = r1.Unlock() //intance 1 unlock
	assert.Nil(t, err)
	assert.True(t, succ)

	succ, err = r2.Lock() // intance 2 lock
	assert.Nil(t, err)
	assert.True(t, succ)

	succ, err = r1.Expire(10) // can't expire same key if other has locked
	assert.Nil(t, err)
	assert.False(t, succ)

	succ, err = r2.Expire(10) // expire key
	assert.Nil(t, err)
	assert.True(t, succ)

	succ, err = r2.Unlock() // unlock
	assert.Nil(t, err)
	assert.True(t, succ)
}
