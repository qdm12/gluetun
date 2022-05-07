package utils

import (
	"strings"

	"github.com/qdm12/gluetun/internal/constants/openvpn"
)

func CipherLines(ciphers []string, version string) (lines []string) {
	if len(ciphers) == 0 {
		return nil
	}

	switch version {
	case openvpn.Openvpn24:
		return []string{
			"cipher " + ciphers[0],
			"ncp-ciphers " + strings.Join(ciphers, ":"),
		}
	default: // 2.5 and above
		return []string{
			"data-ciphers-fallback " + ciphers[0],
			"data-ciphers " + strings.Join(ciphers, ":"),
		}
	}
}
