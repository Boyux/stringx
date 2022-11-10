package stringx

import (
	"bufio"
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

func TestString_Lines(t *testing.T) {
	var s String
	for i := 0; i < 10; i++ {
		str := random(400) + "\r\n" + random(400) + "\r\n"
		s.FromString(str)
		lines, scanner := s.Lines(), bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			if !lines.Next() {
				t.Errorf("extra: Iterator[*String]: Lines: scanner has next token but lines has no value")
			}

			line := lines.Value().UnsafeString()
			expect := scanner.Text()
			if line != expect {
				t.Errorf("extra: Iterator[*String]: Lines: line = %s expect = %s",
					line, expect)
			}
		}
	}
}
