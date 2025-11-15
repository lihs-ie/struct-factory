package hamt

import (
	"testing"
)

// Edge case tests to achieve maximum coverage

func TestBitmapIndexedNodeDeepRecursion(t *testing.T) {
	// Test Set with offset increment (line 43)
	// We need a scenario where BitmapIndexedNode recursively calls Set on a child
	// with an incremented offset

	hash1 := Hash("deep1")
	hash2 := Hash("deep2")

	leaf1 := NewLeafNode(hash1, "deep1", 100)

	bitmap := Initialize()
	pos1 := bitmap.Position(hash1, 0)
	bitmap = bitmap.Next(pos1)

	// Create a BitmapIndexedNode with one child
	var node Node[string, int] = NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf1})

	// Set another value - this will trigger recursive Set on the child
	node = node.Set("deep2", 200, hash2, 0)

	// Verify both values are accessible
	value1, found1 := node.Get(hash1, 0)
	if !found1 || value1 != 100 {
		t.Errorf("Expected to find value 100 for hash1, got %d (found: %v)", value1, found1)
	}

	value2, found2 := node.Get(hash2, 0)
	if !found2 || value2 != 200 {
		t.Errorf("Expected to find value 200 for hash2, got %d (found: %v)", value2, found2)
	}
}

func TestBitmapIndexedNodeSetUnchanged(t *testing.T) {
	// Test the case where next == target (line 45-46)
	// This happens when the child's Set returns itself (no change)

	hash := Hash("unchanged")

	leaf := NewLeafNode(hash, "unchanged", 100)

	bitmap := Initialize()
	pos := bitmap.Position(hash, 0)
	bitmap = bitmap.Next(pos)

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Set the same value - the leaf will return itself
	newNode := node.Set("unchanged", 100, hash, 0)

	// Even though no change, we should still be able to get the value
	value, found := newNode.Get(hash, 0)
	if !found || value != 100 {
		t.Errorf("Expected value 100, got %d (found: %v)", value, found)
	}
}

func TestBitmapIndexedNodeRemoveDeepRecursion(t *testing.T) {
	// Test Remove with offset increment (line 72)
	hash1 := Hash("remove_deep1")
	hash2 := Hash("remove_deep2")
	hash3 := Hash("remove_deep3")

	var root Node[string, int]

	// Build a tree with multiple values
	root = NewLeafNode(hash1, "remove_deep1", 100)
	root = root.Set("remove_deep2", 200, hash2, 0)
	root = root.Set("remove_deep3", 300, hash3, 0)

	// Remove one value from deep in the tree
	newRoot, removed := root.Remove(hash2, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	// Verify the value is gone
	_, found := newRoot.Get(hash2, 0)
	if found {
		t.Error("Expected removed value to be gone")
	}

	// Verify other values still exist
	value1, found1 := newRoot.Get(hash1, 0)
	if !found1 || value1 != 100 {
		t.Error("Expected first value to remain")
	}

	value3, found3 := newRoot.Get(hash3, 0)
	if !found3 || value3 != 300 {
		t.Error("Expected third value to remain")
	}
}

func TestBitmapIndexedNodeRemoveUnchanged(t *testing.T) {
	// Test the case where target == nextNode in Remove (line 78-79)
	hash1 := Hash("remove_unchanged1")
	hash2 := Hash("remove_unchanged2")

	var root Node[string, int]

	root = NewLeafNode(hash1, "remove_unchanged1", 100)
	root = root.Set("remove_unchanged2", 200, hash2, 0)

	// Try to remove a non-existent value at a deeper level
	nonExistentHash := Hash("nonexistent_deep")
	_, removed := root.Remove(nonExistentHash, 0)

	if removed {
		t.Error("Expected removal of non-existent value to fail")
	}

	// Verify original values still exist
	value1, found1 := root.Get(hash1, 0)
	if !found1 || value1 != 100 {
		t.Error("Expected first value to remain unchanged")
	}

	value2, found2 := root.Get(hash2, 0)
	if !found2 || value2 != 200 {
		t.Error("Expected second value to remain unchanged")
	}
}

func TestBitmapIndexedNodeRemoveEmptyResult(t *testing.T) {
	// Test when removing results in empty children (line 86-88)
	hash := Hash("only_value")

	leaf := NewLeafNode(hash, "only_value", 100)

	bitmap := Initialize()
	pos := bitmap.Position(hash, 0)
	bitmap = bitmap.Next(pos)

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

func TestBitmapIndexedNodeReplaceChild(t *testing.T) {
	// Test when nextNode is not nil and we replace a child (line 96)
	hash1 := Hash("replace1")
	hash2 := Hash("replace2")
	hash3 := Hash("replace3")

	var root Node[string, int]

	root = NewLeafNode(hash1, "replace1", 100)
	root = root.Set("replace2", 200, hash2, 0)
	root = root.Set("replace3", 300, hash3, 0)

	// Remove one value - this should replace a child in BitmapIndexedNode
	newRoot, removed := root.Remove(hash2, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	// Verify the structure still works
	value1, found1 := newRoot.Get(hash1, 0)
	if !found1 || value1 != 100 {
		t.Error("Expected first value to remain")
	}

	value3, found3 := newRoot.Get(hash3, 0)
	if !found3 || value3 != 300 {
		t.Error("Expected third value to remain")
	}

	_, found2 := newRoot.Get(hash2, 0)
	if found2 {
		t.Error("Expected removed value to be gone")
	}
}
