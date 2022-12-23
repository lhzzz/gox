package queue

import (
	"container/heap"
	"sync"
	"time"

	"singer.com/util/clock"
	"singer.com/util/recovery"
)

// DelayingQueue is an Interface that can Add an item at a later time. This makes it easier to
// requeue items after failures without ending up in a hot-loop.
type DelayingQueue[T comparable] interface {
	BaseQueue[T]
	// AddAfter adds an item to the workqueue after the indicated duration has passed
	AddAfter(item T, duration time.Duration)
}

// delayingType wraps an Interface and provides delayed re-enquing
type delayingType[T comparable] struct {
	BaseQueue[T]

	// clock tracks time for delayed firing
	clock clock.Clock

	// stopCh lets us signal a shutdown to the waiting loop
	stopCh chan struct{}
	// stopOnce guarantees we only signal shutdown a single time
	stopOnce sync.Once

	// heartbeat ensures we wait no more than maxWait before firing
	heartbeat clock.Ticker

	// waitingForAddCh is a buffered channel that feeds waitingForAdd
	waitingForAddCh chan *waitFor[T]
}

// waitFor holds the data to add and the time it should be added
type waitFor[T comparable] struct {
	data    T
	readyAt time.Time
	// index in the priority queue (heap)
	index int
}

const maxWait = 10 * time.Second

// NewDelayingQueue constructs a new workqueue with delayed queuing ability.
// NewDelayingQueue does not emit metrics. For use with a MetricsProvider, please use
// NewNamedDelayingQueue instead.
func NewDelayingQueue[T comparable]() DelayingQueue[T] {
	return newDelayingQueue(clock.RealClock{}, NewBaseQueue[T]())
}

func newDelayingQueue[T comparable](clock clock.WithTicker, q BaseQueue[T]) *delayingType[T] {
	ret := &delayingType[T]{
		BaseQueue:       q,
		clock:           clock,
		heartbeat:       clock.NewTicker(maxWait),
		stopCh:          make(chan struct{}),
		waitingForAddCh: make(chan *waitFor[T], 1000),
	}

	go ret.waitingLoop()
	return ret
}

// ShutDown stops the queue. After the queue drains, the returned shutdown bool
// on Get() will be true. This method may be invoked more than once.
func (q *delayingType[T]) ShutDown() {
	q.stopOnce.Do(func() {
		q.BaseQueue.ShutDown()
		close(q.stopCh)
		q.heartbeat.Stop()
	})
}

// AddAfter adds the given item to the work queue after the given delay
func (q *delayingType[T]) AddAfter(item T, duration time.Duration) {
	// don't add if we're already shutting down
	if q.ShuttingDown() {
		return
	}

	// immediately add things with no delay
	if duration <= 0 {
		q.Add(item)
		return
	}

	select {
	case <-q.stopCh:
		// unblock if ShutDown() is called
	case q.waitingForAddCh <- &waitFor[T]{data: item, readyAt: q.clock.Now().Add(duration)}:
	}
}

// waitingLoop runs until the workqueue is shutdown and keeps a check on the list of items to be added.
func (q *delayingType[T]) waitingLoop() {
	defer recovery.HandleCrash()

	// Make a placeholder channel to use when there are no items in our list
	never := make(<-chan time.Time)

	// Make a timer that expires when the item at the head of the waiting queue is ready
	var nextReadyAtTimer clock.Timer

	waitingForQueue := &waitForPriorityQueue[T]{}
	heap.Init(waitingForQueue)

	waitingEntryByData := map[T]*waitFor[T]{}

	for {
		if q.BaseQueue.ShuttingDown() {
			return
		}

		now := q.clock.Now()

		// Add ready entries
		for waitingForQueue.Len() > 0 {
			entry := waitingForQueue.Peek().(*waitFor[T])
			if entry.readyAt.After(now) {
				break
			}

			entry = heap.Pop(waitingForQueue).(*waitFor[T])
			q.Add(entry.data)
			delete(waitingEntryByData, entry.data)
		}

		// Set up a wait for the first item's readyAt (if one exists)
		nextReadyAt := never
		if waitingForQueue.Len() > 0 {
			if nextReadyAtTimer != nil {
				nextReadyAtTimer.Stop()
			}
			entry := waitingForQueue.Peek().(*waitFor[T])
			nextReadyAtTimer = q.clock.NewTimer(entry.readyAt.Sub(now))
			nextReadyAt = nextReadyAtTimer.C()
		}

		select {
		case <-q.stopCh:
			return

		case <-q.heartbeat.C():
			// continue the loop, which will add ready items

		case <-nextReadyAt:
			// continue the loop, which will add ready items

		case waitEntry := <-q.waitingForAddCh:
			if waitEntry.readyAt.After(q.clock.Now()) {
				insert(waitingForQueue, waitingEntryByData, waitEntry)
			} else {
				q.Add(waitEntry.data)
			}

			drained := false
			for !drained {
				select {
				case waitEntry := <-q.waitingForAddCh:
					if waitEntry.readyAt.After(q.clock.Now()) {
						insert(waitingForQueue, waitingEntryByData, waitEntry)
					} else {
						q.Add(waitEntry.data)
					}
				default:
					drained = true
				}
			}
		}
	}
}

// insert adds the entry to the priority queue, or updates the readyAt if it already exists in the queue
func insert[T comparable](q *waitForPriorityQueue[T], knownEntries map[T]*waitFor[T], entry *waitFor[T]) {
	// if the entry already exists, update the time only if it would cause the item to be queued sooner
	existing, exists := knownEntries[entry.data]
	if exists {
		if existing.readyAt.After(entry.readyAt) {
			existing.readyAt = entry.readyAt
			heap.Fix(q, existing.index)
		}

		return
	}

	heap.Push(q, entry)
	knownEntries[entry.data] = entry
}

//implement heap
type waitForPriorityQueue[T comparable] []*waitFor[T]

func (pq waitForPriorityQueue[T]) Len() int {
	return len(pq)
}
func (pq waitForPriorityQueue[T]) Less(i, j int) bool {
	return pq[i].readyAt.Before(pq[j].readyAt)
}
func (pq waitForPriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *waitForPriorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*waitFor[T])
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an item from the queue. Pop should not be called directly;
// instead, use `heap.Pop`.
func (pq *waitForPriorityQueue[T]) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	(*pq)[n-1] = nil
	item.index = -1
	*pq = (*pq)[0:(n - 1)]
	return item
}

// Peek returns the item at the beginning of the queue, without removing the
// item or otherwise mutating the queue. It is safe to call directly.
func (pq waitForPriorityQueue[T]) Peek() interface{} {
	return pq[0]
}
