package arraylist

import (
	"singer.com/basic/containers"
)

type iterator[T comparable] struct {
	l     *List[T]
	index int
}

func (l *List[T]) begin() containers.Iterator[T] {
	return iterator[T]{l: l, index: 0}
}

func (l *List[T]) end() containers.Iterator[T] {
	return iterator[T]{l: l, index: l.size}
}

func (it iterator[T]) Next() containers.Iterator[T] {
	index := it.index + 1
	return iterator[T]{l: it.l, index: index}
}

func (it iterator[T]) Value() T {
	return it.l.elems[it.index]
}

func (it iterator[T]) Index() int {
	return it.index
}
