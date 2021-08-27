package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/models"
)

var _ Flusher = (*Storage)(nil)

type Flusher interface {
	FlushToFile(allServers models.AllServers) error
}

func (s *Storage) FlushToFile(allServers models.AllServers) error {
	return flushToFile(s.filepath, allServers)
}

func flushToFile(path string, servers models.AllServers) error {
	dirPath := filepath.Dir(path)
	if err := os.MkdirAll(dirPath, 0644); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(servers); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
