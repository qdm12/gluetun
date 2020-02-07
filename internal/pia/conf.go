package pia

import (
	"fmt"
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) BuildConf(region models.PIARegion, protocol models.NetworkProtocol,
	encryption models.PIAEncryption, uid, gid int) (IPs []net.IP, port uint16, err error) {
	var X509CRL, certificate string // depends on encryption
	var cipherAlgo, authAlgo string // depends on encryption
	if encryption == constants.PIAEncryptionNormal {
		cipherAlgo = "aes-128-cbc"
		authAlgo = "sha1"
		X509CRL = constants.PIAX509CRL_NORMAL
		certificate = constants.PIACertificate_NORMAL
		if protocol == constants.UDP {
			port = 1198
		} else {
			port = 502
		}
	} else { // strong
		cipherAlgo = "aes-256-cbc"
		authAlgo = "sha256"
		X509CRL = constants.PIAX509CRL_STRONG
		certificate = constants.PIACertificate_STRONG
		if protocol == constants.UDP {
			port = 1197
		} else {
			port = 501
		}
	}
	subdomain := constants.PIARegionToSubdomainMapping[region]
	IPs, err = c.lookupIP(subdomain + ".privateinternetaccess.com")
	if err != nil {
		return nil, 0, err
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"tls-client",
		"remote-cert-tls server",
		"compress",
		"verb 1", // TODO env variable
		"reneg-sec 0",
		// Added constant values
		"mute-replay-warnings",
		"user nonrootuser",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"disable-occ",

		// Modified variables
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(protocol)),
		fmt.Sprintf("cipher %s", cipherAlgo),
		fmt.Sprintf("auth %s", authAlgo),
	}
	for _, IP := range IPs {
		lines = append(lines, fmt.Sprintf("remote %s %d", IP.String(), port))
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
	err = c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.FileOwnership(uid, gid), files.FilePermissions(0400))
	return IPs, port, err
}
