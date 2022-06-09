package updater

import "strings"

func codeToCountry(countryCode string, countryCodes map[string]string) (
	country string, warning string) {
	countryCode = strings.ToLower(countryCode)
	country, ok := countryCodes[countryCode]
	if !ok {
		warning = "unknown country code: " + countryCode
		country = countryCode
	}
	return country, warning
}
