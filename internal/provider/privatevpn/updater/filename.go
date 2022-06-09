package updater

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	trailingNumber = regexp.MustCompile(` [0-9]+$`)
)

var (
	errBadPrefix      = errors.New("bad prefix in file name")
	errBadSuffix      = errors.New("bad suffix in file name")
	errNotEnoughParts = errors.New("not enough parts in file name")
)

func parseFilename(fileName string) (
	countryCode, city string, err error,
) {
	fileName = strings.ReplaceAll(fileName, " ", "") // remove spaces

	const prefix = "PrivateVPN-"
	if !strings.HasPrefix(fileName, prefix) {
		return "", "", fmt.Errorf("%w: %s", errBadPrefix, fileName)
	}
	s := strings.TrimPrefix(fileName, prefix)

	const tcpSuffix = "-TUN-443.ovpn"
	const udpSuffix = "-TUN-1194.ovpn"
	switch {
	case strings.HasSuffix(fileName, tcpSuffix):
		s = strings.TrimSuffix(s, tcpSuffix)
	case strings.HasSuffix(fileName, udpSuffix):
		s = strings.TrimSuffix(s, udpSuffix)
	default:
		return "", "", fmt.Errorf("%w: %s", errBadSuffix, fileName)
	}

	s = trailingNumber.ReplaceAllString(s, "")

	parts := strings.Split(s, "-")
	const minParts = 2
	if len(parts) < minParts {
		return "", "", fmt.Errorf("%w: %s", errNotEnoughParts, fileName)
	}
	countryCode, city = parts[0], parts[1]

	countryCode = strings.ToLower(countryCode)
	if countryCode == "co" && strings.HasPrefix(city, "Bogot") {
		city = "Bogota"
	}

	return countryCode, city, nil
}
