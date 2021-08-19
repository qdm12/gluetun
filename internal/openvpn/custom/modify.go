package custom

import (
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func modifyCustomConfig(lines []string, settings configuration.OpenVPN,
	connection models.Connection) (modified []string) {
	// Remove some lines
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "up "),
			strings.HasPrefix(line, "down "),
			strings.HasPrefix(line, "verb "),
			strings.HasPrefix(line, "auth-user-pass "),
			strings.HasPrefix(line, "user "),
			strings.HasPrefix(line, "proto "),
			strings.HasPrefix(line, "remote "),
			settings.Cipher != "" && strings.HasPrefix(line, "cipher "),
			settings.Cipher != "" && strings.HasPrefix(line, "data-ciphers "),
			settings.Auth != "" && strings.HasPrefix(line, "auth "),
			settings.MSSFix > 0 && strings.HasPrefix(line, "mssfix "),
			!settings.IPv6 && strings.HasPrefix(line, "tun-ipv6"):
		default:
			modified = append(modified, line)
		}
	}

	// Add values
	modified = append(modified, connection.OpenVPNProtoLine())
	modified = append(modified, connection.OpenVPNRemoteLine())
	modified = append(modified, "mute-replay-warnings")
	modified = append(modified, "auth-nocache")
	modified = append(modified, "pull-filter ignore \"auth-token\"") // prevent auth failed loop
	modified = append(modified, "auth-retry nointeract")
	modified = append(modified, "suppress-timestamps")
	if settings.User != "" {
		modified = append(modified, "auth-user-pass "+constants.OpenVPNAuthConf)
	}
	modified = append(modified, "verb "+strconv.Itoa(settings.Verbosity))
	if settings.Cipher != "" {
		modified = append(modified, utils.CipherLines(settings.Cipher, settings.Version)...)
	}
	if settings.Auth != "" {
		modified = append(modified, "auth "+settings.Auth)
	}
	if settings.MSSFix > 0 {
		modified = append(modified, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}
	if !settings.IPv6 {
		modified = append(modified, `pull-filter ignore "route-ipv6"`)
		modified = append(modified, `pull-filter ignore "ifconfig-ipv6"`)
	}
	if !settings.Root {
		modified = append(modified, "user "+settings.ProcUser)
	}

	return modified
}
