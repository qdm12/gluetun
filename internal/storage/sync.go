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

func (s *Storage) SyncServers() (err error) {
	serversOnFile, err := s.readFromFile(s.filepath, s.hardcodedServers)
	if err != nil {
		return fmt.Errorf("cannot read servers from file: %w", err)
	}

	hardcodedCount := countServers(s.hardcodedServers)
	countOnFile := countServers(serversOnFile)

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

	if err := flushToFile(s.filepath, &s.mergedServers); err != nil {
		return fmt.Errorf("cannot write servers to file: %w", err)
	}
	return nil
}
