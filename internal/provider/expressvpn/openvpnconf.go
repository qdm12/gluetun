package expressvpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}
	if settings.Auth == "" {
		settings.Auth = constants.SHA512
	}

	if settings.MSSFix == 0 {
		settings.MSSFix = 1200
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(settings.Verbosity),

		// Expressvpn specific
		"fast-io",
		"fragment 1300",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"sndbuf 524288",
		"rcvbuf 524288",
		"verify-x509-name Server name-prefix", // security hole I guess?
		"remote-cert-tls server",              // updated name of ns-cert-type
		"key-direction 1",
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"auth " + settings.Auth,

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		connection.OpenVPNProtoLine(),
		connection.OpenVPNRemoteLine(),
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	if !settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if !settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCert(
		constants.ExpressvpnCert)...)
	lines = append(lines, utils.WrapOpenvpnRSAKey(
		constants.ExpressvpnRSAKey)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.ExpressvpnTLSAuthOpenvpnStaticKeyV1)...)
	lines = append(lines, utils.WrapOpenvpnCA(
		constants.ExpressvpnCA)...)

	lines = append(lines, "")

	return lines, nil
}
