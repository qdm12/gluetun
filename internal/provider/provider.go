package provider

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client.
type Provider interface {
	GetConnection(selection settings.ServerSelection, ipv6Supported bool) (connection models.Connection, err error)
	OpenVPNConfig(connection models.Connection, settings settings.OpenVPN, ipv6Supported bool) (lines []string)
	Name() string
	FetchServers(ctx context.Context, minServers int) (
		servers []models.Server, err error)
}
