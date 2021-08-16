package tools

import (
	"strings"
)

func IsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	return str == ""
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}
