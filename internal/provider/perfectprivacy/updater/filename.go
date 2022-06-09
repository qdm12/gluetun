package updater

import (
	"strings"
	"unicode"
)

func parseFilename(fileName string) (city string) {
	const suffix = ".conf"
	s := strings.TrimSuffix(fileName, suffix)

	for i, r := range s {
		if unicode.IsDigit(r) {
			s = s[:i]
			break
		}
	}

	return s
}
