package arraylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	list.Each(func(index, value int) {
		sum += value
	})
	t.Log(sum)

	list2 := New("a", "b", "c", "d")
	str := ""
	list2.Each(func(index int, value string) {
		str += value
	})
	t.Log(str)
}
