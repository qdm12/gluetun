package ivpn

import (
	"strings"
)

func parseFilename(fileName string) (country, city string) {
	const suffix = ".ovpn"
	fileName = strings.TrimSuffix(fileName, suffix)
	parts := strings.Split(fileName, "-")
	country = strings.ReplaceAll(parts[0], "_", " ")
	if len(parts) > 1 {
		city = strings.ReplaceAll(parts[1], "_", " ")
	}
	return country, city
}
