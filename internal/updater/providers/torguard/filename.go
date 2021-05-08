package torguard

import "strings"

func parseFilename(fileName string) (country, city string) {
	const prefix = "TorGuard."
	const suffix = ".ovpn"
	s := strings.TrimPrefix(fileName, prefix)
	s = strings.TrimSuffix(s, suffix)

	switch {
	case strings.Count(s, ".") == 1 && !strings.HasPrefix(s, "USA"):
		parts := strings.Split(s, ".")
		country = parts[0]
		city = parts[1]

	case strings.HasPrefix(s, "USA"):
		country = "USA"
		s = strings.TrimPrefix(s, "USA-")
		s = strings.ReplaceAll(s, "-", " ")
		s = strings.ReplaceAll(s, ".", " ")
		s = strings.ToLower(s)
		s = strings.Title(s)
		city = s

	default:
		country = s
	}

	return country, city
}
