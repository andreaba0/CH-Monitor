package utils

import (
	"crypto/rand"
	"encoding/hex"
	"unicode"
)

func IsUnicodeDigit(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func RandomString(length int) (string, error) {
	bytes := length / 2
	if length%2 != 0 {
		bytes += 1
	}

	bytesArr := make([]byte, bytes)
	if _, err := rand.Read(bytesArr); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytesArr)[:length], nil
}
