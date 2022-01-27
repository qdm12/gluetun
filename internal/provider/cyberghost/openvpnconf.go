package cyberghost

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/parse"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (c *Cyberghost) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{
			constants.AES256gcm,
			constants.AES256cbc,
			constants.AES128gcm,
		}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = constants.SHA256
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// Cyberghost specific
		"ping 10",
		"remote-cert-tls server",
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

	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}

	if settings.ProcessUser != "root" {
		lines = append(lines, "user "+settings.ProcessUser)
		lines = append(lines, "persist-tun")
		lines = append(lines, "persist-key")
	}

	if *settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(*settings.MSSFix)))
	}

	if !*settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(
		constants.CyberghostCertificate)...)

	certData, err := parse.ExtractCert([]byte(*settings.ClientCrt))
	if err != nil {
		return nil, fmt.Errorf("client cert is not valid: %w", err)
	}
	lines = append(lines, utils.WrapOpenvpnCert(certData)...)

	keyData, err := parse.ExtractPrivateKey([]byte(*settings.ClientKey))
	if err != nil {
		return nil, fmt.Errorf("client key is not valid: %w", err)
	}
	lines = append(lines, utils.WrapOpenvpnKey(keyData)...)

	lines = append(lines, "")

	return lines, nil
}
