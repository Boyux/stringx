package stringx

import (
	"math"
	"unicode/utf8"
	"unsafe"
)

type String struct {
	// nocopy embed this type into a struct, which mustn't be copied,
	// so `go vet` gives a warning if this struct is copied.
	//
	// See https://github.com/golang/go/issues/8005#issuecomment-190753527 for details.
	// and also: https://stackoverflow.com/questions/52494458/nocopy-minimal-example
	nocopy nocopy

	// self represents receiver of this String, to detect copies by value
	// See type definition Builder in strings/builder.go
	self unsafe.Pointer

	mem []byte
	len int
	cap int
}

func (s *String) build(mem []byte, len, cap int) {
	s.mem = mem
	s.len = len
	s.cap = cap
	// s.nocopy = nocopy{}
	s.self = unsafe.Pointer(s)
}

var nullptr = unsafe.Pointer((*String)(nil))

func (s *String) alreadyInit() bool {
	if s.self != nullptr {
		s.copycheck()
		return true
	}

	return false
}

func (s *String) assumeUninit() {
	if s.alreadyInit() {
		panic("String: already initialized")
	}
}

func (s *String) copycheck() {
	if s.self == nullptr {
		panic("String: illegal use of uninitialized value")
	} else if s.self != unsafe.Pointer(s) {
		panic("String: illegal use of copied value")
	}
}

func (s *String) grow(n int) {
	s.copycheck()

	if n >= math.MaxInt32 {
		panic("String.grow: n overflows")
	}

	// next power of 2
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++

	s.mem = append(s.mem, make([]byte, n)...)
	s.cap += n
}

func (s *String) payload() []byte {
	return s.mem[0:s.len]
}

func (s *String) trim(f func(r rune) bool) {
	s.copycheck()

	var start, stop int
	payload := s.payload()

	for start < s.len {
		p := payload[start:]
		r, n := utf8.DecodeRune(p)
		if !f(r) {
			break
		}
		start += n
	}

	stop = s.len
	for stop > start {
		p := payload[start:stop]
		r, n := utf8.DecodeLastRune(p)
		if !f(r) {
			break
		}
		stop -= n
	}

	if start == stop {
		s.len = 0
		return
	}

	copy(payload, payload[start:stop])
	s.len = stop - start
}
