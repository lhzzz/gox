package arraylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArrayList(t *testing.T) {
	l := New(1, 2, 3, 4, 5, 6)

	assert.Equal(t, l.Size(), 6)

	l.Add(7)
	assert.Equal(t, l.Size(), 7)

	l.Clear()
	assert.True(t, l.Empty())

	assert.False(t, l.Contains(1))

	i, b := l.Get(0)
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
