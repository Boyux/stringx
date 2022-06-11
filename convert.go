//go:build !unsafe_convert
// +build !unsafe_convert

package st

func stringToBytes(s string) []byte {
	return []byte(s)
}

func (s *String) toString() string {
	if s.len == 0 {
		return ""
	}

	return string(s.payload())
}
