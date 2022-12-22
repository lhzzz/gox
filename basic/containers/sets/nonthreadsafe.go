package sets

import (
	"fmt"
	"strings"
)

type nonthreadsafeSet[T comparable] struct {
	m map[T]struct{}
}

func newNonThreadSafe[T comparable](values ...T) *nonthreadsafeSet[T] {
	nts := nonthreadsafeSet[T]{
		m: make(map[T]struct{}),
	}
	nts.Add(values...)
	return &nts
}

func (s *nonthreadsafeSet[T]) Add(items ...T) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

func (s *nonthreadsafeSet[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s.m, item)
	}
}

func (s *nonthreadsafeSet[T]) Contains(items ...T) bool {
	has := true
	for _, item := range items {
		if _, has = s.m[item]; !has {
			break
		}
	}
	return has
}

func (s *nonthreadsafeSet[T]) Size() int {
	return len(s.m)
}

func (s *nonthreadsafeSet[T]) Clear() {
	s.m = make(map[T]struct{})
}

func (s *nonthreadsafeSet[T]) Empty() bool {
	return s.Size() == 0
}

func (s *nonthreadsafeSet[T]) IsEqual(t Set[T]) bool {
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item T) bool {
		_, equal = s.m[item]
		return equal
	})
	return equal
}

func (s *nonthreadsafeSet[T]) IsSubset(t Set[T]) bool {
	subset := true
	t.Each(func(item T) bool {
		_, subset = s.m[item]
		return subset
	})

	return subset
}

func (s *nonthreadsafeSet[T]) IsSuperset(t Set[T]) bool {
	return t.IsSubset(s)
}

func (s *nonthreadsafeSet[T]) Each(f func(item T) bool) {
	for item := range s.m {
		if !f(item) {
			break
		}
	}
}

func (s *nonthreadsafeSet[T]) String() string {
	t := make([]string, 0, len(s.Values()))
	for _, item := range s.Values() {
		t = append(t, fmt.Sprintf("%v", item))
	}

	return fmt.Sprintf("[%s]", strings.Join(t, ", "))
}

func (s *nonthreadsafeSet[T]) Values() []T {
	list := make([]T, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *nonthreadsafeSet[T]) Copy() Set[T] {
	u := newNonThreadSafe[T]()
	for item := range s.m {
		u.Add(item)
	}
	return u
}

func (s *nonthreadsafeSet[T]) Merge(t Set[T]) {
	t.Each(func(item T) bool {
		s.m[item] = struct{}{}
		return true
	})
}

func (s *nonthreadsafeSet[T]) Separate(t Set[T]) {
	s.Remove(t.Values()...)
}

func (s *nonthreadsafeSet[T]) Union(sets ...Set[T]) Set[T] {
	u := s.Copy()
	for _, set := range sets {
		set.Each(func(item T) bool {
			u.Add(item)
			return true
		})
	}
	return u
}

func (s *nonthreadsafeSet[T]) Difference(sets ...Set[T]) Set[T] {
	u := s.Copy()
	for _, set := range sets {
		u.Separate(set)
	}
	return u
}

func (s *nonthreadsafeSet[T]) Intersect(t Set[T]) Set[T] {
	result := newNonThreadSafe[T]()
	if s.Size() < t.Size() {
		s.Each(func(item T) bool {
			if t.Contains(item) {
				result.Add(item)
			}
			return true
		})
	} else {
		t.Each(func(item T) bool {
			if s.Contains(item) {
				result.Add(item)
			}
			return true
		})
	}
	return result
}

func (s *nonthreadsafeSet[T]) Intersects(sets ...Set[T]) Set[T] {
	all := s.Union(sets...)
	result := s.Union(sets...)

	all.Each(func(item T) bool {
		for _, set := range sets {
			if !set.Contains(item) {
				result.Remove(item)
			}
		}
		return true
	})
	return result
}
