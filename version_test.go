package strmut

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	expect := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if !Version.EqualToString(expect) {
		t.Errorf("String: expect version(%s), got version(%s)",
			expect, Version.String())
	}
}
