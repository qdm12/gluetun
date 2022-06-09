package updater

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

type Provider interface {
	Name() string
	FetchServers(ctx context.Context, minServers int) (servers []models.Server, err error)
}

func (u *Updater) updateProvider(ctx context.Context, provider Provider) (err error) {
	providerName := provider.Name()
	existingServersCount := u.storage.GetServersCount(providerName)
	minServers := getMinServers(existingServersCount)
	servers, err := provider.FetchServers(ctx, minServers)
	if err != nil {
		return fmt.Errorf("cannot get servers: %w", err)
	}

	if u.storage.ServersAreEqual(providerName, servers) {
		return nil
	}

	// Note the servers variable must NOT BE MUTATED after this call,
	// since the implementation does not deep copy the servers.
	// TODO set in storage in provider updater directly, server by server,
	// to avoid accumulating server data in memory.
	err = u.storage.SetServers(providerName, servers)
	if err != nil {
		return fmt.Errorf("cannot set servers to storage: %w", err)
	}
	return nil
}

func getMinServers(existingServersCount int) (minServers int) {
	const minRatio = 0.8
	return int(minRatio * float64(existingServersCount))
}
