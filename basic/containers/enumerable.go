package containers

type Enumerable[K comparable, T any] interface {
	Begin() Iterator[T]
	End() Iterator[T]
	Each(func(value T) bool)
	Find(value T) Iterator[T]
}
