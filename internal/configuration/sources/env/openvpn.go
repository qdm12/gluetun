package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
	"github.com/qdm12/govalid/binary"
)

func (s *Source) readOpenVPN() (
	openVPN settings.OpenVPN, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"OPENVPN_KEY", "OPENVPN_CERT",
			"OPENVPN_KEY_PASSPHRASE", "OPENVPN_ENCRYPTED_KEY"}, err)
	}()

	openVPN.Version = env.Get("OPENVPN_VERSION")
	openVPN.User = s.readOpenVPNUser()
	openVPN.Password = s.readOpenVPNPassword()
	confFile := env.Get("OPENVPN_CUSTOM_CONFIG")
	if confFile != "" {
		openVPN.ConfFile = &confFile
	}

	ciphersKey, _ := s.getEnvWithRetro("OPENVPN_CIPHERS", []string{"OPENVPN_CIPHER"})
	openVPN.Ciphers = env.CSV(ciphersKey)

	auth := env.Get("OPENVPN_AUTH")
	if auth != "" {
		openVPN.Auth = &auth
	}

	openVPN.Cert = env.StringPtr("OPENVPN_CERT", env.ForceLowercase(false))
	openVPN.Key = env.StringPtr("OPENVPN_KEY", env.ForceLowercase(false))
	openVPN.EncryptedKey = env.StringPtr("OPENVPN_ENCRYPTED_KEY", env.ForceLowercase(false))

	openVPN.KeyPassphrase = s.readOpenVPNKeyPassphrase()

	openVPN.PIAEncPreset = s.readPIAEncryptionPreset()

	openVPN.MSSFix, err = env.Uint16Ptr("OPENVPN_MSSFIX")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_MSSFIX: %w", err)
	}

	_, openVPN.Interface = s.getEnvWithRetro("VPN_INTERFACE",
		[]string{"OPENVPN_INTERFACE"}, env.ForceLowercase(false))

	openVPN.ProcessUser, err = s.readOpenVPNProcessUser()
	if err != nil {
		return openVPN, err
	}

	openVPN.Verbosity, err = env.IntPtr("OPENVPN_VERBOSITY")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_VERBOSITY: %w", err)
	}

	flagsStr := env.Get("OPENVPN_FLAGS", env.ForceLowercase(false))
	if flagsStr != "" {
		openVPN.Flags = strings.Fields(flagsStr)
	}

	return openVPN, nil
}

func (s *Source) readOpenVPNUser() (user *string) {
	user = new(string)
	_, *user = s.getEnvWithRetro("OPENVPN_USER",
		[]string{"USER"}, env.ForceLowercase(false))
	if *user == "" {
		return nil
	}

	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	*user = strings.ReplaceAll(*user, " ", "")
	return user
}

func (s *Source) readOpenVPNPassword() (password *string) {
	password = new(string)
	_, *password = s.getEnvWithRetro("OPENVPN_PASSWORD",
		[]string{"PASSWORD"}, env.ForceLowercase(false))
	if *password == "" {
		return nil
	}

	return password
}

func (s *Source) readOpenVPNKeyPassphrase() (passphrase *string) {
	passphrase = new(string)
	*passphrase = env.Get("OPENVPN_KEY_PASSPHRASE", env.ForceLowercase(false))
	if *passphrase == "" {
		return nil
	}
	return passphrase
}

func (s *Source) readPIAEncryptionPreset() (presetPtr *string) {
	_, preset := s.getEnvWithRetro(
		"PRIVATE_INTERNET_ACCESS_OPENVPN_ENCRYPTION_PRESET",
		[]string{"PIA_ENCRYPTION", "ENCRYPTION"})
	if preset != "" {
		return &preset
	}
	return nil
}

func (s *Source) readOpenVPNProcessUser() (processUser string, err error) {
	key, value := s.getEnvWithRetro("OPENVPN_PROCESS_USER",
		[]string{"OPENVPN_ROOT"})
	if key == "OPENVPN_PROCESS_USER" {
		return value, nil
	}

	// Retro-compatibility
	if value == "" {
		return "", nil
	}
	root, err := binary.Validate(value)
	if err != nil {
		return "", fmt.Errorf("environment variable %s: %w", key, err)
	}
	if *root {
		return "root", nil
	}
	const defaultNonRootUser = "nonrootuser"
	return defaultNonRootUser, nil
}
