package lists

import (
	"singer.com/basic/containers"
)

type List[T comparable] interface {
	containers.Container[T]

	Add(values ...T)
	Remove(index int)
	Get(index int) (T, bool)
	Contains(values ...T) bool
	Sort(comparator containers.Comparator[T])
	Swap(index1, index2 int)
	Insert(index int, values ...T)
}
