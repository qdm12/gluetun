package fastestvpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (f *Fastestvpn) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}
	if settings.Auth == "" {
		settings.Auth = constants.SHA256
	}
	if settings.MSSFix == 0 {
		settings.MSSFix = 1450
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"ping 15",
		"ping-exit 60",
		"tls-exit",

		// Fastestvpn specific
		"ping-restart 0",
		"tls-client",
		"tls-cipher  TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		"key-direction 1",
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)), // defaults to 1450

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
		lines = append(lines, "persist-tun")
	}

	if !settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.FastestvpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.FastestvpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
