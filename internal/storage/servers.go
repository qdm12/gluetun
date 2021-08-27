package storage

import "github.com/qdm12/gluetun/internal/models"

func (s *Storage) GetServers() models.AllServers {
	return s.mergedServers.GetCopy()
}
