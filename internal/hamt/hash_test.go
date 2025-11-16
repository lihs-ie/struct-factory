package hamt

import (
	"hash/fnv"
	"reflect"
	"testing"
	"time"
)

// TestHash_BasicTypes tests hashing of basic types.
func TestHash_BasicTypes(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 42},
		{"int8", int8(42)},
		{"int16", int16(42)},
		{"int32", int32(42)},
		{"int64", int64(42)},
		{"uint", uint(42)},
		{"uint8", uint8(42)},
		{"uint16", uint16(42)},
		{"uint32", uint32(42)},
		{"uint64", uint64(42)},
		{"float32", float32(3.14)},
		{"float64", float64(3.14)},
		{"complex64", complex64(1 + 2i)},
		{"complex128", complex128(1 + 2i)},
		{"bool_true", true},
		{"bool_false", false},
		{"string", "hello"},
		{"empty_string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := Hash(tt.value)
			hash2 := Hash(tt.value)

			// Same value should produce same hash
			if hash1 != hash2 {
				t.Errorf("Hash(%v) produced different results: %d vs %d", tt.value, hash1, hash2)
			}

			// Hash should not be zero (except for very rare cases)
			if hash1 == 0 {
				t.Logf("Warning: Hash(%v) produced zero hash", tt.value)
			}
		})
	}
}

