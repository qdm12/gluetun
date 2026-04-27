package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Storage contains settings to configure the storage.
type Storage struct {
	// ServersPath is the path to the servers files directory.
	// An empty string disables on-disk storage.
	ServersPath *string
	// LegacyServersFilepath is the legacy "fat" JSON filepath to migrate from.
	// TODO v4: remove
	LegacyServersFilepath *string
}

func (s Storage) validate() (err error) {
	if *s.ServersPath != "" { // optional
		_, err := filepath.Abs(*s.ServersPath)
		if err != nil {
			return fmt.Errorf("servers path is not valid: %w", err)
		}
	}
	if *s.LegacyServersFilepath != "" {
		_, err := filepath.Abs(*s.LegacyServersFilepath)
		if err != nil {
			return fmt.Errorf("legacy servers filepath is not valid: %w", err)
		}
	}
	return nil
}

func (s *Storage) copy() (copied Storage) {
	return Storage{
		ServersPath:           gosettings.CopyPointer(s.ServersPath),
		LegacyServersFilepath: gosettings.CopyPointer(s.LegacyServersFilepath),
	}
}

func (s *Storage) overrideWith(other Storage) {
	s.ServersPath = gosettings.OverrideWithPointer(s.ServersPath, other.ServersPath)
	s.LegacyServersFilepath = gosettings.OverrideWithPointer(s.LegacyServersFilepath, other.LegacyServersFilepath)
}

func (s *Storage) SetDefaults() {
	const defaultServersPath = "/gluetun/servers/"
	s.ServersPath = gosettings.DefaultPointer(s.ServersPath, defaultServersPath)
	s.LegacyServersFilepath = gosettings.DefaultPointer(s.LegacyServersFilepath, constants.ServersDataLegacy)
}

func (s Storage) String() string {
	return s.toLinesNode().String()
}

func (s Storage) toLinesNode() (node *gotree.Node) {
	if *s.ServersPath == "" {
		return gotree.New("Storage settings: disabled")
	}
	node = gotree.New("Storage settings:")
	node.Appendf("Servers directory path: %s", *s.ServersPath)
	if *s.LegacyServersFilepath != constants.ServersDataLegacy {
		node.Appendf("Legacy servers filepath: %s", *s.LegacyServersFilepath)
	}
	return node
}

func (s *Storage) Read(r *reader.Reader) (err error) {
	// Retro-compatibility:
	// TODO v4: remove support for STORAGE_FILEPATH
	filePath := r.Get("STORAGE_FILEPATH", reader.AcceptEmpty(true), reader.IsRetro("STORAGE_SERVERS_DIRECTORY_PATH"))
	if filePath != nil {
		s.LegacyServersFilepath = filePath
		if *filePath == "" {
			s.ServersPath = ptrTo("") // skip disk operations
		}
	}
	if s.ServersPath == nil {
		s.ServersPath = r.Get("STORAGE_SERVERS_DIRECTORY_PATH", reader.AcceptEmpty(true))
	}
	return nil
}
