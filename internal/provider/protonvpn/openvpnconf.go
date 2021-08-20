package protonvpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Protonvpn) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}

	if settings.Auth == "" {
		settings.Auth = constants.SHA512
	}

	const defaultMSSFix = 1450
	if settings.MSSFix == 0 {
		settings.MSSFix = defaultMSSFix
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"tls-exit",

		// Protonvpn specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"reneg-sec 0",
		"key-direction 1",
		"pull",
		"comp-lzo no",

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

	if connection.Protocol == constants.UDP {
		lines = append(lines, "fast-io")
	}

	if !settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
	}

	if settings.IPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.ProtonvpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.ProtonvpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines
}
