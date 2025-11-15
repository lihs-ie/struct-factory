package factory

import (
	"github.com/lihs-ie/forge/internal/math"
)

// CharacterSet defines a pool of runes used for random string generation.
type CharacterSet []rune

// Characters provides common rune sets for string generation.
var Characters = struct {
	Alphanumeric CharacterSet
	Alpha        CharacterSet
	Numeric      CharacterSet
	Symbol       CharacterSet
}{
	Alphanumeric: CharacterSet{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	},
	Alpha: CharacterSet{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	},
	Numeric: CharacterSet{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	},
	Symbol: CharacterSet{
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.',
		'/', ':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`',
		'{', '|', '}', '~',
	},
}

// StringProperties carries configuration and generated values for StringFactory.
type StringProperties struct {
	value      string
	min        int
	max        int
	characters CharacterSet
}

// StringFactory generates random strings with configurable constraints.
type StringFactory struct {
	Min        int
	Max        int
	Characters CharacterSet
}

// Instantiate returns the final string value from prepared properties.
func (f *StringFactory) Instantiate(properties StringProperties) string {
	return properties.value
}

// Prepare produces StringProperties using the provided seed and overrides.
func (f *StringFactory) Prepare(overrides Partial[StringProperties], seed int64) StringProperties {
	min := f.Min
	if min <= 0 {
		min = 1
	}

	max := f.Max
	if max <= 0 {
		max = 255
	}
	if max < min {
		max = min
	}

	chars := f.Characters
	if len(chars) == 0 {
		chars = Characters.Alphanumeric
	}

	properties := StringProperties{
		min:        min,
		max:        max,
		characters: chars,
	}

	if overrides != nil {
		overrides(&properties)
	}

	if properties.min <= 0 {
		properties.min = 1
	}
	if properties.max < properties.min {
		properties.max = properties.min
	}
	if len(properties.characters) == 0 {
		properties.characters = Characters.Alphanumeric
	}

	if properties.value == "" {
		offset := seed % int64(properties.max-properties.min+1)
		length := properties.min + int(offset)

		value := make([]rune, length)
		for index := range length {
			scrambled := math.Scramble(uint32(seed + int64(index)))
			characterIndex := int(scrambled) % len(properties.characters)
			value[index] = properties.characters[characterIndex]
		}

		properties.value = string(value)
	}

	return properties
}

// Retrieve converts a string instance back into StringProperties.
func (f *StringFactory) Retrieve(instance string) StringProperties {
	return StringProperties{
		value: instance,
	}
}
