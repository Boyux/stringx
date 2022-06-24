package strmut

import (
	"fmt"
)

func Format(format string, args ...any) *String {
	var s String
	s.FromString(fmt.Sprintf(format, args...))
	return &s
}
