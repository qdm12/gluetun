package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
)

func (s *Source) readOpenVPN() (
	openVPN settings.OpenVPN, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"OPENVPN_KEY", "OPENVPN_CERT",
			"OPENVPN_KEY_PASSPHRASE", "OPENVPN_ENCRYPTED_KEY"}, err)
	}()

	openVPN.Version = getCleanedEnv("OPENVPN_VERSION")
	openVPN.User = s.readOpenVPNUser()
	openVPN.Password = s.readOpenVPNPassword()
	confFile := getCleanedEnv("OPENVPN_CUSTOM_CONFIG")
	if confFile != "" {
		openVPN.ConfFile = &confFile
	}

	ciphersKey, _ := s.getEnvWithRetro("OPENVPN_CIPHERS", "OPENVPN_CIPHER")
	openVPN.Ciphers = envToCSV(ciphersKey)

	auth := getCleanedEnv("OPENVPN_AUTH")
	if auth != "" {
		openVPN.Auth = &auth
	}

	openVPN.Cert = envToStringPtr("OPENVPN_CERT")
	openVPN.Key = envToStringPtr("OPENVPN_KEY")
	openVPN.EncryptedKey = envToStringPtr("OPENVPN_ENCRYPTED_KEY")

	openVPN.KeyPassphrase = s.readOpenVPNKeyPassphrase()

	openVPN.PIAEncPreset = s.readPIAEncryptionPreset()

	openVPN.MSSFix, err = envToUint16Ptr("OPENVPN_MSSFIX")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_MSSFIX: %w", err)
	}

	_, openVPN.Interface = s.getEnvWithRetro("VPN_INTERFACE", "OPENVPN_INTERFACE")

	openVPN.ProcessUser, err = s.readOpenVPNProcessUser()
	if err != nil {
		return openVPN, err
	}

	openVPN.Verbosity, err = envToIntPtr("OPENVPN_VERBOSITY")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_VERBOSITY: %w", err)
	}

	flagsStr := getCleanedEnv("OPENVPN_FLAGS")
	if flagsStr != "" {
		openVPN.Flags = strings.Fields(flagsStr)
	}

	return openVPN, nil
}

func (s *Source) readOpenVPNUser() (user *string) {
	user = new(string)
	_, *user = s.getEnvWithRetro("OPENVPN_USER", "USER")
	if *user == "" {
		return nil
	}

	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	*user = strings.ReplaceAll(*user, " ", "")
	return user
}

func (s *Source) readOpenVPNPassword() (password *string) {
	password = new(string)
	_, *password = s.getEnvWithRetro("OPENVPN_PASSWORD", "PASSWORD")
	if *password == "" {
		return nil
	}

	return password
}

func (s *Source) readOpenVPNKeyPassphrase() (passphrase *string) {
	passphrase = new(string)
	*passphrase = getCleanedEnv("OPENVPN_KEY_PASSPHRASE")
	if *passphrase == "" {
		return nil
	}
	return passphrase
}

func (s *Source) readPIAEncryptionPreset() (presetPtr *string) {
	_, preset := s.getEnvWithRetro(
		"PRIVATE_INTERNET_ACCESS_OPENVPN_ENCRYPTION_PRESET",
		"PIA_ENCRYPTION", "ENCRYPTION")
	if preset != "" {
		return &preset
	}
	return nil
}

func (s *Source) readOpenVPNProcessUser() (processUser string, err error) {
	key, value := s.getEnvWithRetro("OPENVPN_PROCESS_USER", "OPENVPN_ROOT")
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
	if root {
		return "root", nil
	}
	const defaultNonRootUser = "nonrootuser"
	return defaultNonRootUser, nil
}
