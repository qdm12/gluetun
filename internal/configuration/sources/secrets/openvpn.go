package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readOpenVPN() (
	settings settings.OpenVPN, err error) {
	settings.User, err = s.readSecretFileAsStringPtr(
		"OPENVPN_USER_SECRETFILE",
		"/run/secrets/openvpn_user",
	)
	if err != nil {
		return settings, fmt.Errorf("reading user file: %w", err)
	}

	settings.Password, err = s.readSecretFileAsStringPtr(
		"OPENVPN_PASSWORD_SECRETFILE",
		"/run/secrets/openvpn_password",
	)
	if err != nil {
		return settings, fmt.Errorf("reading password file: %w", err)
	}

	settings.Key, err = s.readPEMSecretFile(
		"OPENVPN_CLIENTKEY_SECRETFILE",
		"/run/secrets/openvpn_clientkey",
	)
	if err != nil {
		return settings, fmt.Errorf("reading client key file: %w", err)
	}

	settings.EncryptedKey, err = s.readPEMSecretFile(
		"OPENVPN_ENCRYPTED_KEY_SECRETFILE",
		"/run/secrets/openvpn_encrypted_key",
	)
	if err != nil {
		return settings, fmt.Errorf("reading encrypted key file: %w", err)
	}

	settings.KeyPassphrase, err = s.readSecretFileAsStringPtr(
		"OPENVPN_KEY_PASSPHRASE_SECRETFILE",
		"/run/secrets/openvpn_key_passphrase",
	)
	if err != nil {
		return settings, fmt.Errorf("reading key passphrase file: %w", err)
	}

	settings.Cert, err = s.readPEMSecretFile(
		"OPENVPN_CLIENTCRT_SECRETFILE",
		"/run/secrets/openvpn_clientcrt",
	)
	if err != nil {
		return settings, fmt.Errorf("reading client certificate file: %w", err)
	}

	return settings, nil
}
