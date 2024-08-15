package settings

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// StorageSettings contains settings to configure the storage.
type StorageSettings struct {
	// Filepath is the path to the servers.json file. An empty string disables on-disk storage.
	Filepath *string
}

func (s StorageSettings) validate() (err error) {
	return nil
}

func (s *StorageSettings) copy() (copied StorageSettings) {
	return StorageSettings{
		Filepath: s.Filepath,
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (s *StorageSettings) overrideWith(other StorageSettings) {
	s.Filepath = gosettings.OverrideWithPointer(s.Filepath, other.Filepath)
}

func (s *StorageSettings) setDefaults() {
	s.Filepath = gosettings.DefaultPointer(s.Filepath, constants.ServersData)
}

func (s StorageSettings) String() string {
	return s.toLinesNode().String()
}

func (s StorageSettings) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Storage settings:")
	node.Appendf("Filepath: %s", *s.Filepath)
	return node
}

func (s *StorageSettings) read(r *reader.Reader) (err error) {
	s.Filepath = r.Get("STORAGE_FILEPATH", reader.AcceptEmpty(true))
	return nil
}
