package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Retry interface {
	Do(func() error) error
	DoWithContext(ctx context.Context, fn func() error) error
}

type retry struct {
	ebo *backoff.ExponentialBackOff
}

var kRetry = retry{ebo: &backoff.ExponentialBackOff{
	InitialInterval:     backoff.DefaultInitialInterval,
	RandomizationFactor: backoff.DefaultRandomizationFactor,
	Multiplier:          backoff.DefaultMultiplier,
	MaxInterval:         backoff.DefaultMaxInterval,
	MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
	Stop:                backoff.Stop,
	Clock:               backoff.SystemClock,
}}

/*
Note: it's not thread-safe
 next call time:  random(current * (1- 0.5), current * (1 + 0.5)) * 1.5
 exit: current >= MaxElapsedTime
*/
func NewRetry(initInterval, maxInterval, maxElapsedTime time.Duration) Retry {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = initInterval
	bo.MaxInterval = maxInterval
	bo.MaxElapsedTime = maxElapsedTime
	bo.Reset()
	r := retry{ebo: bo}
	return &r
}

func (r *retry) Do(fn func() error) error {
	r.ebo.Reset()
	return backoff.Retry(fn, r.ebo)
}

func (r *retry) DoWithContext(ctx context.Context, fn func() error) error {
	bctx := backoff.WithContext(r.ebo, ctx)
	r.ebo.Reset()
	return backoff.Retry(fn, bctx)
}

func Do(fn func() error) error {
	return kRetry.Do(fn)
}

func DoWithContext(ctx context.Context, fn func() error) error {
	return kRetry.DoWithContext(ctx, fn)
}
