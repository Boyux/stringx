package strmut

import (
	"unicode/utf8"
)

type String struct {
	mem []byte
	len int
	cap int
}

func (s *String) grow(n int) {
	s.mem = append(s.mem, make([]byte, n)...)
	s.cap += n
}

func (s *String) payload() []byte {
	return s.mem[0:s.len]
}

func (s *String) block() []byte {
	return s.mem
}

func (s *String) trim(f func(r rune) bool) {
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
