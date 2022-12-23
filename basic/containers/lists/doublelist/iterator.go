package doublelist

import "singer.com/basic/containers"

type iterator[T comparable] struct {
	l     *List[T]
	index int
	node  *Node[T]
}

func (l *List[T]) begin() containers.Iterator[T] {
	return iterator[T]{l: l, index: 0, node: l.head}
}

func (l *List[T]) end() containers.Iterator[T] {
	return iterator[T]{l: l, index: l.size, node: nil}
}

func (it iterator[T]) Next() containers.Iterator[T] {
	index := it.index + 1
	return iterator[T]{l: it.l, index: index, node: it.node.next}
}

func (it iterator[T]) Prev() containers.Iterator[T] {
	if it.index == it.l.size {
		return iterator[T]{l: it.l, index: it.l.size - 1, node: it.l.tail}
	}
	index := it.index - 1
	return iterator[T]{l: it.l, index: index, node: it.node.prev}
}

func (it iterator[T]) Value() T {
	return it.node.value
}

func (it iterator[T]) Index() int {
	return it.index
}
