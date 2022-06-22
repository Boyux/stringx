package strmut

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func testStringRunes(t *testing.T, data string) {
	s := From(data)
	cvt, tgt := s.Runes(), []rune(s.String())
	for i := 0; cvt.Next(); i++ {
		if r := cvt.Value(); r != tgt[i] {
			t.Errorf("Runes: Iterator.Value() = '%v', []rune(string)[%d] = '%v'",
				r, i, tgt[i])
		}
	}
}

var runeData = []string{
	"abc123abc",
	"你好世界",
	random(10),
	random(100),
	random(1000),
}

func TestString_Runes(t *testing.T) {
	for _, data := range runeData {
		testStringRunes(t, data)
	}
}

func TestRuneCount(t *testing.T) {
	for _, data := range runeData {
		s := From(data)
		runes, rev := s.Runes(), s.Runes().Reverse()
		if runes.Size() != rev.Size() {
			t.Errorf("Runes: count is not equal: runes=%d rev=%d",
				runes.Size(), rev.Size())
		}
	}
}

func TestRuneRev(t *testing.T) {
	for _, data := range runeData {
		s := From(data)
		runes, rev := s.Runes().Consume(), s.Runes().Reverse().Consume()
		for i := 0; i < len(runes); i++ {
			if runes[i] != rev[len(rev)-1-i] {
				t.Errorf("Runes: reverse failed: runes[%d]=%d rev[%d]=%d",
					i, runes[i], len(rev)-1-i, rev[len(rev)-1-i])
			}
		}
	}
}

func testStringDrain(t *testing.T, data struct {
	s      string
	r1, r2 int
}) {
	var l, r int
	if data.r1 < data.r2 {
		l, r = data.r1, data.r2
	} else {
		l, r = data.r2, data.r1
	}

	s := From(data.s)
	before := s.Clone()
	s.Drain(l, r)
	expect := data.s[:l] + data.s[r:]
	if !s.EqualToString(expect) {
		t.Errorf("String: drain failed: before=%s after=%s expect=%s",
			before.String(), s.String(), expect)
	}
}

var drainData = []struct {
	s  string
	r1 int
	r2 int
}{
	{random(0), 0, 0},
	{random(10), 0, 0},
	{random(10), rand.Intn(10), rand.Intn(10)},
	{random(10), rand.Intn(10), rand.Intn(10)},
	{random(10), rand.Intn(10), rand.Intn(10)},
	{random(100), rand.Intn(100), rand.Intn(100)},
	{random(100), rand.Intn(100), rand.Intn(100)},
	{random(100), rand.Intn(100), rand.Intn(100)},
}

func TestString_Drain(t *testing.T) {
	for _, data := range drainData {
		testStringDrain(t, data)
	}
}

func testStringReplace(t *testing.T, data []string) {
	str, from, to, exp := From(data[0]), data[1], data[2], data[3]
	before := str.Clone()
	str.Replace(from, to)
	if !str.EqualToString(exp) {
		t.Errorf("String: replacing pattern failed: before=%s after=%s old=%s new=%s expect=%s",
			before.String(), str.String(), from, to, exp)
	}
}

var replaceData = [][]string{
	{"abc123abc", "123", "abc", "abcabcabc"},
	{"abcAAAabcAAA", "AAA", "123", "abc123abc123"},
	{"abcAAAabcAAA", "AAA", "AA", "abcAAabcAA"},
	{"abcAAAabcAAA", "AAA", "AAAA", "abcAAAAabcAAAA"},
	{"你好大世界", "大", "", "你好世界"},
	{"你好大世界", "小", "大", "你好大世界"},
	{"世界真大", "大", "小", "世界真小"},
}

func TestString_Replace(t *testing.T) {
	for _, data := range replaceData {
		testStringReplace(t, data)
	}
}

func BenchmarkString_Replace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range replaceData {
			str, from, to := From(data[0]), data[1], data[2]
			str.Replace(from, to)
		}
	}
}

func BenchmarkString_ReplaceToNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range replaceData {
			str, from, to := From(data[0]), data[1], data[2]
			str = str.ReplaceToNew(from, to)
		}
	}
}

func BenchmarkStdReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range replaceData {
			str, from, to := data[0], data[1], data[2]
			str = strings.ReplaceAll(str, from, to)
		}
	}
}

func testStringTrimSpace(t *testing.T, data []string) {
	str, exp := From(data[0]), data[1]
	before := str.Clone()
	str.TrimSpace()
	if !str.EqualToString(exp) {
		t.Errorf("String: triming space failed: before=%s after=%s expect=%s",
			before.String(), str.String(), exp)
	}
}

func testStringTrimSpaceSlow(t *testing.T, data []string) {
	str, exp := From(data[0]), data[1]
	before := str.Clone()
	str.TrimSpaceSlow()
	if !str.EqualToString(exp) {
		t.Errorf("String: triming space failed: before=%s after=%s expect=%s",
			before.String(), str.String(), exp)
	}
}

var trimSpaceData = [][]string{
	{"\t\n\r     ", ""},
	{"\t aaa \n", "aaa"},
	{"\t 你好世界 \n", "你好世界"},
}

func TestString_TrimSpace(t *testing.T) {
	for _, data := range trimSpaceData {
		testStringTrimSpace(t, data)
	}
}

func TestString_TrimSpaceSlow(t *testing.T) {
	for _, data := range trimSpaceData {
		testStringTrimSpaceSlow(t, data)
	}
}

func BenchmarkString_TrimSpace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range trimSpaceData {
			str := From(data[0])
			str.TrimSpace()
		}
	}
}

func BenchmarkStdTrimSpace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range trimSpaceData {
			str := data[0]
			str = strings.TrimSpace(str)
		}
	}
}

func TestString_ParseInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := strconv.Itoa(i)
		s := From(n)
		cvt, _ := s.ParseInt()
		if cvt != int64(i) {
			t.Errorf("String: parse integer failed: parsed=%d expect=%d",
				cvt, i)
		}
	}
}

type integer int64

func (i *integer) FromString(s *String) error {
	n, err := s.ParseInt()
	if err != nil {
		return err
	}
	*i = integer(n)
	return nil
}

func TestString_ParseTo(t *testing.T) {
	for i := 0; i < 100; i++ {
		var n integer
		s := From(strconv.Itoa(i))
		if err := s.ParseTo(&n); err != nil {
			t.Errorf("String.ParseTo: %s", err.Error())
			return
		}
		if n != integer(i) {
			t.Errorf("String: parse to type failed: parsed=%d expect=%d",
				n, i)
		}
	}
}

func testStringReverse(t *testing.T, data []string) {
	src, tgt := From(data[0]), data[1]
	before := src.Clone()
	src.Reverse()
	if !src.EqualToString(tgt) {
		t.Errorf("String: reverse failed: before=%s after=%s expect=%s",
			before.String(), src.String(), tgt)
	}
}

var reverseData = [][]string{
	{"123456789", "987654321"},
	{"abcdefghi", "ihgfedcba"},
	{"你好", "好你"},
	{"123你好", "好你321"},
	{"1234你好", "好你4321"},
	{"你好世界", "界世好你"},
	{"你好世界👋", "👋界世好你"},
	{"💯", "💯"},
	{"👋💯", "💯👋"},
}

func TestString_Reverse(t *testing.T) {
	for _, data := range reverseData {
		testStringReverse(t, data)
	}
}

func TestString_ToUpper(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := random(10)
		s := From(str)
		before := s.Clone()
		s.ToUpper()
		expect := strings.ToUpper(str)
		if !s.EqualToString(expect) {
			t.Errorf("String: to upper failed: before=%s after=%s expect=%s",
				before.String(), s.String(), expect)
		}
	}
}

func TestString_ToLower(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := random(10)
		s := From(str)
		before := s.Clone()
		s.ToLower()
		expect := strings.ToLower(str)
		if !s.EqualToString(expect) {
			t.Errorf("String: to lower failed: before=%s after=%s expect=%s",
				before.String(), s.String(), expect)
		}
	}
}
