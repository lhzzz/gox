package sets

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadsafeSet(t *testing.T) {
	s := New(true)

	s.Add(1, 1, 2, 3, 4, 5)
	assert.Equal(t, s.Size(), 5)

	s.Remove(1, 2) // s [3,4,5]
	assert.False(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	assert.True(t, s.Contains(3))

	assert.Equal(t, s.Size(), 3)

	s.Clear() // s nil
	assert.True(t, s.IsEmpty())

	s.Add(3, 4, 5) // s [3,4,5]
	s2 := New(true)
	s2.Add(3, 4, 5) // s2 [3,4,5]

	assert.True(t, s.IsEqual(s2))

	s2.Remove(3) // s2 [4,5]

	assert.True(t, s.IsSubset(s2))
	assert.True(t, s2.IsSuperset(s))
	s.Each(func(i interface{}) bool {
		t.Log(i)
		return true
	})
	t.Log(s2.List()...)
	c := s.Copy()
	assert.True(t, s.IsEqual(c))

	s2.Add(7)   // s2: [4,5,7]
	s.Merge(s2) // s: [3,4,5,7]
	assert.True(t, s.Contains(7))

	s.Separate(s2) // s: [3]
	assert.True(t, s.Size() == 1)

	s = s.Union(s2) // s: [3,4,5,7]
	us := New(true)
	us.Add(3, 4, 5, 7)
	assert.True(t, s.IsEqual(us))

	ds := s.Difference(s2) // ds: 3
	t.Log(ds.String())
	dss := New(true)
	dss.Add(3)
	assert.True(t, ds.IsEqual(dss))

	is := s.Intersect(ds) // s [3,4,5,7] ds[3]
	s3 := New(true)
	s3.Add(3)
	assert.True(t, is.IsEqual(s3))

	iss := s.Intersects(s2, ds) //s [3,4,5,7] s2 [4,5,7] ds [3] = iss nil
	assert.True(t, iss.IsEmpty())

	//test concurency
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(3)
		go func(idx int) {
			s.Add(idx*10, idx)
			wg.Done()
		}(i)

		go func(idx int) {
			s.Remove(idx)
			wg.Done()
		}(i)

		go func(idx int) {
			s.Add(idx * 11)
			wg.Done()
		}(i)
	}
	wg.Wait()
	t.Log(s.String())
}

func TestNonThreadsafeSet(t *testing.T) {
	s := New(false)

	s.Add(1, 1, 2, 3, 4, 5)
	assert.Equal(t, s.Size(), 5)

	s.Remove(1, 2) // s [3,4,5]
	assert.False(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	assert.True(t, s.Contains(3))

	assert.Equal(t, s.Size(), 3)

	s.Clear() // s nil
	assert.True(t, s.IsEmpty())

	s.Add(3, 4, 5) // s [3,4,5]
	s2 := New(false)
	s2.Add(3, 4, 5) // s2 [3,4,5]

	assert.True(t, s.IsEqual(s2))

	s2.Remove(3) // s2 [4,5]

	assert.True(t, s.IsSubset(s2))
	assert.True(t, s2.IsSuperset(s))
	s.Each(func(i interface{}) bool {
		t.Log(i)
		return false
	})
	t.Log(s2.List()...)
	c := s.Copy()
	assert.True(t, s.IsEqual(c))

	s2.Add(7)   // s2: [4,5,7]
	s.Merge(s2) // s: [3,4,5,7]
	assert.True(t, s.Contains(7))

	s.Separate(s2) // s: [3]
	assert.True(t, s.Size() == 1)

	s = s.Union(s2) // s: [3,4,5,7]
	us := New(false)
	us.Add(3, 4, 5, 7)
	assert.True(t, s.IsEqual(us))

	ds := s.Difference(s2) // ds: 3
	t.Log(ds.String())
	dss := New(false)
	dss.Add(3)
	assert.True(t, ds.IsEqual(dss))

	is := s.Intersect(ds) // s [3,4,5,7] ds[3]
	s3 := New(false)
	s3.Add(3)
	assert.True(t, is.IsEqual(s3))

	iss := s.Intersects(s2, ds) //s [3,4,5,7] s2 [4,5,7] ds [3] = iss nil
	assert.True(t, iss.IsEmpty())

	t.Log(s)
	//will carsh
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(3)
		go func(idx int) {
			s.Add(idx*10, idx)
			wg.Done()
		}(i)

		go func(idx int) {
			s.Remove(idx)
			wg.Done()
		}(i)

		go func(idx int) {
			s.Add(idx * 11)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
