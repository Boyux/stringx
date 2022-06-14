//go:build unsafe_convert
// +build unsafe_convert

package st

import (
	"bytes"
	"testing"
)

var str = "abcdefghijklmnopqrst1234567890你好世界👋"

func TestStringToBytes(t *testing.T) {
	b1, b2 := stringToBytesSlow(str), stringToBytes(str)
	if !bytes.Equal(b1, b2) {
		t.Errorf("unsafe_convert: error converting string to bytes, safe_version=%s unsafe_version=%s\n",
			string(b1), string(b2))
	}
}

func BenchmarkUnsafeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringToBytes(str)
	}
}

func BenchmarkSafeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringToBytesSlow(str)
	}
}

var s = From(str)

func TestBytesToString(t *testing.T) {
	s1, s2 := s.toString(), s.toStringUnsafe()
	if s1 != s2 {
		t.Errorf("unsafe_convert: error converting bytes to string, safe_version=%s unsafe_version=%s\n",
			s1, s2)
	}
}

func BenchmarkToStringUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = s.toStringUnsafe()
	}
}

func BenchmarkSafeToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = s.toString()
	}
}
