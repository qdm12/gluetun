package purevpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Purevpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{constants.AES256gcm}
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Purevpn specific
		"ping 10",
		"remote-cert-tls server",
		"key-direction 1",
		"auth-user-pass " + constants.OpenVPNAuthConf,

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

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	if *settings.Auth != "" {
		lines = append(lines, "auth "+*settings.Auth)
	}

	if *settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(*settings.MSSFix)))
	}

	if settings.ProcessUser != "root" {
		lines = append(lines, "user "+settings.ProcessUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if !*settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.PurevpnCA)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.PurevpnCert)...)
	lines = append(lines, utils.WrapOpenvpnKey(
		constants.PurevpnKey)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.PurevpnTLSAuth)...)

	lines = append(lines, "")

	return lines, nil
}
