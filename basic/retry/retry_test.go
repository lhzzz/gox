package retry

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	err := Do(func() error {
		time.Sleep(3 * time.Second)
		t.Log(time.Now())
		return fmt.Errorf("do fail")
	})

	t.Log(err)
}

func TestRetry(t *testing.T) {
	r := NewRetry(200*time.Millisecond, 2*time.Second, 10*time.Second)
	err := r.Do(func() error {
		time.Sleep(3 * time.Second)
		t.Log(time.Now())
		return fmt.Errorf("do fail")
	})
	t.Log(err)
}

func TestRetryContext(t *testing.T) {
	r := NewRetry(200*time.Millisecond, 2*time.Second, 10*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.DoWithContext(ctx, func() error {
		time.Sleep(3 * time.Second)
		t.Log(time.Now())
		return fmt.Errorf("do fail")
	})
	t.Log(err)
}
