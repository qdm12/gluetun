package updater

import (
	"errors"
	"fmt"
	"strings"
)

var errNotOvpnExt = errors.New("filename does not have the openvpn file extension")

func parseFilename(fileName string) (
	region string, err error,
) {
	const suffix = ".ovpn"
	if !strings.HasSuffix(fileName, suffix) {
		return "", fmt.Errorf("%w: %s", errNotOvpnExt, fileName)
	}

	region = strings.TrimSuffix(fileName, suffix)
	region = strings.ReplaceAll(region, " - ", " ")
	return region, nil
}
