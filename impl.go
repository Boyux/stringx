package st

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (s *String) Write(p []byte) (n int, err error) {
	s.PushBytes(p)
	return len(p), nil
}

func (s *String) Scan(src any) error {
	if src == nil {
		s.len = 0
		return nil
	}

	if str, ok := src.(string); ok {
		mem := stringToBytesSlow(str)
		s.mem = mem
		s.len = len(mem)
		s.cap = len(mem)
		return nil
	}

	if bytes, ok := src.([]byte); ok {
		mem := make([]byte, len(bytes))
		copy(mem, bytes)
		s.mem = mem
		s.len = len(mem)
		s.cap = len(mem)
		return nil
	}

	return fmt.Errorf("sql: cannot assign type %s to String", reflect.TypeOf(src).String())
}

func (s *String) MarshalJSON() ([]byte, error) {
	dst := make([]byte, s.len+2)
	copy(dst[1:len(dst)-1], s.payload())
	dst[0] = '"'
	dst[len(dst)-1] = '"'
	return dst, nil
}

func (s *String) UnmarshalJSON(src []byte) (err error) {
	dst, ok := unquoteBytes(src)
	if !ok {
		return &json.UnmarshalTypeError{
			Value: "string",
			Type:  reflect.TypeOf(s),
		}
	}
	s.mem = dst
	s.len = len(dst)
	s.cap = len(dst)
	return nil
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
// copy from encoding/json/decode.go
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

// copy from encoding/json/decode.go
func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8, and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}
