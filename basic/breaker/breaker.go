package breaker

// Acceptable is the func to check if the error can be accepted.
type Acceptable func(err error) bool

type Promise interface {
	// Accept tells the Breaker that the call is successful.
	Accept()
	// Reject tells the Breaker that the call is failed.
	Reject(reason string)
}

type Breaker interface {
	// Name returns the name of the Breaker.
	Name() string

	// Allow checks if the request is allowed.
	// If allowed, a promise will be returned, the caller needs to call promise.Accept()
	// on success, or call promise.Reject() on failure.
	// If not allow, ErrServiceUnavailable will be returned.
	Allow() (Promise, error)

	// Do runs the given request if the Breaker accepts it.
	// Do returns an error instantly if the Breaker rejects the request.
	// If a panic occurs in the request, the Breaker handles it as an error
	// and causes the same panic again.
	Do(req func() error) error

	// DoWithAcceptable runs the given request if the Breaker accepts it.
	// DoWithAcceptable returns an error instantly if the Breaker rejects the request.
	// If a panic occurs in the request, the Breaker handles it as an error
	// and causes the same panic again.
	// acceptable checks if it's a successful call, even if the err is not nil.
	DoWithAcceptable(req func() error, acceptable Acceptable) error

	// DoWithFallback runs the given request if the Breaker accepts it.
	// DoWithFallback runs the fallback if the Breaker rejects the request.
	// If a panic occurs in the request, the Breaker handles it as an error
	// and causes the same panic again.
	DoWithFallback(req func() error, fallback func(err error) error) error

	// DoWithFallbackAcceptable runs the given request if the Breaker accepts it.
	// DoWithFallbackAcceptable runs the fallback if the Breaker rejects the request.
	// If a panic occurs in the request, the Breaker handles it as an error
	// and causes the same panic again.
	// acceptable checks if it's a successful call, even if the err is not nil.
	DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
}
