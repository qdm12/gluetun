package ipvanish

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ipvanish) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}
	if settings.Auth == "" {
		settings.Auth = constants.SHA256
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping-timer-rem",
		"tls-exit",

		// Ipvanish specific
		"comp-lzo",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		"keysize 256",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"proto " + connection.Protocol,
		connection.RemoteLine(),
		"verify-x509-name " + connection.Hostname, //  + " name"
		"auth " + settings.Auth,
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if !settings.Root {
		lines = append(lines, "user "+username)
	}

	lines = append(lines, utils.WrapOpenvpnCA(constants.IpvanishCA)...)

	lines = append(lines, "")

	return lines
}
