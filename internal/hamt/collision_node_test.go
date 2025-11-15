package hamt

import (
	"testing"
)

func TestNewCollisionNode(t *testing.T) {
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
	}

	node := NewCollisionNode[string, int](12345, entries)
	if node == nil {
		t.Fatal("Expected non-nil node")
	}

	if node.hash != 12345 {
		t.Errorf("Expected hash 12345, got %d", node.hash)
	}

	if len(node.entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(node.entries))
	}
}

func TestCollisionNodeKey(t *testing.T) {
	entries := []Entry[string, int]{
		{Key: "apple", Value: 100},
		{Key: "banana", Value: 200},
	}

	node := NewCollisionNode[string, int](12345, entries)

	key := node.Key()
	if key != "apple" {
		t.Errorf("Expected key 'apple', got '%s'", key)
	}
}

func TestCollisionNodeKeyEmpty(t *testing.T) {
	node := NewCollisionNode[string, int](12345, []Entry[string, int]{})

	key := node.Key()
	if key != "" {
		t.Errorf("Expected zero value for key, got '%s'", key)
	}
}

func TestCollisionNodeValue(t *testing.T) {
	entries := []Entry[string, int]{
		{Key: "apple", Value: 100},
		{Key: "banana", Value: 200},
	}

	node := NewCollisionNode[string, int](12345, entries)

	value := node.Value()
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}
}

func TestCollisionNodeValueEmpty(t *testing.T) {
	node := NewCollisionNode[string, int](12345, []Entry[string, int]{})

	value := node.Value()
	if value != 0 {
		t.Errorf("Expected zero value, got %d", value)
	}
}

func TestCollisionNodeGet(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
	}

	node := NewCollisionNode(hash, entries)

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

func TestCollisionNodeGetEmpty(t *testing.T) {
	node := NewCollisionNode[string, int](12345, []Entry[string, int]{})

	_, found := node.Get(12345, 0)
	if found {
		t.Error("Expected not to find value in empty collision node")
	}
}

func TestCollisionNodeSetSameHash(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
	}

	node := NewCollisionNode(hash, entries)

	// Add another entry with same hash
	newNode := node.Set("key2", 200, hash, 0)

	collision, ok := newNode.(*CollisionNode[string, int])
	if !ok {
		t.Fatal("Expected result to be CollisionNode")
	}

	if len(collision.entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(collision.entries))
	}

	// Verify first entry is preserved
	if collision.entries[0].Key != "key1" || collision.entries[0].Value != 100 {
		t.Error("Expected first entry to be preserved")
	}

	// Verify new entry is added
	if collision.entries[1].Key != "key2" || collision.entries[1].Value != 200 {
		t.Error("Expected new entry to be added")
	}
}

func TestCollisionNodeSetDifferentHash(t *testing.T) {
	hash1 := uint64(12345)
	hash2 := uint64(67890)

	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
	}

	node := NewCollisionNode(hash1, entries)

	// Set with different hash should create a new leaf
	newNode := node.Set("key2", 200, hash2, 0)

	leaf, ok := newNode.(*LeafNode[string, int])
	if !ok {
		t.Fatal("Expected result to be LeafNode")
	}

	if leaf.hash != hash2 {
		t.Errorf("Expected hash %d, got %d", hash2, leaf.hash)
	}

	if leaf.Key() != "key2" {
		t.Errorf("Expected key 'key2', got '%s'", leaf.Key())
	}
}

func TestCollisionNodeRemoveSingleEntry(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
	}

	node := NewCollisionNode(hash, entries)

	// Remove the only entry
	newNode, removed := node.Remove(hash, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	if newNode != nil {
		t.Error("Expected nil node after removing last entry")
	}
}

func TestCollisionNodeRemoveTwoEntries(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
	}

	node := NewCollisionNode(hash, entries)

	// Remove one entry, should get a leaf back
	newNode, removed := node.Remove(hash, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	leaf, ok := newNode.(*LeafNode[string, int])
	if !ok {
		t.Fatal("Expected result to be LeafNode")
	}

	if leaf.Key() != "key2" || leaf.Value() != 200 {
		t.Error("Expected remaining entry to be key2")
	}
}

func TestCollisionNodeRemoveMultipleEntries(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
		{Key: "key3", Value: 300},
	}

	node := NewCollisionNode(hash, entries)

	// Remove one entry
	newNode, removed := node.Remove(hash, 0)
	if !removed {
		t.Error("Expected removal to succeed")
	}

	collision, ok := newNode.(*CollisionNode[string, int])
	if !ok {
		t.Fatal("Expected result to be CollisionNode")
	}

	if len(collision.entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(collision.entries))
	}

	// Should have removed the first entry
	if collision.entries[0].Key != "key2" || collision.entries[0].Value != 200 {
		t.Error("Expected first remaining entry to be key2")
	}

	if collision.entries[1].Key != "key3" || collision.entries[1].Value != 300 {
		t.Error("Expected second remaining entry to be key3")
	}
}

func TestCollisionNodeRemoveNonMatching(t *testing.T) {
	hash := uint64(12345)
	entries := []Entry[string, int]{
		{Key: "key1", Value: 100},
	}

	node := NewCollisionNode(hash, entries)

	// Try to remove with different hash
	newNode, removed := node.Remove(99999, 0)
	if removed {
		t.Error("Expected removal to fail")
	}

	if newNode != node {
		t.Error("Expected node to be unchanged")
	}
}
