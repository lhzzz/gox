package containers

// Container is base interface that all data structures implement.
type Container[T any] interface {
	Empty() bool
	Size() int
	Clear()
	Values() []T
	String() string
}
