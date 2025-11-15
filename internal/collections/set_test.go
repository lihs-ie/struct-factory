package collections

import (
	"testing"
)

func TestNewSet(t *testing.T) {
	set := NewSet[int](nil)
	if set == nil {
		t.Error("Expected non-nil set")
	}
	if !set.IsEmpty() {
		t.Error("Expected empty set")
	}
}

func TestNewFromSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	set := NewFromSlice(items)

	if set.IsEmpty() {
		t.Error("Expected non-empty set")
	}

	for _, item := range items {
		if !set.Has(item) {
			t.Errorf("Expected set to contain %d", item)
		}
	}

	if set.Size() != len(items) {
		t.Errorf("Expected size %d, got %d", len(items), set.Size())
	}
}

func TestNewFromSliceWithDuplicates(t *testing.T) {
	items := []int{1, 2, 2, 3, 3, 3}
	set := NewFromSlice(items)

	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	if !set.Has(1) || !set.Has(2) || !set.Has(3) {
		t.Error("Expected set to contain 1, 2, 3")
	}
}

func TestNewFromEmptySlice(t *testing.T) {
	set := NewFromSlice([]int{})
	if !set.IsEmpty() {
		t.Error("Expected empty set")
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0, got %d", set.Size())
	}
}

func TestSetAdd(t *testing.T) {
	set := NewSet[string](nil)

	set.Set("apple")
	if !set.Has("apple") {
		t.Error("Expected set to contain 'apple'")
	}

	set.Set("banana")
	if !set.Has("banana") {
		t.Error("Expected set to contain 'banana'")
	}

	if set.Size() != 2 {
		t.Errorf("Expected size 2, got %d", set.Size())
	}
}

func TestSetAddDuplicate(t *testing.T) {
	set := NewSet[int](nil)

	set.Set(42)
	set.Set(42)

	if set.Size() != 1 {
		t.Errorf("Expected size 1, got %d", set.Size())
	}
}

func TestRemove(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})

	set.Remove(3)
	if set.Has(3) {
		t.Error("Expected 3 to be removed")
	}

	if set.Size() != 4 {
		t.Errorf("Expected size 4, got %d", set.Size())
	}

	if !set.Has(1) || !set.Has(2) || !set.Has(4) || !set.Has(5) {
		t.Error("Expected remaining items to still be in set")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3})

	set.Remove(999)
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}
}

func TestRemoveFromEmptySet(t *testing.T) {
	set := NewSet[int](nil)
	set.Remove(1)

	if !set.IsEmpty() {
		t.Error("Expected set to remain empty")
	}
}

func TestHas(t *testing.T) {
	set := NewFromSlice([]string{"apple", "banana", "cherry"})

	if !set.Has("apple") {
		t.Error("Expected set to contain 'apple'")
	}

	if !set.Has("banana") {
		t.Error("Expected set to contain 'banana'")
	}

	if set.Has("orange") {
		t.Error("Expected set to not contain 'orange'")
	}
}

func TestHasEmptySet(t *testing.T) {
	set := NewSet[int](nil)

	if set.Has(1) {
		t.Error("Expected empty set to not contain any items")
	}
}

func TestIsEmpty(t *testing.T) {
	set := NewSet[int](nil)
	if !set.IsEmpty() {
		t.Error("Expected new set to be empty")
	}

	set.Set(1)
	if set.IsEmpty() {
		t.Error("Expected set with items to not be empty")
	}

	set.Remove(1)
	if !set.IsEmpty() {
		t.Error("Expected set to be empty after removing all items")
	}
}

func TestSize(t *testing.T) {
	set := NewSet[int](nil)
	if set.Size() != 0 {
		t.Errorf("Expected size 0, got %d", set.Size())
	}

	for i := 1; i <= 10; i++ {
		set.Set(i)
		if set.Size() != i {
			t.Errorf("Expected size %d, got %d", i, set.Size())
		}
	}

	for i := 10; i >= 1; i-- {
		set.Remove(i)
		if set.Size() != i-1 {
			t.Errorf("Expected size %d, got %d", i-1, set.Size())
		}
	}
}

func TestToSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	set := NewFromSlice(items)

	slice := set.ToSlice()
	if len(slice) != len(items) {
		t.Errorf("Expected slice length %d, got %d", len(items), len(slice))
	}

	resultSet := NewFromSlice(slice)
	for _, item := range items {
		if !resultSet.Has(item) {
			t.Errorf("Expected slice to contain %d", item)
		}
	}
}

func TestToSliceEmpty(t *testing.T) {
	set := NewSet[int](nil)
	slice := set.ToSlice()

	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}
}

func TestSetWithStrings(t *testing.T) {
	set := NewSet[string](nil)

	set.Set("hello")
	set.Set("world")
	set.Set("hello")

	if set.Size() != 2 {
		t.Errorf("Expected size 2, got %d", set.Size())
	}

	if !set.Has("hello") || !set.Has("world") {
		t.Error("Expected set to contain 'hello' and 'world'")
	}
}

func TestSetWithStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	set := NewSet[Person](nil)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	set.Set(p1)
	set.Set(p2)

	if set.Size() != 2 {
		t.Errorf("Expected size 2, got %d", set.Size())
	}

	if !set.Has(p1) {
		t.Error("Expected set to contain p1")
	}

	if !set.Has(p2) {
		t.Error("Expected set to contain p2")
	}
}
