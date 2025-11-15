package hamt

import (
	"testing"
)

func TestNewBitmapIndexedNode(t *testing.T) {
	bitmap := Initialize().Next(1 << 5)
	child := NewLeafNode[string, int](12345, "key", 42)
	children := []Node[string, int]{child}

	node := NewBitmapIndexedNode(bitmap, children)
	if node == nil {
		t.Fatal("Expected non-nil node")
	}

	if node.bitmap != bitmap {
		t.Errorf("Expected bitmap %d, got %d", bitmap, node.bitmap)
	}

	if len(node.children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(node.children))
	}
}

func TestBitmapIndexedNodeKeyValue(t *testing.T) {
	bitmap := Initialize().Next(1 << 5)
	child := NewLeafNode[string, int](12345, "key", 42)
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{child})

	// BitmapIndexedNode doesn't store key/value directly
	key := node.Key()
	if key != "" {
		t.Errorf("Expected zero value for key, got '%s'", key)
	}

	value := node.Value()
	if value != 0 {
		t.Errorf("Expected zero value for value, got %d", value)
	}
}

func TestBitmapIndexedNodeGet(t *testing.T) {
	// Create a simple tree with two leaf nodes
	hash1 := uint64(0b000001) // position 1 at offset 0
	hash2 := uint64(0b000010) // position 2 at offset 0

	leaf1 := NewLeafNode(hash1, "key1", 100)
	leaf2 := NewLeafNode(hash2, "key2", 200)

	bitmap := Initialize()
	bitmap = bitmap.Next(bitmap.Position(hash1, 0))
	bitmap = bitmap.Next(bitmap.Position(hash2, 0))

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1, leaf2})

	// Get existing values
	value, found := node.Get(hash1, 0)
	if !found {
		t.Error("Expected to find first value")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}

	value, found = node.Get(hash2, 0)
	if !found {
		t.Error("Expected to find second value")
	}
	if value != 200 {
		t.Errorf("Expected value 200, got %d", value)
	}

	// Get non-existing value
	_, found = node.Get(99999, 0)
	if found {
		t.Error("Expected not to find non-existing value")
	}
}

func TestBitmapIndexedNodeSet(t *testing.T) {
	// Start with a single leaf
	hash1 := uint64(0b000001)
	leaf1 := NewLeafNode(hash1, "key1", 100)

	bm := Initialize()
	bitmap := bm.Next(bm.Position(hash1, 0))
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1})

	// Add another value with different position
	hash2 := uint64(0b000010)
	newNode := node.Set("key2", 200, hash2, 0)

	// Both values should be accessible
	value1, found1 := newNode.Get(hash1, 0)
	value2, found2 := newNode.Get(hash2, 0)

	if !found1 || !found2 {
		t.Error("Expected to find both values")
	}

	if value1 != 100 || value2 != 200 {
		t.Errorf("Expected values 100 and 200, got %d and %d", value1, value2)
	}
}

func TestBitmapIndexedNodeSetUpdate(t *testing.T) {
	hash := uint64(0b000001)
	leaf := NewLeafNode(hash, "key1", 100)

	bm := Initialize()
	bitmap := bm.Next(bm.Position(hash, 0))
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Update existing value
	newNode := node.Set("key1", 999, hash, 0)

	value, found := newNode.Get(hash, 0)
	if !found {
		t.Error("Expected to find updated value")
	}
	if value != 999 {
		t.Errorf("Expected value 999, got %d", value)
	}
}

func TestBitmapIndexedNodeRemove(t *testing.T) {
	hash1 := uint64(0b000001)
	hash2 := uint64(0b000010)

	leaf1 := NewLeafNode(hash1, "key1", 100)
	leaf2 := NewLeafNode(hash2, "key2", 200)

	bitmap := Initialize()
	pos1 := bitmap.Position(hash1, 0)
	pos2 := bitmap.Position(hash2, 0)
	bitmap = bitmap.Next(pos1).Next(pos2)

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1, leaf2})

	// Remove first value
	newNode, removed := node.Remove(hash1, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	_, found := newNode.Get(hash1, 0)
	if found {
		t.Error("Expected first value to be removed")
	}

	value, found := newNode.Get(hash2, 0)
	if !found {
		t.Error("Expected second value to still exist")
	}
	if value != 200 {
		t.Errorf("Expected value 200, got %d", value)
	}
}

