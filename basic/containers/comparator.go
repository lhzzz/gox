package containers

type Comparator[T comparable] func(a, b T) int
