package app

import (
	"math/rand"

	"github.com/lihs-ie/forge/internal/collections"
)

const maxSafeInteger = 1<<53 - 1

type BuilderInstance[T any, P any] struct {
	Build         func(overrides Partial[P]) T
	BuildList     func(size int, overrides Partial[P]) []T
	BuildWith     func(seed int64, overrides Partial[P]) T
	BuildListWith func(size int, seed int64, overrides Partial[P]) []T
	Duplicate     func(instance T, overrides Partial[P]) T
}

func Builder[T any, P any](factory Factory[T, P]) *BuilderInstance[T, P] {
	seeds := collections.NewSet[int64](nil)

	var nextSeed func() int64
	var nextSeeds func(size int) []int64

	nextSeeds = func(size int) []int64 {
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

	nextSeed = func() int64 {
		return nextSeeds(1)[0]
	}

	return &BuilderInstance[T, P]{
		Build: func(overrides Partial[P]) T {
			seed := nextSeed()
			return Create(factory, overrides, seed)
		},

		BuildList: func(size int, overrides Partial[P]) []T {
			seedList := nextSeeds(size)
			results := make([]T, 0, size)

			for _, seed := range seedList {
				results = append(results, Create(factory, overrides, seed))
			}

			return results
		},

		BuildWith: func(seed int64, overrides Partial[P]) T {
			return Create(factory, overrides, seed)
		},

		BuildListWith: func(size int, seed int64, overrides Partial[P]) []T {
			results := make([]T, 0, size)

			for i := 0; i < size; i++ {
				results = append(results, Create(factory, overrides, seed+int64(i)))
			}

			return results
		},

		Duplicate: func(instance T, overrides Partial[P]) T {
			return Duplicate(factory, instance, overrides)
		},
	}
}
