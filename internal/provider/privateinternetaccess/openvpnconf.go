package privateinternetaccess

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	var defaultCipher, defaultAuth, X509CRL, certificate string
	switch settings.EncPreset {
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

	if settings.Cipher == "" {
		settings.Cipher = defaultCipher
	}

	if settings.Auth == "" {
		settings.Auth = defaultAuth
	}

	lines = []string{
		"client",
		"dev " + settings.Interface,
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// PIA specific
		"reneg-sec 0",
		"disable-occ",
		"compress",    // allow PIA server to choose the compression to use
		"ncp-disable", // prevent from auto-upgrading cipher to aes-256-gcm

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
	}

	if settings.Cipher != "" {
		lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)
	}

	if settings.Auth != "" {
		lines = append(lines, "auth "+settings.Auth)
	}

	if !settings.Root {
		lines = append(lines, "user "+settings.ProcUser)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if !settings.IPv6 {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(certificate)...)
	lines = append(lines, utils.WrapOpenvpnCRLVerify(X509CRL)...)

	lines = append(lines, "")

	return lines, nil
}
