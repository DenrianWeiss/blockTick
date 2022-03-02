package utils

import "strings"

func AddPrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return "0x" + s
}

func CheckInputAddress(address string) bool {
	ia := AddPrefix(address)
	if len(ia) != 42 {
		return false
	}
	return true
}
