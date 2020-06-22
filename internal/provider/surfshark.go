package provider

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type surfshark struct {
	fileManager files.FileManager
	lookupIP    func(host string) ([]net.IP, error)
}

func newSurfshark(fileManager files.FileManager) *surfshark {
	return &surfshark{fileManager, net.LookupIP}
}

func (s *surfshark) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.SurfsharkServers() {
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
	switch {
	case selection.Protocol == constants.TCP:
		port = 1443
	case selection.Protocol == constants.UDP:
		port = 1194
	default:
		return nil, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
	}
	return connections, nil
}

func (s *surfshark) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (err error) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "SHA512"
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Surfshark specific
		"resolv-retry infinite",
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix 1450",
		"ping 15",
		"ping-restart 0",
		"ping-timer-rem",
		"reneg-sec 0",
		"fast-io",
		"key-direction 1",
		"script-security 2",

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
	if !root {
		lines = append(lines, "user nonrootuser")
	}
	for _, connection := range connections {
		lines = append(lines, fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port))
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.SurfsharkCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.SurfsharkOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return s.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}

func (s *surfshark) GetPortForward() (port uint16, err error) {
	panic("port forwarding is not supported for surfshark")
}

func (s *surfshark) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	panic("port forwarding is not supported for surfshark")
}

func (s *surfshark) AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error) {
	panic("port forwarding is not supported for surfshark")
}
