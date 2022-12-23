package containers

type Enumerable[K comparable, T any] interface {
	Each(func(value T) bool)
	Find(value T) Iterator[T]
}
