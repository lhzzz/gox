package limit

import (
	"context"
	"strings"
	"time"

	"github.com/juju/ratelimit"
)

type RateLimitStrategy int32

const (
	RateLimitStrategy_REJECT RateLimitStrategy = iota
	RateLimitStrategy_WAIT
	RateLimitStrategy_BLOCK
)

const infinityDuration time.Duration = 1<<63 - 1

var (
	defaultLimitOptions = limitOptions{
		strategy:     RateLimitStrategy_REJECT,
		fillInterval: time.Millisecond,
		capacity:     1000,
		quantum:      1,
	}
)

type limitOptions struct {
	strategy     RateLimitStrategy
	fillInterval time.Duration
	capacity     int64
	quantum      int64
	waitDuration time.Duration
}

type LimitOption func(lo *limitOptions)

func RejectStrategy() LimitOption {
	return func(lo *limitOptions) {
		lo.strategy = RateLimitStrategy_REJECT
	}
}

func WaitStrategy(d time.Duration) LimitOption {
	return func(lo *limitOptions) {
		lo.strategy = RateLimitStrategy_WAIT
		lo.waitDuration = d
	}
}

func BlockStrategy() LimitOption {
	return func(lo *limitOptions) {
		lo.strategy = RateLimitStrategy_BLOCK
		lo.waitDuration = infinityDuration
	}
}

func Capcaity(cap int64) LimitOption {
	return func(lo *limitOptions) {
		lo.capacity = cap
	}
}

func Quantum(qu int64) LimitOption {
	return func(lo *limitOptions) {
		lo.quantum = qu
	}
}

func FillInterval(interval time.Duration) LimitOption {
	return func(lo *limitOptions) {
		lo.fillInterval = interval
	}
}

type prefixTokenLimiter struct {
	prefix string
	bucket *ratelimit.Bucket
	opts   limitOptions
}

func NewPrefixTokenLimiter(prefix string, opts ...LimitOption) Limiter {
	opt := defaultLimitOptions
	for _, o := range opts {
		o(&opt)
	}
	bucket := ratelimit.NewBucketWithQuantum(opt.fillInterval, opt.capacity, opt.quantum)
	return &prefixTokenLimiter{
		prefix: prefix,
		bucket: bucket,
		opts:   opt,
	}
}

func (pl *prefixTokenLimiter) Limit(method string) bool {
	if !strings.HasPrefix(method, pl.prefix) {
		return false
	}
	return !pl.bucket.WaitMaxDuration(1, pl.opts.waitDuration)
}

func (pl *prefixTokenLimiter) LimitWithContext(ctx context.Context, method string) bool {
	if !strings.HasPrefix(method, pl.prefix) {
		return false
	}
	wait := pl.opts.waitDuration
	if pl.opts.strategy != RateLimitStrategy_REJECT {
		deadline, ok := ctx.Deadline()
		if ok {
			until := time.Until(deadline)
			if until < wait {
				wait = until
			}
		}
	}
	return !pl.bucket.WaitMaxDuration(1, wait)
}

type mutiTokenLimiter struct {
	muti map[string]struct {
		bucket *ratelimit.Bucket
		opts   limitOptions
	}
}

func NewMutiTokenLimiter(configs map[string][]LimitOption) Limiter {
	l := mutiTokenLimiter{
		muti: make(map[string]struct {
			bucket *ratelimit.Bucket
			opts   limitOptions
		}),
	}
	for method, opts := range configs {
		opt := defaultLimitOptions
		for _, o := range opts {
			o(&opt)
		}
		b := ratelimit.NewBucketWithQuantum(opt.fillInterval, opt.capacity, opt.quantum)
		l.muti[method] = struct {
			bucket *ratelimit.Bucket
			opts   limitOptions
		}{
			bucket: b,
			opts:   opt,
		}
	}
	return &l
}

func (ml *mutiTokenLimiter) Limit(method string) bool {
	bo, ok := ml.muti[method]
	if !ok {
		return false
	}
	return !bo.bucket.WaitMaxDuration(1, bo.opts.waitDuration)
}

func (ml *mutiTokenLimiter) LimitWithContext(ctx context.Context, method string) bool {
	bo, ok := ml.muti[method]
	if !ok {
		return false
	}
	wait := bo.opts.waitDuration
	if bo.opts.strategy != RateLimitStrategy_REJECT {
		deadline, ok := ctx.Deadline()
		if ok {
			until := time.Until(deadline)
			if until < wait {
				wait = until
			}
		}
	}
	return !bo.bucket.WaitMaxDuration(1, wait)
}
