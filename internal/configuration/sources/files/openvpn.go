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
	openVPNEncryptedKey          = "/gluetun/openvpn_encrypted_key"
)

func (s *Source) readOpenVPN() (settings settings.OpenVPN, err error) {
	settings.Key, err = readPEMFile(OpenVPNClientKeyPath)
	if err != nil {
		return settings, fmt.Errorf("client key: %w", err)
	}

	settings.Cert, err = readPEMFile(OpenVPNClientCertificatePath)
	if err != nil {
		return settings, fmt.Errorf("client certificate: %w", err)
	}
	settings.EncryptedKey, err = readPEMFile(openVPNEncryptedKey)
	if err != nil {
		return settings, fmt.Errorf("reading encrypted key file: %w", err)
	}

	return settings, nil
}
