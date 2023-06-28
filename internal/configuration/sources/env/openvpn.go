package env

import (
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readOpenVPN() (
	openVPN settings.OpenVPN, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"OPENVPN_KEY", "OPENVPN_CERT",
			"OPENVPN_KEY_PASSPHRASE", "OPENVPN_ENCRYPTED_KEY"}, err)
	}()

	openVPN.Version = s.env.String("OPENVPN_VERSION")
	openVPN.User = s.env.Get("OPENVPN_USER",
		env.RetroKeys("USER"), env.ForceLowercase(false))
	openVPN.Password = s.env.Get("OPENVPN_PASSWORD",
		env.RetroKeys("PASSWORD"), env.ForceLowercase(false))
	openVPN.ConfFile = s.env.Get("OPENVPN_CUSTOM_CONFIG", env.ForceLowercase(false))
	openVPN.Ciphers = s.env.CSV("OPENVPN_CIPHERS", env.RetroKeys("OPENVPN_CIPHER"))
	openVPN.Auth = s.env.Get("OPENVPN_AUTH")
	openVPN.Cert = s.env.Get("OPENVPN_CERT", env.ForceLowercase(false))
	openVPN.Key = s.env.Get("OPENVPN_KEY", env.ForceLowercase(false))
	openVPN.EncryptedKey = s.env.Get("OPENVPN_ENCRYPTED_KEY", env.ForceLowercase(false))
	openVPN.KeyPassphrase = s.env.Get("OPENVPN_KEY_PASSPHRASE", env.ForceLowercase(false))

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
