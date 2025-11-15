package factory

func create[T any, P any](factory Factory[T, P], overrides Partial[P], seed int64) T {
	properties := factory.Prepare(overrides, seed)
	return factory.Instantiate(properties)
}

func duplicate[T any, P any](factory Factory[T, P], instance T, overrides Partial[P]) T {
	properties := factory.Retrieve(instance)
	if overrides != nil {
		overrides(&properties)
	}
	return factory.Instantiate(properties)
}
