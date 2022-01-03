package perfectprivacy

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Perfectprivacy) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{constants.AES256cbc, constants.AES256gcm}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = constants.SHA512
	}

	mssFix := *settings.MSSFix
	if mssFix == 0 {
		mssFix = 1450
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Perfect Privacy specific
		"ping 5",
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(mssFix)),
		"reneg-sec 3600",
		"key-direction 1",
		"tls-cipher TLS_CHACHA20_POLY1305_SHA256:TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-128-GCM-SHA256:TLS-DHE-RSA-WITH-AES-128-CBC-SHA:TLS_AES_256_GCM_SHA384:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"auth " + auth,

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		connection.OpenVPNProtoLine(),
		connection.OpenVPNRemoteLine(),
	}

	lines = append(lines, utils.CipherLines(settings.Ciphers, settings.Version)...)

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	if !*settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if !*settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
		// Perfect Privacy specific IPv6
		lines = append(lines, "redirect-gateway def1")
		lines = append(lines, `pull-filter ignore "redirect-gateway def1 ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.PerfectprivacyCA)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.PerfectprivacyCert)...)
	lines = append(lines, utils.WrapOpenvpnKey(
		constants.PerfectprivacyPrivateKey)...)
	lines = append(lines, utils.WrapOpenvpnTLSCrypt(
		constants.PerfectprivacyTLSCryptOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
