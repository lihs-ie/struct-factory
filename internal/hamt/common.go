package hamt

type Entry[K any, V any] struct {
	Key   K
	Value V
}

type Node[K any, V any] interface {
	Key() K
	Value() V
	Get(hash uint64, offset int) (V, bool)
	Set(key K, value V, hash uint64, offset int) Node[K, V]
	Remove(hash uint64, offset int) (Node[K, V], bool)
	ToSlice() []Entry[K, V]
}
