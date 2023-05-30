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

	openVPN.Version = env.String("OPENVPN_VERSION")
	_, openVPN.User = s.getEnvWithRetro("OPENVPN_USER",
		[]string{"USER"}, env.ForceLowercase(false))
	_, openVPN.Password = s.getEnvWithRetro("OPENVPN_PASSWORD",
		[]string{"PASSWORD"}, env.ForceLowercase(false))
	openVPN.ConfFile = env.Get("OPENVPN_CUSTOM_CONFIG")

	ciphersKey, _ := s.getEnvWithRetro("OPENVPN_CIPHERS", []string{"OPENVPN_CIPHER"})
	openVPN.Ciphers = env.CSV(ciphersKey)

	openVPN.Auth = env.Get("OPENVPN_AUTH")
	openVPN.Cert = env.Get("OPENVPN_CERT", env.ForceLowercase(false))
	openVPN.Key = env.Get("OPENVPN_KEY", env.ForceLowercase(false))
	openVPN.EncryptedKey = env.Get("OPENVPN_ENCRYPTED_KEY", env.ForceLowercase(false))
	openVPN.KeyPassphrase = env.Get("OPENVPN_KEY_PASSPHRASE", env.ForceLowercase(false))

	openVPN.PIAEncPreset = s.readPIAEncryptionPreset()

	openVPN.MSSFix, err = env.Uint16Ptr("OPENVPN_MSSFIX")
	if err != nil {
		return openVPN, err
	}

	_, openvpnInterface := s.getEnvWithRetro("VPN_INTERFACE",
		[]string{"OPENVPN_INTERFACE"}, env.ForceLowercase(false))
	if openvpnInterface != nil {
		openVPN.Interface = *openvpnInterface
	}

	openVPN.ProcessUser, err = s.readOpenVPNProcessUser()
	if err != nil {
		return openVPN, err
	}

	openVPN.Verbosity, err = env.IntPtr("OPENVPN_VERBOSITY")
	if err != nil {
		return openVPN, err
	}

	flagsPtr := env.Get("OPENVPN_FLAGS", env.ForceLowercase(false))
	if flagsPtr != nil {
		openVPN.Flags = strings.Fields(*flagsPtr)
	}

	return openVPN, nil
}

func (s *Source) readPIAEncryptionPreset() (presetPtr *string) {
	_, presetPtr = s.getEnvWithRetro(
		"PRIVATE_INTERNET_ACCESS_OPENVPN_ENCRYPTION_PRESET",
		[]string{"PIA_ENCRYPTION", "ENCRYPTION"})
	return presetPtr
}

func (s *Source) readOpenVPNProcessUser() (processUser string, err error) {
	key, value := s.getEnvWithRetro("OPENVPN_PROCESS_USER",
		[]string{"OPENVPN_ROOT"})
	if value == nil {
		return "", nil
	} else if key == "OPENVPN_PROCESS_USER" {
		return *value, nil
	}

	// Retro-compatibility
	if *value == "" {
		return "", nil
	}
	root, err := binary.Validate(*value)
	if err != nil {
		return "", fmt.Errorf("environment variable %s: %w", key, err)
	}
	if *root {
		return "root", nil
	}
	const defaultNonRootUser = "nonrootuser"
	return defaultNonRootUser, nil
}
