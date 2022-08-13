package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

const (
	// OpenVPNClientKeyPath is the OpenVPN client key filepath.
	OpenVPNClientKeyPath = "/gluetun/client.key"
	// OpenVPNClientCertificatePath is the OpenVPN client certificate filepath.
	OpenVPNClientCertificatePath = "/gluetun/client.crt"
)

func (r *Reader) readOpenVPN() (settings settings.OpenVPN, err error) {
	settings.Key, err = ReadFromFile(OpenVPNClientKeyPath)
	if err != nil {
		return settings, fmt.Errorf("client key: %w", err)
	}

	settings.ClientCrt, err = ReadFromFile(OpenVPNClientCertificatePath)
	if err != nil {
		return settings, fmt.Errorf("client certificate: %w", err)
	}

	return settings, nil
}
