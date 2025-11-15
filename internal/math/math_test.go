package math

import (
	"testing"
)

func TestAsUint32(t *testing.T) {
	tests := []struct {
		input    uint64
		expected uint32
	}{
		{0, 0},
		{1, 1},
		{0xFFFFFFFF, 0xFFFFFFFF},
		{0x100000000, 0},
		{0x100000001, 1},
		{0x123456789ABCDEF0, 0x9ABCDEF0},
	}

	for _, test := range tests {
		result := asUint32(test.input)
		if result != test.expected {
			t.Errorf("asUint32(0x%X) = 0x%X, expected 0x%X", test.input, result, test.expected)
		}
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{0x00000000, 0x00000000},
		{0xFFFFFFFF, 0xFFFFFFFF},
		{0x12345678, 0x1E6A2C48},
		{0x00000001, 0x80000000},
		{0x80000000, 0x00000001},
	}

	for _, test := range tests {
		result := invert(test.input)
		if result != test.expected {
			t.Errorf("invert(0x%08X) = 0x%08X, expected 0x%08X", test.input, result, test.expected)
		}
	}
}

func TestExtendedGCD(t *testing.T) {
	tests := []struct {
		a           uint32
		b           uint32
		expectedGCD uint32
	}{
		{0, 5, 5},
		{5, 0, 5},
		{10, 15, 5},
		{17, 13, 1},
		{252, 105, 21},
		{1, 1, 1},
	}

	for _, test := range tests {
		gcd, _, _ := extendedGCD(test.a, test.b)
		if gcd != test.expectedGCD {
			t.Errorf("extendedGCD(%d, %d) GCD = %d, expected %d", test.a, test.b, gcd, test.expectedGCD)
		}
	}
}

func TestExtendedGCDCoefficients(t *testing.T) {
	a := uint32(17)
	b := uint32(13)

	gcd, x, y := extendedGCD(a, b)

	if gcd != 1 {
		t.Errorf("extendedGCD(%d, %d) GCD = %d, expected 1", a, b, gcd)
	}

	result := asUint32(uint64(a)*uint64(x) + uint64(b)*uint64(y))
	if result != gcd {
		t.Errorf("Bezout identity failed: %d*%d + %d*%d = %d, expected %d", a, x, b, y, result, gcd)
	}
}

func TestModularInverse(t *testing.T) {
	tests := []struct {
		a uint32
		b uint32
	}{
		{3, 11},
		{7, 26},
		{17, 43},
	}

	for _, test := range tests {
		inverse := modularInverse(test.a, test.b)
		result := asUint32(uint64(test.a) * uint64(inverse) % uint64(test.b))

		if result != 1 {
			t.Errorf("modularInverse(%d, %d) = %d, but %d * %d mod %d = %d, expected 1",
				test.a, test.b, inverse, test.a, inverse, test.b, result)
		}
	}
}

func TestModularInversePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for numbers with no inverse")
		}
	}()

	modularInverse(4, 8)
}

func TestInvertedSaltConstant(t *testing.T) {
	modulus := uint64(0x100000000)

	gcd, x, _ := extendedGCD64(uint64(salt), modulus)

	if gcd != 1 {
		t.Errorf("salt and modulus are not coprime")
	}

	computed := uint32((x%int64(modulus) + int64(modulus)) % int64(modulus))

	if computed != invertedSalt {
		t.Errorf("invertedSalt constant is incorrect: got 0x%08X, expected 0x%08X", invertedSalt, computed)
	}

	verification := asUint32(uint64(salt) * uint64(computed))
	if verification != 1 {
		t.Errorf("salt * invertedSalt mod 2^32 = %d, expected 1", verification)
	}
}

func extendedGCD64(a, b uint64) (gcd uint64, x int64, y int64) {
	if a == 0 {
		return b, 0, 1
	}

	g, yTemp, xTemp := extendedGCD64(b%a, a)

	return g, xTemp - int64(b/a)*yTemp, yTemp
}

func TestScramble(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{0, 0},
		{1, 0x2F14B1E8},
		{2, 0x178A58F4},
		{100, 0x1FFFFD44},
		{1000, 0xEFFFFF06},
	}

	for _, test := range tests {
		result := Scramble(test.input)
		if result != test.expected {
			t.Errorf("Scramble(%d) = 0x%08X, expected 0x%08X", test.input, result, test.expected)
		}
	}
}

func TestScrambleRoundTrip(t *testing.T) {
	original := uint32(12345)

	scrambled := Scramble(original)

	if scrambled == original {
		t.Error("Scrambled value should differ from original")
	}
}

func TestScrambleDifferentInputs(t *testing.T) {
	inputs := []uint32{0, 1, 2, 100, 1000, 0xFFFFFFFF}
	results := make(map[uint32]bool)

	for _, input := range inputs {
		result := Scramble(input)
		if results[result] {
			t.Errorf("Scramble produced duplicate result for different inputs")
		}
		results[result] = true
	}
}

func TestInvertInvolutive(t *testing.T) {
	tests := []uint32{0, 1, 0x12345678, 0xFFFFFFFF, 0x80000000}

	for _, value := range tests {
		inverted := invert(value)
		doubleInverted := invert(inverted)

		if doubleInverted != value {
			t.Errorf("invert(invert(0x%08X)) = 0x%08X, expected 0x%08X", value, doubleInverted, value)
		}
	}
}

func TestExtendedGCDSymmetry(t *testing.T) {
	a := uint32(35)
	b := uint32(15)

	gcd1, _, _ := extendedGCD(a, b)
	gcd2, _, _ := extendedGCD(b, a)

	if gcd1 != gcd2 {
		t.Errorf("extendedGCD is not symmetric: GCD(%d, %d) = %d, GCD(%d, %d) = %d",
			a, b, gcd1, b, a, gcd2)
	}
}
