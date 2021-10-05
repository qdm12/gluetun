package vpnunlimited

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(settings.Verbosity),

		// VPNUnlimited specific
		"ping 5",
		"remote-cert-tls server",
		"reneg-sec 0",
		"route-metric 1",

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

	if len(settings.Ciphers) > 0 {
		lines = append(lines, utils.CipherLines(settings.Ciphers, settings.Version)...)
	}

	if settings.Auth != "" {
		lines = append(lines, "auth "+settings.Auth)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

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

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.VPNUnlimitedCertificateAuthority)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		settings.ClientCrt)...)
	lines = append(lines, utils.WrapOpenvpnKey(
		settings.ClientKey)...)

	lines = append(lines, "")

	return lines, nil
}
