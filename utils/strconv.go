package utils

import (
	"strconv"
)

func ParseUint(s string, base int, bitSize int) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseUint(s, 10, 64)
}
