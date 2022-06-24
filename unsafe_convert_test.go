//go:build unsafe_convert

package strmut

import (
	"bytes"
	"testing"
)

var unsafeDataStr = "abcdefghijklmnopqrst1234567890你好世界👋"

func TestStringToBytes(t *testing.T) {
	b1, b2 := stringToBytesSlow(unsafeDataStr), stringToBytes(unsafeDataStr)
	if !bytes.Equal(b1, b2) {
		t.Errorf("unsafe_convert: error converting string to bytes, safe_version=%s unsafe_version=%s",
			string(b1), string(b2))
	}
}

func BenchmarkUnsafeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringToBytes(unsafeDataStr)
	}
}

func BenchmarkSafeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringToBytesSlow(unsafeDataStr)
	}
}

var unsafeDataS String

func init() {
	unsafeDataS.FromString(unsafeDataStr)
}

func TestBytesToString(t *testing.T) {
	s1, s2 := unsafeDataS.toString(), unsafeDataS.toStringUnsafe()
	if s1 != s2 {
		t.Errorf("unsafe_convert: error converting bytes to string, safe_version=%s unsafe_version=%s",
			s1, s2)
	}
}

func BenchmarkToStringUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = unsafeDataS.toStringUnsafe()
	}
}

func BenchmarkSafeToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = unsafeDataS.toString()
	}
}
