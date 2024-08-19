package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Storage contains settings to configure the storage.
type Storage struct {
	// Filepath is the path to the servers.json file. An empty string disables on-disk storage.
	Filepath *string
}

func (s Storage) validate() (err error) {
	if *s.Filepath != "" { // optional
		_, err := filepath.Abs(*s.Filepath)
		if err != nil {
			return fmt.Errorf("filepath is not valid: %w", err)
		}
	}
	return nil
}

func (s *Storage) copy() (copied Storage) {
	return Storage{
		Filepath: gosettings.CopyPointer(s.Filepath),
	}
}

func (s *Storage) overrideWith(other Storage) {
	s.Filepath = gosettings.OverrideWithPointer(s.Filepath, other.Filepath)
}

func (s *Storage) setDefaults() {
	const defaultFilepath = "/gluetun/servers.json"
	s.Filepath = gosettings.DefaultPointer(s.Filepath, defaultFilepath)
}

func (s Storage) String() string {
	return s.toLinesNode().String()
}

func (s Storage) toLinesNode() (node *gotree.Node) {
	if *s.Filepath == "" {
		return gotree.New("Storage settings: disabled")
	}
	node = gotree.New("Storage settings:")
	node.Appendf("Filepath: %s", *s.Filepath)
	return node
}

func (s *Storage) read(r *reader.Reader) (err error) {
	s.Filepath = r.Get("STORAGE_FILEPATH", reader.AcceptEmpty(true))
	return nil
}
