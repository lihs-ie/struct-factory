package math

func asUint32(original uint64) uint32 {
	return uint32(original & 0xFFFFFFFF)
}

func invert(original uint32) uint32 {
	masks := []uint32{0x55555555, 0x33333333, 0x0F0F0F0F, 0x00FF00FF, 0xFFFFFFFF}
	carry := original

	for index, mask := range masks {
		padding := uint32(1 << index)
		left := (carry >> padding) & mask
		right := (carry & mask) << padding
		carry = asUint32(uint64(left | right))
	}

	return carry
}

func extendedGCD(a, b uint32) (gcd uint32, x int64, y int64) {
	if a == 0 {
		return b, 0, 1
	}

	g, yTemp, xTemp := extendedGCD(asUint32(uint64(b%a)), a)

	return g, xTemp - int64(b/a)*yTemp, yTemp
}

func modularInverse(a, b uint32) uint32 {
	gcd, x, _ := extendedGCD(a, b)

	if gcd != 1 {
		panic("no inverse found")
	}

	result := x % int64(b)
	if result < 0 {
		result += int64(b)
	}

	return uint32(result)
}

const salt uint32 = 0x17654321
const invertedSalt uint32 = 0x700000E1

func Scramble(original uint32) uint32 {
	normalized := asUint32(uint64(original))
	base := asUint32(uint64(normalized) * uint64(salt))
	inverted := invert(base)

	return asUint32(uint64(inverted) * uint64(invertedSalt))
}
