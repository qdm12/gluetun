package updater

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
)

type Providers interface {
	Get(providerName string) provider.Provider
}

type Storage interface {
	SetServers(provider string, servers []models.Server) (err error)
	GetServersCount(provider string) (count int)
	ServersAreEqual(provider string, servers []models.Server) (equal bool)
	// Extra methods to match the provider.New storage interface
	FilterServers(provider string, selection settings.ServerSelection) (filtered []models.Server, err error)
	GetServerByName(provider string, name string) (server models.Server, ok bool)
}

type Unzipper interface {
	FetchAndExtract(ctx context.Context, url string) (
		contents map[string][]byte, err error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
