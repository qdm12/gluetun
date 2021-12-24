package secrets

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func (r *Reader) readOpenVPN() (
	settings settings.OpenVPN, err error) {
	settings.User, err = r.readOpenVPNUser()
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN user file: %w", err)
	}

	settings.Password, err = r.readOpenVPNPassword()
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN password file: %w", err)
	}

	settings.ClientKey, err = r.readOpenVPNClientKey()
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN client key file: %w", err)
	}

	settings.ClientCrt, err = r.readOpenVPNClientCrt()
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN client certificate file: %w", err)
	}

	return settings, nil
}

func (r *Reader) readOpenVPNUser() (user string, err error) {
	// TODO have as part of env reader in settings somewhere
	const envKey = "OPENVPN_USER_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/openvpn_user"
	}
	stringPtr, err := files.ReadFromFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read OpenVPN user file: %w", err)
	} else if stringPtr != nil {
		user = *stringPtr
	}
	return user, nil
}

func (r *Reader) readOpenVPNPassword() (password string, err error) {
	// TODO have as part of env reader in settings somewhere
	const envKey = "OPENVPN_PASSWORD_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/openvpn_password"
	}
	stringPtr, err := files.ReadFromFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read OpenVPN password file: %w", err)
	} else if stringPtr != nil {
		password = *stringPtr
	}
	return password, nil
}

func (r *Reader) readOpenVPNClientKey() (clientKey *string, err error) {
	// TODO have as part of env reader in settings somewhere
	const envKey = "OPENVPN_CLIENTKEY_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/openvpn_clientkey"
	}
	return files.ReadFromFile(path)
}

func (r *Reader) readOpenVPNClientCrt() (clientCrt *string, err error) {
	// TODO have as part of env reader in settings somewhere
	const envKey = "OPENVPN_CLIENTCRT_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/openvpn_clientcrt"
	}
	return files.ReadFromFile(path)
}
