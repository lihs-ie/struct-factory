package factory

// Partial mutates factory properties prior to instantiation.
type Partial[P any] func(*P)

// Factory describes the lifecycle hooks required to build values of type T.
type Factory[T any, P any] interface {
	Instantiate(properties P) T
	Prepare(overrides Partial[P], seed int64) P
	Retrieve(instance T) P
}
