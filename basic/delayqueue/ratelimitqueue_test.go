package queue

import (
	"testing"
	"time"

	"singer.com/util/clock"
)

func TestRateLimitQueue(t *testing.T) {
	limiter := NewExponentialFailureRateLimiter[string](1*time.Second, 10*time.Second)
	queue := NewRateLimitingQueue(limiter).(*rateLimitingType[string])
	real := clock.RealClock{}
	delayingQueue := &delayingType[string]{
		BaseQueue:       NewBaseQueue[string](),
		clock:           real,
		heartbeat:       real.NewTicker(maxWait),
		stopCh:          make(chan struct{}),
		waitingForAddCh: make(chan *waitFor[string], 1000),
	}
	queue.DelayingQueue = delayingQueue

	delay := queue.DelayingQueue.(*delayingType[string])
	queue.AddRateLimited("one")
	waitEntry := <-delay.waitingForAddCh
	t.Log(waitEntry.readyAt)

	queue.AddRateLimited("one")
	waitEntry = <-delay.waitingForAddCh
	t.Log(waitEntry.readyAt)

	queue.AddRateLimited("one")
	waitEntry = <-delay.waitingForAddCh
	t.Log(waitEntry.readyAt)

	queue.AddRateLimited("two")
	waitEntry = <-delay.waitingForAddCh
	t.Log(waitEntry.readyAt)

	queue.Forget("one")
	t.Log(queue.NumRequeues("one"))
}
