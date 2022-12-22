package sets

import "sync"

type threadsafeSet[T comparable] struct {
	nts nonthreadsafeSet[T]
	mtx sync.RWMutex
}

func newThreadSafe[T comparable](values ...T) *threadsafeSet[T] {
	ts := threadsafeSet[T]{
		nts: nonthreadsafeSet[T]{
			m: make(map[T]struct{}),
		},
	}
	ts.Add(values...)
	return &ts
}

func (s *threadsafeSet[T]) Add(items ...T) {
	s.mtx.Lock()
	s.nts.Add(items...)
	s.mtx.Unlock()
}

func (s *threadsafeSet[T]) Remove(items ...T) {
	s.mtx.Lock()
	s.nts.Remove(items...)
	s.mtx.Unlock()
}

func (s *threadsafeSet[T]) Contains(items ...T) bool {
	s.mtx.RLock()
	b := s.nts.Contains(items...)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet[T]) Size() int {
	s.mtx.RLock()
	size := s.nts.Size()
	s.mtx.RUnlock()
	return size
}

func (s *threadsafeSet[T]) Clear() {
	s.mtx.Lock()
	s.nts.Clear()
	s.mtx.Unlock()
}

func (s *threadsafeSet[T]) Empty() bool {
	s.mtx.Lock()
	b := s.nts.Size() == 0
	s.mtx.Unlock()
	return b
}

func (s *threadsafeSet[T]) IsEqual(t Set[T]) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if conv, ok := t.(*threadsafeSet[T]); ok {
		conv.mtx.RLock()
		defer conv.mtx.RUnlock()
	}

	return s.nts.IsEqual(t)
}

func (s *threadsafeSet[T]) IsSubset(t Set[T]) (subset bool) {
	s.mtx.RLock()
	b := s.nts.IsSubset(t)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet[T]) IsSuperset(t Set[T]) (subset bool) {
	s.mtx.RLock()
	b := s.nts.IsSuperset(t)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet[T]) Each(f func(item T) bool) {
	s.mtx.RLock()
	s.nts.Each(f)
	s.mtx.RUnlock()
}

func (s *threadsafeSet[T]) String() string {
	s.mtx.RLock()
	str := s.nts.String()
	s.mtx.RUnlock()
	return str
}

func (s *threadsafeSet[T]) Values() []T {
	s.mtx.RLock()
	l := s.nts.Values()
	s.mtx.RUnlock()
	return l
}

func (s *threadsafeSet[T]) Copy() Set[T] {
	s.mtx.RLock()
	us := s.nts.Copy().(*nonthreadsafeSet[T])
	s.mtx.RUnlock()

	return &threadsafeSet[T]{nts: *us}
}

func (s *threadsafeSet[T]) Merge(t Set[T]) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.nts.Merge(t)
}

func (s *threadsafeSet[T]) Separate(t Set[T]) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.nts.Remove(t.Values()...)
}

func (s *threadsafeSet[T]) Union(sets ...Set[T]) Set[T] {
	s.mtx.RLock()
	us := s.nts.Union(sets...).(*nonthreadsafeSet[T])
	s.mtx.RUnlock()
	return &threadsafeSet[T]{nts: *us}
}

func (s *threadsafeSet[T]) Difference(sets ...Set[T]) Set[T] {
	s.mtx.RLock()
	us := s.nts.Difference(sets...).(*nonthreadsafeSet[T])
	s.mtx.RUnlock()

	return &threadsafeSet[T]{nts: *us}
}

func (s *threadsafeSet[T]) Intersect(t Set[T]) Set[T] {
	s.mtx.RLock()
	us := s.nts.Intersect(t).(*nonthreadsafeSet[T])
	s.mtx.RUnlock()
	return &threadsafeSet[T]{nts: *us}
}

func (s *threadsafeSet[T]) Intersects(sets ...Set[T]) Set[T] {
	s.mtx.RLock()
	us := s.nts.Intersects(sets...).(*nonthreadsafeSet[T])
	s.mtx.RUnlock()
	return &threadsafeSet[T]{nts: *us}
}
