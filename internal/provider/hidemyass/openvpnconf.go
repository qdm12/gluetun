package hidemyass

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (h *HideMyAss) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if settings.Cipher == "" {
		settings.Cipher = constants.AES256cbc
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping 5",
		"ping-exit 30",
		"ping-timer-rem",
		"tls-exit",

		// HideMyAss specific
		"remote-cert-tls server", // updated name of ns-cert-type
		// "route-metric 1",
		"comp-lzo yes",
		"comp-noadapt",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"proto " + connection.Protocol,
		connection.RemoteLine(),
	}

	lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)

	if settings.Auth != "" {
		lines = append(lines, "auth "+settings.Auth)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if !settings.Root {
		lines = append(lines, "user "+username)
	}

	if settings.Provider.ExtraConfigOptions.OpenVPNIPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.HideMyAssCA)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.HideMyAssCertificate)...)
	lines = append(lines, utils.WrapOpenvpnRSAKey(
		constants.HideMyAssRSAPrivateKey)...)

	lines = append(lines, "")

	return lines
}