func TestBitmapIndexedNodeRemoveNonExistent(t *testing.T) {
	hash := uint64(0b000001)
	leaf := NewLeafNode(hash, "key1", 100)

	bm := Initialize()
	bitmap := bm.Next(bm.Position(hash, 0))
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Try to remove non-existent value
	newNode, removed := node.Remove(99999, 0)
	if removed {
		t.Error("Expected removal to fail")
	}

	if newNode != node {
		t.Error("Expected node to be unchanged")
	}
}

func TestBitmapIndexedNodeRemoveAll(t *testing.T) {
	hash := uint64(0b000001)
	leaf := NewLeafNode(hash, "key1", 100)

	bm := Initialize()
	bitmap := bm.Next(bm.Position(hash, 0))
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Remove the only value
	newNode, removed := node.Remove(hash, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	if newNode != nil {
		t.Error("Expected nil node after removing all children")
	}
}

func TestBitmapIndexedNodeToSlice(t *testing.T) {
	hash1 := uint64(0b000001)
	hash2 := uint64(0b000010)

	leaf1 := NewLeafNode(hash1, "apple", 100)
	leaf2 := NewLeafNode(hash2, "banana", 200)

	bitmap := Initialize()
	bitmap = bitmap.Next(bitmap.Position(hash1, 0))
	bitmap = bitmap.Next(bitmap.Position(hash2, 0))

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1, leaf2})

	slice := node.ToSlice()
	if len(slice) != 2 {
		t.Errorf("Expected slice length 2, got %d", len(slice))
	}

	// Verify all entries are present (order may vary)
	found := make(map[string]int)
	for _, entry := range slice {
		found[entry.Key] = entry.Value
	}

	if found["apple"] != 100 {
		t.Error("Expected to find 'apple' with value 100")
	}

	if found["banana"] != 200 {
		t.Error("Expected to find 'banana' with value 200")
	}
}

func TestReplaceNode(t *testing.T) {
	leaf1 := NewLeafNode[string, int](1, "key1", 100)
	leaf2 := NewLeafNode[string, int](2, "key2", 200)
	leaf3 := NewLeafNode[string, int](3, "key3", 300)

	children := []Node[string, int]{leaf1, leaf2}

	newChildren := replaceNode(children, 1, leaf3)

	if len(newChildren) != 2 {
		t.Errorf("Expected 2 children, got %d", len(newChildren))
	}

	if newChildren[0] != leaf1 {
		t.Error("Expected first child to be unchanged")
	}

	if newChildren[1] != leaf3 {
		t.Error("Expected second child to be replaced")
	}
}

func TestInsertNode(t *testing.T) {
	leaf1 := NewLeafNode[string, int](1, "key1", 100)
	leaf2 := NewLeafNode[string, int](2, "key2", 200)
	leaf3 := NewLeafNode[string, int](3, "key3", 300)

	children := []Node[string, int]{leaf1, leaf3}

	// Insert at index 1
	newChildren := insertNode(children, 1, leaf2)

	if len(newChildren) != 3 {
		t.Errorf("Expected 3 children, got %d", len(newChildren))
	}

	if newChildren[0] != leaf1 {
		t.Error("Expected first child to be leaf1")
	}

	if newChildren[1] != leaf2 {
		t.Error("Expected second child to be leaf2")
	}

	if newChildren[2] != leaf3 {
		t.Error("Expected third child to be leaf3")
	}
}

func TestBitmapIndexedNodeRemoveNodeMethod(t *testing.T) {
	leaf1 := NewLeafNode[string, int](1, "key1", 100)
	leaf2 := NewLeafNode[string, int](2, "key2", 200)
	leaf3 := NewLeafNode[string, int](3, "key3", 300)

	bitmap := Initialize()
	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1, leaf2, leaf3})

	newChildren := node.removeNode(1)

	if len(newChildren) != 2 {
		t.Errorf("Expected 2 children, got %d", len(newChildren))
	}

	if newChildren[0] != leaf1 {
		t.Error("Expected first child to be leaf1")
	}

	if newChildren[1] != leaf3 {
		t.Error("Expected second child to be leaf3")
	}
}
