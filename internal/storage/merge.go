package storage

import (
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
)

func (s *Storage) mergeServers(hardcoded, persisted models.AllServers) models.AllServers {
	allProviders := providers.All()
	merged := models.AllServers{
		Version:           hardcoded.Version,
		ProviderToServers: make(map[string]models.Servers, len(allProviders)),
	}

	for _, provider := range allProviders {
		hardcodedServers := hardcoded.ProviderToServers[provider]
		persistedServers := persisted.ProviderToServers[provider]
		merged.ProviderToServers[provider] = s.mergeProviderServers(provider,
			hardcodedServers, persistedServers)
	}

	return merged
}

func (s *Storage) mergeProviderServers(provider string,
	hardcoded, persisted models.Servers) (merged models.Servers) {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	diff := time.Unix(persisted.Timestamp, 0).Sub(time.Unix(hardcoded.Timestamp, 0))
	if diff < 0 {
		diff = -diff
	}
	diff = diff.Truncate(time.Second)
	message := "Using " + provider + " servers from file which are " +
		diff.String() + " more recent"
	s.logger.Info(message)

	return persisted
}
