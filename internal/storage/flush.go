package storage

import "github.com/qdm12/gluetun/internal/models"

func (s *storage) FlushToFile(allServers models.AllServers) error {
	return flushToFile(s.filepath, allServers)
}
