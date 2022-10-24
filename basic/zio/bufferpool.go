package zio

import (
	"bytes"
	"sync"
)

var (
	bufferPool BufferPool
)

type BufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}

type defaultPool struct {
	pool *sync.Pool
}

func init() {
	SetBufferPool(&defaultPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	})
}

func (p *defaultPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

func (p *defaultPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func SetBufferPool(bp BufferPool) {
	bufferPool = bp
}
