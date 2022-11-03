package single

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func sleeping() int {
	time.Sleep(3 * time.Second)
	fmt.Println("someone has entry this")
	return rand.Int()
}

func TestDo(t *testing.T) {
	g := New()

	//concurency only call once
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			v, err, shared := g.Do("caceh-miss", func() (interface{}, error) {
				res := sleeping()
				return res, nil
			})
			t.Log(v, err, shared, i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	v, err, shared := g.Do("caceh-miss", func() (interface{}, error) {
		res := sleeping()
		return res, nil
	})
	t.Log(" loop out: ", v, err, shared)

	v, err, shared = g.Do("caceh-miss", func() (interface{}, error) {
		res := sleeping()
		return res, nil
	})
	t.Log(" loop out: ", v, err, shared)
}
