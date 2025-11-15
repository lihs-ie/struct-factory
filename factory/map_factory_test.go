package factory

import (
	"testing"
)

type IntProperties struct {
	Value int
}

type IntFactory struct{}

func (f *IntFactory) Instantiate(properties IntProperties) int {
	return properties.Value
}

func (f *IntFactory) Prepare(overrides Partial[IntProperties], seed int64) IntProperties {
	properties := IntProperties{
		Value: int(seed),
	}

	if overrides != nil {
		overrides(&properties)
	}

	return properties
}

func (f *IntFactory) Retrieve(instance int) IntProperties {
	return IntProperties{
		Value: instance,
	}
}

func TestNewMapFactory(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &StringFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	if mapFactory == nil {
		t.Fatal("Expected non-nil MapFactory")
	}

	if mapFactory.keyFactory != keyFactory {
		t.Error("Expected keyFactory to be set")
	}

	if mapFactory.valueFactory != valueFactory {
		t.Error("Expected valueFactory to be set")
	}
}

func TestMapFactoryInstantiate(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	properties := MapProperties[int, int]{
		entries: []MapEntry[int, int]{
			{Key: 1, Value: 100},
			{Key: 2, Value: 200},
			{Key: 3, Value: 300},
		},
	}

	result := mapFactory.Instantiate(properties)

	if len(result) != 3 {
		t.Errorf("Expected map with 3 entries, got %d", len(result))
	}

	if result[1] != 100 {
		t.Errorf("Expected result[1] = 100, got %d", result[1])
	}

	if result[2] != 200 {
		t.Errorf("Expected result[2] = 200, got %d", result[2])
	}

	if result[3] != 300 {
		t.Errorf("Expected result[3] = 300, got %d", result[3])
	}
}

func TestMapFactoryInstantiateEmpty(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	properties := MapProperties[int, int]{entries: []MapEntry[int, int]{}}

	result := mapFactory.Instantiate(properties)

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(result))
	}
}

func TestMapFactoryPrepare(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	properties := mapFactory.Prepare(nil, 0)

	if len(properties.entries) < 1 || len(properties.entries) > 10 {
		t.Errorf("Expected between 1 and 10 entries, got %d", len(properties.entries))
	}
}

func TestMapFactoryPrepareDeterministic(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	seed := int64(5)
	properties1 := mapFactory.Prepare(nil, seed)
	properties2 := mapFactory.Prepare(nil, seed)

	if len(properties1.entries) != len(properties2.entries) {
		t.Error("Expected same number of entries for same seed")
	}

	for i := range properties1.entries {
		if properties1.entries[i].Key != properties2.entries[i].Key {
			t.Errorf("Expected same key at index %d", i)
		}
		if properties1.entries[i].Value != properties2.entries[i].Value {
			t.Errorf("Expected same value at index %d", i)
		}
	}
}

func TestMapFactoryPrepareVariousSizes(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	sizes := make(map[int]bool)

	for seed := int64(0); seed < 100; seed++ {
		properties := mapFactory.Prepare(nil, seed)
		size := len(properties.entries)

		if size < 1 || size > 10 {
			t.Errorf("Expected size between 1 and 10, got %d for seed %d", size, seed)
		}

		sizes[size] = true
	}

	if len(sizes) < 5 {
		t.Errorf("Expected at least 5 different sizes, got %d", len(sizes))
	}
}

func TestMapFactoryPrepareWithOverrides(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	customEntries := []MapEntry[int, int]{
		{Key: 999, Value: 1000},
	}

	properties := mapFactory.Prepare(Override[MapProperties[int, int]](map[string]any{
		"entries": customEntries,
	}).Func(), 5)

	if len(properties.entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(properties.entries))
	}

	if properties.entries[0].Key != 999 {
		t.Errorf("Expected key 999, got %d", properties.entries[0].Key)
	}

	if properties.entries[0].Value != 1000 {
		t.Errorf("Expected value 1000, got %d", properties.entries[0].Value)
	}
}

func TestMapFactoryRetrieve(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	instance := map[int]int{
		1: 100,
		2: 200,
		3: 300,
	}

	properties := mapFactory.Retrieve(instance)

	if len(properties.entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(properties.entries))
	}

	found := make(map[int]int)
	for _, entry := range properties.entries {
		found[entry.Key] = entry.Value
	}

	if found[1] != 100 {
		t.Errorf("Expected value 100 for key 1, got %d", found[1])
	}

	if found[2] != 200 {
		t.Errorf("Expected value 200 for key 2, got %d", found[2])
	}

	if found[3] != 300 {
		t.Errorf("Expected value 300 for key 3, got %d", found[3])
	}
}

func TestMapFactoryRetrieveEmpty(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	instance := map[int]int{}

	properties := mapFactory.Retrieve(instance)

	if len(properties.entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(properties.entries))
	}
}

