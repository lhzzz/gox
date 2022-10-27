package sets

import (
	"fmt"
	"strings"
)

type nonthreadsafeSet struct {
	m map[interface{}]struct{}
}

func newNonThreadSafe() *nonthreadsafeSet {
	nts := nonthreadsafeSet{
		m: make(map[interface{}]struct{}),
	}
	return &nts
}

func (s *nonthreadsafeSet) Add(items ...interface{}) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

func (s *nonthreadsafeSet) Remove(items ...interface{}) {
	for _, item := range items {
		delete(s.m, item)
	}
}

func (s *nonthreadsafeSet) Contains(items ...interface{}) bool {
	has := true
	for _, item := range items {
		if _, has = s.m[item]; !has {
			break
		}
	}
	return has
}

func (s *nonthreadsafeSet) Size() int {
	return len(s.m)
}

func (s *nonthreadsafeSet) Clear() {
	s.m = make(map[interface{}]struct{})
}

func (s *nonthreadsafeSet) IsEmpty() bool {
	return s.Size() == 0
}

func (s *nonthreadsafeSet) IsEqual(t Set) bool {
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item interface{}) bool {
		_, equal = s.m[item]
		return equal
	})
	return equal
}

func (s *nonthreadsafeSet) IsSubset(t Set) bool {
	subset := true
	t.Each(func(item interface{}) bool {
		_, subset = s.m[item]
		return subset
	})

	return subset
}

func (s *nonthreadsafeSet) IsSuperset(t Set) bool {
	return t.IsSubset(s)
}

func (s *nonthreadsafeSet) Each(f func(item interface{}) bool) {
	for item := range s.m {
		if !f(item) {
			break
		}
	}
}

func (s *nonthreadsafeSet) String() string {
	t := make([]string, 0, len(s.List()))
	for _, item := range s.List() {
		t = append(t, fmt.Sprintf("%v", item))
	}

	return fmt.Sprintf("[%s]", strings.Join(t, ", "))
}

func (s *nonthreadsafeSet) List() []interface{} {
	list := make([]interface{}, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *nonthreadsafeSet) Copy() Set {
	u := newNonThreadSafe()
	for item := range s.m {
		u.Add(item)
	}
	return u
}

func (s *nonthreadsafeSet) Merge(t Set) {
	t.Each(func(item interface{}) bool {
		s.m[item] = struct{}{}
		return true
	})
}

func (s *nonthreadsafeSet) Separate(t Set) {
	s.Remove(t.List()...)
}

func (s *nonthreadsafeSet) Union(sets ...Set) Set {
	u := s.Copy()
	for _, set := range sets {
		set.Each(func(item interface{}) bool {
			u.Add(item)
			return true
		})
	}
	return u
}

func (s *nonthreadsafeSet) Difference(sets ...Set) Set {
	u := s.Copy()
	for _, set := range sets {
		u.Separate(set)
	}
	return u
}

func (s *nonthreadsafeSet) Intersect(t Set) Set {
	result := newNonThreadSafe()
	if s.Size() < t.Size() {
		s.Each(func(item interface{}) bool {
			if t.Contains(item) {
				result.Add(item)
			}
			return true
		})
	} else {
		t.Each(func(item interface{}) bool {
			if s.Contains(item) {
				result.Add(item)
			}
			return true
		})
	}
	return result
}

func (s *nonthreadsafeSet) Intersects(sets ...Set) Set {
	all := s.Union(sets...)
	result := s.Union(sets...)

	all.Each(func(item interface{}) bool {
		for _, set := range sets {
			if !set.Contains(item) {
				result.Remove(item)
			}
		}
		return true
	})
	return result
}
