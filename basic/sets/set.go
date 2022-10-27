package sets

type Set interface {
	Add(items ...interface{})
	Remove(items ...interface{})
	Contains(items ...interface{}) bool
	Size() int
	Clear()
	IsEmpty() bool
	IsEqual(t Set) bool
	// check whether the t is the subset of caller
	IsSubset(t Set) bool
	IsSuperset(t Set) bool
	Each(func(interface{}) bool)
	String() string
	List() []interface{}
	Copy() Set
	Merge(t Set)
	Separate(t Set)
	Union(sets ...Set) Set
	Difference(sets ...Set) Set
	Intersect(t Set) Set
	Intersects(sets ...Set) Set
}

func New(threadsafe bool) Set {
	if !threadsafe {
		return newNonThreadSafe()
	}
	return newThreadSafe()
}