func TestMapFactoryWithBuilder(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	builder := Builder(mapFactory)

	result := builder.Build(nil)

	if len(result) < 1 || len(result) > 10 {
		t.Errorf("Expected between 1 and 10 entries, got %d", len(result))
	}
}

func TestMapFactoryBuildWith(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	builder := Builder(mapFactory)

	seed := int64(7)
	result1 := builder.BuildWith(seed, nil)
	result2 := builder.BuildWith(seed, nil)

	if len(result1) != len(result2) {
		t.Error("Expected same size for same seed")
	}

	for key, value := range result1 {
		if result2[key] != value {
			t.Errorf("Expected same value for key %d", key)
		}
	}
}

func TestMapFactoryBuildList(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	builder := Builder(mapFactory)

	results := builder.BuildList(5, nil)

	if len(results) != 5 {
		t.Errorf("Expected 5 maps, got %d", len(results))
	}

	for _, result := range results {
		if len(result) < 1 || len(result) > 10 {
			t.Errorf("Expected between 1 and 10 entries, got %d", len(result))
		}
	}
}

func TestMapFactoryWithStringKeys(t *testing.T) {
	minimum := 5
	maximum := 10
	keyFactory := &StringFactory{Min: minimum, Max: maximum}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	builder := Builder(mapFactory)

	result := builder.Build(nil)

	if len(result) < 1 || len(result) > 10 {
		t.Errorf("Expected between 1 and 10 entries, got %d", len(result))
	}

	for key := range result {
		if len(key) < 5 || len(key) > 10 {
			t.Errorf("Expected key length between 5 and 10, got %d", len(key))
		}
	}
}

func TestMapFactoryWithStringValues(t *testing.T) {
	minimum := 3
	maximum := 8
	keyFactory := &IntFactory{}
	valueFactory := &StringFactory{Min: minimum, Max: maximum}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	builder := Builder(mapFactory)

	result := builder.Build(nil)

	if len(result) < 1 || len(result) > 10 {
		t.Errorf("Expected between 1 and 10 entries, got %d", len(result))
	}

	for _, value := range result {
		if len(value) < 3 || len(value) > 8 {
			t.Errorf("Expected value length between 3 and 8, got %d", len(value))
		}
	}
}

func TestMapFactoryPrepareEntryCount(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)

	testCases := []struct {
		seed          int64
		expectedCount int
	}{
		{0, 1},
		{1, 2},
		{5, 6},
		{9, 10},
		{10, 1},
		{19, 10},
	}

	for _, tc := range testCases {
		properties := mapFactory.Prepare(nil, tc.seed)
		if len(properties.entries) != tc.expectedCount {
			t.Errorf("For seed %d, expected %d entries, got %d",
				tc.seed, tc.expectedCount, len(properties.entries))
		}
	}
}

func TestMapFactoryWithOverrideLiteral(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)
	builder := Builder(mapFactory)

	customEntries := []MapEntry[int, int]{
		{Key: 1, Value: 100},
		{Key: 2, Value: 200},
	}

	result := builder.Build(Override[MapProperties[int, int]](map[string]any{
		"Entries": customEntries,
	}))

	if len(result) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(result))
	}

	if result[1] != 100 {
		t.Errorf("Expected result[1] = 100, got %d", result[1])
	}

	if result[2] != 200 {
		t.Errorf("Expected result[2] = 200, got %d", result[2])
	}
}

func TestMapFactoryWithLiteralAndFuncOverrides(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)
	builder := Builder(mapFactory)

	result := builder.Build(Override[MapProperties[int, int]](map[string]any{
		"Entries": []MapEntry[int, int]{
			{Key: 5, Value: 500},
			{Key: 6, Value: 600},
		},
	}))

	if len(result) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(result))
	}

	if result[5] != 500 {
		t.Errorf("Expected result[5] = 500, got %d", result[5])
	}

	if result[6] != 600 {
		t.Errorf("Expected result[6] = 600, got %d", result[6])
	}
}

func TestMapFactoryWithInlineOverride(t *testing.T) {
	keyFactory := &IntFactory{}
	valueFactory := &IntFactory{}

	mapFactory := NewMapFactory(keyFactory, valueFactory)
	builder := Builder(mapFactory)

	result := builder.Build(Override[MapProperties[int, int]](map[string]any{
		"Entries": []MapEntry[int, int]{
			{Key: 10, Value: 1000},
			{Key: 20, Value: 2000},
			{Key: 30, Value: 3000},
		},
	}))

	if len(result) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(result))
	}

	expectedValues := map[int]int{
		10: 1000,
		20: 2000,
		30: 3000,
	}

	for key, expectedValue := range expectedValues {
		if result[key] != expectedValue {
			t.Errorf("Expected result[%d] = %d, got %d", key, expectedValue, result[key])
		}
	}
}
