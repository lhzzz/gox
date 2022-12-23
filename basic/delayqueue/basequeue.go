package queue

import (
	"sync"
	"time"

	"singer.com/util/clock"
)

type BaseQueue[T comparable] interface {
	Add(item T)
	Len() int
	Get() (item T, shutdown bool)
	Done(item T)
	ShutDown()
	ShutDownWithDrain()
	ShuttingDown() bool
}

const defaultUnfinishedWorkUpdatePeriod = 500 * time.Millisecond

/*
 queue 		[[1],[2],[3]]    	queue 		[[2],[3]]		queue 		[[2],[3]]
 dirty 		[[1],[2],[3]]   --> dirty 		[[2],[3]]  --> 	dirty		[[2],[3]]
 processing []					processing 	[[1]]			processing	[]
*/
// baseQueue is a work queue.
type baseQueue[T comparable] struct {
	// queue defines the order in which we will work on items. Every
	// element of queue should be in the dirty set and not in the
	// processing set.
	queue []T

	// dirty defines all of the items that need to be processed.
	dirty set[T]

	// Things that are currently being processed are in the processing set.
	// These things may be simultaneously in the dirty set. When we finish
	// processing something and remove it from this set, we'll check if
	// it's in the dirty set, and if so, add it to the queue.
	processing set[T]

	cond *sync.Cond

	shuttingDown bool
	drain        bool

	unfinishedWorkUpdatePeriod time.Duration
	clock                      clock.WithTicker
}

type empty struct{}
type set[T comparable] map[T]empty

func NewBaseQueue[T comparable]() BaseQueue[T] {
	rc := clock.RealClock{}
	return newQueue[T](rc, defaultUnfinishedWorkUpdatePeriod)
}

func newQueue[T comparable](c clock.WithTicker, updatePeriod time.Duration) *baseQueue[T] {
	t := &baseQueue[T]{
		clock:                      c,
		dirty:                      set[T]{},
		processing:                 set[T]{},
		cond:                       sync.NewCond(&sync.Mutex{}),
		unfinishedWorkUpdatePeriod: updatePeriod,
	}
	return t
}

// Add marks item as needing processing.
func (q *baseQueue[T]) Add(item T) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return
	}
	if q.dirty.has(item) {
		return
	}

	q.dirty.insert(item)
	if q.processing.has(item) {
		return
	}

	q.queue = append(q.queue, item)
	q.cond.Signal()
}

// Len returns the current queue length, for informational purposes only. You
// shouldn't e.g. gate a call to Add() or Get() on Len() being a particular
// value, that can't be synchronized properly.
func (q *baseQueue[T]) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return len(q.queue)
}

// Get blocks until it can return an item to be processed. If shutdown = true,
// the caller should end their goroutine. You must call Done with item when you
// have finished processing it.
func (q *baseQueue[T]) Get() (item T, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
	var zero T
	if len(q.queue) == 0 {
		// We must be shutting down.
		return zero, true
	}

	item = q.queue[0]
	// The underlying array still exists and reference this object, so the object will not be garbage collected.
	q.queue[0] = zero
	q.queue = q.queue[1:]

	q.processing.insert(item)
	q.dirty.delete(item)

	return item, false
}

// Done marks item as done processing, and if it has been marked as dirty again
// while it was being processed, it will be re-added to the queue for
// re-processing.
func (q *baseQueue[T]) Done(item T) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.processing.delete(item)
	if q.dirty.has(item) {
		q.queue = append(q.queue, item)
		q.cond.Signal()
	} else if q.processing.len() == 0 {
		q.cond.Signal()
	}
}

// ShutDown will cause q to ignore all new items added to it and
// immediately instruct the worker goroutines to exit.
func (q *baseQueue[T]) ShutDown() {
	q.setDrain(false)
	q.shutdown()
}

// ShutDownWithDrain will cause q to ignore all new items added to it. As soon
// as the worker goroutines have "drained", i.e: finished processing and called
// Done on all existing items in the queue; they will be instructed to exit and
// ShutDownWithDrain will return. Hence: a strict requirement for using this is;
// your workers must ensure that Done is called on all items in the queue once
// the shut down has been initiated, if that is not the case: this will block
// indefinitely. It is, however, safe to call ShutDown after having called
// ShutDownWithDrain, as to force the queue shut down to terminate immediately
// without waiting for the drainage.
func (q *baseQueue[T]) ShutDownWithDrain() {
	q.setDrain(true)
	q.shutdown()
	for q.isProcessing() && q.shouldDrain() {
		q.waitForProcessing()
	}
}

// isProcessing indicates if there are still items on the work queue being
// processed. It's used to drain the work queue on an eventual shutdown.
func (q *baseQueue[T]) isProcessing() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.processing.len() != 0
}

// waitForProcessing waits for the worker goroutines to finish processing items
// and call Done on them.
func (q *baseQueue[T]) waitForProcessing() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	// Ensure that we do not wait on a queue which is already empty, as that
	// could result in waiting for Done to be called on items in an empty queue
	// which has already been shut down, which will result in waiting
	// indefinitely.
	if q.processing.len() == 0 {
		return
	}
	q.cond.Wait()
}

func (q *baseQueue[T]) setDrain(shouldDrain bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.drain = shouldDrain
}

func (q *baseQueue[T]) shouldDrain() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.drain
}

func (q *baseQueue[T]) shutdown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.shuttingDown = true
	q.cond.Broadcast()
}

func (q *baseQueue[T]) ShuttingDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return q.shuttingDown
}

func (s set[T]) has(item T) bool {
	_, exists := s[item]
	return exists
}

func (s set[T]) insert(item T) {
	s[item] = empty{}
}

func (s set[T]) delete(item T) {
	delete(s, item)
}

func (s set[T]) len() int {
	return len(s)
}
