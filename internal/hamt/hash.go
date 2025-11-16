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
	hashResult, _ := hashValue(hasher, reflect.ValueOf(value))
	return hashResult
}

var timeType = reflect.TypeOf(time.Time{})

// unwrapValue removes interface and pointer wrapping from a reflect.Value.
func unwrapValue(value reflect.Value) reflect.Value {
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
	return value
}

// tryHashable checks if the value implements the Hashable interface and returns its hash.
func tryHashable(value reflect.Value) (hashValue uint64, found bool, err error) {
	if value.CanInterface() {
		if hashable, ok := value.Interface().(Hashable); ok {
			hashValue, err = hashable.Hash()
			return hashValue, true, err
		}
	}
	return 0, false, nil
}

// normalizeValue converts platform-dependent types (int, uint, bool) to fixed-size types.
func normalizeValue(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Int:
		return reflect.ValueOf(value.Int())
	case reflect.Uint:
		return reflect.ValueOf(value.Uint())
	case reflect.Bool:
		var temp int8
		if value.Bool() {
			temp = 1
		}
		return reflect.ValueOf(temp)
	}
	return value
}

// hashNil returns the hash for nil values.
func hashNil(hasher hash.Hash64) uint64 {
	hasher.Reset()
	return hasher.Sum64()
}

// hashNumeric returns the hash for numeric types (int8-64, uint8-64, float32-64, complex64-128).
func hashNumeric(hasher hash.Hash64, value reflect.Value) (uint64, error) {
	hasher.Reset()
	if err := binary.Write(hasher, binary.LittleEndian, value.Interface()); err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

// hashString returns the hash for string values.
func hashString(hasher hash.Hash64, value reflect.Value) uint64 {
	hasher.Reset()
	hasher.Write([]byte(value.String()))
	return hasher.Sum64()
}

// hashTime returns the hash for time.Time values.
func hashTime(hasher hash.Hash64, value reflect.Value) (uint64, error) {
	hasher.Reset()
	bytes, err := value.Interface().(time.Time).MarshalBinary()
	if err != nil {
		return 0, err
	}
	hasher.Write(bytes)
	return hasher.Sum64(), nil
}

// hashSequence returns the hash for arrays and slices (order-dependent).
func hashSequence(hasher hash.Hash64, value reflect.Value) (uint64, error) {
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
}

// hashMap returns the hash for maps (order-independent using XOR).
func hashMap(hasher hash.Hash64, value reflect.Value) (uint64, error) {
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
}

// hashStruct returns the hash for struct values.
func hashStruct(hasher hash.Hash64, value reflect.Value) (uint64, error) {
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
}

func hashValue(hasher hash.Hash64, value reflect.Value) (uint64, error) {
	// Unwrap pointers and interfaces
	value = unwrapValue(value)

	// Handle invalid values (nil)
	if !value.IsValid() {
		return hashNil(hasher), nil
	}

	// Check if the value implements Hashable interface
	if result, ok, err := tryHashable(value); ok {
		return result, err
	}

	// Normalize platform-dependent types
	value = normalizeValue(value)

	// Handle time.Time specially (must be before numeric check)
	if value.Type() == timeType {
		return hashTime(hasher, value)
	}

	// Dispatch based on kind
	kind := value.Kind()

	// Handle numeric types
	if kind >= reflect.Int && kind <= reflect.Complex128 {
		return hashNumeric(hasher, value)
	}

	// Handle other types
	switch kind {
	case reflect.String:
		return hashString(hasher, value), nil

	case reflect.Array, reflect.Slice:
		return hashSequence(hasher, value)

	case reflect.Map:
		return hashMap(hasher, value)

	case reflect.Struct:
		return hashStruct(hasher, value)

	default:
		// For unsupported types (chan, func, etc.), return a default hash
		return hashNil(hasher), nil
	}
}

// hashUpdateOrdered combines two hash values in an order-dependent way.
func hashUpdateOrdered(hasher hash.Hash64, a, b uint64) uint64 {
	hasher.Reset()

	// Convert uint64 values to bytes and write to hasher.
	// Using a byte buffer avoids the need for error checking.
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[0:8], a)
	binary.LittleEndian.PutUint64(buf[8:16], b)
	hasher.Write(buf[:])

	return hasher.Sum64()
}
