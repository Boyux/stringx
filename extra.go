package strmut

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type Iterator[T any] interface {
	Next() bool
	Value() T
}

var _ Iterator[*String] = (*Lines)(nil)

type Lines struct {
	mem []byte
	idx int
}

func (l *Lines) Next() bool {
	return l.idx < len(l.mem)
}

func (l *Lines) Value() *String {
	dropCR := func(data []byte) []byte {
		if len(data) > 0 && data[len(data)-1] == '\r' {
			return data[0 : len(data)-1]
		}
		return data
	}

	var next String

	loc := bytes.IndexByte(l.mem[l.idx:], '\n')

	if loc < 0 {
		next.FromBytes(dropCR(l.mem[l.idx:]))
		l.idx = len(l.mem)
		return &next
	}

	next.FromBytes(dropCR(l.mem[l.idx : l.idx+loc]))
	l.idx += loc + 1

	return &next
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

func (r *Runes) Reverse() *ReverseRunes {
	return &ReverseRunes{
		runes: r,
		last:  len(r.mem),
	}
}

type ReverseRunes struct {
	runes *Runes
	last  int
}

func (r *ReverseRunes) Next() bool {
	return r.runes.idx < r.last
}

func (r *ReverseRunes) Value() rune {
	next, n := utf8.DecodeLastRune(r.runes.mem[r.runes.idx:r.last])
	r.last -= n
	return next
}

func (r *ReverseRunes) Size() int {
	return utf8.RuneCount(r.runes.mem[r.runes.idx:r.last])
}

func (r *ReverseRunes) Consume() []rune {
	slice := make([]rune, 0, r.Size())

	for r.Next() {
		slice = append(slice, r.Value())
	}

	return slice
}

var _ Iterator[*String] = (*Split)(nil)

type Split struct {
	mem []byte
	idx int
	sep []byte
}

func (s *Split) Next() bool {
	return s.idx < len(s.mem)
}

func (s *Split) Value() *String {
	var next String

	loc := bytes.Index(s.mem[s.idx:], s.sep)

	if loc < 0 {
		next.FromBytes(s.mem[s.idx:])
		s.idx = len(s.mem)
		return &next
	}

	next.FromBytes(s.mem[s.idx : s.idx+loc])
	s.idx += loc + len(s.sep)

	return &next
}

func (s *Split) Size() int {
	return bytes.Count(s.mem[s.idx:], s.sep)
}

func (s *Split) Consume() []*String {
	slice := make([]*String, 0, s.Size())

	for s.Next() {
		slice = append(slice, s.Value())
	}

	return slice
}

type FromString interface {
	FromString(*String) error
}

type ToString interface {
	ToString() *String
}

type Int int

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

func (i Int) ToString() *String {
	var s String
	s.FromString(i.String())
	return &s
}

type Str string

func (str Str) String() string {
	return string(str)
}

func (str Str) ToString() *String {
	var s String
	s.FromString(str.String())
	return &s
}

func (str Str) Len() int {
	return len(str)
}

type Initializer[T any] interface {
	Initialize(T)
}

var _ Initializer[*String] = StringInitializer("StringInitializer")
var _ Initializer[*String] = BytesInitializer("BytesInitializer")
var _ Initializer[*String] = RunesInitializer("RunesInitializer")
var _ Initializer[*String] = (*String)(nil)

type StringInitializer string

func (str StringInitializer) Initialize(s *String) {
	if s.cap < len(str) {
		s.grow(len(str))
	}

	copy(s.mem[0:], str)
	s.len = len(str)
}

type BytesInitializer []byte

func (b BytesInitializer) Initialize(s *String) {
	if s.cap < len(b) {
		s.grow(len(b))
	}

	copy(s.mem[0:], b)
	s.len = len(b)
}

type RunesInitializer []rune

func (r RunesInitializer) Initialize(s *String) {
	l := len(r) * utf8.UTFMax
	if s.cap < l {
		s.grow(l)
	}

	var n int
	for _, rr := range r {
		n += utf8.EncodeRune(s.mem[n:], rr)
	}
	s.len = n
}

type From[Self any, T Initializer[Self]] interface {
	From(T) Self
}

var _ From[*String, Initializer[*String]] = (*String)(nil)

type List[S fmt.Stringer] []S

var _ = List[*String]{(*String)(nil)}
var _ interface{ Len() int } = (*String)(nil)
var _ interface{ Len() int } = Str("Str.Len")

func (l List[S]) Join(sep string) *String {
	var s String

	if len(l) == 0 {
		s.Init()
		return &s
	}

	if len(l) == 1 {
		s.FromString(l[0].String())
		return &s
	}

	var head any = l[0]
	if _, ok := head.(interface{ Len() int }); ok {
		n := len(sep) * (len(l) - 1)
		for i := 0; i < len(l); i++ {
			var ele any = l[i]
			n += ele.(interface{ Len() int }).Len()
		}
		s.SetCapacity(n)
	}

	if !s.alreadyInit() {
		s.Init()
	}

	s.PushString(l[0].String())
	for i := 1; i < len(l); i++ {
		s.PushString(sep)
		s.PushString(l[i].String())
	}

	return &s
}
