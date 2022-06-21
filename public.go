package strmut

import (
	"bytes"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func (s *String) Bytes() []byte {
	return s.payload()
}

func (s *String) Runes() *Runes {
	return &Runes{
		mem: s.payload(),
		idx: 0,
	}
}

func (s *String) RuneSlice() []rune {
	return s.Runes().Consume()
}

func (s *String) Length() int {
	return s.len
}

func (s *String) Capacity() int {
	return s.cap
}

func (s *String) Reset() {
	s.len = 0
}

func (s *String) IsEmpty() bool {
	return s.len == 0
}

func (s *String) Clone() String {
	mem := make([]byte, s.len)
	copy(mem, s.payload())

	return String{
		mem: mem,
		len: s.len,
		cap: s.len,
	}
}

func (s *String) Insert(i int, b byte) {
	if s.len >= s.cap {
		s.grow(s.cap)
	}

	copy(s.mem[i+1:s.len+1], s.mem[i:s.len])
	s.mem[i] = b
	s.len += 1
}

func (s *String) InsertString(i int, str string) {
	l := len(str)

	if s.len+l >= s.cap {
		s.grow(s.len + l - s.cap)
	}

	copy(s.mem[i+l:s.len+l], s.mem[i:s.len])
	copy(s.mem[i:i+l], str)
	s.len += l
}

func (s *String) Push(b byte) {
	if s.len >= s.cap {
		s.grow(s.cap)
	}

	s.mem[s.len] = b
	s.len += 1
}

func (s *String) PushRune(r rune) {
	if r < utf8.RuneSelf {
		s.Push(byte(r))
		return
	}

	if s.len+utf8.UTFMax > s.cap {
		s.grow(utf8.UTFMax)
	}

	n := utf8.EncodeRune(s.mem[s.len:s.cap], r)
	s.len += n
}

func (s *String) PushString(str string) {
	l := len(str)

	if s.len+l >= s.cap {
		s.grow(s.len + l - s.cap)
	}

	copy(s.mem[s.len:s.cap], str)
	s.len += l
}

func (s *String) PushBytes(bytes []byte) {
	l := len(bytes)

	if s.len+l >= s.cap {
		s.grow(s.len + l - s.cap)
	}

	copy(s.mem[s.len:s.cap], bytes)
	s.len += l
}

func (s *String) PushRunes(runes []rune) {
	for _, r := range runes {
		s.PushRune(r)
	}
}

func (s *String) Drain(l, r int) {
	copy(s.mem[l:s.len], s.mem[r:s.len])
	s.len -= r - l
}

func (s *String) Get(i int) byte {
	payload := s.payload()
	return payload[i]
}

func (s *String) Index(l, r int) String {
	mem := make([]byte, r-l)
	copy(mem, s.mem[l:r])

	return String{
		mem: mem,
		len: r - l,
		cap: r - l,
	}
}

func (s *String) EqualTo(other String) bool {
	if s.len != other.len {
		return false
	}

	return bytes.Equal(s.payload(), other.mem[0:other.len])
}

func (s *String) EqualToString(str string) bool {
	if s.len != len(str) {
		return false
	}

	return bytes.Equal(s.payload(), stringToBytes(str))
}

func (s *String) CompareTo(other String) int {
	return bytes.Compare(s.payload(), other.payload())
}

func (s *String) CompareToString(str string) int {
	return bytes.Compare(s.payload(), stringToBytes(str))
}

func (s *String) Contains(sub string) bool {
	return bytes.Contains(s.payload(), stringToBytes(sub))
}

func (s *String) StartsWith(pat string) bool {
	b := stringToBytes(pat)
	return s.len >= len(pat) && bytes.Equal(s.mem[0:len(b)], b)
}

func (s *String) EndsWith(pat string) bool {
	b := stringToBytes(pat)
	return s.len >= len(pat) && bytes.Equal(s.mem[s.len-len(b):s.len], b)
}

func (s *String) Split(sep string) *Split {
	return &Split{
		mem: s.payload(),
		idx: 0,
		sep: stringToBytes(sep),
	}
}

func (s *String) SplitSlice(sep string) []String {
	return s.Split(sep).Consume()
}

func (s *String) Find(pat string) int {
	return bytes.Index(s.payload(), stringToBytes(pat))
}

func (s *String) Replace(from, to string) {
	oldsl, newsl := stringToBytes(from), stringToBytes(to)

	payload := s.payload()
	points := make([]int, 0, 8)
	var point, size int
	for {
		index := bytes.Index(payload, oldsl)
		if index < 0 {
			break
		}
		payload = payload[index+len(oldsl):]
		size += len(newsl) - len(oldsl)
		point += index
		points = append(points, point)
		point += len(oldsl)
	}

	// SAFETY: payload drop here
	payload = nil

	if len(points) == 0 {
		return
	}

	if size > 0 {
		s.grow(size)
	}

	block := s.block()
	var offset int
	for _, point = range points {
		loc := point + offset
		copy(block[loc+len(newsl):s.len+len(newsl)-len(oldsl)], block[loc+len(oldsl):s.len])
		copy(block[loc:loc+len(newsl)], newsl)
		offset += len(newsl) - len(oldsl)
		s.len += len(newsl) - len(oldsl)
	}
}

func (s *String) ReplaceToNew(from, to string) String {
	mem := bytes.ReplaceAll(s.payload(), stringToBytes(from), stringToBytes(to))

	return String{
		mem: mem,
		len: len(mem),
		cap: len(mem),
	}
}

// TrimSpaceSlow benchmark: 90.12 ns/op
func (s *String) TrimSpaceSlow() {
	tgt := bytes.TrimSpace(s.payload())
	s.mem = tgt
	s.len = len(tgt)
	s.cap = len(tgt)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// TrimSpace benchmark: 72.55 ns/op
func (s *String) TrimSpace() {
	var start, stop int
	for ; start < s.len; start++ {
		c := s.Get(start)
		if c >= utf8.RuneSelf {
			s.trim(unicode.IsSpace)
			return
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	stop = s.len
	for ; stop > start; stop-- {
		c := s.Get(stop - 1)
		if c >= utf8.RuneSelf {
			s.trim(unicode.IsSpace)
			return
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	if start == stop {
		s.len = 0
		return
	}

	payload := s.payload()
	copy(payload, payload[start:stop])
	s.len = stop - start
}

func (s *String) ParseInt() (int64, error) {
	return strconv.ParseInt(s.UnsafeString(), 10, 64)
}

func (s *String) ParseTo(to FromString) error {
	return to.FromString(s)
}

func (s *String) Reverse() {
	if s.len < 2 {
		return
	}

	payload := s.payload()

	// slower case
	for i := 0; i < s.len; i++ {
		if c := payload[i]; c >= utf8.RuneSelf {
			// NOTE: cloning (*s) is necessary, since changing the memory in (*s) would cause
			// reverse problem, the Runes iterator shares the same memory owned by (*s)
			cl := s.Clone()

			if max := utf8.RuneCount(payload) * utf8.UTFMax; s.cap < max {
				cl.grow(max - s.cap)
			}

			var n int
			for rev := cl.Runes().Reverse(); rev.Next(); {
				n += utf8.EncodeRune(s.mem[n:], rev.Value())
			}

			return
		}
	}

	for i := 0; i < s.len/2; i++ {
		payload[i], payload[s.len-1-i] = payload[s.len-1-i], payload[i]
	}
}

func (s *String) ToUpper() {
	isASCII, hasLower := true, false
	for i := 0; i < s.len; i++ {
		c := s.Get(i)
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasLower = hasLower || ('a' <= c && c <= 'z')
	}

	p := s.payload()
	if isASCII {
		if !hasLower {
			return
		}
		for i := 0; i < s.len; i++ {
			if 'a' <= p[i] && p[i] <= 'z' {
				p[i] -= 'a' - 'A'
			}
		}
		return
	}

	if max := utf8.RuneCount(p) * utf8.UTFMax; s.cap < max {
		s.grow(max - s.cap)
	}

	var n int
	for runes := s.Runes(); runes.Next(); {
		n += utf8.EncodeRune(s.mem[n:], unicode.ToUpper(runes.Value()))
	}

	return
}

func (s *String) ToLower() {
	isASCII, hasUpper := true, false
	for i := 0; i < s.len; i++ {
		c := s.Get(i)
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}

	p := s.payload()
	if isASCII {
		if !hasUpper {
			return
		}
		for i := 0; i < s.len; i++ {
			if 'A' <= p[i] && p[i] <= 'Z' {
				p[i] += 'a' - 'A'
			}
		}
		return
	}

	if max := utf8.RuneCount(p) * utf8.UTFMax; s.cap < max {
		s.grow(max - s.cap)
	}

	var n int
	for runes := s.Runes(); runes.Next(); {
		n += utf8.EncodeRune(s.mem[n:], unicode.ToLower(runes.Value()))
	}

	return
}
