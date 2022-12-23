package containers

type Iterator[T any] interface {
	Value() T
	Index() int
	Next() Iterator[T]
	Prev() Iterator[T]
}
