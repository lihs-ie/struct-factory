package factory

// MapEntry represents a single key/value pair in MapFactory output.
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// MapProperties stores the entries prepared for MapFactory.
type MapProperties[K comparable, V any] struct {
	entries []MapEntry[K, V]
}

// MapFactory builds maps using dedicated key/value factories.
type MapFactory[K comparable, KP any, V any, VP any] struct {
	keyFactory   Factory[K, KP]
	valueFactory Factory[V, VP]
}

// NewMapFactory wires key and value factories into a MapFactory.
func NewMapFactory[K comparable, KP any, V any, VP any](
	keyFactory Factory[K, KP],
	valueFactory Factory[V, VP],
) *MapFactory[K, KP, V, VP] {
	return &MapFactory[K, KP, V, VP]{
		keyFactory:   keyFactory,
		valueFactory: valueFactory,
	}
}

// Instantiate converts prepared map properties into a concrete map.
func (f *MapFactory[K, KP, V, VP]) Instantiate(properties MapProperties[K, V]) map[K]V {
	result := make(map[K]V, len(properties.entries))

	for _, entry := range properties.entries {
		result[entry.Key] = entry.Value
	}

	return result
}

// Prepare builds random entries (optionally overridden) for later instantiation.
func (f *MapFactory[K, KP, V, VP]) Prepare(overrides Partial[MapProperties[K, V]], seed int64) MapProperties[K, V] {
	count := int((seed % 10) + 1)
	entries := make([]MapEntry[K, V], count)

	for index := range count {
		keyInstance := create(f.keyFactory, nil, seed+int64(index))
		valueInstance := create(f.valueFactory, nil, seed+int64(index))

		entries[index] = MapEntry[K, V]{
			Key:   keyInstance,
			Value: valueInstance,
		}
	}

	properties := MapProperties[K, V]{
		entries: entries,
	}

	if overrides != nil {
		overrides(&properties)
	}

	return properties
}

// Retrieve converts an existing map into MapProperties for duplication/override.
func (f *MapFactory[K, KP, V, VP]) Retrieve(instance map[K]V) MapProperties[K, V] {
	entries := make([]MapEntry[K, V], 0, len(instance))
	for key, value := range instance {
		entries = append(entries, MapEntry[K, V]{
			Key:   key,
			Value: value,
		})
	}

	return MapProperties[K, V]{
		entries: entries,
	}
}
