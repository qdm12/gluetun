package privatevpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privatevpn) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES128gcm
	}

	if settings.Auth == "" {
		settings.Auth = constants.SHA256
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"tls-exit",

		// Privatevpn specific

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
		lines = append(lines, "key-direction 1")
	}

	if !settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
		lines = append(lines, "persist-tun")
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if !settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.PrivatevpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSCrypt(
		constants.PrivatevpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
