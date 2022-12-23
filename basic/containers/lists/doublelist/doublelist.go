package doublelist

import (
	"fmt"
	"strings"

	"singer.com/basic/containers"
	"singer.com/basic/containers/lists"
)

type List[T comparable] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

type Node[T comparable] struct {
	prev  *Node[T]
	next  *Node[T]
	value T
}

func New[T comparable](values ...T) lists.List[T] {
	list := &List[T]{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (l *List[T]) Add(values ...T) {
	for _, v := range values {
		node := &Node[T]{prev: l.tail, value: v}
		if l.size == 0 {
			l.head = node
		} else {
			l.tail.next = node
		}
		l.tail = node
		l.size++
	}
}

func (l *List[T]) Empty() bool {
	return l.size == 0
}

func (l *List[T]) Size() int {
	return l.size
}

func (l *List[T]) Get(index int) (T, bool) {
	var zero T
	if !l.withinRange(index) {
		return zero, false
	}

	//尾->头
	if l.size-index < index {
		node := l.tail
		for tailIdx := l.size - 1; tailIdx != index; tailIdx, node = tailIdx-1, node.prev {
		}
		return node.value, true
	}
	//头->尾
	node := l.head
	for i := 0; i < index; i, node = i+1, node.next {
	}
	return node.value, true
}

func (l *List[T]) Insert(index int, values ...T) {
	if !l.withinRange(index) {
		if index == l.size {
			l.Add(values...)
		}
		return
	}
	l.size += len(values)

	prefound := (*Node[T])(nil)
	foundNode := (*Node[T])(nil)
	if l.size-index < index {
		foundNode = l.tail
		for tailIdx := l.size - 1; tailIdx != index; tailIdx, foundNode = tailIdx-1, foundNode.prev {
			prefound = foundNode.prev
		}
	} else {
		foundNode = l.head
		for i := 0; i < index; i, foundNode = i+1, foundNode.next {
			prefound = foundNode
		}
	}

	if foundNode == l.head {
		old := l.head
		for i, v := range values {
			node := &Node[T]{value: v}
			if i == 0 {
				l.head = node
			} else {
				node.prev = prefound
				prefound.next = node
			}
			prefound = node
		}
		old.prev = prefound
		prefound.next = old
	} else {
		old := foundNode
		for _, v := range values {
			node := &Node[T]{value: v}
			node.prev = prefound
			prefound.next = node
			prefound = node
		}
		old.prev = prefound
		prefound.next = old
	}
}

func (l *List[T]) Remove(index int) {
	if !l.withinRange(index) {
		return
	}

	foundNode := (*Node[T])(nil)
	if l.size-index < index {
		foundNode = l.tail
		for tailIdx := l.size - 1; tailIdx != index; tailIdx, foundNode = tailIdx-1, foundNode.prev {
		}
	} else {
		foundNode = l.head
		for i := 0; i < index; i, foundNode = i+1, foundNode.next {
		}
	}

	if foundNode == l.head {
		l.head = l.head.next
	}
	if foundNode == l.tail {
		l.tail = l.tail.prev
	}
	if foundNode.prev != nil {
		foundNode.prev.next = foundNode.next
	}
	if foundNode.next != nil {
		foundNode.next.prev = foundNode.prev
	}
	foundNode = nil
	l.size--
}

func (l *List[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}

func (l *List[T]) Values() []T {
	list := make([]T, l.size, l.size)
	for i, node := 0, l.head; i < l.size; i, node = i+1, node.next {
		list[i] = node.value
	}
	return list
}

func (l *List[T]) Contains(values ...T) bool {
	if len(values) == 0 {
		return true
	}

	if l.size == 0 {
		return false
	}

	for _, v := range values {
		flag := false
		for node := l.head; node != nil; node = node.next {
			if node.value == v {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func (l *List[T]) String() string {
	str := "DoubeList\n"
	values := []string{}
	for node := l.head; node != nil; node = node.next {
		values = append(values, fmt.Sprintf("%v", node.value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (l *List[T]) Swap(index1 int, index2 int) {
	if !l.withinRange(index1) || !l.withinRange(index2) || index1 == index2 {
		return
	}
	node := l.head
	node1 := (*Node[T])(nil)
	node2 := (*Node[T])(nil)
	for i := 0; i < l.size; i++ {
		if i == index1 {
			node1 = node
		} else if i == index2 {
			node2 = node
		}
		if node1 != nil && node2 != nil {
			break
		}
		node = node.next
	}
	node1.value, node2.value = node2.value, node1.value
}

func (l *List[T]) Sort(comparator containers.Comparator[T]) {
	if l.size <= 1 {
		return
	}

	array := l.Values()
	containers.Sort(array, comparator)
	l.Clear()
	l.Add(array...)
}

func (list *List[T]) withinRange(index int) bool {
	return index >= 0 && index < list.size
}
