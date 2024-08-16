package storage

import (
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
)

// SetServers sets the given servers for the given provider
// in the storage in-memory map and saves all the servers
// to file.
// Note the servers given are not copied so the caller must
// NOT MUTATE them after calling this method.
func (s *Storage) SetServers(provider string, servers []models.Server) (err error) {
	if provider == providers.Custom {
		return
	}

	s.mergedMutex.Lock()
	defer s.mergedMutex.Unlock()

	serversObject := s.getMergedServersObject(provider)
	serversObject.Timestamp = time.Now().Unix()
	serversObject.Servers = servers
	s.mergedServers.ProviderToServers[provider] = serversObject

	err = s.flushToFile(s.filepath)
	if err != nil {
		return fmt.Errorf("saving servers to file: %w", err)
	}
	return nil
}

// GetServersCount returns the number of servers for the provider given.
func (s *Storage) GetServersCount(provider string) (count int) {
	if provider == providers.Custom {
		return 0
	}

	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()

	serversObject := s.getMergedServersObject(provider)
	return len(serversObject.Servers)
}

// Format formats the servers for the provider using the format given
// and returns the resulting string.
func (s *Storage) Format(provider, format string) (formatted string, err error) {
	if provider == providers.Custom {
		return "", nil
	}

	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()

	serversObject := s.getMergedServersObject(provider)
	return serversObject.Format(provider, format)
}

// GetServersCount returns the number of servers for the provider given.
func (s *Storage) ServersAreEqual(provider string, servers []models.Server) (equal bool) {
	if provider == providers.Custom {
		return true
	}

	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()

	serversObject := s.getMergedServersObject(provider)
	existingServers := serversObject.Servers

	if len(existingServers) != len(servers) {
		return false
	}

	for i := range existingServers {
		if !existingServers[i].Equal(servers[i]) {
			return false
		}
	}

	return true
}

func (s *Storage) getMergedServersObject(provider string) (serversObject models.Servers) {
	serversObject, ok := s.mergedServers.ProviderToServers[provider]
	if !ok {
		panicOnProviderMissingHardcoded(provider)
	}
	return serversObject
}
