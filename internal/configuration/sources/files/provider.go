package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readProvider() (provider settings.Provider, err error) {
	provider.ServerSelection, err = s.readServerSelection()
	if err != nil {
		return provider, fmt.Errorf("server selection: %w", err)
	}

	return provider, nil
}
