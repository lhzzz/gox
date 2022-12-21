package arraylist

import (
	"fmt"
	"strings"

	"singer.com/basic/containers"
	"singer.com/basic/containers/lists"
)

type List[T comparable] struct {
	elems []T
	size  int
}

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

func New[T comparable](values ...T) lists.List[T] {
	list := &List[T]{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (l *List[T]) Empty() bool {
	return l.size == 0
}

func (l *List[T]) Size() int {
	return l.size
}

func (l *List[T]) Clear() {
	l.size = 0
	l.elems = make([]T, 0, 0)
}

func (l *List[T]) Values() []T {
	newElems := make([]T, l.size, l.size)
	copy(newElems, l.elems[:l.size])
	return newElems
}

func (l *List[T]) String() string {
	str := "ArrayList\n"
	values := []string{}
	for _, value := range l.elems[:l.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (l *List[T]) Add(values ...T) {
	l.growBy(len(values))
	for _, value := range values {
		l.elems[l.size] = value
		l.size++
	}
}

func (l *List[T]) Remove(index int) {
	if !l.withinRange(index) {
		return
	}
	var zero T
	l.elems[index] = zero
	copy(l.elems[index:], l.elems[index+1:l.size])
	l.size--
	l.shrink()
}

func (l *List[T]) Get(index int) (T, bool) {
	var zero T
	if !l.withinRange(index) {
		return zero, false
	}
	return l.elems[index], true
}

func (l *List[T]) Contains(values ...T) bool {
	for _, v := range values {
		found := false
		for i := 0; i < l.size; i++ {
			if l.elems[i] == v {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (l *List[T]) Sort(comparator containers.Comparator[T]) {
	if len(l.elems) < 2 {
		return
	}
	containers.Sort(l.elems[:l.size], comparator)
}

// Swap swaps the two values at the specified positions.
func (list *List[T]) Swap(i, j int) {
	if list.withinRange(i) && list.withinRange(j) {
		list.elems[i], list.elems[j] = list.elems[j], list.elems[i]
	}
}

//if index larger than the list size, it will append the values to list tail
func (l *List[T]) Insert(index int, values ...T) {
	if !l.withinRange(index) {
		if index == l.size {
			l.Add(values...)
		}
		return
	}

	length := len(values)
	l.growBy(length)
	l.size += length
	copy(l.elems[index+length:], l.elems[index:l.size-length])
	copy(l.elems[index:], values)
}

func (l *List[T]) withinRange(index int) bool {
	return index >= 0 && index < l.size
}

func (l *List[T]) growBy(n int) {
	currCap := cap(l.elems)
	if l.size+n >= currCap {
		newCap := int(growthFactor * float32(currCap+n))
		l.resize(newCap)
	}
}

func (list *List[T]) resize(cap int) {
	newElements := make([]T, cap, cap)
	copy(newElements, list.elems)
	list.elems = newElements
}

func (l *List[T]) shrink() {
	if shrinkFactor-0 < 0.00001 {
		return
	}
	currCap := cap(l.elems)
	if l.size <= int(float32(currCap)*shrinkFactor) {
		l.resize(l.size)
	}
}
