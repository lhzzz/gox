package doublelist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArrayList(t *testing.T) {
	l := New(1, 2, 3, 4, 5, 6)

	assert.Equal(t, l.Size(), 6)

	l.Add(7)
	assert.Equal(t, l.Size(), 7)

	i, b := l.Get(0)
	t.Log(i)
	assert.True(t, b)

	l.Clear()
	assert.True(t, l.Empty())

	assert.False(t, l.Contains(1))

	i, b = l.Get(0)
	t.Log(i)
	assert.False(t, b)

	l.Insert(0, 10, 11, 12, 13)
	l.Remove(0)
	assert.Equal(t, l.Values(), []int{11, 12, 13})

	l.Swap(0, 1)
	assert.Equal(t, l.Values(), []int{12, 11, 13})

	l.Sort(func(a, b int) int {
		if a > b {
			return 1
		} else if a < b {
			return -1
		}
		return 0
	})
	t.Log(l.String())
}

func TestBegin(t *testing.T) {
	list := New("1", "2", "3")

	it := list.Begin()
	assert.Equal(t, it.Index(), 0)
	assert.Equal(t, it.Value(), "1")
}

func TestEnd(t *testing.T) {
	list := New("1", "2", "3")

	it := list.End()
	assert.Equal(t, it.Index(), list.Size())
}

func TestFind(t *testing.T) {
	list := New("1", "2", "3")

	assert.NotEqual(t, list.Find("3"), list.End())
	assert.Equal(t, list.Find("10"), list.End())
}

func TestEach(t *testing.T) {
	list := New(1, 2, 3, 4, 5, 6, 7)
	sum := 0
	list.Each(func(value int) bool {
		sum += value
		return true
	})
	t.Log(sum)

	list2 := New("a", "b", "c", "d")
	str := ""
	list2.Each(func(value string) bool {
		str += value
		return true
	})
	t.Log(str)
}

func TestIterator(t *testing.T) {
	list := New(6, 5, 4, 3, 2, 1)

	for it := list.Begin(); it != list.End(); it = it.Next() {
		t.Log(it.Index(), it.Value())
	}

	it := list.End()
	for it = it.Prev(); it != list.Begin(); it = it.Prev() {
		t.Log(it.Index(), it.Value())
	}
	t.Log(it.Index(), it.Value())
}