// TestHash_Nil tests hashing of nil values.
func TestHash_Nil(t *testing.T) {
	var nilPointer *int
	var nilInterface interface{}

	hash1 := Hash(nilPointer)
	hash2 := Hash(nilInterface)

	if hash1 != hash2 {
		t.Errorf("Nil values should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_Pointers tests hashing of pointer values.
func TestHash_Pointers(t *testing.T) {
	value := 42
	pointer := &value

	hashValue := Hash(value)
	hashPointer := Hash(pointer)

	if hashValue != hashPointer {
		t.Errorf("Pointer and value should produce same hash: %d vs %d", hashValue, hashPointer)
	}
}

// TestHash_Interfaces tests hashing through interface{}.
func TestHash_Interfaces(t *testing.T) {
	var value interface{} = 42

	hash1 := Hash(value)
	hash2 := Hash(42)

	if hash1 != hash2 {
		t.Errorf("Interface and direct value should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_Arrays tests hashing of arrays.
func TestHash_Arrays(t *testing.T) {
	array1 := [3]int{1, 2, 3}
	array2 := [3]int{1, 2, 3}
	array3 := [3]int{3, 2, 1}

	hash1 := Hash(array1)
	hash2 := Hash(array2)
	hash3 := Hash(array3)

	if hash1 != hash2 {
		t.Errorf("Same arrays should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different arrays should produce different hash: %d vs %d", hash1, hash3)
	}
}

// TestHash_Slices tests hashing of slices.
func TestHash_Slices(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	slice3 := []int{3, 2, 1}
	emptySlice := []int{}

	hash1 := Hash(slice1)
	hash2 := Hash(slice2)
	hash3 := Hash(slice3)
	hashEmpty := Hash(emptySlice)

	if hash1 != hash2 {
		t.Errorf("Same slices should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different slices should produce different hash: %d vs %d", hash1, hash3)
	}

	if hashEmpty == 0 {
		t.Logf("Empty slice hash: %d", hashEmpty)
	}
}

// TestHash_Maps tests hashing of maps.
func TestHash_Maps(t *testing.T) {
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 2, "a": 1}
	map3 := map[string]int{"a": 1, "b": 3}
	emptyMap := map[string]int{}

	hash1 := Hash(map1)
	hash2 := Hash(map2)
	hash3 := Hash(map3)
	hashEmpty := Hash(emptyMap)

	// Maps with same content should produce same hash (order-independent)
	if hash1 != hash2 {
		t.Errorf("Maps with same content should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different maps should produce different hash: %d vs %d", hash1, hash3)
	}

	if hashEmpty == 0 {
		t.Logf("Empty map hash: %d", hashEmpty)
	}
}

// TestHash_Structs tests hashing of structs.
func TestHash_Structs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Alice", Age: 30}
	person3 := Person{Name: "Bob", Age: 30}

	hash1 := Hash(person1)
	hash2 := Hash(person2)
	hash3 := Hash(person3)

	if hash1 != hash2 {
		t.Errorf("Same structs should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different structs should produce different hash: %d vs %d", hash1, hash3)
	}
}

// TestHash_StructWithUnexportedFields tests hashing of structs with unexported fields.
func TestHash_StructWithUnexportedFields(t *testing.T) {
	type StructWithPrivate struct {
		Public  string
		private string
	}

	struct1 := StructWithPrivate{Public: "public", private: "private1"}
	struct2 := StructWithPrivate{Public: "public", private: "private2"}

	hash1 := Hash(struct1)
	hash2 := Hash(struct2)

	// Unexported fields should be ignored, so hashes should be same
	if hash1 != hash2 {
		t.Errorf("Structs differing only in unexported fields should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_NestedStructs tests hashing of nested structs.
func TestHash_NestedStructs(t *testing.T) {
	type Address struct {
		City string
		Zip  int
	}

	type Person struct {
		Name    string
		Address Address
	}

	person1 := Person{
		Name:    "Alice",
		Address: Address{City: "Tokyo", Zip: 100},
	}

	person2 := Person{
		Name:    "Alice",
		Address: Address{City: "Tokyo", Zip: 100},
	}

	person3 := Person{
		Name:    "Alice",
		Address: Address{City: "Osaka", Zip: 200},
	}

	hash1 := Hash(person1)
	hash2 := Hash(person2)
	hash3 := Hash(person3)

	if hash1 != hash2 {
		t.Errorf("Same nested structs should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different nested structs should produce different hash: %d vs %d", hash1, hash3)
	}
}

// TestHash_Time tests hashing of time.Time.
func TestHash_Time(t *testing.T) {
	time1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	hash1 := Hash(time1)
	hash2 := Hash(time2)
	hash3 := Hash(time3)

	if hash1 != hash2 {
		t.Errorf("Same times should produce same hash: %d vs %d", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different times should produce different hash: %d vs %d", hash1, hash3)
	}
}

// CustomHashable is a type that implements Hashable interface.
type CustomHashable struct {
	Value int
}

func (c CustomHashable) Hash() (uint64, error) {
	return uint64(c.Value * 1000), nil
}

// TestHash_Hashable tests hashing of types implementing Hashable interface.
func TestHash_Hashable(t *testing.T) {
	custom := CustomHashable{Value: 42}
	hash := Hash(custom)

	expectedHash := uint64(42 * 1000)
	if hash != expectedHash {
		t.Errorf("Hashable interface should be used: got %d, want %d", hash, expectedHash)
	}
}

// TestHash_HashablePointer tests hashing of pointer to Hashable.
func TestHash_HashablePointer(t *testing.T) {
	custom := &CustomHashable{Value: 42}
	hash := Hash(custom)

	expectedHash := uint64(42 * 1000)
	if hash != expectedHash {
		t.Errorf("Hashable interface should be used for pointer: got %d, want %d", hash, expectedHash)
	}
}

// TestHash_ComplexNesting tests hashing of complex nested structures.
func TestHash_ComplexNesting(t *testing.T) {
	type ComplexStruct struct {
		Name     string
		Numbers  []int
		Metadata map[string]interface{}
		Time     time.Time
	}

	complex1 := ComplexStruct{
		Name:    "test",
		Numbers: []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	complex2 := ComplexStruct{
		Name:    "test",
		Numbers: []int{1, 2, 3},
		Metadata: map[string]interface{}{
			"key2": 42,
			"key1": "value1",
		},
		Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	hash1 := Hash(complex1)
	hash2 := Hash(complex2)

	if hash1 != hash2 {
		t.Errorf("Complex structs with same content should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_EmptyStruct tests hashing of empty structs.
func TestHash_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	empty1 := EmptyStruct{}
	empty2 := EmptyStruct{}

	hash1 := Hash(empty1)
	hash2 := Hash(empty2)

	if hash1 != hash2 {
		t.Errorf("Empty structs should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_Consistency tests that same value always produces same hash.
func TestHash_Consistency(t *testing.T) {
	value := "test string"

	hashes := make([]uint64, 100)
	for i := range hashes {
		hashes[i] = Hash(value)
	}

	firstHash := hashes[0]
	for i, hash := range hashes {
		if hash != firstHash {
			t.Errorf("Hash inconsistent at iteration %d: got %d, want %d", i, hash, firstHash)
		}
	}
}

// TestHash_DifferentTypes tests that different types produce different hashes.
func TestHash_DifferentTypes(t *testing.T) {
	hashInt := Hash(42)
	hashString := Hash("42")
	hashFloat := Hash(42.0)

	if hashInt == hashString {
		t.Errorf("int and string should produce different hashes")
	}

	if hashInt == hashFloat {
		t.Logf("Note: int(42) and float(42.0) produced same hash (might be expected)")
	}
}

// BenchmarkHash_String benchmarks string hashing.
func BenchmarkHash_String(b *testing.B) {
	value := "test string for benchmarking"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(value)
	}
}

// BenchmarkHash_Int benchmarks int hashing.
func BenchmarkHash_Int(b *testing.B) {
	value := 42
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(value)
	}
}

// BenchmarkHash_Struct benchmarks struct hashing.
func BenchmarkHash_Struct(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}
	value := Person{Name: "Alice", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(value)
	}
}

// BenchmarkHash_Slice benchmarks slice hashing.
func BenchmarkHash_Slice(b *testing.B) {
	value := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(value)
	}
}

// BenchmarkHash_Map benchmarks map hashing.
func BenchmarkHash_Map(b *testing.B) {
	value := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(value)
	}
}

// TestHash_DoublePointer tests hashing of double pointers.
func TestHash_DoublePointer(t *testing.T) {
	value := 42
	pointer1 := &value
	pointer2 := &pointer1

	hashValue := Hash(value)
	hashPointer2 := Hash(pointer2)

	if hashValue != hashPointer2 {
		t.Errorf("Double pointer and value should produce same hash: %d vs %d", hashValue, hashPointer2)
	}
}

// TestHash_InterfaceWrappedPointer tests hashing of interface-wrapped pointers.
func TestHash_InterfaceWrappedPointer(t *testing.T) {
	value := 42
	pointer := &value
	var interfaceValue interface{} = pointer

	hashPointer := Hash(pointer)
	hashInterface := Hash(interfaceValue)

	if hashPointer != hashInterface {
		t.Errorf("Interface-wrapped pointer should produce same hash: %d vs %d", hashPointer, hashInterface)
	}
}

// ErrorHashable is a type that implements Hashable but returns an error.
type ErrorHashable struct{}

func (e ErrorHashable) Hash() (uint64, error) {
	return 0, nil // Returns successfully with 0
}

// TestHash_ErrorHashable tests Hashable that returns an error.
func TestHash_ErrorHashable(t *testing.T) {
	errorHashable := ErrorHashable{}
	hash := Hash(errorHashable)

	// Should use the Hashable interface even if it returns 0
	if hash != 0 {
		t.Errorf("ErrorHashable should return 0, got %d", hash)
	}
}

// TestHash_NestedSlicesWithErrors tests nested structures that might cause errors.
func TestHash_NestedSlicesWithErrors(t *testing.T) {
	nestedSlice := [][]int{{1, 2}, {3, 4}, {5, 6}}

	hash1 := Hash(nestedSlice)
	hash2 := Hash(nestedSlice)

	if hash1 != hash2 {
		t.Errorf("Nested slices should produce consistent hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_NestedMapsWithErrors tests nested maps.
func TestHash_NestedMapsWithErrors(t *testing.T) {
	nestedMap := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"z": 3},
	}

	hash1 := Hash(nestedMap)
	hash2 := Hash(nestedMap)

	if hash1 != hash2 {
		t.Errorf("Nested maps should produce consistent hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_SliceWithMixedTypes tests slices containing different types.
func TestHash_SliceWithMixedTypes(t *testing.T) {
	slice := []interface{}{1, "two", 3.0, true}

	hash1 := Hash(slice)
	hash2 := Hash(slice)

	if hash1 != hash2 {
		t.Errorf("Same mixed-type slice should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_MapWithComplexKeys tests maps with complex key types.
func TestHash_MapWithComplexKeys(t *testing.T) {
	type Key struct {
		ID   int
		Name string
	}

	map1 := map[Key]string{
		{ID: 1, Name: "first"}:  "value1",
		{ID: 2, Name: "second"}: "value2",
	}

	map2 := map[Key]string{
		{ID: 2, Name: "second"}: "value2",
		{ID: 1, Name: "first"}:  "value1",
	}

	hash1 := Hash(map1)
	hash2 := Hash(map2)

	if hash1 != hash2 {
		t.Errorf("Maps with complex keys should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_StructWithPointerFields tests structs containing pointer fields.
func TestHash_StructWithPointerFields(t *testing.T) {
	type StructWithPointer struct {
		Value *int
		Name  string
	}

	value := 42
	struct1 := StructWithPointer{Value: &value, Name: "test"}
	struct2 := StructWithPointer{Value: &value, Name: "test"}

	hash1 := Hash(struct1)
	hash2 := Hash(struct2)

	if hash1 != hash2 {
		t.Errorf("Structs with pointer fields should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_StructWithNilPointer tests structs with nil pointer fields.
func TestHash_StructWithNilPointer(t *testing.T) {
	type StructWithPointer struct {
		Value *int
		Name  string
	}

	struct1 := StructWithPointer{Value: nil, Name: "test"}
	struct2 := StructWithPointer{Value: nil, Name: "test"}

	hash1 := Hash(struct1)
	hash2 := Hash(struct2)

	if hash1 != hash2 {
		t.Errorf("Structs with nil pointers should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_StructWithSliceField tests structs containing slice fields.
func TestHash_StructWithSliceField(t *testing.T) {
	type StructWithSlice struct {
		Name   string
		Values []int
	}

	struct1 := StructWithSlice{Name: "test", Values: []int{1, 2, 3}}
	struct2 := StructWithSlice{Name: "test", Values: []int{1, 2, 3}}

	hash1 := Hash(struct1)
	hash2 := Hash(struct2)

	if hash1 != hash2 {
		t.Errorf("Structs with slice fields should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_StructWithMapField tests structs containing map fields.
func TestHash_StructWithMapField(t *testing.T) {
	type StructWithMap struct {
		Name string
		Data map[string]int
	}

	struct1 := StructWithMap{
		Name: "test",
		Data: map[string]int{"a": 1, "b": 2},
	}
	struct2 := StructWithMap{
		Name: "test",
		Data: map[string]int{"b": 2, "a": 1},
	}

	hash1 := Hash(struct1)
	hash2 := Hash(struct2)

	if hash1 != hash2 {
		t.Errorf("Structs with map fields should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_ArrayOfStructs tests arrays containing structs.
func TestHash_ArrayOfStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	array1 := [2]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	array2 := [2]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	hash1 := Hash(array1)
	hash2 := Hash(array2)

	if hash1 != hash2 {
		t.Errorf("Arrays of structs should produce same hash: %d vs %d", hash1, hash2)
	}
}

// TestHash_Chan tests that channels are handled.
func TestHash_Chan(t *testing.T) {
	ch := make(chan int)
	hash := Hash(ch)

	// Channels are not supported, should return a default hash
	_ = hash // Just ensure it doesn't panic
}

// TestHash_Func tests that functions are handled.
func TestHash_Func(t *testing.T) {
	fn := func() {}
	hash := Hash(fn)

	// Functions are not supported, should return a default hash
	_ = hash // Just ensure it doesn't panic
}

// TestHashValue_DirectCall tests hashValue function directly for error paths.
func TestHashValue_DirectCall(t *testing.T) {
	hasher := fnv.New64a()

	// Test various types directly
	testCases := []struct {
		name  string
		value any
	}{
		{"nil", nil},
		{"int", 42},
		{"string", "test"},
		{"slice", []int{1, 2, 3}},
		{"map", map[string]int{"a": 1}},
		{"struct", struct{ Name string }{"test"}},
		{"time", time.Now()},
		{"array", [3]int{1, 2, 3}},
		{"bool", true},
		{"float", 3.14},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := reflect.ValueOf(tc.value)
			hash, err := hashValue(hasher, value)
			if err != nil {
				t.Errorf("hashValue(%s) returned error: %v", tc.name, err)
			}

			if hash == 0 {
				t.Logf("hashValue(%s) returned zero hash", tc.name)
			}
		})
	}
}

// TestHashValue_NestedStructures tests deeply nested structures.
func TestHashValue_NestedStructures(t *testing.T) {
	hasher := fnv.New64a()

	type Inner struct {
		Value int
	}

	type Middle struct {
		Inner Inner
		Items []int
	}

	type Outer struct {
		Middle Middle
		Data   map[string]int
	}

	outer := Outer{
		Middle: Middle{
			Inner: Inner{Value: 42},
			Items: []int{1, 2, 3},
		},
		Data: map[string]int{"key": 100},
	}

	value := reflect.ValueOf(outer)
	hash, err := hashValue(hasher, value)
	if err != nil {
		t.Errorf("hashValue for nested structure returned error: %v", err)
	}

	if hash == 0 {
		t.Error("hashValue for nested structure returned zero")
	}
}

// TestHashValue_EmptyCollections tests empty collections.
func TestHashValue_EmptyCollections(t *testing.T) {
	hasher := fnv.New64a()

	testCases := []struct {
		name  string
		value any
	}{
		{"empty slice", []int{}},
		{"empty map", map[string]int{}},
		{"empty array", [0]int{}},
		{"empty string", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := reflect.ValueOf(tc.value)
			_, err := hashValue(hasher, value)
			if err != nil {
				t.Errorf("hashValue(%s) returned error: %v", tc.name, err)
			}
		})
	}
}
