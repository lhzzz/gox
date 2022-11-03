package single

import (
	"golang.org/x/sync/singleflight"
)

type Group struct {
	sf singleflight.Group
}

type Result singleflight.Result

func New() Group {
	return Group{sf: singleflight.Group{}}
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return g.sf.Do(key, fn)
}

func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	return convertResultChan(g.sf.DoChan(key, fn))
}

func (g *Group) Forget(key string) {
	g.sf.Forget(key)
}

func convertResultChan(in <-chan singleflight.Result) <-chan Result {
	out := make(chan Result, 1)

	go func() {
		for r := range in {
			out <- Result(r)
		}
	}()

	return out
}
