package circulequeue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPushQueue(t *testing.T) {
	q := New[int](3)
	assert.True(t, q.Empty())
	assert.Equal(t, q.Size(), 0)

	q.Push(1)
	t.Log("q push 1")
	q.Push(2)
	t.Log("q push 2")
	q.Push(3)
	t.Log("q push 3")

	go func() {
		time.Sleep(time.Second * 3)
		i := q.Pop()
		t.Log("q pop:", i)
	}()

	q.Push(4)
	t.Log("q push 4")
	t.Log("q size:", q.Size())
	t.Log("q string:", q.String())
	t.Log("q value:", q.Values())
}

func TestPop(t *testing.T) {
	q := New[int](3)

	go func() {
		time.Sleep(2 * time.Second)
		q.Push(10)
	}()

	i := q.Pop()
	t.Log("pop:", i)
	assert.True(t, q.Empty())
	assert.Equal(t, q.Size(), 0)
}
