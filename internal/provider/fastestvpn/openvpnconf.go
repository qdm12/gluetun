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
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(settings.Verbosity),

		// Fastestvpn specific
		"mssfix " + strconv.Itoa(int(settings.MSSFix)), // defaults to 1450
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		"key-direction 1",
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"auth " + settings.Auth,
		"comp-lzo",
		"reneg-sec 0",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		// "pull-filter ignore \"auth-token\"", // needed for FastestVPN
		"auth-retry nointeract",
		"suppress-timestamps",

		// Connection variables
		connection.OpenVPNProtoLine(),
		connection.OpenVPNRemoteLine(),
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
		lines = append(lines, "tun-mtu 1500")     // FastestVPN specific
		lines = append(lines, "tun-mtu-extra 32") // FastestVPN specific
		lines = append(lines, "ping 15")          // FastestVPN specific
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

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.FastestvpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.FastestvpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
