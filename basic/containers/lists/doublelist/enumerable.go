package doublelist

import "singer.com/basic/containers"

func (l *List[T]) Begin() containers.Iterator[T] {
	return l.begin()
}

func (l *List[T]) End() containers.Iterator[T] {
	return l.end()
}

func (l *List[T]) Each(f func(value T) bool) {
	for it := l.Begin(); it != l.End(); it = it.Next() {
		if !f(it.Value()) {
			break
		}
	}
}

func (l *List[T]) Find(value T) containers.Iterator[T] {
	for it := l.Begin(); it != l.End(); it = it.Next() {
		if it.Value() == value {
			return it
		}
	}
	return l.End()
}
