package hamt

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"reflect"
	"time"
)

// Hashable is an interface that can optionally be implemented by a value
// to override the default hashing behavior.
type Hashable interface {
	Hash() (uint64, error)
}

// Hash returns the hash value of an arbitrary value using FNV-1a algorithm.
// This function uses reflection to handle any Go type.
func Hash(value any) uint64 {
	hasher := fnv.New64a()
	hash, _ := hashValue(hasher, reflect.ValueOf(value))
	return hash
}

var timeType = reflect.TypeOf(time.Time{})

func hashValue(hasher hash.Hash64, value reflect.Value) (uint64, error) {
	// Handle interface and pointer dereferencing
	for {
		if value.Kind() == reflect.Interface {
			value = value.Elem()
			continue
		}

		if value.Kind() == reflect.Ptr {
			value = reflect.Indirect(value)
			continue
		}

		break
	}

	// Handle invalid values (nil)
	if !value.IsValid() {
		hasher.Reset()
		return hasher.Sum64(), nil
	}

	// Check if the value implements Hashable interface
	if value.CanInterface() {
		if hashable, ok := value.Interface().(Hashable); ok {
			return hashable.Hash()
		}
	}

	// Convert int/uint/bool to fixed-size types
	switch value.Kind() {
	case reflect.Int:
		value = reflect.ValueOf(int64(value.Int()))
	case reflect.Uint:
		value = reflect.ValueOf(uint64(value.Uint()))
	case reflect.Bool:
		var temp int8
		if value.Bool() {
			temp = 1
		}
		value = reflect.ValueOf(temp)
	}

	kind := value.Kind()

	// Handle numeric types directly
	if kind >= reflect.Int && kind <= reflect.Complex128 {
		hasher.Reset()
		binary.Write(hasher, binary.LittleEndian, value.Interface())
		return hasher.Sum64(), nil
	}

	// Handle time.Time specially
	if value.Type() == timeType {
		hasher.Reset()
		bytes, err := value.Interface().(time.Time).MarshalBinary()
		if err != nil {
			return 0, err
		}
		hasher.Write(bytes)
		return hasher.Sum64(), nil
	}

	// Handle different kinds
	switch kind {
	case reflect.String:
		hasher.Reset()
		hasher.Write([]byte(value.String()))
		return hasher.Sum64(), nil

	case reflect.Array, reflect.Slice:
		var result uint64
		length := value.Len()
		for i := 0; i < length; i++ {
			elementHash, err := hashValue(hasher, value.Index(i))
			if err != nil {
				return 0, err
			}
			result = hashUpdateOrdered(hasher, result, elementHash)
		}
		return result, nil

	case reflect.Map:
		var result uint64
		for _, key := range value.MapKeys() {
			keyHash, err := hashValue(hasher, key)
			if err != nil {
				return 0, err
			}

			valueHash, err := hashValue(hasher, value.MapIndex(key))
			if err != nil {
				return 0, err
			}

			// XOR for unordered combination (maps are unordered)
			fieldHash := hashUpdateOrdered(hasher, keyHash, valueHash)
			result ^= fieldHash
		}
		return result, nil

	case reflect.Struct:
		typeInfo := value.Type()
		typeNameHash, _ := hashValue(hasher, reflect.ValueOf(typeInfo.Name()))
		result := typeNameHash

		fieldCount := value.NumField()
		for i := 0; i < fieldCount; i++ {
			field := typeInfo.Field(i)

			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			fieldValue := value.Field(i)

			// Hash field name
			fieldNameHash, err := hashValue(hasher, reflect.ValueOf(field.Name))
			if err != nil {
				return 0, err
			}

			// Hash field value
			fieldValueHash, err := hashValue(hasher, fieldValue)
			if err != nil {
				return 0, err
			}

			// Combine field name and value hashes
			fieldHash := hashUpdateOrdered(hasher, fieldNameHash, fieldValueHash)
			result ^= fieldHash
		}

		return result, nil

	default:
		// For unsupported types, return a default hash
		hasher.Reset()
		return hasher.Sum64(), nil
	}
}

// hashUpdateOrdered combines two hash values in an order-dependent way
func hashUpdateOrdered(hasher hash.Hash64, a, b uint64) uint64 {
	hasher.Reset()
	binary.Write(hasher, binary.LittleEndian, a)
	binary.Write(hasher, binary.LittleEndian, b)
	return hasher.Sum64()
}
