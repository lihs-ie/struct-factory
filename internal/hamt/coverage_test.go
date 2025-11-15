package hamt

import (
	"testing"
)

// Tests to achieve 100% coverage

func TestCollisionNodeToSlice(t *testing.T) {
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
		{Key: "key3", Value: 300},
	}

	node := NewCollisionNode[string, int](12345, entries)

	slice := node.ToSlice()

	if len(slice) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(slice))
	}

	for i, entry := range slice {
		if entry.Key != entries[i].Key || entry.Value != entries[i].Value {
			t.Errorf("Expected entry %d to be %v, got %v", i, entries[i], entry)
		}
	}
}

func TestLeafNodeSetPositionCollision(t *testing.T) {
	// Create hashes that collide at offset 0 but differ at offset 1
	// We need to find two hashes where Position(hash, 0) is the same
	// but Position(hash, 1) is different

	// Using simple hashes for testing
	hash1 := uint64(0b000001) // bits 0-5: 1
	hash2 := uint64(0b000001 | (0b000001 << 6)) // bits 0-5: 1, bits 6-11: 1

	leaf1 := NewLeafNode(hash1, "key1", 100)

	// Check if they have the same position at offset 0
	bitmap := Initialize()
	pos1 := bitmap.Position(hash1, 0)
	pos2 := bitmap.Position(hash2, 0)

	if pos1 != pos2 {
		// Find different hashes for this test
		t.Skip("Need to find hashes that collide at offset 0")
	}

	newNode := leaf1.Set("key2", 200, hash2, 0)

	// Should be able to get both values
	value1, found1 := newNode.Get(hash1, 0)
	value2, found2 := newNode.Get(hash2, 0)

	if !found1 || !found2 {
		t.Error("Expected to find both values after position collision")
	}

	if value1 != 100 {
		t.Errorf("Expected value 100 for hash1, got %d", value1)
	}

	if value2 != 200 {
		t.Errorf("Expected value 200 for hash2, got %d", value2)
	}
}

func TestBitmapIndexedNodeSetNoChange(t *testing.T) {
	// Test the case where Set returns the same node (no change)
	hash := uint64(0b000001)
	leaf := NewLeafNode(hash, "key1", 100)

	bitmap := Initialize()
	pos := bitmap.Position(hash, 0)
	bitmap = bitmap.Next(pos)

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Set the same value - should return a different node but with same structure
	newNode := node.Set("key1", 100, hash, 0)

	// Verify the value is still accessible
	value, found := newNode.Get(hash, 0)
	if !found {
		t.Error("Expected to find value")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}
}

func TestBitmapIndexedNodeRemoveEdgeCases(t *testing.T) {
	// Test removal when target == nextNode (no change)
	hash := uint64(0b000001)
	leaf := NewLeafNode(hash, "key1", 100)

	bitmap := Initialize()
	pos := bitmap.Position(hash, 0)
	bitmap = bitmap.Next(pos)

	node := NewBitmapIndexedNode(bitmap, []Node[string, int]{leaf})

	// Try to remove a non-existent deep value
	deepHash := uint64(0b000001 | (0b000010 << 6))
	_, removed := node.Remove(deepHash, 0)

	if removed {
		t.Error("Expected removal of non-existent value to fail")
	}
}

func TestCompleteHAMTWorkflow(t *testing.T) {
	// Build a comprehensive tree that exercises all code paths
	var root Node[string, int]

	// Test data with various hash patterns
	testData := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
		"g": 7,
		"h": 8,
		"i": 9,
		"j": 10,
	}

	// Insert all values
	for key, value := range testData {
		hash := Hash(key)
		if root == nil {
			root = NewLeafNode(hash, key, value)
		} else {
			root = root.Set(key, value, hash, 0)
		}
	}

	// Verify all values are present
	for key, expectedValue := range testData {
		hash := Hash(key)
		value, found := root.Get(hash, 0)
		if !found {
			t.Errorf("Expected to find key '%s'", key)
		}
		if value != expectedValue {
			t.Errorf("Expected value %d for key '%s', got %d", expectedValue, key, value)
		}
	}

	// Update some values
	for key := range testData {
		if key <= "e" {
			hash := Hash(key)
			root = root.Set(key, 999, hash, 0)
		}
	}

	// Verify updates
	for key, originalValue := range testData {
		hash := Hash(key)
		value, found := root.Get(hash, 0)
		if !found {
			t.Errorf("Expected to find key '%s' after update", key)
		}
		if key <= "e" {
			if value != 999 {
				t.Errorf("Expected updated value 999 for key '%s', got %d", key, value)
			}
		} else {
			if value != originalValue {
				t.Errorf("Expected original value %d for key '%s', got %d", originalValue, key, value)
			}
		}
	}

	// Remove half the values
	keysToRemove := []string{"a", "c", "e", "g", "i"}
	for _, key := range keysToRemove {
		hash := Hash(key)
		var removed bool
		root, removed = root.Remove(hash, 0)
		if !removed {
			t.Errorf("Expected removal of key '%s' to succeed", key)
		}
	}

	// Verify removals
	for _, key := range keysToRemove {
		hash := Hash(key)
		_, found := root.Get(hash, 0)
		if found {
			t.Errorf("Expected key '%s' to be removed", key)
		}
	}

	// Verify remaining values
	remainingKeys := []string{"b", "d", "f", "h", "j"}
	for _, key := range remainingKeys {
		hash := Hash(key)
		value, found := root.Get(hash, 0)
		if !found {
			t.Errorf("Expected key '%s' to still exist", key)
		}
		expectedValue := testData[key]
		if key <= "e" {
			expectedValue = 999
		}
		if value != expectedValue {
			t.Errorf("Expected value %d for key '%s', got %d", expectedValue, key, value)
		}
	}

	// Test ToSlice
	slice := root.ToSlice()
	if len(slice) != len(remainingKeys) {
		t.Errorf("Expected %d entries in slice, got %d", len(remainingKeys), len(slice))
	}
}
