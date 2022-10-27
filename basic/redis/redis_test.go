package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	r, err := Client("127.0.0.1:6379")
	assert.Nil(t, err)

	_, err = r.Ping(context.Background()).Result()
	assert.Nil(t, err)
}

func TestFailoverClient(t *testing.T) {
	r, err := FailoverClient("mymaster", "127.0.0.1:26379")
	assert.Nil(t, err)

	_, err = r.Ping(context.Background()).Result()
	assert.Nil(t, err)

	s, err := r.Keys(context.Background(), "*").Result()
	t.Log(s, err)
}
