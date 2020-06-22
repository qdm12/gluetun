package provider

import (
	"context"
	"fmt"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type mullvad struct {
	fileManager files.FileManager
	logger      logging.Logger
}

func newMullvad(fileManager files.FileManager, logger logging.Logger) *mullvad {
	return &mullvad{
		fileManager: fileManager,
		logger:      logger.WithPrefix("Mullvad configurator: "),
	}
}

func (m *mullvad) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	servers := constants.MullvadServerFilter(selection.Country, selection.City, selection.ISP)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for country %q, city %q and ISP %q", selection.Country, selection.City, selection.ISP)
	}
	for _, server := range servers {
		port := server.DefaultPort
		if selection.CustomPort > 0 {
			port = selection.CustomPort
		}
		for _, IP := range server.IPs {
			if selection.TargetIP != nil {
				if selection.TargetIP.Equal(IP) {
					return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: selection.Protocol}}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
			}
		}
	}
	if selection.TargetIP != nil {
		return nil, fmt.Errorf("target IP address %q not found in IP addresses", selection.TargetIP)
	}
	return connections, nil
}

func (m *mullvad) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (err error) {
	if len(connections) == 0 {
		return fmt.Errorf("at least one connection string is expected")
	}
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Mullvad specific
		"ping 10",
		"ping-restart 60",
		"sndbuf 524288",
		"rcvbuf 524288",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",
		"tun-ipv6",
		"fast-io",
		"script-security 2",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
		fmt.Sprintf("cipher %s", cipher),
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
		constants.MullvadCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
		"",
	}...)
	return m.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}

func (m *mullvad) GetPortForward() (port uint16, err error) {
	panic("port forwarding is not supported for mullvad")
}

func (m *mullvad) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	panic("port forwarding is not supported for mullvad")
}

func (m *mullvad) AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error) {
	panic("port forwarding is not supported for mullvad")
}
