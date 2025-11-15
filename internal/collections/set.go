package collections

import "github.com/lihs-ie/struct-factory/internal/hamt"

type void = struct{}

type Set[T any] struct {
	root hamt.Node[T, void]
}

func NewSet[T any](root hamt.Node[T, void]) *Set[T] {
	return &Set[T]{root: root}
}

func NewFromSlice[T any](items []T) *Set[T] {
	var root hamt.Node[T, void]

	for _, item := range items {
		hash := hamt.Hash(item)
		if root == nil {
			root = hamt.NewLeafNode(hash, item, void{})
		} else {
			root = root.Set(item, void{}, hash, 0)
		}
	}

	return NewSet(root)
}

func (set *Set[T]) Set(item T) {
	hash := hamt.Hash(item)
	if set.root == nil {
		set.root = hamt.NewLeafNode(hash, item, void{})
	} else {
		set.root = set.root.Set(item, void{}, hash, 0)
	}
}

func (set *Set[T]) Remove(item T) {
	if set.root == nil {
		return
	}
	hash := hamt.Hash(item)
	newRoot, _ := set.root.Remove(hash, 0)
	set.root = newRoot
}

func (set *Set[T]) Has(item T) bool {
	if set.root == nil {
		return false
	}
	hash := hamt.Hash(item)
	_, found := set.root.Get(hash, 0)
	return found
}

func (set *Set[T]) IsEmpty() bool {
	return set.root == nil
}

func (set *Set[T]) Size() int {
	if set.root == nil {
		return 0
	}
	return len(set.root.ToSlice())
}

func (set *Set[T]) ToSlice() []T {
	if set.root == nil {
		return []T{}
	}
	entries := set.root.ToSlice()
	result := make([]T, len(entries))
	for i, entry := range entries {
		result[i] = entry.Key
	}
	return result
}
