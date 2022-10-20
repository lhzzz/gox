package breaker

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func TestPrefixBreaker(t *testing.T) {
	cmd := NewHystrixConfig(4, 50, nil, nil)
	b := NewPrefixHystrixBreaker("/pb", cmd)
	breaker, exist, err := hystrix.GetCircuit("/pb")
	t.Log(breaker, exist, err)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			method := fmt.Sprintf("/pb.%d", idx)
			err := b.Do(method, func() error {
				if idx%2 == 0 {
					return fmt.Errorf("random error:%d", idx)
				}
				return nil
			})
			t.Log(method, idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				isRun = true
				t.Log("Success") // not execute
				return nil
			})
			t.Log("pb", idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestPrefixBreakerWithFallback(t *testing.T) {
	cmd := NewHystrixConfig(4, 40, nil, func(err error) error {
		t.Log("breaker fallback error", err)
		return nil
	})
	b := NewPrefixHystrixBreaker("/pb", cmd)
	breaker, exist, err := hystrix.GetCircuit("/pb")
	t.Log(breaker, exist, err)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				isRun = true
				if idx%2 == 0 {
					return fmt.Errorf("random error:%d", idx)
				}
				return nil
			})
			t.Log(idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				isRun = true
				t.Log("Success") // not execute
				return nil
			})
			t.Log(idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()
}

//accept的测试用例：accept如果判断某个error是不需要被熔断的，则该error不会被统计到熔断中，且会直接返回该error
func TestPrefixBreakerWithAccept(t *testing.T) {
	cmd := NewHystrixConfig(4, 50, func(err error) bool {
		if strings.HasPrefix(err.Error(), "random") {
			t.Log("accept random error")
			return true
		}
		return false
	}, nil)
	b := NewPrefixHystrixBreaker("/pb", cmd)
	breaker, exist, err := hystrix.GetCircuit("/pb")
	t.Log(breaker, exist, err)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				return fmt.Errorf("random error:%d", idx)
			})
			t.Log("pb", idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				isRun = true
				t.Log("Success") // will execute
				return nil
			})
			t.Log("pb", idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()

	//上面已经有20个对于熔断来说是成功的请求，下面至少要20个错误请求来达到50%错误率
	t.Log("------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 20; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				return fmt.Errorf("error:%d", idx)
			})
			t.Log("pb", idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/pb.Ex", func() error {
				isRun = true
				t.Log("Success") //  not execute
				return nil
			})
			t.Log("pb", idx, isRun, err, breaker.IsOpen(), breaker.AllowRequest())
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestMutiHystrix(t *testing.T) {
	config := NewHystrixConfig(4, 50, nil, nil)
	b := NewMutiHystrixBreaker(map[string]HystrixConfig{
		"/api/Create": config,
		"/api/Update": config,
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/api/Create", func() error {
				return fmt.Errorf("create error:%d", idx)
			})
			t.Log("/api/Create", idx, isRun, err)
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/api/Create", func() error {
				isRun = true
				t.Log("Success") //  not execute
				return nil
			})
			t.Log("/api/Create", idx, isRun, err)
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/api/Update", func() error {
				isRun = true
				t.Log("Success") // execute
				return nil
			})
			t.Log("/api/Update", idx, isRun, err)
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/api/Update", func() error {
				isRun = true
				return fmt.Errorf("update error") //execute
			})
			t.Log("/api/Update", idx, isRun, err)
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log("-------------------------")
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		idx := i
		wg.Add(1)
		go func() {
			isRun := false
			err := b.Do("/api/Update", func() error {
				isRun = true
				t.Log("i am running") // not execute
				return fmt.Errorf("update error")
			})
			t.Log("/api/Update", idx, isRun, err)
			wg.Done()
		}()
	}
	wg.Wait()
}
