package arraylist

import "singer.com/basic/containers"

func (l *List[T]) Begin() containers.Iterator[T] {
	return l.begin()
}

func (l *List[T]) End() containers.Iterator[T] {
	return l.end()
}

func (l *List[T]) Find(value T) containers.Iterator[T] {
	for it := l.Begin(); it != l.End(); it = it.Next() {
		if it.Value() == value {
			return it
		}
	}
	return l.End()
}

func (l *List[T]) Each(f func(index int, value T)) {
	for it := l.Begin(); it != l.End(); it = it.Next() {
		f(it.Index(), it.Value())
	}
}
