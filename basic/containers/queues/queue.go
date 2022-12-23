package queues

import "singer.com/basic/containers"

type Queue[T comparable] interface {
	containers.Container[T]

	Push(value T)
	Pop() T
}
