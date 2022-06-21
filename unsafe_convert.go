//go:build unsafe_convert

package strmut

import (
	"reflect"
	"unsafe"
)

// stringToBytesSlow benchmark: 2.716 ns/op
func stringToBytesSlow(s string) (b []byte) {
	return []byte(s)
}

// stringToBytes benchmarks:
// BenchmarkUnsafeToBytes-8        1000000000               0.3147 ns/op
// BenchmarkSafeToBytes-8          443243308                2.716 ns/op
//
// SAFETY: byte slice converted by this function is immutable, don't mutate
// those bytes, keep readonly in mind
func stringToBytes(s string) (b []byte) {
	st := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sl.Data = st.Data
	sl.Len = st.Len
	sl.Cap = st.Len
	return b
}

// toStringUnsafe benchmark: 0.6994 ns/op
func (s *String) toStringUnsafe() (dst string) {
	src := s.payload()
	st := (*reflect.StringHeader)(unsafe.Pointer(&dst))
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	st.Data = sl.Data
	st.Len = s.len
	return dst
}

// UnsafeString is a faster way to convert String to primitive string by unsafe.Pointer,
// it takes no extra cost but may cause memory issue if caller use UnsafeString incorrectly
func (s *String) UnsafeString() string {
	return s.toStringUnsafe()
}

// toString benchmark: 3.191 ns/op
func (s *String) toString() string {
	if s.len == 0 {
		return ""
	}

	return string(s.payload())
}
