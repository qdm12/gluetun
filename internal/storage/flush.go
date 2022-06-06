package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

// FlushToFile flushes the merged servers data to the file
// specified by path, as indented JSON.
func (s *Storage) FlushToFile(path string) error {
	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()

	return s.flushToFile(path)
}

// flushToFile flushes the merged servers data to the file
// specified by path, as indented JSON. It is not thread-safe.
func (s *Storage) flushToFile(path string) error {
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

	for _, obj := range s.mergedServers.ProviderToServers {
		sort.Sort(models.SortableServers(obj.Servers))
	}

	err = encoder.Encode(&s.mergedServers)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
