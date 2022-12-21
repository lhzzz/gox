package containers

type Iterator[T any] interface {
	Next() bool

	Value() T

	Index() int

	Begin()

	First() bool
}
