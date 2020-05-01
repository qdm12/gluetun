package pia

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(region models.PIARegion, protocol models.NetworkProtocol, encryption models.PIAEncryption, targetIP net.IP) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.PIAServers() {
		if strings.EqualFold(string(server.Region), string(region)) {
			IPs = server.IPs
		}
	}
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IP found for region %q", region)
	}
	if targetIP != nil {
		found := false
		for i := range IPs {
			if IPs[i].Equal(targetIP) {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("target IP address %q not found in IP addresses", targetIP)
		}
		IPs = []net.IP{targetIP}
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

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, encryption models.PIAEncryption, verbosity, uid, gid int, root bool, cipher, auth string) (err error) {
	var X509CRL, certificate string
	if encryption == constants.PIAEncryptionNormal {
		if len(cipher) == 0 {
			cipher = "aes-128-cbc"
		}
		if len(auth) == 0 {
			auth = "sha1"
		}
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	} else { // strong encryption
		if len(cipher) == 0 {
			cipher = "aes-256-cbc"
		}
		if len(auth) == 0 {
			auth = "sha256"
		}
		X509CRL = constants.PiaX509CRLStrong
		certificate = constants.PIACertificateStrong
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"tls-client",
		"remote-cert-tls server",
		"ping 300", // Ping every 5 minutes to prevent a timeout error

		// PIA specific
		"reneg-sec 0",
		"compress", // allow PIA server to choose the compression to use

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if strings.HasSuffix(cipher, "-gcm") {
		lines = append(lines, "ncp-disable")
	}
	if !root {
		lines = append(lines, "user nonrootuser")
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
