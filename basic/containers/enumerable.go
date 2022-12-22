package containers

type Enumerable[K comparable, T any] interface {
	Begin() Iterator[T]

	End() Iterator[T]

	Each(func(index K, value T))

	Find(value T) Iterator[T]
}
