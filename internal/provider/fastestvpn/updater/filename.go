package updater

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var errFilenameNoProtocolSuffix = errors.New("filename does not have a protocol suffix")

var trailNumberExp = regexp.MustCompile(`[0-9]+$`)

func parseFilename(fileName string) (
	country string, tcp, udp bool, err error,
) {
	const (
		tcpSuffix = "-tcp.ovpn"
		udpSuffix = "-udp.ovpn"
	)
	var suffix string
	switch {
	case strings.HasSuffix(strings.ToLower(fileName), tcpSuffix):
		suffix = tcpSuffix
		tcp = true
	case strings.HasSuffix(strings.ToLower(fileName), udpSuffix):
		suffix = udpSuffix
		udp = true
	default:
		return "", false, false, fmt.Errorf("%w: %s",
			errFilenameNoProtocolSuffix, fileName)
	}

	countryWithNumber := strings.TrimSuffix(fileName, suffix)
	number := trailNumberExp.FindString(countryWithNumber)
	country = countryWithNumber[:len(countryWithNumber)-len(number)]

	return country, tcp, udp, nil
}
