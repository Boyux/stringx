package st

import "unicode/utf8"

func New() String {
	return WithCapacity(0)
}

func WithCapacity(capacity int) String {
	return String{
		mem: make([]byte, capacity),
		len: 0,
		cap: capacity,
	}
}

func From(in string) String {
	// convert input to bytes, there is a faster way by using
	// unsafe.Pointer, which is not recommended
	//
	// SAFETY
	// because String is a mutable type, so String.mem should
	// also be mutable. However, byte slice converted by unsafe
	// function stringToBytes is immutable, mutating those bytes
	// would cause 'unexpected fault address' error
	mem := stringToBytesSlow(in)

	return String{
		mem: mem,
		len: len(mem),
		// we use slice length as capacity but not the slice cap
		cap: len(mem),
	}
}

func FromBytes(in []byte) String {
	mem := make([]byte, len(in))
	copy(mem, in)

	return String{
		mem: mem,
		len: len(mem),
		cap: len(mem),
	}
}

func FromRune(in []rune) String {
	mem := make([]byte, len(in)*utf8.UTFMax)

	var n int
	for _, r := range in {
		n += utf8.EncodeRune(mem[n:], r)
	}

	return String{
		mem: mem,
		len: n,
		cap: len(mem),
	}
}
