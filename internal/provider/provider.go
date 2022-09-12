package provider

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client.
type Provider interface {
	GetConnection(selection settings.ServerSelection, ipv6Supported bool) (connection models.Connection, err error)
	OpenVPNConfig(connection models.Connection, settings settings.OpenVPN, ipv6Supported bool) (lines []string)
	Name() string
	PortForwarder
	FetchServers(ctx context.Context, minServers int) (
		servers []models.Server, err error)
}

type PortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger utils.Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, client *http.Client,
		port uint16, gateway net.IP, serverName string) (err error)
}
