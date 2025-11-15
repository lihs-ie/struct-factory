package hamt

type BitmapIndexedNode[K any, V any] struct {
	Node[K, V]
	bitmap   Bitmap
	children []Node[K, V]
}

func NewBitmapIndexedNode[K any, V any](bitmap Bitmap, children []Node[K, V]) *BitmapIndexedNode[K, V] {
	return &BitmapIndexedNode[K, V]{
		bitmap:   bitmap,
		children: children,
	}
}

func (node *BitmapIndexedNode[K, V]) Key() K {
	return *new(K)
}

func (node *BitmapIndexedNode[K, V]) Value() V {
	return *new(V)
}

func (node *BitmapIndexedNode[K, V]) Get(hash uint64, offset int) (V, bool) {
	position := node.bitmap.Position(hash, offset)

	if !node.bitmap.Has(position) {
		return *new(V), false
	}

	index, _ := node.bitmap.Index(position)

	return node.children[index].Get(hash, offset+1)
}

func (node *BitmapIndexedNode[K, V]) Set(key K, value V, hash uint64, offset int) Node[K, V] {
	position := node.bitmap.Position(hash, offset)

	index, _ := node.bitmap.Index(position)

	if node.bitmap.Has(position) {
		target := node.children[index]
		next := target.Set(key, value, hash, offset)

		if next == target {
			return node
		}

		return &BitmapIndexedNode[K, V]{
			bitmap:   node.bitmap,
			children: replaceNode(node.children, index, next),
		}
	}

	nextChildren := insertNode(node.children, index, NewLeafNode(hash, key, value))

	return NewBitmapIndexedNode(
		node.bitmap.Next(position),
		nextChildren,
	)
}

func (node *BitmapIndexedNode[K, V]) Remove(hash uint64, offset int) (Node[K, V], bool) {
	position := node.bitmap.Position(hash, offset)

	if !node.bitmap.Has(position) {
		return node, false
	}

	index, _ := node.bitmap.Index(position)
	target := node.children[index]
	nextNode, exists := target.Remove(hash, offset+1)

	if !exists {
		return node, false
	}

	if target == nextNode {
		return node, false
	}

	if nextNode == nil {
		nextBitmap := node.bitmap.Without(position)
		nextChildren := node.removeNode(index)

		if len(nextChildren) == 0 {
			return nil, true
		}

		return NewBitmapIndexedNode(
			nextBitmap,
			nextChildren,
		), true
	}

	return NewBitmapIndexedNode(node.bitmap, replaceNode(node.children, index, nextNode)), true
}

func (node *BitmapIndexedNode[K, V]) ToSlice() []Entry[K, V] {
	var entries []Entry[K, V]

	for _, child := range node.children {
		entries = append(entries, child.ToSlice()...)
	}

	return entries
}

func replaceNode[K any, V any](children []Node[K, V], index int, node Node[K, V]) []Node[K, V] {
	newChildren := make([]Node[K, V], len(children))

	copy(newChildren, children)
	newChildren[index] = node

	return newChildren
}

func insertNode[K any, V any](children []Node[K, V], index int, node Node[K, V]) []Node[K, V] {
	nextChildren := make([]Node[K, V], len(children)+1)

	copy(nextChildren[:index], children[:index])
	nextChildren[index] = node
	copy(nextChildren[index+1:], children[index:])

	return nextChildren
}

func (node *BitmapIndexedNode[K, V]) removeNode(index int) []Node[K, V] {
	newChildren := make([]Node[K, V], len(node.children)-1)

	copy(newChildren[:index], node.children[:index])
	copy(newChildren[index:], node.children[index+1:])

	return newChildren
}
