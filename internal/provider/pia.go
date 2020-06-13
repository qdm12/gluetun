package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type pia struct {
	client      network.Client
	fileManager files.FileManager
	firewall    firewall.Configurator
	random      random.Random
	verifyPort  func(port string) error
	lookupIP    func(host string) ([]net.IP, error)
}

func newPrivateInternetAccess(client network.Client, fileManager files.FileManager, firewall firewall.Configurator) *pia {
	return &pia{
		client:      client,
		fileManager: fileManager,
		firewall:    firewall,
		random:      random.NewRandom(),
		verifyPort:  verification.NewVerifier().VerifyPort,
		lookupIP:    net.LookupIP}
}

func (p *pia) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.PIAServers() {
		if strings.EqualFold(server.Region, selection.Region) {
			IPs = server.IPs
		}
	}
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IP found for region %q", selection.Region)
	}
	if selection.TargetIP != nil {
		found := false
		for i := range IPs {
			if IPs[i].Equal(selection.TargetIP) {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("target IP address %q not found in IP addresses", selection.TargetIP)
		}
		IPs = []net.IP{selection.TargetIP}
	}
	var port uint16
	switch selection.Protocol {
	case constants.TCP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 502
		case constants.PIAEncryptionPresetStrong:
			port = 501
		}
	case constants.UDP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 1198
		case constants.PIAEncryptionPresetStrong:
			port = 1197
		}
	}
	if port == 0 {
		return nil, fmt.Errorf("combination of protocol %q and encryption %q does not yield any port number", selection.Protocol, selection.EncryptionPreset)
	}
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
	}
	return connections, nil
}

func (p *pia) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (err error) {
	var X509CRL, certificate string
	if extras.EncryptionPreset == constants.PIAEncryptionPresetNormal {
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
			cipher = aes256cbc
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
	return p.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}

func (p *pia) GetPortForward() (port uint16, err error) {
	b, err := p.random.GenerateRandomBytes(32)
	if err != nil {
		return 0, err
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	content, status, err := p.client.GetContent(url)
	switch {
	case err != nil:
		return 0, err
	case status != http.StatusOK:
		return 0, fmt.Errorf("status is %d for %s; does your PIA server support port forwarding?", status, url)
	case len(content) == 0:
		return 0, fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding")
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(content, &body); err != nil {
		return 0, fmt.Errorf("port forwarding response: %w", err)
	}
	return body.Port, nil
}

func (p *pia) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	return p.fileManager.WriteLinesToFile(
		string(filepath),
		[]string{fmt.Sprintf("%d", port)},
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func (p *pia) AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error) {
	return p.firewall.AllowInputTrafficOnPort(ctx, device, port)
}
