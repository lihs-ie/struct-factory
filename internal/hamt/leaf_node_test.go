package hamt

import (
	"testing"
)

func TestNewLeafNode(t *testing.T) {
	node := NewLeafNode[string, int](12345, "key", 42)
	if node == nil {
		t.Fatal("Expected non-nil node")
	}

	if node.Key() != "key" {
		t.Errorf("Expected key 'key', got '%s'", node.Key())
	}

	if node.Value() != 42 {
		t.Errorf("Expected value 42, got %d", node.Value())
	}
}

func TestLeafNodeGet(t *testing.T) {
	hash := uint64(12345)
	node := NewLeafNode(hash, "apple", 100)

	// Get with matching hash
	value, found := node.Get(hash, 0)
	if !found {
		t.Error("Expected to find value with matching hash")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}

	// Get with non-matching hash
	_, found = node.Get(99999, 0)
	if found {
		t.Error("Expected not to find value with non-matching hash")
	}
}

func TestLeafNodeSetSameHash(t *testing.T) {
	hash := uint64(12345)
	node := NewLeafNode(hash, "key1", 100)

	// Set with same hash (update)
	newNode := node.Set("key1", 200, hash, 0)
	if newNode == nil {
		t.Fatal("Expected non-nil node")
	}

	value, found := newNode.Get(hash, 0)
	if !found {
		t.Error("Expected to find updated value")
	}
	if value != 200 {
		t.Errorf("Expected value 200, got %d", value)
	}
}

func TestLeafNodeSetDifferentHash(t *testing.T) {
	hash1 := uint64(12345)
	hash2 := uint64(67890)

	node := NewLeafNode(hash1, "key1", 100)

	// Set with different hash should create a new structure
	newNode := node.Set("key2", 200, hash2, 0)
	if newNode == nil {
		t.Fatal("Expected non-nil node")
	}

	// Both values should be accessible
	value1, found1 := newNode.Get(hash1, 0)
	value2, found2 := newNode.Get(hash2, 0)

	if !found1 {
		t.Error("Expected to find first value")
	}
	if !found2 {
		t.Error("Expected to find second value")
	}

	if value1 != 100 {
		t.Errorf("Expected first value 100, got %d", value1)
	}
	if value2 != 200 {
		t.Errorf("Expected second value 200, got %d", value2)
	}
}

func TestLeafNodeRemove(t *testing.T) {
	hash := uint64(12345)
	node := NewLeafNode(hash, "key", 42)

	// Remove with matching hash
	newNode, removed := node.Remove(hash, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}
	if newNode != nil {
		t.Error("Expected nil node after removal")
	}

	// Remove with non-matching hash
	newNode, removed = node.Remove(99999, 0)
	if removed {
		t.Error("Expected removal to fail with non-matching hash")
	}
	if newNode != node {
		t.Error("Expected node to be unchanged")
	}
}

func TestLeafNodeToSlice(t *testing.T) {
	node := NewLeafNode[string, int](12345, "apple", 100)

	slice := node.ToSlice()
	if len(slice) != 1 {
		t.Errorf("Expected slice length 1, got %d", len(slice))
	}

	if slice[0].Key != "apple" {
		t.Errorf("Expected key 'apple', got '%s'", slice[0].Key)
	}

	if slice[0].Value != 100 {
		t.Errorf("Expected value 100, got %d", slice[0].Value)
	}
}
