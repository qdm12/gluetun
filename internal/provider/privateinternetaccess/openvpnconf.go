package privateinternetaccess

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string, err error) {
	var defaultCipher, defaultAuth, X509CRL, certificate string
	switch *settings.PIAEncPreset {
	case constants.PIAEncryptionPresetNormal:
		defaultCipher = constants.AES128cbc
		defaultAuth = constants.SHA1
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	case constants.PIAEncryptionPresetNone:
		defaultCipher = "none"
		defaultAuth = "none"
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	default: // strong
		defaultCipher = constants.AES256cbc
		defaultAuth = constants.SHA256
		X509CRL = constants.PiaX509CRLStrong
		certificate = constants.PIACertificateStrong
	}

	if len(settings.Ciphers) == 0 {
		settings.Ciphers = []string{defaultCipher}
	}

	auth := *settings.Auth
	if auth == "" {
		auth = defaultAuth
	}

	lines = []string{
		"client",
		"nobind",
		"tls-exit",
		"dev " + settings.Interface,
		"verb " + strconv.Itoa(*settings.Verbosity),

		// PIA specific
		"remote-cert-tls server",
		"reneg-sec 0",
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

	if len(settings.Ciphers) > 0 {
		lines = append(lines, utils.CipherLines(settings.Ciphers, settings.Version)...)
	}

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

	lines = append(lines, utils.WrapOpenvpnCA(certificate)...)
	lines = append(lines, utils.WrapOpenvpnCRLVerify(X509CRL)...)

	lines = append(lines, "")

	return lines, nil
}
