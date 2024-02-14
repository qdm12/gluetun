package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readOpenVPN() (
	openVPN settings.OpenVPN, err error) {

	openVPN.Version = s.env.String("OPENVPN_VERSION")

	openVPN.User, err = s.readSecretFileAsStringPtr(
		"OPENVPN_USER",
		"/run/secrets/openvpn_user",
		[]string{"OPENVPN_USER_SECRETFILE", "USER"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading user file: %w", err)
	}

	openVPN.Key, err = s.readPEMSecretFile(
		"OPENVPN_CLIENTKEY",
		"/run/secrets/openvpn_clientkey",
		[]string{"OPENVPN_CLIENTKEY_SECRETFILE"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading client key file: %w", err)
	}

	openVPN.Password, err = s.readSecretFileAsStringPtr(
		"OPENVPN_PASSWORD",
		"/run/secrets/openvpn_password",
		[]string{"OPENVPN_PASSWORD_SECRETFILE", "PASSWORD"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading password file: %w", err)
	}

	openVPN.ConfFile = s.env.Get("OPENVPN_CUSTOM_CONFIG", env.ForceLowercase(false))
	openVPN.Ciphers = s.env.CSV("OPENVPN_CIPHERS", env.RetroKeys("OPENVPN_CIPHER"))
	openVPN.Auth = s.env.Get("OPENVPN_AUTH")
	openVPN.Cert = s.env.Get("OPENVPN_CERT", env.ForceLowercase(false))
	openVPN.Key = s.env.Get("OPENVPN_KEY", env.ForceLowercase(false))

	openVPN.EncryptedKey, err = s.readPEMSecretFile(
		"OPENVPN_ENCRYPTED_KEY",
		"/run/secrets/openvpn_encrypted_key",
		[]string{"OPENVPN_ENCRYPTED_KEY_SECRETFILE"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading encrypted key file: %w", err)
	}

	openVPN.KeyPassphrase, err = s.readSecretFileAsStringPtr(
		"OPENVPN_KEY_PASSPHRASE",
		"/run/secrets/openvpn_key_passphrase",
		[]string{"OPENVPN_KEY_PASSPHRASE_SECRETFILE"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading key passphrase file: %w", err)
	}

	openVPN.Cert, err = s.readPEMSecretFile(
		"OPENVPN_CLIENTCRT",
		"/run/secrets/openvpn_clientcrt",
		[]string{"OPENVPN_CLIENTCRT_SECRETFILE"},
	)
	if err != nil {
		return openVPN, fmt.Errorf("reading client certificate file: %w", err)
	}

	openVPN.PIAEncPreset = s.readPIAEncryptionPreset()

	openVPN.MSSFix, err = s.env.Uint16Ptr("OPENVPN_MSSFIX")
	if err != nil {
		return openVPN, err
	}

	openVPN.Interface = s.env.String("VPN_INTERFACE",
		env.RetroKeys("OPENVPN_INTERFACE"), env.ForceLowercase(false))

	openVPN.ProcessUser, err = s.readOpenVPNProcessUser()
	if err != nil {
		return openVPN, err
	}

	openVPN.Verbosity, err = s.env.IntPtr("OPENVPN_VERBOSITY")
	if err != nil {
		return openVPN, err
	}

	flagsPtr := s.env.Get("OPENVPN_FLAGS", env.ForceLowercase(false))
	if flagsPtr != nil {
		openVPN.Flags = strings.Fields(*flagsPtr)
	}

	return openVPN, nil
}

func (s *Source) readPIAEncryptionPreset() (presetPtr *string) {
	return s.env.Get(
		"PRIVATE_INTERNET_ACCESS_OPENVPN_ENCRYPTION_PRESET",
		env.RetroKeys("ENCRYPTION", "PIA_ENCRYPTION"))
}

func (s *Source) readOpenVPNProcessUser() (processUser string, err error) {
	value, err := s.env.BoolPtr("OPENVPN_ROOT") // Retro-compatibility
	if err != nil {
		return "", err
	} else if value != nil {
		if *value {
			return "root", nil
		}
		const defaultNonRootUser = "nonrootuser"
		return defaultNonRootUser, nil
	}

	return s.env.String("OPENVPN_PROCESS_USER"), nil
}
