package containers

import (
	"sort"
)

func Sort[T comparable](values []T, comparator Comparator[T]) {
	sort.Sort(sortable[T]{values, comparator})
}

type sortable[T comparable] struct {
	values     []T
	comparator Comparator[T]
}

func (s sortable[T]) Len() int {
	return len(s.values)
}
func (s sortable[T]) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}
func (s sortable[T]) Less(i, j int) bool {
	return s.comparator(s.values[i], s.values[j]) < 0
}
