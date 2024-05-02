package utils

import (
	"strings"
)

func CipherLines(ciphers []string) (lines []string) {
	if len(ciphers) == 0 {
		return nil
	}

	return []string{
		"cipher " + ciphers[0],
		"ncp-ciphers " + strings.Join(ciphers, ":"),
	}
}
