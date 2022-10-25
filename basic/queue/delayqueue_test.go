package queue

import (
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	dq := NewDelayingQueue()

	dq.Add(1)
	dq.AddAfter(2, 3*time.Second)
	dq.Add(3)

	go func() {
		it, _ := dq.Get()
		t.Log(it)

		it, _ = dq.Get()
		t.Log(it)

		cur := time.Now()
		it, _ = dq.Get()
		cost := time.Since(cur)
		t.Log(it, cost)
	}()

	time.Sleep(5 * time.Second)
}
