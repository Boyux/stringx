//go:build !unsafe_convert

package strmut

func stringToBytesSlow(s string) (b []byte) {
	return []byte(s)
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

func (s *String) toString() string {
	if s.len == 0 {
		return ""
	}

	return string(s.payload())
}

// UnsafeString in 'convert.go' is safe, it is just for preventing compile issue
// while disabling 'unsafe_convert' tag
func (s *String) UnsafeString() string {
	return s.toString()
}
