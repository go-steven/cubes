package utils

import (
	"testing"
)

func TestInArrayStr(t *testing.T) {
	list := []string{
		"aa", "bb", "cc",
	}
	if !InArrayStr("aa", list) {
		t.Error("InArrayStr() failed. Got false, expected true.")
	}
	if InArrayStr("Aa", list) {
		t.Error("InArrayStr() failed. Got true, expected false.")
	}
}
