package torguard

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (t *Torguard) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{constants.AES256gcm}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = constants.SHA256
	}

	mssFix := *settings.MSSFix
	if mssFix == 0 {
		const defaultMSSFix = 1450
		mssFix = defaultMSSFix
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Torguard specific
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(mssFix)),
		"sndbuf 393216",
		"rcvbuf 393216",
		"ping 5",
		"remote-cert-tls server",
		"reneg-sec 0",
		"key-direction 1",
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

	if !*settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if connection.Protocol == constants.UDP {
		lines = append(lines, "fast-io")
		lines = append(lines, "explicit-exit-notify")
	}

	if !*settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.TorguardCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSAuth(
		constants.TorguardOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
