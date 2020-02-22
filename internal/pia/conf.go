package pia

import (
	"fmt"
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(region models.PIARegion, protocol models.NetworkProtocol, encryption models.PIAEncryption) (connections []models.OpenVPNConnection, err error) {
	geoMapping := constants.PIAGeoToSubdomainMapping()
	var subdomain string
	for r, s := range geoMapping {
		if strings.ToLower(string(region)) == strings.ToLower(string(r)) {
			subdomain = s
			break
		}
	}
	if len(subdomain) == 0 {
		return nil, fmt.Errorf("region %q has no associated PIA subdomain", region)
	}
	if err != nil {
		return nil, err
	}
	IPs, err := c.lookupIP(subdomain + ".privateinternetaccess.com")
	if err != nil {
		return nil, err
	}
	var port uint16
	switch protocol {
	case constants.TCP:
		switch encryption {
		case constants.PIAEncryptionNormal:
			port = 502
		case constants.PIAEncryptionStrong:
			port = 501
		}
	case constants.UDP:
		switch encryption {
		case constants.PIAEncryptionNormal:
			port = 1198
		case constants.PIAEncryptionStrong:
			port = 1197
		}
	}
	if port == 0 {
		return nil, fmt.Errorf("combination of protocol %q and encryption %q does not yield any port number", protocol, encryption)
	}
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: protocol})
	}
	return connections, nil
}

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, encryption models.PIAEncryption, uid, gid int) (err error) {
	var X509CRL, certificate, cipherAlgo, authAlgo string
	if encryption == constants.PIAEncryptionNormal {
		cipherAlgo = "aes-128-cbc"
		authAlgo = "sha1"
		X509CRL = constants.PIAX509CRL_NORMAL
		certificate = constants.PIACertificate_NORMAL
	} else { // strong encryption
		cipherAlgo = "aes-256-cbc"
		authAlgo = "sha256"
		X509CRL = constants.PIAX509CRL_STRONG
		certificate = constants.PIACertificate_STRONG
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"remote-cert-tls server",
		"ping 300", // Ping every 5 minutes to prevent a timeout error
		"verb 1",   // TODO env variable

		// PIA specific
		"reneg-sec 0",

		// Added constant values
		"mute-replay-warnings",
		"user nonrootuser",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",

		// Modified variables
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
		fmt.Sprintf("cipher %s", cipherAlgo),
		fmt.Sprintf("auth %s", authAlgo),
	}
	for _, connection := range connections {
		lines = append(lines, fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port))
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
	return c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}
