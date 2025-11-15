package factory

import (
	"strings"
	"testing"
)

func TestStringFactoryOverrideConfig(t *testing.T) {
	factory := &StringFactory{}

	min := 10
	max := 12
	customChars := CharacterSet{'x', 'y', 'z'}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min":        min,
		"max":        max,
		"characters": customChars,
	}).Func(), 0)

	if properties.min != min {
		t.Fatalf("expected min %d, got %d", min, properties.min)
	}
	if properties.max != max {
		t.Fatalf("expected max %d, got %d", max, properties.max)
	}
	if len(properties.characters) != len(customChars) {
		t.Fatalf("expected custom charset")
	}
}

func TestStringFactoryInstantiate(t *testing.T) {
	factory := &StringFactory{}

	properties := StringProperties{value: "test string"}

	result := factory.Instantiate(properties)

	if result != "test string" {
		t.Errorf("Expected 'test string', got '%s'", result)
	}
}

func TestStringFactoryPrepare(t *testing.T) {
	minimum := 5
	maximum := 10
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min": minimum,
		"max": maximum,
	}).Func(), 0)

	if len(properties.value) < 5 || len(properties.value) > 10 {
		t.Errorf("Expected length between 5 and 10, got %d", len(properties.value))
	}
}

func TestStringFactoryPrepareDeterministic(t *testing.T) {
	factory := &StringFactory{}

	seed := int64(12345)
	properties1 := factory.Prepare(nil, seed)
	properties2 := factory.Prepare(nil, seed)

	if properties1.value != properties2.value {
		t.Error("Expected same value for same seed")
	}
}

func TestStringFactoryPrepareDifferentSeeds(t *testing.T) {
	factory := &StringFactory{}

	properties1 := factory.Prepare(nil, 1)
	properties2 := factory.Prepare(nil, 2)

	if properties1.value == properties2.value {
		t.Error("Expected different values for different seeds")
	}
}

func TestStringFactoryPrepareWithOverrides(t *testing.T) {
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"value": "overridden",
	}).Func(), 0)

	if properties.value != "overridden" {
		t.Errorf("Expected 'overridden', got '%s'", properties.value)
	}
}

func TestStringFactoryPrepareLength(t *testing.T) {
	minimum := 10
	maximum := 10
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min": minimum,
		"max": maximum,
	}).Func(), 0)

	if len(properties.value) != 10 {
		t.Errorf("Expected length 10, got %d", len(properties.value))
	}
}

func TestStringFactoryWithMaxLessThanMin(t *testing.T) {
	factory := &StringFactory{
		Min: 10,
		Max: 5,
	}

	properties := factory.Prepare(nil, 0)

	if properties.max < properties.min {
		t.Errorf("Expected max >= min, got min=%d, max=%d", properties.min, properties.max)
	}
}

func TestStringFactoryWithNegativeMin(t *testing.T) {
	factory := &StringFactory{
		Min: -5,
		Max: 10,
	}

	properties := factory.Prepare(nil, 0)

	if properties.min <= 0 {
		t.Errorf("Expected min > 0, got %d", properties.min)
	}
}

func TestStringFactoryWithZeroMax(t *testing.T) {
	factory := &StringFactory{
		Min: 0,
		Max: 0,
	}

	properties := factory.Prepare(nil, 0)

	if len(properties.value) < 1 {
		t.Errorf("Expected length >= 1, got %d", len(properties.value))
	}
}

func TestStringFactoryWithOverrideMaxLessThanMin(t *testing.T) {
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min": 10,
		"max": 5,
	}).Func(), 0)

	if properties.max < properties.min {
		t.Errorf("Expected max >= min, got min=%d, max=%d", properties.min, properties.max)
	}
}

func TestStringFactoryWithOverrideNegativeMin(t *testing.T) {
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min": -5,
	}).Func(), 0)

	if properties.min <= 0 {
		t.Errorf("Expected min > 0, got %d", properties.min)
	}
}

func TestStringFactoryWithOverrideEmptyCharacters(t *testing.T) {
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"characters": CharacterSet{},
	}).Func(), 0)

	if len(properties.characters) == 0 {
		t.Errorf("Expected characters to be set to default")
	}
}

func TestStringFactoryPrepareWithNumericCharacters(t *testing.T) {
	minimum := 5
	maximum := 5
	numericChars := Characters.Numeric
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min":        minimum,
		"max":        maximum,
		"characters": numericChars,
	}).Func(), 0)

	for _, char := range properties.value {
		if char < '0' || char > '9' {
			t.Errorf("Expected only numeric characters, got '%c'", char)
		}
	}
}

func TestStringFactoryPrepareWithAlphaCharacters(t *testing.T) {
	minimum := 5
	maximum := 5
	alphaChars := Characters.Alpha
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min":        minimum,
		"max":        maximum,
		"characters": alphaChars,
	}).Func(), 100)

	for _, char := range properties.value {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			t.Errorf("Expected only alpha characters, got '%c'", char)
		}
	}
}

