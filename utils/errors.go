package utils

import (
	"errors"
	"fmt"
)

func Error(err string) error {
	return errors.New(err)
}

// New returns an error that formats as the given text.
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
