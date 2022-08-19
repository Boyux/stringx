package strmut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (s *String) TryFrom(from any) error {
	switch src := from.(type) {
	case bool:
		s.FromString(strconv.FormatBool(src))
	case int, int8, int16, int32, int64:
		s.FromString(strconv.FormatInt(reflect.ValueOf(src).Int(), 10))
	case uint, uint8, uint16, uint32, uint64:
		s.FromString(strconv.FormatUint(reflect.ValueOf(src).Uint(), 10))
	case float32:
		s.FromString(strconv.FormatFloat(reflect.ValueOf(src).Float(), 'g', -1, 32))
	case float64:
		s.FromString(strconv.FormatFloat(src, 'g', -1, 64))
	case string:
		s.FromString(src)
	case []byte:
		s.FromBytes(src)
	case []rune:
		s.FromRunes(src)
	case fmt.Stringer:
		s.FromString(src.String())
	case ToString:
		s.From(src.ToString())
	case Initializer[*String]:
		s.From(src)
	default:
		return fmt.Errorf("string: cannot convert type %s to String", reflect.TypeOf(from).String())
	}
	return nil
}

func (s *String) From(ini Initializer[*String]) *String {
	s.build(nil, 0, 0)
	ini.Initialize(s)
	s.copycheck()
	return s
}

func (s *String) Initialize(target *String) {
	s.CloneInto(target)
}

func (s *String) String() string {
	return s.toString()
}

func (s *String) GoString() string {
	return "\"" + s.toString() + "\""
}

func (s *String) ToString() *String {
	return s.Clone()
}

// Error transform *String as an error, however, this method make IDE like Goland
// unhappy because all errors should be handled, so code like
// 		var s String
// 		s.From(StringInitialier("init"))
// would cause warning messages, so remove it temporary.
// uncomment codes below to enable Error method
// func (s *String) Error() string {
// 	return s.toString()
// }

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
		s.build(mem, len(mem), len(mem))
		return nil
	}

	if b, ok := src.([]byte); ok {
		mem := make([]byte, len(b))
		copy(mem, b)
		s.build(mem, len(mem), len(mem))
		return nil
	}

	return fmt.Errorf("sql: cannot assign type %s to String", reflect.TypeOf(src).String())
}

func (s *String) MarshalJSON() ([]byte, error) {
	return encodeJSON(s.payload()), nil
}

var hex = "0123456789abcdef"

// encodeJSON copy from encoding/json/encode.go encodeState.stringBytes
func encodeJSON(s []byte) []byte {
	var e bytes.Buffer

	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteByte('\\')
			switch b {
			case '\\', '"':
				e.WriteByte(b)
			case '\n':
				e.WriteByte('n')
			case '\r':
				e.WriteByte('r')
			case '\t':
				e.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.WriteString(`u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\u202`)
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.Write(s[start:])
	}
	e.WriteByte('"')

	return e.Bytes()
}

func (s *String) UnmarshalJSON(src []byte) (err error) {
	dst, ok := unquoteBytes(src)
	if !ok {
		return &json.UnmarshalTypeError{
			Value: "string",
			Type:  reflect.TypeOf(s),
		}
	}
	s.build(dst, len(dst), len(dst))
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
