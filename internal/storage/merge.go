package storage

import (
	"strconv"
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
)

func (s *Storage) logVersionDiff(provider string, hardcodedVersion, persistedVersion uint16) {
	message := provider + " servers from file discarded because they have version " +
		strconv.Itoa(int(persistedVersion)) +
		" and hardcoded servers have version " +
		strconv.Itoa(int(hardcodedVersion))
	s.logger.Info(message)
}

func (s *Storage) logTimeDiff(provider string, persistedUnix, hardcodedUnix int64) {
	diff := time.Unix(persistedUnix, 0).Sub(time.Unix(hardcodedUnix, 0))
	if diff < 0 {
		diff = -diff
	}
	diff = diff.Truncate(time.Second)
	message := "Using " + provider + " servers from file which are " +
		diff.String() + " more recent"
	s.logger.Info(message)
}

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

	s.logTimeDiff(provider, persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}
