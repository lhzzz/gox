package circulequeue

import (
	"fmt"
	"strings"
	"sync"
)

type Queue[T comparable] struct {
	values   []T
	writeIdx int
	readIdx  int
	full     bool
	cap      int
	size     int
	notempty *sync.Cond
	notfull  *sync.Cond
	lock     sync.RWMutex
}

func New[T comparable](maxSize int) *Queue[T] {
	if maxSize == 0 {
		return nil
	}
	q := &Queue[T]{
		values: make([]T, maxSize, maxSize),
		cap:    maxSize,
	}
	q.notempty = sync.NewCond(&q.lock)
	q.notfull = sync.NewCond(&q.lock)
	return q
}

func (q *Queue[T]) Empty() bool {
	return q.size == 0
}

func (q *Queue[T]) Full() bool {
	return q.size == q.cap
}

// Size returns number of elements within the queue.
func (q *Queue[T]) Size() int {
	return q.size
}

// Clear removes all elements from the queue.
func (q *Queue[T]) Clear() {
	q.values = make([]T, q.cap, q.cap)
	q.writeIdx = 0
	q.readIdx = 0
	q.full = false
	q.size = 0
}

func (q *Queue[T]) Values() []T {
	values := make([]T, q.Size(), q.Size())
	for i := 0; i < q.Size(); i++ {
		values[i] = q.values[(q.readIdx+i)%q.cap]
	}
	return values
}

// String returns a string representation of container
func (q *Queue[T]) String() string {
	str := "CircularQueue\n"
	var values []string
	for _, value := range q.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (q *Queue[T]) Push(value T) {
	if q.Full() {
		q.Pop()
	}
	q.values[q.writeIdx] = value
	q.writeIdx = q.writeIdx + 1
	if q.writeIdx >= q.cap {
		q.writeIdx = 0
	}
	if q.writeIdx == q.readIdx {
		q.full = true
	}
	q.size = q.calculateSize()
}

func (q *Queue[T]) Pop() (value interface{}, ok bool) {
	if q.Empty() {
		return nil, false
	}

	value, ok = q.values[q.readIdx], true
	if value != nil {
		var zero T
		q.values[q.readIdx] = zero
		q.readIdx = q.readIdx + 1
		if q.readIdx >= q.cap {
			q.readIdx = 0
		}
		q.full = false
	}
	q.size = q.size - 1
	return
}

func (q *Queue[T]) calculateSize() int {
	if q.writeIdx < q.readIdx {
		return q.cap - q.readIdx + q.writeIdx
	} else if q.writeIdx == q.readIdx {
		if q.full {
			return q.cap
		}
		return 0
	}
	return q.writeIdx - q.readIdx
}
