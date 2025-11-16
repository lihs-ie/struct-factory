package factory

import (
	"fmt"
	"testing"
)

type stubProps struct {
	Value string
	Seed  int64
}

func (p *stubProps) SetValue(v string) {
	p.Value = v
}

type stubInstance struct {
	Value string
	Seed  int64
}

type stubFactory struct {
	prepareSeeds []int64
}

func (f *stubFactory) Instantiate(props stubProps) stubInstance {
	return stubInstance(props)
}

func (f *stubFactory) Prepare(overrides Partial[stubProps], seed int64) stubProps {
	props := stubProps{Value: fmt.Sprintf("seed-%d", seed), Seed: seed}
	if overrides != nil {
		overrides(&props)
	}
	f.prepareSeeds = append(f.prepareSeeds, seed)
	return props
}

func (f *stubFactory) Retrieve(instance stubInstance) stubProps {
	return stubProps(instance)
}

func TestBuilderBuildAppliesLiteralOverride(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	result := builder.Build(Override[stubProps](map[string]any{
		"Value": "literal",
	}))

	if result.Value != "literal" {
		t.Fatalf("expected literal override, got %s", result.Value)
	}

	if len(factory.prepareSeeds) != 1 {
		t.Fatalf("expected one prepare call, got %d", len(factory.prepareSeeds))
	}
}

type stubPropsOverride struct {
	Value string
}

func TestBuilderBuildWithStructOverride(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	result := builder.Build(Override[stubProps](stubPropsOverride{Value: "struct-value"}))

	if result.Value != "struct-value" {
		t.Fatalf("expected struct-value, got %s", result.Value)
	}
}

func TestBuilderBuildListGeneratesUniqueSeeds(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	results := builder.BuildList(5, nil)
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	seen := make(map[int64]struct{})
	for _, instance := range results {
		seen[instance.Seed] = struct{}{}
	}

	if len(seen) != len(results) {
		t.Fatalf("expected unique seeds, got %d unique for %d items", len(seen), len(results))
	}
}

func TestBuilderBuildWithSeed(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	result := builder.BuildWith(42, nil)
	if result.Seed != 42 {
		t.Fatalf("expected seed 42, got %d", result.Seed)
	}

	if result.Value != fmt.Sprintf("seed-%d", result.Seed) {
		t.Fatalf("unexpected value %s", result.Value)
	}
}

func TestBuilderBuildListWithSequentialSeeds(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	results := builder.BuildListWith(3, 100, nil)

	for i, instance := range results {
		expectedSeed := int64(100 + i)
		if instance.Seed != expectedSeed {
			t.Fatalf("expected seed %d, got %d", expectedSeed, instance.Seed)
		}
		actualValue := fmt.Sprintf("seed-%d", expectedSeed)
		if instance.Value != actualValue {
			t.Fatalf("expected %s, got %s", actualValue, instance.Value)
		}
	}
}

func TestBuilderPanicsOnUnsupportedOverride(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unsupported override type")
		}
	}()

	builder.Build(func() {})
}

func TestBuilderDuplicateAppliesOverride(t *testing.T) {
	factory := &stubFactory{}
	builder := Builder(factory)

	base := builder.BuildWith(55, Override[stubProps](map[string]any{
		"Value": "base",
	}))

	dup := builder.Duplicate(base, Override[stubProps](map[string]any{
		"Value": "copied",
	}))

	if dup.Seed != 55 {
		t.Fatalf("expected seed 55, got %d", dup.Seed)
	}

	if dup.Value != "copied" {
		t.Fatalf("expected value 'copied', got %s", dup.Value)
	}

	if base.Value != "base" {
		t.Fatalf("base should remain unchanged, got %s", base.Value)
	}
}
