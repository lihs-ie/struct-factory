package factory

import "github.com/lihs-ie/forge/internal/collections"

// EnumProperties captures the selected value and exclusions for EnumFactory.
type EnumProperties[T comparable] struct {
	value      T
	exclusions []T
}

// EnumFactory selects values from a predefined candidate set.
type EnumFactory[T comparable] struct {
	candidates *collections.Set[T]
}

// NewEnumFactory constructs an EnumFactory with the provided candidates.
func NewEnumFactory[T comparable](candidates []T) *EnumFactory[T] {
	return &EnumFactory[T]{
		candidates: collections.NewFromSlice(candidates),
	}
}

// Instantiate returns the chosen enum value.
func (f *EnumFactory[T]) Instantiate(properties EnumProperties[T]) T {
	return properties.value
}

// Prepare applies overrides and exclusions before choosing a value.
func (f *EnumFactory[T]) Prepare(overrides Partial[EnumProperties[T]], seed int64) EnumProperties[T] {
	properties := EnumProperties[T]{
		exclusions: []T{},
	}

	if overrides != nil {
		overrides(&properties)
	}

	actuals := f.filterExclusions(properties.exclusions)

	if len(actuals) == 0 {
		panic("no candidates available after exclusions")
	}

	index := int(seed % int64(len(actuals)))

	var zero T
	if properties.value == zero {
		properties.value = actuals[index]
	}

	return properties
}

// Retrieve wraps an existing instance into EnumProperties.
func (f *EnumFactory[T]) Retrieve(instance T) EnumProperties[T] {
	return EnumProperties[T]{
		value:      instance,
		exclusions: []T{},
	}
}

func (f *EnumFactory[T]) filterExclusions(exclusions []T) []T {
	exclusionSet := collections.NewFromSlice(exclusions)
	result := make([]T, 0)

	candidates := f.candidates.ToSlice()
	for _, candidate := range candidates {
		if !exclusionSet.Has(candidate) {
			result = append(result, candidate)
		}
	}

	return result
}
