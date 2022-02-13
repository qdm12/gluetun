package privado

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privado) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{constants.AES256cbc}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = constants.SHA256
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Privado specific
		"ping 10",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		"verify-x509-name " + connection.Hostname + " name",
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"auth " + auth,

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Connection variables
		connection.OpenVPNProtoLine(),
		connection.OpenVPNRemoteLine(),
	}

	lines = append(lines, utils.CipherLines(settings.Ciphers, settings.Version)...)

	if settings.ProcessUser != "root" {
		lines = append(lines, "user "+settings.ProcessUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if *settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(*settings.MSSFix)))
	}

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	if !*settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.PrivadoCA)...)

	lines = append(lines, "")

	return lines, nil
}
