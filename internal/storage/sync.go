package storage

import (
	"fmt"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
)

func countServers(allServers models.AllServers) (count int) {
	for _, servers := range allServers.ProviderToServers {
		count += len(servers.Servers)
	}
	return count
}

// syncServers merges the hardcoded servers with the ones from the file.
func (s *Storage) syncServers() (err error) {
	hardcodedVersions := make(map[string]uint16, len(s.hardcodedServers.ProviderToServers))
	for provider, servers := range s.hardcodedServers.ProviderToServers {
		hardcodedVersions[provider] = servers.Version
	}

	serversOnFile, err := s.readFromFile(s.filepath, hardcodedVersions)
	if err != nil {
		return fmt.Errorf("cannot read servers from file: %w", err)
	}

	hardcodedCount := countServers(s.hardcodedServers)
	countOnFile := countServers(serversOnFile)

	s.mergedMutex.Lock()
	defer s.mergedMutex.Unlock()

	if countOnFile == 0 {
		s.logger.Info(fmt.Sprintf(
			"creating %s with %d hardcoded servers",
			s.filepath, hardcodedCount))
		s.mergedServers = s.hardcodedServers
	} else {
		s.logger.Info(fmt.Sprintf(
			"merging by most recent %d hardcoded servers and %d servers read from %s",
			hardcodedCount, countOnFile, s.filepath))

		s.mergedServers = s.mergeServers(s.hardcodedServers, serversOnFile)
	}

	// Eventually write file
	if s.filepath == "" || reflect.DeepEqual(serversOnFile, s.mergedServers) {
		return nil
	}

	err = s.flushToFile(s.filepath)
	if err != nil {
		return fmt.Errorf("cannot write servers to file: %w", err)
	}
	return nil
}
