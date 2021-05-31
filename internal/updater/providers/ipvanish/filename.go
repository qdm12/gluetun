package ipvanish

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

var errCountryCodeUnknown = errors.New("country code is unknown")

func parseFilename(fileName, hostname string) (
	country, city string, err error) {
	const prefix = "ipvanish-"
	s := strings.TrimPrefix(fileName, prefix)

	const ext = ".ovpn"
	host := strings.Split(hostname, ".")[0]
	suffix := "-" + host + ext
	s = strings.TrimSuffix(s, suffix)

	parts := strings.Split(s, "-")

	countryCodes := constants.CountryCodes()
	countryCode := strings.ToLower(parts[0])
	country, ok := countryCodes[countryCode]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", errCountryCodeUnknown, countryCode)
	}
	country = strings.Title(country)

	if len(parts) > 1 {
		city = strings.Join(parts[1:], " ")
		city = strings.Title(city)
	}

	return country, city, nil
}
