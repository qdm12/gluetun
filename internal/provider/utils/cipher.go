package utils

import (
	"strings"
)

func CipherLines(ciphers []string) (lines []string) {
	if len(ciphers) == 0 {
		return nil
	}

	return []string{
		"data-ciphers-fallback " + ciphers[0],
		"data-ciphers " + strings.Join(ciphers, ":"),
	}
}
