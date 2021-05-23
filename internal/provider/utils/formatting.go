package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

func commaJoin(slice []string) string {
	return strings.Join(slice, ", ")
}

var ErrNoServerFound = errors.New("no server found")

func NoServerFoundError(selection configuration.ServerSelection) (err error) {
	var messageParts []string

	protocol := constants.UDP
	if selection.TCP {
		protocol = constants.TCP
	}
	messageParts = append(messageParts, "protocol "+protocol)

	if selection.Group != "" {
		part := "group " + selection.Group
		messageParts = append(messageParts, part)
	}

	switch len(selection.Countries) {
	case 0:
	case 1:
		part := "country " + selection.Countries[0]
		messageParts = append(messageParts, part)
	default:
		part := "countries " + commaJoin(selection.Countries)
		messageParts = append(messageParts, part)
	}

	switch len(selection.Regions) {
	case 0:
	case 1:
		part := "region " + selection.Regions[0]
		messageParts = append(messageParts, part)
	default:
		part := "regions " + commaJoin(selection.Regions)
		messageParts = append(messageParts, part)
	}

	switch len(selection.Cities) {
	case 0:
	case 1:
		part := "city " + selection.Cities[0]
		messageParts = append(messageParts, part)
	default:
		part := "cities " + commaJoin(selection.Cities)
		messageParts = append(messageParts, part)
	}

	if selection.Owned {
		messageParts = append(messageParts, "owned servers only")
	}

	switch len(selection.ISPs) {
	case 0:
	case 1:
		part := "ISP " + selection.ISPs[0]
		messageParts = append(messageParts, part)
	default:
		part := "ISPs " + commaJoin(selection.ISPs)
		messageParts = append(messageParts, part)
	}

	switch len(selection.Hostnames) {
	case 0:
	case 1:
		part := "hostname " + selection.Hostnames[0]
		messageParts = append(messageParts, part)
	default:
		part := "hostnames " + commaJoin(selection.Hostnames)
		messageParts = append(messageParts, part)
	}

	switch len(selection.Names) {
	case 0:
	case 1:
		part := "name " + selection.Names[0]
		messageParts = append(messageParts, part)
	default:
		part := "names " + commaJoin(selection.Names)
		messageParts = append(messageParts, part)
	}

	switch len(selection.Numbers) {
	case 0:
	case 1:
		part := "server number " + strconv.Itoa(int(selection.Numbers[0]))
		messageParts = append(messageParts, part)
	default:
		serverNumbers := make([]string, len(selection.Numbers))
		for i := range selection.Numbers {
			serverNumbers[i] = strconv.Itoa(int(selection.Numbers[i]))
		}
		part := "server numbers " + commaJoin(serverNumbers)
		messageParts = append(messageParts, part)
	}

	if selection.EncryptionPreset != "" {
		part := "encryption preset " + selection.EncryptionPreset
		messageParts = append(messageParts, part)
	}

	if selection.FreeOnly {
		messageParts = append(messageParts, "free tier only")
	}

	message := "for " + strings.Join(messageParts, "; ")

	return fmt.Errorf("%w: %s", ErrNoServerFound, message)
}
