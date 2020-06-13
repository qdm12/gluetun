package openvpn

import (
	"context"
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/openvpn/pia"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(region models.PIARegion, protocol models.NetworkProtocol,
		encryption models.PIAEncryption, targetIP net.IP) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, encryption models.PIAEncryption, verbosity, uid, gid int, root bool, cipher, auth string) (err error)
	GetPortForward() (port uint16, err error)
	WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error)
	AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error)
}

func newPrivateInternetAccess(client network.Client, fileManager files.FileManager, firewall firewall.Configurator) {
	return pia.New(client, fileManager, firewall)
}
