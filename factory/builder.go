package factory

import (
	"math/rand"

	"github.com/lihs-ie/struct-factory/internal/collections"
)

const maxSafeInteger = 1<<53 - 1

// BuilderHandle exposes the supported build operations for a factory.
type BuilderHandle[T any, P any] interface {
	Build(overrides any) T
	BuildList(size int, overrides any) []T
	BuildWith(seed int64, overrides any) T
	BuildListWith(size int, seed int64, overrides any) []T
	Duplicate(instance T, overrides any) T
}

type builderInstance[T any, P any] struct {
	factory         Factory[T, P]
	nextSeed        func() int64
	nextSeeds       func(size int) []int64
	convertOverride func(any) Partial[P]
}

var _ BuilderHandle[any, any] = (*builderInstance[any, any])(nil)

// Builder creates a BuilderHandle for the provided Factory.
func Builder[T any, P any](factory Factory[T, P]) BuilderHandle[T, P] {
	seeds := collections.NewSet[int64](nil)

	nextSeeds := func(size int) []int64 {
		next := make([]int64, 0, size)

		for len(next) < size {
			seed := rand.Int63n(maxSafeInteger)
			if !seeds.Has(seed) {
				next = append(next, seed)
				seeds.Set(seed)
			}
		}

		return next
	}

	nextSeed := func() int64 {
		return nextSeeds(1)[0]
	}

	convertOverride := func(override any) Partial[P] {
		if override == nil {
			return nil
		}

		overrider, ok := override.(Overrider[P])
		if !ok {
			panic("builder: overrides must be generated via Override()")
		}
		return overrider.Func()
	}

	return &builderInstance[T, P]{
		factory:         factory,
		nextSeed:        nextSeed,
		nextSeeds:       nextSeeds,
		convertOverride: convertOverride,
	}
}

func (b *builderInstance[T, P]) Build(overrides any) T {
	seed := b.nextSeed()
	return create(b.factory, b.convertOverride(overrides), seed)
}

func (b *builderInstance[T, P]) BuildList(size int, overrides any) []T {
	seedList := b.nextSeeds(size)
	results := make([]T, 0, size)
	converted := b.convertOverride(overrides)

	for _, seed := range seedList {
		results = append(results, create(b.factory, converted, seed))
	}

	return results
}

func (b *builderInstance[T, P]) BuildWith(seed int64, overrides any) T {
	return create(b.factory, b.convertOverride(overrides), seed)
}

func (b *builderInstance[T, P]) BuildListWith(size int, seed int64, overrides any) []T {
	results := make([]T, 0, size)
	converted := b.convertOverride(overrides)

	for i := range size {
		results = append(results, create(b.factory, converted, seed+int64(i)))
	}

	return results
}

func (b *builderInstance[T, P]) Duplicate(instance T, overrides any) T {
	return duplicate(b.factory, instance, b.convertOverride(overrides))
}
