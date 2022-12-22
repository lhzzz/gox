package containers

type Iterator[T any] interface {
	Next() Iterator[T]
	Value() T
	Index() int
}
