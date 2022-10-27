package sets

import "sync"

type threadsafeSet struct {
	nts nonthreadsafeSet
	mtx sync.RWMutex
}

func newThreadSafe() *threadsafeSet {
	ts := threadsafeSet{
		nts: nonthreadsafeSet{
			m: make(map[interface{}]struct{}),
		},
	}
	return &ts
}

func (s *threadsafeSet) Add(items ...interface{}) {
	s.mtx.Lock()
	s.nts.Add(items...)
	s.mtx.Unlock()
}

func (s *threadsafeSet) Remove(items ...interface{}) {
	s.mtx.Lock()
	s.nts.Remove(items...)
	s.mtx.Unlock()
}

func (s *threadsafeSet) Contains(items ...interface{}) bool {
	s.mtx.RLock()
	b := s.nts.Contains(items...)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet) Size() int {
	s.mtx.RLock()
	size := s.nts.Size()
	s.mtx.RUnlock()
	return size
}

func (s *threadsafeSet) Clear() {
	s.mtx.Lock()
	s.nts.Clear()
	s.mtx.Unlock()
}

func (s *threadsafeSet) IsEmpty() bool {
	s.mtx.Lock()
	b := s.nts.Size() == 0
	s.mtx.Unlock()
	return b
}

func (s *threadsafeSet) IsEqual(t Set) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if conv, ok := t.(*threadsafeSet); ok {
		conv.mtx.RLock()
		defer conv.mtx.RUnlock()
	}

	return s.nts.IsEqual(t)
}

func (s *threadsafeSet) IsSubset(t Set) (subset bool) {
	s.mtx.RLock()
	b := s.nts.IsSubset(t)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet) IsSuperset(t Set) (subset bool) {
	s.mtx.RLock()
	b := s.nts.IsSuperset(t)
	s.mtx.RUnlock()
	return b
}

func (s *threadsafeSet) Each(f func(item interface{}) bool) {
	s.mtx.RLock()
	s.nts.Each(f)
	s.mtx.RUnlock()
}

func (s *threadsafeSet) String() string {
	s.mtx.RLock()
	str := s.nts.String()
	s.mtx.RUnlock()
	return str
}

func (s *threadsafeSet) List() []interface{} {
	s.mtx.RLock()
	l := s.nts.List()
	s.mtx.RUnlock()
	return l
}

func (s *threadsafeSet) Copy() Set {
	s.mtx.RLock()
	us := s.nts.Copy().(*nonthreadsafeSet)
	s.mtx.RUnlock()

	return &threadsafeSet{nts: *us}
}

func (s *threadsafeSet) Merge(t Set) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.nts.Merge(t)
}

func (s *threadsafeSet) Separate(t Set) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.nts.Remove(t.List()...)
}

func (s *threadsafeSet) Union(sets ...Set) Set {
	s.mtx.RLock()
	us := s.nts.Union(sets...).(*nonthreadsafeSet)
	s.mtx.RUnlock()
	return &threadsafeSet{nts: *us}
}

func (s *threadsafeSet) Difference(sets ...Set) Set {
	s.mtx.RLock()
	us := s.nts.Difference(sets...).(*nonthreadsafeSet)
	s.mtx.RUnlock()

	return &threadsafeSet{nts: *us}
}

func (s *threadsafeSet) Intersect(t Set) Set {
	s.mtx.RLock()
	us := s.nts.Intersect(t).(*nonthreadsafeSet)
	s.mtx.RUnlock()
	return &threadsafeSet{nts: *us}
}

func (s *threadsafeSet) Intersects(sets ...Set) Set {
	s.mtx.RLock()
	us := s.nts.Intersects(sets...).(*nonthreadsafeSet)
	s.mtx.RUnlock()
	return &threadsafeSet{nts: *us}
}
