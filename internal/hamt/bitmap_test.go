package hamt

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	bitmap := Initialize()
	if bitmap != 0 {
		t.Errorf("Expected initialized bitmap to be 0, got %d", bitmap)
	}
}

func TestNewBitmap(t *testing.T) {
	bitmap := NewBitmap(12345)
	if uint64(bitmap) != 12345 {
		t.Errorf("Expected bitmap value 12345, got %d", bitmap)
	}
}

func TestBitmapPosition(t *testing.T) {
	bitmap := Initialize()

	// Test position calculation at offset 0
	hash := uint64(0b111111) // 63 in binary
	position := bitmap.Position(hash, 0)
	expected := uint64(1 << 63)
	if position != expected {
		t.Errorf("Expected position %d, got %d", expected, position)
	}

	// Test position calculation at offset 1
	hash = uint64(0b111111 << 6) // shifted to second band
	position = bitmap.Position(hash, 1)
	expected = uint64(1 << 63)
	if position != expected {
		t.Errorf("Expected position %d for offset 1, got %d", expected, position)
	}
}

func TestBitmapNext(t *testing.T) {
	bitmap := Initialize()

	bitmap = bitmap.Next(1 << 5)
	if uint64(bitmap) != (1 << 5) {
		t.Errorf("Expected bitmap %d, got %d", 1<<5, bitmap)
	}

	bitmap = bitmap.Next(1 << 10)
	expected := (1 << 5) | (1 << 10)
	if uint64(bitmap) != uint64(expected) {
		t.Errorf("Expected bitmap %d, got %d", expected, bitmap)
	}
}

func TestBitmapWithout(t *testing.T) {
	bitmap := Initialize()
	bitmap = bitmap.Next(1 << 5)
	bitmap = bitmap.Next(1 << 10)
	bitmap = bitmap.Next(1 << 15)

	bitmap = bitmap.Without(1 << 10)
	if bitmap.Has(1 << 10) {
		t.Error("Expected position to be removed")
	}

	if !bitmap.Has(1 << 5) {
		t.Error("Expected other positions to remain")
	}

	if !bitmap.Has(1 << 15) {
		t.Error("Expected other positions to remain")
	}
}

func TestBitmapHas(t *testing.T) {
	bitmap := Initialize()

	if bitmap.Has(1 << 5) {
		t.Error("Expected empty bitmap to not have any positions")
	}

	bitmap = bitmap.Next(1 << 5)
	if !bitmap.Has(1 << 5) {
		t.Error("Expected bitmap to have position")
	}

	if bitmap.Has(1 << 10) {
		t.Error("Expected bitmap to not have other positions")
	}
}

func TestBitmapIndex(t *testing.T) {
	bitmap := Initialize()
	bitmap = bitmap.Next(1 << 2)  // position 0
	bitmap = bitmap.Next(1 << 5)  // position 1
	bitmap = bitmap.Next(1 << 10) // position 2

	// Index of first position
	index, ok := bitmap.Index(1 << 2)
	if !ok {
		t.Error("Expected index calculation to succeed")
	}
	if index != 0 {
		t.Errorf("Expected index 0, got %d", index)
	}

	// Index of second position
	index, ok = bitmap.Index(1 << 5)
	if !ok {
		t.Error("Expected index calculation to succeed")
	}
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Index of third position
	index, ok = bitmap.Index(1 << 10)
	if !ok {
		t.Error("Expected index calculation to succeed")
	}
	if index != 2 {
		t.Errorf("Expected index 2, got %d", index)
	}

	// Index of position 0 should fail
	_, ok = bitmap.Index(0)
	if ok {
		t.Error("Expected index calculation to fail for position 0")
	}
}

func TestMsbMask(t *testing.T) {
	// Test with 0
	mask := msbMask(0)
	if mask != ^uint64(0) {
		t.Errorf("Expected mask ^uint64(0) for input 0, got %d", mask)
	}

	// Test with power of 2
	mask = msbMask(8) // 0b1000
	expected := uint64(0b0111)
	if mask != expected {
		t.Errorf("Expected mask %b, got %b", expected, mask)
	}

	// Test with non-power of 2
	mask = msbMask(10) // 0b1010
	expected = uint64(0b0111)
	if mask != expected {
		t.Errorf("Expected mask %b, got %b", expected, mask)
	}
}

func TestSplitMix64(t *testing.T) {
	// Test that it produces deterministic output
	input := uint64(12345)
	output1 := splitMix64(input)
	output2 := splitMix64(input)

	if output1 != output2 {
		t.Error("Expected deterministic output from splitMix64")
	}

	// Test that different inputs produce different outputs
	output3 := splitMix64(67890)
	if output1 == output3 {
		t.Error("Expected different outputs for different inputs")
	}
}

func TestBandMask(t *testing.T) {
	mask := bandMask()
	expected := uint64(63) // 2^6 - 1 = 63
	if mask != expected {
		t.Errorf("Expected mask %d, got %d", expected, mask)
	}
}
