package utils

import "strings"

func CipherLines(cipher, version string) (lines []string) {
	switch {
	case strings.HasPrefix(version, "2.4"):
		return []string{"cipher " + cipher}
	default: // 2.5 and above
		return []string{
			"data-ciphers-fallback " + cipher,
			"data-ciphers " + cipher,
		}
	}
}
