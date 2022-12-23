package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialFailureRateLimiter(t *testing.T) {
	limiter := NewExponentialFailureRateLimiter[string](1*time.Millisecond, 1*time.Second)

	assert.EqualValues(t, limiter.When("one"), 1*time.Millisecond)
	assert.EqualValues(t, limiter.When("one"), 2*time.Millisecond)
	assert.EqualValues(t, limiter.When("one"), 4*time.Millisecond)
	assert.EqualValues(t, limiter.When("one"), 8*time.Millisecond)
	assert.EqualValues(t, limiter.When("one"), 16*time.Millisecond)
	assert.EqualValues(t, limiter.NumRequeues("one"), 5)

	assert.EqualValues(t, limiter.When("two"), 1*time.Millisecond)
	assert.EqualValues(t, limiter.When("two"), 2*time.Millisecond)
	assert.EqualValues(t, limiter.NumRequeues("two"), 2)

	limiter.Forget("one")
	assert.EqualValues(t, limiter.NumRequeues("one"), 0)
	assert.EqualValues(t, limiter.When("one"), 1*time.Millisecond)
}

func TestExponentialFailureRateLimiterOverFlow(t *testing.T) {
	limiter := NewExponentialFailureRateLimiter[string](1*time.Millisecond, 1000*time.Second)
	for i := 0; i < 5; i++ {
		limiter.When("one")
	}
	assert.EqualValues(t, limiter.When("one"), 32*time.Millisecond)

	for i := 0; i < 1000; i++ {
		limiter.When("overflow1")
	}
	assert.EqualValues(t, limiter.When("overflow1"), 1000*time.Second)

	limiter = NewExponentialFailureRateLimiter[string](1*time.Minute, 1000*time.Hour)
	for i := 0; i < 2; i++ {
		limiter.When("two")
	}
	assert.EqualValues(t, limiter.When("two"), 4*time.Minute)

	for i := 0; i < 1000; i++ {
		limiter.When("overflow2")
	}
	assert.EqualValues(t, limiter.When("overflow2"), 1000*time.Hour)
}

func TestItemFastSlowRateLimiter(t *testing.T) {
	limiter := NewFastSlowRateLimiter[string](5*time.Millisecond, 10*time.Second, 3)

	assert.EqualValues(t, limiter.When("one"), 5*time.Millisecond)
	assert.EqualValues(t, limiter.When("one"), 5*time.Millisecond)

	assert.EqualValues(t, limiter.When("one"), 5*time.Millisecond)

	assert.EqualValues(t, limiter.When("one"), 10*time.Second)
	assert.EqualValues(t, limiter.When("one"), 10*time.Second)
	assert.EqualValues(t, limiter.NumRequeues("one"), 5)

	assert.EqualValues(t, limiter.When("two"), 5*time.Millisecond)
	assert.EqualValues(t, limiter.When("two"), 5*time.Millisecond)
	assert.EqualValues(t, limiter.NumRequeues("two"), 2)

	limiter.Forget("one")
	assert.EqualValues(t, limiter.NumRequeues("one"), 0)
	assert.EqualValues(t, limiter.When("one"), 5*time.Millisecond)
}
