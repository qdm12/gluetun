package vpnunlimited

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"tls-exit",
		"remote-cert-tls server",

		// VPNUnlimited specific
		"reneg-sec 0",
		"ping 5",
		"ping-exit 30",
		"comp-lzo no",
		"route-metric 1",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		// "auth-user-pass " + constants.OpenVPNAuthConf,
		connection.ProtoLine(),
		connection.RemoteLine(),
	}

	if settings.Cipher != "" {
		lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)
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
		constants.VPNUnlimitedCertificateAuthority)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		settings.Provider.ExtraConfigOptions.ClientCertificate)...)
	lines = append(lines, utils.WrapOpenvpnKey(
		settings.Provider.ExtraConfigOptions.ClientKey)...)

	lines = append(lines, "")

	return lines
}
