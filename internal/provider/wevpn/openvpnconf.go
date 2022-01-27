package wevpn

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/parse"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Wevpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{constants.AES256gcm}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = constants.SHA512
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Wevpn specific
		"ping 30",
		"remote-cert-tls server",
		"redirect-gateway def1 bypass-dhcp",
		"reneg-sec 0",
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

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	lines = append(lines, utils.CipherLines(settings.Ciphers, settings.Version)...)

	if settings.ProcessUser != "root" {
		lines = append(lines, "user "+settings.ProcessUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if *settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(*settings.MSSFix)))
	}

	if *settings.IPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	keyData, err := parse.ExtractPrivateKey([]byte(*settings.ClientKey))
	if err != nil {
		return nil, fmt.Errorf("client key is not valid: %w", err)
	}
	lines = append(lines, utils.WrapOpenvpnKey(keyData)...)

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.WevpnCA)...)
	lines = append(lines, utils.WrapOpenvpnCert(
		constants.WevpnCertificate)...)
	lines = append(lines, utils.WrapOpenvpnTLSCrypt(
		constants.WevpnOpenvpnStaticKeyV1)...)

	lines = append(lines, "")

	return lines, nil
}
