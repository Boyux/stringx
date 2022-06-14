package st

import (
	"math/rand"
	"strings"
	"testing"
)

var elements = []rune{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '!', '@', '#', '$', '%', '^', '&',
	'*', '(', ')', '{', '}', '[', ']', '\'', '\'', '"', '"', '你', '好', '世', '界', '💰', '🐱',
}

func random(n int) string {
	slice := make([]rune, n)
	for i := 0; i < n/2; i++ {
		slice[i] = elements[rand.Intn(len(elements))]
	}
	for j := n / 2; j < n; j++ {
		slice[j] = elements[len(elements)-6:][rand.Intn(len(elements[len(elements)-6:]))]
	}
	return string(elements)
}

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

func testStringReplace(t *testing.T, data []string) {
	str, from, to, exp := From(data[0]), data[1], data[2], data[3]
	before := str.Clone()
	str.Replace(from, to)
	if !str.EqualToString(exp) {
		t.Errorf("String: replacing pattern failed: before=%s after=%s old=%s new=%s expect=%s\n",
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
		t.Errorf("String: triming space failed: before=%s after=%s expect=%s\n",
			before.String(), str.String(), exp)
	}
}

func testStringTrimSpaceSlow(t *testing.T, data []string) {
	str, exp := From(data[0]), data[1]
	before := str.Clone()
	str.TrimSpaceSlow()
	if !str.EqualToString(exp) {
		t.Errorf("String: triming space failed: before=%s after=%s expect=%s\n",
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
