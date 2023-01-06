package circulequeue

import (
	"fmt"
	"strings"
	"sync"

	"singer.com/basic/containers/queues"
)

type Queue[T comparable] struct {
	values   []T
	writeIdx int
	readIdx  int
	cap      int
	size     int
	notempty *sync.Cond
	notfull  *sync.Cond
	lock     sync.RWMutex
}

func New[T comparable](maxSize int) queues.Queue[T] {
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
	q.lock.Lock()
	for q.cap == q.size {
		q.notfull.Wait()
	}
	q.values[q.writeIdx] = value
	q.writeIdx++
	if q.writeIdx >= q.cap {
		q.writeIdx = 0
	}
	q.size++
	q.notempty.Signal()
	q.lock.Unlock()
}

func (q *Queue[T]) Pop() (value T) {
	q.lock.Lock()
	for q.size == 0 {
		q.notempty.Wait()
	}

	value = q.values[q.readIdx]
	q.readIdx++
	if q.readIdx >= q.cap {
		q.readIdx = 0
	}
	q.size--
	q.notfull.Signal()
	q.lock.Unlock()
	return
}
