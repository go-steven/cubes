package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Sha1(s string) string {
	t := sha1.New()
	io.WriteString(t, s)
	return fmt.Sprintf("%x", t.Sum(nil))
}
