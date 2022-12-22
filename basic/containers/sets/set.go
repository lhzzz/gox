package sets

import "singer.com/basic/containers"

type Set[T comparable] interface {
	containers.Container[T]

	Add(items ...T)
	Remove(items ...T)
	Contains(items ...T) bool
	IsEqual(t Set[T]) bool
	// check whether the t is the subset of caller
	IsSubset(t Set[T]) bool
	IsSuperset(t Set[T]) bool
	Each(func(T) bool)
	Copy() Set[T]
	Merge(t Set[T])
	Separate(t Set[T])
	Union(sets ...Set[T]) Set[T]
	Difference(sets ...Set[T]) Set[T]
	Intersect(t Set[T]) Set[T]
	Intersects(sets ...Set[T]) Set[T]
}

func New[T comparable](threadsafe bool, values ...T) Set[T] {
	if !threadsafe {
		return newNonThreadSafe(values...)
	}
	return newThreadSafe(values...)
}
