package strmut

import (
	"strings"
	"testing"
)

func TestList_Join(t *testing.T) {
	for i := 0; i < 100; i++ {
		strslices := make([]string, 10)
		strlist := make([]Str, 10)
		for j := 0; j < 10; j++ {
			str := random(i + 10)
			strslices[j] = str
			strlist[j] = Str(str)
		}
		s := List[Str](strlist).Join("-")
		exp := strings.Join(strslices, "-")
		if !s.EqualToString(exp) {
			t.Errorf("extra: List[*String]: joint=%s expect=%s",
				s.String(), exp)
		}
	}
}
