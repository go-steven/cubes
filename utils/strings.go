package utils

import (
	"sort"
	"strings"
)

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func LowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func UpperTrim(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func InArrayStr(s string, list []string) bool {
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	index := sort.SearchStrings(list, s)
	return index >= 0 && index < len(list) && list[index] == s
}
