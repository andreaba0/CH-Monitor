package utils

import "unicode"

func IsUnicodeDigit(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
