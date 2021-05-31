package utils

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func CipherLines(cipher, version string) (lines []string) {
	switch version {
	case constants.Openvpn24:
		return []string{"cipher " + cipher}
	default: // 2.5 and above
		return []string{
			"data-ciphers-fallback " + cipher,
			"data-ciphers " + cipher,
		}
	}
}
