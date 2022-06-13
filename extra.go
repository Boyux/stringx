package st

import (
	"bytes"
	"unicode/utf8"
)

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
	mem []byte
	idx int
	sep []byte
}

func (s *Split) Next() bool {
	return s.idx < len(s.mem)
}

func (s *Split) Value() String {
	loc := bytes.Index(s.mem[s.idx:], s.sep)

	if loc < 0 {
		s.idx = len(s.mem)
		return FromBytes(s.mem[s.idx:])
	}

	next := FromBytes(s.mem[s.idx : s.idx+loc])
	s.idx += loc + len(s.sep)

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

type ToString interface {
	ToString() String
}
