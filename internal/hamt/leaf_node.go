package hamt

type LeafNode[K any, V any] struct {
	Node[K, V]
	hash  uint64
	key   K
	value V
}

func NewLeafNode[K any, V any](hash uint64, key K, value V) *LeafNode[K, V] {
	return &LeafNode[K, V]{
		hash:  hash,
		key:   key,
		value: value,
	}
}

func (leaf *LeafNode[K, V]) Key() K {
	return leaf.key
}

func (leaf *LeafNode[K, V]) Value() V {
	return leaf.value
}

func (leaf *LeafNode[K, V]) Get(hash uint64, offset int) (V, bool) {
	if leaf.hash == hash {
		return leaf.value, true
	}

	var zero V

	return zero, false
}

func (leaf *LeafNode[K, V]) Set(key K, value V, hash uint64, offset int) Node[K, V] {
	if leaf.hash == hash {
		// Same hash - update the value
		return NewLeafNode(hash, key, value)
	}

	// Different hashes - create a BitmapIndexedNode
	bitmap := Initialize()
	position1 := bitmap.Position(leaf.hash, offset)
	position2 := bitmap.Position(hash, offset)

	if position1 == position2 {
		// Positions collide at this level, need to go deeper
		nextNode := leaf.Set(key, value, hash, offset+1)
		return NewBitmapIndexedNode(
			bitmap.Next(position1),
			[]Node[K, V]{nextNode},
		)
	}

	// Different positions - create bitmap with both nodes
	bitmapNode := NewBitmapIndexedNode(
		bitmap.Next(position1),
		[]Node[K, V]{leaf},
	)

	return bitmapNode.Set(key, value, hash, offset)
}

func (leaf *LeafNode[K, V]) Remove(hash uint64, offset int) (Node[K, V], bool) {
	if leaf.hash == hash {
		return nil, true
	}

	return leaf, false
}

func (leaf *LeafNode[K, V]) ToSlice() []Entry[K, V] {
	return []Entry[K, V]{
		{Key: leaf.key, Value: leaf.value},
	}
}
