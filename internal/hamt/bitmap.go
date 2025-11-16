package hamt

import "math/bits"

type Bitmap uint64

const (
	shiftWidth = 6
	maxOffset  = 10 // Maximum HAMT depth: 64-bit hash / 6-bit width ≈ 10 levels
	mask64     = ^uint64(0)
	mixMul1    = 0xBF58476D1CE4E589
	mixMul2    = 0x94D049BB133111EB
)

func Initialize() Bitmap {
	return Bitmap(0)
}

func NewBitmap(value uint64) Bitmap {
	return Bitmap(value)
}

func bandMask() uint64 {
	return (1 << shiftWidth) - 1
}

func (bitmap Bitmap) Position(hash uint64, offset int) uint64 {
	// Safety: offset is constrained by HAMT depth (0-10), making conversion to uint safe.
	// HAMT uses 6-bit chunks from a 64-bit hash, limiting depth to ~10 levels.
	if offset < 0 || offset > maxOffset {
		panic("bitmap: offset out of valid HAMT range")
	}
	// Conversion is safe: offset bounds-checked above [0, 10]
	shift := uint(offset * shiftWidth)

	shifted := hash >> shift
	masked := shifted & bandMask()

	return 1 << masked
}

func (bitmap Bitmap) Next(position uint64) Bitmap {
	return Bitmap(uint64(bitmap) | position)
}

func (bitmap Bitmap) Without(position uint64) Bitmap {
	return Bitmap(uint64(bitmap) &^ position)
}

func (bitmap Bitmap) Has(position uint64) bool {
	return (uint64(bitmap) & position) != 0
}

func (bitmap Bitmap) Index(position uint64) (count int, ok bool) {
	msbMask := msbMask(position)

	if msbMask == ^uint64(0) {
		return 0, false
	}

	masked := bitmap & Bitmap(msbMask)

	return bits.OnesCount64(uint64(masked)), true
}

func msbMask(value uint64) uint64 {
	if value == 0 {
		return ^uint64(0)
	}

	x := value
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32

	return x >> 1
}

func splitMix64(x uint64) uint64 {
	y := x & mask64
	y ^= y >> 30
	y *= mixMul1
	y &= mask64

	y ^= y >> 27
	y *= mixMul2
	y &= mask64

	y ^= y >> 31
	y &= mask64

	return y
}
