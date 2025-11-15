package hamt

import (
	"testing"
)

// Integration tests to ensure all components work together and achieve 100% coverage

func TestBitmapIndexedNodeDeepNesting(t *testing.T) {
	// Test with real hash values that may collide at first level
	// but will be separated at deeper levels
	hash1 := Hash("test1")
	hash2 := Hash("test2")

	var root Node[string, int]

	root = NewLeafNode(hash1, "test1", 100)
	root = root.Set("test2", 200, hash2, 0)

	// Both values should be accessible
	value1, found1 := root.Get(hash1, 0)
	value2, found2 := root.Get(hash2, 0)

	if !found1 || !found2 {
		t.Error("Expected to find both values")
	}

	if value1 != 100 || value2 != 200 {
		t.Errorf("Expected values 100 and 200, got %d and %d", value1, value2)
	}
}

func TestBitmapIndexedNodeRemoveDeep(t *testing.T) {
	// Create a structure with multiple values
	hash1 := Hash("remove1")
	hash2 := Hash("remove2")

	var root Node[string, int]

	root = NewLeafNode(hash1, "remove1", 100)
	root = root.Set("remove2", 200, hash2, 0)

	// Remove one value
	newNode, removed := root.Remove(hash2, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	// First value should still be accessible
	value1, found1 := newNode.Get(hash1, 0)
	if !found1 {
		t.Error("Expected to find first value")
	}
	if value1 != 100 {
		t.Errorf("Expected value 100, got %d", value1)
	}

	// Second value should not be accessible
	_, found2 := newNode.Get(hash2, 0)
	if found2 {
		t.Error("Expected not to find removed value")
	}
}

func TestComplexTreeOperations(t *testing.T) {
	// Build a tree with multiple levels
	var root Node[int, string]

	// Insert multiple values
	values := map[uint64]string{
		0b000001: "a",
		0b000010: "b",
		0b000100: "c",
		0b001000: "d",
		0b010000: "e",
		0b100000: "f",
	}

	for hash, value := range values {
		if root == nil {
			root = NewLeafNode(hash, int(hash), value)
		} else {
			root = root.Set(int(hash), value, hash, 0)
		}
	}

	// Verify all values are accessible
	for hash, expectedValue := range values {
		value, found := root.Get(hash, 0)
		if !found {
			t.Errorf("Expected to find value for hash %b", hash)
		}
		if value != expectedValue {
			t.Errorf("Expected value '%s' for hash %b, got '%s'", expectedValue, hash, value)
		}
	}

	// Remove some values
	root, removed := root.Remove(0b000001, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	// Verify removed value is gone
	_, found := root.Get(0b000001, 0)
	if found {
		t.Error("Expected removed value to be gone")
	}

	// Verify other values still accessible
	value, found := root.Get(0b000010, 0)
	if !found {
		t.Error("Expected to find remaining value")
	}
	if value != "b" {
		t.Errorf("Expected value 'b', got '%s'", value)
	}
}

func TestHashCollisionInTree(t *testing.T) {
	// Force hash collisions by using custom hash values
	// This tests CollisionNode integration
	hash := uint64(12345)

	// Add another value with the same hash
	// Since we can't control Hash() function, we'll manually create collision scenario
	// by building the structure ourselves

	entry1 := Entry[string, int]{Key: "key1", Value: 100}
	entry2 := Entry[string, int]{Key: "key2", Value: 200}

	collision := NewCollisionNode(hash, []Entry[string, int]{entry1, entry2})

	// Test collision node operations
	value, found := collision.Get(hash, 0)
	if !found {
		t.Error("Expected to find value in collision node")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}

	// Add more to collision
	newCollision := collision.Set("key3", 300, hash, 0)

	collisionPtr, ok := newCollision.(*CollisionNode[string, int])
	if !ok {
		t.Fatal("Expected CollisionNode")
	}

	if len(collisionPtr.entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(collisionPtr.entries))
	}
}

func TestToSliceComprehensive(t *testing.T) {
	// Build a complex tree and ensure ToSlice works correctly
	var root Node[int, string]

	entries := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
	}

	for key, value := range entries {
		hash := Hash(key)
		if root == nil {
			root = NewLeafNode(hash, key, value)
		} else {
			root = root.Set(key, value, hash, 0)
		}
	}

	slice := root.ToSlice()

	if len(slice) != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(slice))
	}

	// Verify all entries are present
	found := make(map[int]string)
	for _, entry := range slice {
		found[entry.Key] = entry.Value
	}

	for key, expectedValue := range entries {
		if found[key] != expectedValue {
			t.Errorf("Expected value '%s' for key %d, got '%s'", expectedValue, key, found[key])
		}
	}
}