func TestStringFactoryPrepareWithSymbolCharacters(t *testing.T) {
	minimum := 5
	maximum := 5
	symbolChars := Characters.Symbol
	factory := &StringFactory{}

	properties := factory.Prepare(Override[StringProperties](map[string]any{
		"min":        minimum,
		"max":        maximum,
		"characters": symbolChars,
	}).Func(), 200)

	if len(properties.value) != 5 {
		t.Errorf("Expected length 5, got %d", len(properties.value))
	}

	for _, char := range properties.value {
		found := false
		for _, symbol := range Characters.Symbol {
			if char == symbol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected only symbol characters, got '%c'", char)
		}
	}
}

func TestStringFactoryRetrieve(t *testing.T) {
	factory := &StringFactory{}

	instance := "retrieved string"
	properties := factory.Retrieve(instance)

	if properties.value != "retrieved string" {
		t.Errorf("Expected 'retrieved string', got '%s'", properties.value)
	}
}

func TestStringFactoryWithBuilder(t *testing.T) {
	minimum := 10
	maximum := 20
	factory := &StringFactory{}

	builder := Builder(factory)

	str1 := builder.Build(Override[StringProperties](map[string]any{
		"min": minimum,
		"max": maximum,
	}))
	if len(str1) < 10 || len(str1) > 20 {
		t.Errorf("Expected length between 10 and 20, got %d", len(str1))
	}

	str2 := builder.Build(Override[StringProperties](map[string]any{
		"min": minimum,
		"max": maximum,
	}))
	if str1 == str2 {
		t.Error("Expected different strings from consecutive builds")
	}
}

func TestStringFactoryBuildList(t *testing.T) {
	minimum := 5
	maximum := 15
	factory := &StringFactory{}

	builder := Builder(factory)

	strings := builder.BuildList(10, Override[StringProperties](map[string]any{
		"min": minimum,
		"max": maximum,
	}))

	if len(strings) != 10 {
		t.Errorf("Expected 10 strings, got %d", len(strings))
	}

	for _, str := range strings {
		if len(str) < 5 || len(str) > 15 {
			t.Errorf("Expected length between 5 and 15, got %d", len(str))
		}
	}
}

func TestStringFactoryBuildWith(t *testing.T) {
	factory := &StringFactory{}

	builder := Builder(factory)

	seed := int64(99999)
	str1 := builder.BuildWith(seed, nil)
	str2 := builder.BuildWith(seed, nil)

	if str1 != str2 {
		t.Error("Expected same string for same seed")
	}
}

func TestCharactersAlphanumeric(t *testing.T) {
	if len(Characters.Alphanumeric) != 62 {
		t.Errorf("Expected 62 alphanumeric characters, got %d", len(Characters.Alphanumeric))
	}

	alphanumericString := string(Characters.Alphanumeric)
	if !strings.Contains(alphanumericString, "a") ||
		!strings.Contains(alphanumericString, "Z") ||
		!strings.Contains(alphanumericString, "0") {
		t.Error("Expected Alphanumeric to contain letters and numbers")
	}
}

func TestCharactersAlpha(t *testing.T) {
	if len(Characters.Alpha) != 52 {
		t.Errorf("Expected 52 alpha characters, got %d", len(Characters.Alpha))
	}

	for _, char := range Characters.Alpha {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			t.Errorf("Expected only alpha characters, got '%c'", char)
		}
	}
}

func TestCharactersNumeric(t *testing.T) {
	if len(Characters.Numeric) != 10 {
		t.Errorf("Expected 10 numeric characters, got %d", len(Characters.Numeric))
	}

	for _, char := range Characters.Numeric {
		if char < '0' || char > '9' {
			t.Errorf("Expected only numeric characters, got '%c'", char)
		}
	}
}

func TestCharactersSymbol(t *testing.T) {
	if len(Characters.Symbol) != 32 {
		t.Errorf("Expected 32 symbol characters, got %d", len(Characters.Symbol))
	}

	symbolString := string(Characters.Symbol)
	if !strings.Contains(symbolString, "!") ||
		!strings.Contains(symbolString, "@") ||
		!strings.Contains(symbolString, "#") {
		t.Error("Expected Symbol to contain common symbols")
	}
}

func TestStringFactoryPrepareVariousLengths(t *testing.T) {
	minimum := 1
	maximum := 100
	factory := &StringFactory{}

	lengths := make(map[int]bool)

	for seed := int64(0); seed < 200; seed++ {
		properties := factory.Prepare(Override[StringProperties](map[string]any{
			"min": minimum,
			"max": maximum,
		}).Func(), seed)
		length := len(properties.value)

		if length < 1 || length > 100 {
			t.Errorf("Expected length between 1 and 100, got %d", length)
		}

		lengths[length] = true
	}

	if len(lengths) < 50 {
		t.Errorf("Expected at least 50 different lengths, got %d", len(lengths))
	}
}
