package limit

import "context"

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter interface {
	Limit(method string) bool

	LimitWithContext(ctx context.Context, method string) bool
}
