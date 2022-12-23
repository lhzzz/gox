package queue

import (
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter[T comparable] interface {
	// When gets an item and gets to decide how long that item should wait
	When(item T) time.Duration
	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for failing
	// or for success, we'll stop tracking it
	Forget(item T)
	// NumRequeues returns back how many failures the item has had
	NumRequeues(item T) int
}

// DefaultControllerRateLimiter is a no-arg constructor for a default rate limiter for a workqueue.  It has
// both overall and per-item rate limiting.  The overall is a token bucket and the per-item is exponential
func DefaultControllerRateLimiter[T comparable]() RateLimiter[T] {
	return NewMaxOfRateLimiter(
		NewExponentialFailureRateLimiter[T](5*time.Millisecond, 1000*time.Second),
		// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		NewBucketRateLimiter[T](10, 100),
	)
}

// BucketRateLimiter adapts a standard bucket to the workqueue ratelimiter API
type BucketRateLimiter[T comparable] struct {
	*rate.Limiter
}

func NewBucketRateLimiter[T comparable](qps, bucket int) RateLimiter[T] {
	return &BucketRateLimiter[T]{Limiter: rate.NewLimiter(rate.Limit(qps), bucket)}
}

func (r *BucketRateLimiter[T]) When(item T) time.Duration {
	return r.Limiter.Reserve().Delay()
}

func (r *BucketRateLimiter[T]) NumRequeues(item T) int {
	return 0
}

func (r *BucketRateLimiter[T]) Forget(item T) {
}

// exponentialFailureRateLimiter does a simple baseDelay*2^<num-failures> limit
// dealing with max failures and expiration are up to the caller
type exponentialFailureRateLimiter[T comparable] struct {
	failuresLock sync.Mutex
	failures     map[T]int

	baseDelay time.Duration
	maxDelay  time.Duration
}

func NewExponentialFailureRateLimiter[T comparable](baseDelay time.Duration, maxDelay time.Duration) RateLimiter[T] {
	return &exponentialFailureRateLimiter[T]{
		failures:  map[T]int{},
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}

func DefaultBasedRateLimiter[T comparable]() RateLimiter[T] {
	return NewExponentialFailureRateLimiter[T](time.Millisecond, 1000*time.Second)
}

func (r *exponentialFailureRateLimiter[T]) When(item T) time.Duration {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	fs := r.failures[item]
	r.failures[item]++

	backoff := float64(r.baseDelay.Nanoseconds()) * math.Pow(2, float64(fs))
	if backoff > math.MaxInt64 {
		return r.maxDelay
	}

	d := time.Duration(backoff)
	if d > r.maxDelay {
		return r.maxDelay
	}
	return d
}

func (r *exponentialFailureRateLimiter[T]) Forget(item T) {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	delete(r.failures, item)
}

func (r *exponentialFailureRateLimiter[T]) NumRequeues(item T) int {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	return r.failures[item]
}

// FastSlowRateLimiter does a quick retry for a certain number of attempts, then a slow retry after that
type FastSlowRateLimiter[T comparable] struct {
	failuresLock sync.Mutex
	failures     map[T]int

	maxFastAttempts int
	fastDelay       time.Duration
	slowDelay       time.Duration
}

func NewFastSlowRateLimiter[T comparable](fastDelay, slowDelay time.Duration, maxFastAttempts int) RateLimiter[T] {
	return &FastSlowRateLimiter[T]{
		failures:        map[T]int{},
		fastDelay:       fastDelay,
		slowDelay:       slowDelay,
		maxFastAttempts: maxFastAttempts,
	}
}

func (r *FastSlowRateLimiter[T]) When(item T) time.Duration {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	r.failures[item] = r.failures[item] + 1

	if r.failures[item] <= r.maxFastAttempts {
		return r.fastDelay
	}

	return r.slowDelay
}

func (r *FastSlowRateLimiter[T]) Forget(item T) {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	delete(r.failures, item)
}

func (r *FastSlowRateLimiter[T]) NumRequeues(item T) int {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	return r.failures[item]
}

// MaxOfRateLimiter calls every RateLimiter and returns the worst case response
// When used with a token bucket limiter, the burst could be apparently exceeded in cases where particular items
// were separately delayed a longer time.
type MaxOfRateLimiter[T comparable] struct {
	limiters []RateLimiter[T]
}

func (r *MaxOfRateLimiter[T]) When(item T) time.Duration {
	ret := time.Duration(0)
	for _, limiter := range r.limiters {
		curr := limiter.When(item)
		if curr > ret {
			ret = curr
		}
	}

	return ret
}

func NewMaxOfRateLimiter[T comparable](limiters ...RateLimiter[T]) RateLimiter[T] {
	return &MaxOfRateLimiter[T]{limiters: limiters}
}

func (r *MaxOfRateLimiter[T]) NumRequeues(item T) int {
	ret := 0
	for _, limiter := range r.limiters {
		curr := limiter.NumRequeues(item)
		if curr > ret {
			ret = curr
		}
	}

	return ret
}

func (r *MaxOfRateLimiter[T]) Forget(item T) {
	for _, limiter := range r.limiters {
		limiter.Forget(item)
	}
}
