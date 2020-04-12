package utils

import (
	"path"
	"strings"
)

// e.g .vscode .a.html.swap
func IsSpecialDorf(dorf string) bool {
	return strings.HasPrefix(path.Base(dorf), ".")
}

// check if an array contains some string item
func ContainsStr(arr []string, key string) bool {
	for _, cur := range arr {
		if cur == key {
			return true
		}
	}
	return false
}
