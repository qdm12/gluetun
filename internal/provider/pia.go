package provider

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func buildPIAConf(connection models.OpenVPNConnection, verbosity int, root bool, cipher, auth string,
	extras models.ExtraConfigOptions) (lines []string) {
	var X509CRL, certificate string
	var defaultCipher, defaultAuth string
	if extras.EncryptionPreset == constants.PIAEncryptionPresetNormal {
		defaultCipher = "aes-128-cbc"
		defaultAuth = "sha1"
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	} else { // strong encryption
		defaultCipher = aes256cbc
		defaultAuth = "sha256"
		X509CRL = constants.PiaX509CRLStrong
		certificate = constants.PIACertificateStrong
	}
	if len(cipher) == 0 {
		cipher = defaultCipher
	}
	if len(auth) == 0 {
		auth = defaultAuth
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// PIA specific
		"ping 300", // Ping every 5 minutes to prevent a timeout error
		"reneg-sec 0",
		"compress", // allow PIA server to choose the compression to use

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if strings.HasSuffix(cipher, "-gcm") {
		lines = append(lines, "ncp-disable")
	}
	if !root {
		lines = append(lines, "user nonrootuser")
	}
	lines = append(lines, []string{
		"<crl-verify>",
		"-----BEGIN X509 CRL-----",
		X509CRL,
		"-----END X509 CRL-----",
		"</crl-verify>",
	}...)
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		certificate,
		"-----END CERTIFICATE-----",
		"</ca>",
		"",
	}...)
	return lines
}
