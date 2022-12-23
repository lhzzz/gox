package queue

// RateLimitingInterface is an interface that rate limits items being added to the queue.
type RateLimitingQueue[T comparable] interface {
	DelayingQueue[T]

	// AddRateLimited adds an item to the workqueue after the rate limiter says it's ok
	AddRateLimited(item T)

	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for perm failing
	// or for success, we'll stop the rate limiter from tracking it.  This only clears the `rateLimiter`, you
	// still have to call `Done` on the queue.
	Forget(item T)

	// NumRequeues returns back how many times the item was requeued
	NumRequeues(item T) int
}

// rateLimitingType wraps an Interface and provides rateLimited re-enquing
type rateLimitingType[T comparable] struct {
	DelayingQueue[T]

	rateLimiter RateLimiter[T]
}

func NewRateLimitingQueue[T comparable](rateLimiter RateLimiter[T]) RateLimitingQueue[T] {
	return &rateLimitingType[T]{
		DelayingQueue: NewDelayingQueue[T](),
		rateLimiter:   rateLimiter,
	}
}

// AddRateLimited AddAfter's the item based on the time when the rate limiter says it's ok
func (q *rateLimitingType[T]) AddRateLimited(item T) {
	q.DelayingQueue.AddAfter(item, q.rateLimiter.When(item))
}

func (q *rateLimitingType[T]) NumRequeues(item T) int {
	return q.rateLimiter.NumRequeues(item)
}

func (q *rateLimitingType[T]) Forget(item T) {
	q.rateLimiter.Forget(item)
}
