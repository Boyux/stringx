package stringx

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func testStringRunes(t *testing.T, data string) {
	var s String
	s.FromString(data)
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
	var s String
	for _, data := range runeData {
		s.FromString(data)
		runes, rev := s.Runes(), s.Runes().Reverse()
		if runes.Size() != rev.Size() {
			t.Errorf("Runes: count is not equal: runes=%d rev=%d",
				runes.Size(), rev.Size())
		}
	}
}

func TestRuneRev(t *testing.T) {
	var s String
	for _, data := range runeData {
		s.FromString(data)
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

	var s String
	s.FromString(data.s)
	var before String
	s.CloneInto(&before)
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
	var str String
	str.FromString(data[0])
	from, to, exp := data[1], data[2], data[3]
	var before String
	str.CloneInto(&before)
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
	var str String
	for i := 0; i < b.N; i++ {
		for _, data := range replaceData {
			str.FromString(data[0])
			from, to := data[1], data[2]
			str.Replace(from, to)
		}
	}
}

func BenchmarkString_ReplaceToNew(b *testing.B) {
	str := New()
	for i := 0; i < b.N; i++ {
		for _, data := range replaceData {
			str.FromString(data[0])
			from, to := data[1], data[2]
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

var trimPrefixData = [][]string{
	{"abc123", "abc"},
	{"🐱🐶", "🐱"},
}

func TestString_TrimPrefix(t *testing.T) {
	for _, data := range trimPrefixData {
		after, exp := New().
			FromString(data[0]).
			TrimPrefix(data[1]), strings.TrimPrefix(data[0], data[1])
		if !after.EqualToString(exp) {
			t.Errorf("String: triming prefix failed: before=%s after=%s expect=%s",
				data[0], after.UnsafeString(), exp)
		}
	}
}

func BenchmarkString_TrimPrefix(b *testing.B) {
	str := New()
	for i := 0; i < b.N; i++ {
		for _, data := range trimPrefixData {
			str.
				FromString(data[0]).
				TrimPrefix(data[1])
		}
	}
}

func BenchmarkStdTrimPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range trimPrefixData {
			strings.TrimPrefix(data[0], data[1])
		}
	}
}

var trimSuffixData = [][]string{
	{"abc123", "123"},
	{"🐶🐱", "🐱"},
}

func BenchmarkString_TrimSuffix(b *testing.B) {
	str := New()
	for i := 0; i < b.N; i++ {
		for _, data := range trimSuffixData {
			str.
				FromString(data[0]).
				TrimSuffix(data[1])
		}
	}
}

func BenchmarkStdTrimSuffix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, data := range trimSuffixData {
			strings.TrimSuffix(data[0], data[1])
		}
	}
}

func TestString_TrimSuffix(t *testing.T) {
	for _, data := range trimSuffixData {
		after, exp := New().
			FromString(data[0]).
			TrimSuffix(data[1]), strings.TrimSuffix(data[0], data[1])
		if !after.EqualToString(exp) {
			t.Errorf("String: triming suffix failed: before=%s after=%s expect=%s",
				data[0], after.UnsafeString(), exp)
		}
	}
}

func testStringTrimSpace(t *testing.T, data []string) {
	var str String
	str.FromString(data[0])
	exp := data[1]
	var before String
	str.CloneInto(&before)
	str.TrimSpace()
	if !str.EqualToString(exp) {
		t.Errorf("String: triming space failed: before=%s after=%s expect=%s",
			before.String(), str.String(), exp)
	}
}

func testStringTrimSpaceSlow(t *testing.T, data []string) {
	var str String
	str.FromString(data[0])
	exp := data[1]
	var before String
	str.CloneInto(&before)
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
	var str String
	for i := 0; i < b.N; i++ {
		for _, data := range trimSpaceData {
			str.FromString(data[0])
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
	var s String
	for i := 0; i < 100; i++ {
		n := strconv.Itoa(i)
		s.FromString(n)
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
	var s String
	for i := 0; i < 100; i++ {
		var n integer
		s.FromString(strconv.Itoa(i))
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
	var src String
	src.FromString(data[0])
	tgt := data[1]
	var before String
	src.CloneInto(&before)
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
	var s String
	for i := 0; i < 100; i++ {
		str := random(10)
		s.FromString(str)
		var before String
		s.CloneInto(&before)
		s.ToUpper()
		expect := strings.ToUpper(str)
		if !s.EqualToString(expect) {
			t.Errorf("String: to upper failed: before=%s after=%s expect=%s",
				before.String(), s.String(), expect)
		}
	}
}

func TestString_ToLower(t *testing.T) {
	var s String
	for i := 0; i < 100; i++ {
		str := random(10)
		s.FromString(str)
		before := s.Clone()
		s.ToLower()
		expect := strings.ToLower(str)
		if !s.EqualToString(expect) {
			t.Errorf("String: to lower failed: before=%s after=%s expect=%s",
				before.String(), s.String(), expect)
		}
	}
}

func BenchmarkString_ToUpper(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			s.FromString(random(10)).
				ToUpper()
		}
	}
}

func BenchmarkStdToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			strings.ToUpper(random(10))
		}
	}
}

func BenchmarkString_ToLower(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			s.FromString(random(10)).
				ToLower()
		}
	}
}

func BenchmarkStdToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			strings.ToLower(random(10))
		}
	}
}

func BenchmarkString_PushString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := New()
		s.SetCapacity(100 * 10)
		for j := 0; j < 100; j++ {
			s.PushString(random(10))
		}
	}
}

func BenchmarkStdPushString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += random(10)
		}
	}
}

func BenchmarkStdStringBuilderPushString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := new(strings.Builder)
		for j := 0; j < 100; j++ {
			s.WriteString(random(10))
		}
	}
}
