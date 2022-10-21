package limit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrefixLimiter_Reject(t *testing.T) {
	pl := NewPrefixTokenLimiter("/pb", Capcaity(10), Quantum(1), FillInterval(time.Second))

	var wg sync.WaitGroup
	limitCh := make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			name := "/pb" + fmt.Sprint(idx)
			limit := pl.Limit(name)
			if limit {
				limitCh <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(limitCh)
	limitCount := 0
	for range limitCh {
		limitCount++
	}
	assert.EqualValues(t, limitCount, 10)
}

func TestPrefixLimiter_Block(t *testing.T) {
	pl := NewPrefixTokenLimiter("/pb", Capcaity(10), Quantum(1), FillInterval(time.Second), BlockStrategy())

	var wg sync.WaitGroup
	limitCh := make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			name := "/pb" + fmt.Sprint(idx)
			limit := pl.Limit(name)
			if limit {
				limitCh <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(limitCh)
	limitCount := 0
	for range limitCh {
		limitCount++
	}
	assert.EqualValues(t, limitCount, 0)
}

func TestPrefixLimiter_Wait(t *testing.T) {
	pl := NewPrefixTokenLimiter("/pb", Capcaity(10), Quantum(1), FillInterval(time.Second), WaitStrategy(2*time.Second))

	var wg sync.WaitGroup
	limitCh := make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			name := "/pb" + fmt.Sprint(idx)
			limit := pl.Limit(name)
			if limit {
				limitCh <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(limitCh)

	limitCount := 0
	for range limitCh {
		limitCount++
	}
	assert.EqualValues(t, limitCount, 8)
}

func TestPrefixLimiter_Wait_Context(t *testing.T) {
	pl := NewPrefixTokenLimiter("/pb", Capcaity(10), Quantum(1), FillInterval(time.Second), WaitStrategy(5*time.Second))
	wg := sync.WaitGroup{}
	limitCh := make(chan struct{}, 20)
	ctx := context.Background()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			name := "/pb" + fmt.Sprint(idx)
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			limit := pl.LimitWithContext(ctx, name)
			if limit {
				limitCh <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(limitCh)

	limitCount := 0
	for range limitCh {
		limitCount++
	}
	assert.EqualValues(t, limitCount, 5)

	pl2 := NewPrefixTokenLimiter("/pb", Capcaity(10), Quantum(1), FillInterval(time.Second), WaitStrategy(10*time.Second))
	wg2 := sync.WaitGroup{}
	limitCh2 := make(chan struct{}, 20)
	ctx2 := context.Background()
	for i := 0; i < 20; i++ {
		wg2.Add(1)
		idx := i
		go func() {
			defer wg2.Done()
			name := "/pb" + fmt.Sprint(idx)
			ctx, cancel := context.WithTimeout(ctx2, time.Second*5)
			defer cancel()
			limit := pl2.LimitWithContext(ctx, name)
			if limit {
				limitCh2 <- struct{}{}
			}
		}()
	}
	wg2.Wait()
	close(limitCh2)

	limitCount2 := 0
	for range limitCh2 {
		limitCount2++
	}
	assert.EqualValues(t, limitCount2, 5)
}

func TestMutiLimiter_Reject(t *testing.T) {
	pl := NewMutiTokenLimiter(map[string][]LimitOption{
		"/api/Create": {Capcaity(10), Quantum(1), FillInterval(time.Second)},
		"/api/Update": {Capcaity(10), Quantum(1), FillInterval(time.Second)},
	})

	var wg sync.WaitGroup
	createLimit := make(chan struct{}, 20)
	updateLimit := make(chan struct{}, 20)
	deleteLimit := make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limit := pl.Limit("/api/Create")
			if limit {
				createLimit <- struct{}{}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			limit := pl.Limit("/api/Update")
			if limit {
				updateLimit <- struct{}{}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			limit := pl.Limit("/api/Delete")
			if limit {
				deleteLimit <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(createLimit)
	close(updateLimit)
	close(deleteLimit)
	createCount := 0
	updateCount := 0
	deleteCount := 0
	for range createLimit {
		createCount++
	}
	for range updateLimit {
		updateCount++
	}
	for range deleteLimit {
		deleteCount++
	}
	assert.EqualValues(t, createCount, 10)
	assert.EqualValues(t, updateCount, 10)
	assert.EqualValues(t, deleteCount, 0)
}
