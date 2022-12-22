package arraylist

import (
	"testing"
)

func TestIterator(t *testing.T) {
	list := New(6, 5, 4, 3, 2, 1)

	for it := list.Begin(); it != list.End(); it = it.Next() {
		t.Log(it.Index(), it.Value())
	}
}
