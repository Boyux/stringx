//go:build unsafe_convert
// +build unsafe_convert

package st

import (
	"reflect"
	"runtime"
	"unsafe"
)

// FIXME
// segment fault
func stringToBytes(s string) (b []byte) {
	st := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sl.Data = st.Data
	sl.Len = st.Len
	sl.Cap = st.Len
	runtime.KeepAlive(s)
	return b
}

func (s *String) toString() (dst string) {
	src := make([]byte, s.len)
	copy(src, s.payload())
	st := (*reflect.StringHeader)(unsafe.Pointer(&dst))
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	st.Data = sl.Data
	st.Len = s.len
	runtime.KeepAlive(src)
	return dst
}
