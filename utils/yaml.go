package utils

import (
	"github.com/go-yaml/yaml"
)

func Yaml(obj interface{}) string {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}
