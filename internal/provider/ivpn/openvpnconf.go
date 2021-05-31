package ivpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ivpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping 5",
		"ping-exit 30",
		"ping-timer-rem",
		"tls-exit",

		// IVPN specific
		"remote-cert-tls server", // updated name of ns-cert-type
		"comp-lzo no",
		"key-direction 1",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",

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
		"verify-x509-name " + connection.Hostname, //  + " name-prefix"
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if settings.Auth != "" {
		lines = append(lines, "auth "+settings.Auth)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if !settings.Root {
		lines = append(lines, "user "+username)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.IvpnCA)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.IvpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines
}
