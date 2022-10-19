package breaker

// Acceptable is the func to check if the error can be accepted.
type Acceptable func(err error) bool

type Breaker interface {
	// Do runs the given request if the Breaker accepts it.
	// Do returns an error instantly if the Breaker rejects the request.
	// If a panic occurs in the request, the Breaker handles it as an error
	// and causes the same panic again.
	Do(method string, req func() error) error
}
