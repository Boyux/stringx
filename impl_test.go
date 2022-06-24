package strmut

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestString_Init(t *testing.T) {
	var s String
	for i := 0; i < 100; i++ {
		str := random(i + 10)
		s.Reset()
		s.From(StringInitializer(str))
		if !s.EqualToString(str) {
			t.Errorf("String: impl From[string]: String=%s string=%s",
				s.String(), str)
		}
	}
}

func TestString_MarshalJSON(t *testing.T) {
	var (
		str string
		s   String
	)
	for i := 0; i < 100; i++ {
		str = random(i + 10)
		s.FromString(str)
		cvt, _ := s.MarshalJSON()
		exp, _ := json.Marshal(str)
		if !bytes.Equal(cvt, exp) {
			t.Errorf("String: impl MarshalJSON: convert=%s expect=%s",
				string(cvt), string(exp))
		}
	}
}

func TestString_MarshalJSON2(t *testing.T) {
	var (
		str string
		s   String
	)
	for i := 0; i < 100; i++ {
		str = random(i + 10)
		s.FromString(str)
		// *String implements MarshalJSON, but String doesn't
		cvt, _ := json.Marshal(&s)
		exp, _ := json.Marshal(str)
		if !bytes.Equal(cvt, exp) {
			t.Errorf("String: impl MarshalJSON: convert=%s expect=%s",
				string(cvt), string(exp))
		}
	}
}
