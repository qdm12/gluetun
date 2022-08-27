package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readServerSelection() (selection settings.ServerSelection, err error) {
	selection.Wireguard, err = s.readWireguardSelection()
	if err != nil {
		return selection, fmt.Errorf("wireguard: %w", err)
	}

	return selection, nil
}
