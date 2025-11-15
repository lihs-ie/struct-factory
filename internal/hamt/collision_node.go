package hamt

type CollisionNode[K any, V any] struct {
	Node[K, V]
	hash    uint64
	entries []Entry[K, V]
}

func NewCollisionNode[K any, V any](hash uint64, entries []Entry[K, V]) *CollisionNode[K, V] {
	return &CollisionNode[K, V]{
		hash:    hash,
		entries: entries,
	}
}

func (node *CollisionNode[K, V]) Key() K {
	if len(node.entries) > 0 {
		return node.entries[0].Key
	}
	return *new(K)
}

func (node *CollisionNode[K, V]) Value() V {
	if len(node.entries) > 0 {
		return node.entries[0].Value
	}
	return *new(V)
}

func (node *CollisionNode[K, V]) Get(hash uint64, offset int) (V, bool) {
	if node.hash != hash {
		return *new(V), false
	}

	if len(node.entries) > 0 {
		return node.entries[0].Value, true
	}

	return *new(V), false
}

func (node *CollisionNode[K, V]) Set(key K, value V, hash uint64, offset int) Node[K, V] {
	if node.hash == hash {

		newEntries := make([]Entry[K, V], len(node.entries)+1)
		copy(newEntries, node.entries)
		newEntries[len(node.entries)] = Entry[K, V]{Key: key, Value: value}
		return NewCollisionNode(hash, newEntries)
	}

	return NewLeafNode(hash, key, value)
}

func (node *CollisionNode[K, V]) Remove(hash uint64, offset int) (Node[K, V], bool) {
	if node.hash != hash {
		return node, false
	}

	if len(node.entries) == 1 {
		return nil, true
	}

	if len(node.entries) == 2 {
		remaining := node.entries[1]
		return NewLeafNode(node.hash, remaining.Key, remaining.Value), true
	}

	newEntries := make([]Entry[K, V], len(node.entries)-1)
	copy(newEntries, node.entries[1:])
	return NewCollisionNode(node.hash, newEntries), true
}

func (node *CollisionNode[K, V]) ToSlice() []Entry[K, V] {
	result := make([]Entry[K, V], len(node.entries))
	copy(result, node.entries)
	return result
}
