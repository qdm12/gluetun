package purevpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Purevpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Purevpn specific
		"key-direction 1",
		"remote-cert-tls server",
		"cipher AES-256-CBC",
		"route-method exe",
		"route-delay 0",
		"script-security 2",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		connection.ProtoLine(),
		connection.RemoteLine(),
		"data-ciphers-fallback " + settings.Cipher,
		"data-ciphers " + settings.Cipher,
	}

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

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
		constants.PurevpnCertificateAuthority)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.PurevpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnKey(
		constants.PurevpnKey)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.PurevpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines
}
