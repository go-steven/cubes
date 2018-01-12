package utils

import (
	"testing"
)

func TestJson(t *testing.T) {
	obj := []string{
		"aa", "bb", "cc",
	}
	ret := Json(obj)
	expect := `["aa", "bb", "cc"]`
	if ret != expect {
		t.Errorf("Json failed. Got %s, expected %s.", ret, expect)
	}
}
