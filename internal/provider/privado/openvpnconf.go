package privado

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privado) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}

	if settings.Auth == "" {
		settings.Auth = constants.SHA256
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Privado specific
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		"verify-x509-name " + connection.Hostname + " name",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		connection.OpenVPNProtoLine(),
		connection.OpenVPNRemoteLine(),
		"auth " + settings.Auth,
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if !settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if settings.IPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.PrivadoCertificate)...)

	lines = append(lines, "")

	return lines
}
