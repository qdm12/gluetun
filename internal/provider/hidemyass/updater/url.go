package updater

import (
	"strings"
	"unicode"
)

func parseOpenvpnURL(url, protocol string) (country, region, city string) {
	lastSlashIndex := strings.LastIndex(url, "/")
	url = url[lastSlashIndex+1:]

	suffix := "." + strings.ToUpper(protocol) + ".ovpn"
	url = strings.TrimSuffix(url, suffix)

	parts := strings.Split(url, ".")

	switch len(parts) {
	case 1:
		country = parts[0]
		return country, "", ""
	case 2: //nolint:gomnd
		country = parts[0]
		city = parts[1]
	default:
		country = parts[0]
		region = parts[1]
		city = parts[2]
	}

	country = camelCaseToWords(country)
	region = camelCaseToWords(region)
	city = camelCaseToWords(city)

	country = mutateSpecialCountryCases(country)

	return country, region, city
}

func camelCaseToWords(camelCase string) (words string) {
	wasLowerCase := false
	for _, r := range camelCase {
		if wasLowerCase && unicode.IsUpper(r) {
			words += " "
		}
		wasLowerCase = unicode.IsLower(r)
		words += string(r)
	}
	return words
}

func mutateSpecialCountryCases(country string) string {
	switch country {
	case "Coted`Ivoire":
		return "Cote d'Ivoire"
	default:
		return country
	}
}
