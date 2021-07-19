package privateinternetaccess

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	var defaultCipher, defaultAuth, X509CRL, certificate string
	switch settings.Provider.ExtraConfigOptions.EncryptionPreset {
	case constants.PIAEncryptionPresetNormal:
		defaultCipher = constants.AES128cbc
		defaultAuth = constants.SHA1
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	case constants.PIAEncryptionPresetStrong:
		defaultCipher = constants.AES256cbc
		defaultAuth = constants.SHA256
		X509CRL = constants.PiaX509CRLStrong
		certificate = constants.PIACertificateStrong
	default: // no encryption preset
		defaultCipher = "none"
		defaultAuth = "none"
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	}

	if settings.Cipher == "" {
		settings.Cipher = defaultCipher
	}

	if settings.Auth == "" {
		settings.Auth = defaultAuth
	}

	lines = []string{
		"client",
		"dev tun",
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
		connection.ProtoLine(),
		connection.RemoteLine(),
	}

	if settings.Cipher != "" {
		lines = append(lines, utils.CipherLines(settings.Cipher, settings.Version)...)
	}

	if settings.Auth != "" {
		lines = append(lines, "auth "+settings.Auth)
	}

	if !settings.Root {
		lines = append(lines, "user "+username)
	}

	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}

	if settings.Provider.ExtraConfigOptions.OpenVPNIPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}

	lines = append(lines, utils.WrapOpenvpnCA(certificate)...)
	lines = append(lines, utils.WrapOpenvpnCRLVerify(X509CRL)...)

	lines = append(lines, "")

	return lines
}
