package provider

import (
	"context"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (err error)
	GetPortForward() (port uint16, err error)
	WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error)
	AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error)
}

func New(provider models.VPNProvider, logger logging.Logger, client network.Client, fileManager files.FileManager, firewall firewall.Configurator) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccess(client, fileManager, firewall)
	case constants.Mullvad:
		return newMullvad(fileManager, logger)
	case constants.Windscribe:
		return newWindscribe(fileManager)
	case constants.Surfshark:
		return newSurfshark(fileManager)
	case constants.Cyberghost:
		return newCyberghost(fileManager)
	default:
		return nil // should never occur
	}
}
