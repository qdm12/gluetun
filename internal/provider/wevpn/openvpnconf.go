package wevpn

import (
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Wevpn) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256gcm
	}

	if settings.Auth == "" {
		settings.Auth = constants.SHA512
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"	-exit",

		// Wevpn specific
		"redirect-gateway def1 bypass-dhcp",
		"route-delay 0",
		"reneg-sec 0",

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

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if strings.HasSuffix(settings.Cipher, "-gcm") {
		lines = append(lines, "ncp-ciphers AES-256-GCM:AES-256-CBC:AES-128-GCM")
	}

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

	lines = append(lines, utils.WrapOpenvpnKey(
		settings.ClientKey)...)
	lines = append(lines, utils.WrapOpenvpnCA(
		constants.WevpnCA)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.WevpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.WevpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines
}
