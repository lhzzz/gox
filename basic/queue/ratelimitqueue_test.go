package queue

import (
	"testing"
	"time"

	"singer.com/util/clock"
)

func TestRateLimitQueue(t *testing.T) {
	limiter := NewExponentialFailureRateLimiter(1*time.Second, 10*time.Second)
	queue := NewRateLimitingQueue(limiter).(*rateLimitingType)
	real := clock.RealClock{}
	delayingQueue := &delayingType{
		BaseQueue:       NewBaseQueue(),
		clock:           real,
		heartbeat:       real.NewTicker(maxWait),
		stopCh:          make(chan struct{}),
		waitingForAddCh: make(chan *waitFor, 1000),
	}
	queue.DelayingQueue = delayingQueue

	delay := queue.DelayingQueue.(*delayingType)
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
