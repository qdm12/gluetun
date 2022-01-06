package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

func (r *Reader) readOpenVPN() (settings settings.OpenVPN, err error) {
	settings.ClientKey, err = ReadFromFile(constants.ClientKey)
	if err != nil {
		return settings, fmt.Errorf("cannot read client key: %w", err)
	}

	settings.ClientCrt, err = ReadFromFile(constants.ClientCertificate)
	if err != nil {
		return settings, fmt.Errorf("cannot read client certificate: %w", err)
	}

	return settings, nil
}
