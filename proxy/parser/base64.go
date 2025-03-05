package parser

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"
)

func IsBase64String(s string) bool {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return false
	}

	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}

	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	sWithoutPadding := strings.TrimRight(s, "=")
	for _, c := range sWithoutPadding {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}

	return utf8.Valid(decoded)
}

func DecodeBase64(s string) string {
	if !IsBase64String(s) {
		return s
	}

	s = strings.TrimSpace(s)
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}

	if !utf8.Valid(decoded) {
		return s
	}

	return string(decoded)
}
