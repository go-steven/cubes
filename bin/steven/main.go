package main

import (
	"strings"
	"fmt"
)

func trim_zero(v string) string {
	vals := strings.Split(v, ".")
	if len(vals) < 2 {
		return v
	}
	println(vals[0])
	println(vals[1])
	
	str := vals[1]
	i := len(str) - 1
	for ; i >= 0; i-- {
		if string(str[i]) != "0" {
			break
		}
	}

	if i == 0 {
		return vals[0]
	} else {
		s := ""
		for j := 0; j < i+1; j++ {
			s += string(str[j])
		}
		return fmt.Sprintf("%s.%s", vals[0], s)
	}
}

func main() {
	s := "13418.000000"
	println(trim_zero(s))
}
