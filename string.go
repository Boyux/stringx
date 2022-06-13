package st

import (
	"bytes"
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

type Iterator[T any] interface {
	Next() bool
	Value() T
}

var _ Iterator[rune] = (*Runes)(nil)

type Runes struct {
	mem []byte
	idx int
}

func (r *Runes) Next() bool {
	return r.idx < len(r.mem)
}

func (r *Runes) Value() rune {
	next, n := utf8.DecodeRune(r.mem[r.idx:])
	r.idx += n
	return next
}

func (r *Runes) Nth(i int) rune {
	for j := 0; j < i; j++ {
		_ = r.Value()
	}

	return r.Value()
}

func (r *Runes) Size() int {
	return utf8.RuneCount(r.mem[r.idx:])
}

func (r *Runes) Consume() []rune {
	slice := make([]rune, 0, r.Size())

	for r.Next() {
		slice = append(slice, r.Value())
	}

	return slice
}

var _ Iterator[String] = (*Split)(nil)

type Split struct {
	mem        []byte
	idx        int
	sep        []byte
	cacheIndex int
}

func (s *Split) Next() bool {
	s.cacheIndex = bytes.Index(s.mem[s.idx:], s.sep)
	return s.cacheIndex >= 0
}

func (s *Split) Value() String {
	if s.cacheIndex < 0 {
		if s.cacheIndex = bytes.Index(s.mem[s.idx:], s.sep); s.cacheIndex < 0 {
			return FromBytes(s.mem[s.idx:])
		}
	}

	next := FromBytes(s.mem[s.idx:s.cacheIndex])
	s.idx += s.cacheIndex + len(s.sep)
	s.cacheIndex = -1 // reset cache

	return next
}

func (s *Split) Size() int {
	return bytes.Count(s.mem[s.idx:], s.sep)
}

func (s *Split) Consume() []String {
	slice := make([]String, 0, s.Size())

	for s.Next() {
		slice = append(slice, s.Value())
	}

	return slice
}

type FromString interface {
	FromString(*String) error
}
