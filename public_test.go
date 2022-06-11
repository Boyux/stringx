package st

import (
	"strings"
	"testing"
)

func testStringReplace(t *testing.T, data []string) {
	str, from, to, exp := From(data[0]), data[1], data[2], data[3]
	before := str.Clone()
	str.Replace(from, to)
	if !str.EqualToString(exp) {
		t.Errorf("replacing String failed: before=%s after=%s old=%s new=%s expected=%s\n",
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
