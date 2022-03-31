package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readOpenVPN() (
	openVPN settings.OpenVPN, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"OPENVPN_CLIENTKEY", "OPENVPN_CLIENTCRT"}, err)
	}()

	openVPN.Version = os.Getenv("OPENVPN_VERSION")
	openVPN.User = r.readOpenVPNUser()
	openVPN.Password = r.readOpenVPNPassword()
	confFile := os.Getenv("OPENVPN_CUSTOM_CONFIG")
	if confFile != "" {
		openVPN.ConfFile = &confFile
	}

	openVPN.Ciphers = envToCSV("OPENVPN_CIPHER")
	auth := os.Getenv("OPENVPN_AUTH")
	if auth != "" {
		openVPN.Auth = &auth
	}

	openVPN.ClientCrt, err = readBase64OrNil("OPENVPN_CLIENTCRT")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_CLIENTCRT: %w", err)
	}

	openVPN.ClientKey, err = readBase64OrNil("OPENVPN_CLIENTKEY")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_CLIENTKEY: %w", err)
	}

	openVPN.PIAEncPreset = r.readPIAEncryptionPreset()

	openVPN.IPv6, err = envToBoolPtr("OPENVPN_IPV6")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_IPV6: %w", err)
	}

	openVPN.MSSFix, err = envToUint16Ptr("OPENVPN_MSSFIX")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_MSSFIX: %w", err)
	}

	openVPN.Interface = os.Getenv("OPENVPN_INTERFACE")

	openVPN.Root, err = envToBoolPtr("OPENVPN_ROOT")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_ROOT: %w", err)
	}

	// TODO ProcUser once Root is deprecated.

	openVPN.Verbosity, err = envToIntPtr("OPENVPN_VERBOSITY")
	if err != nil {
		return openVPN, fmt.Errorf("environment variable OPENVPN_VERBOSITY: %w", err)
	}

	flagsStr := os.Getenv("OPENVPN_FLAGS")
	if flagsStr != "" {
		openVPN.Flags = strings.Fields(flagsStr)
	}

	return openVPN, nil
}

func (r *Reader) readOpenVPNUser() (user string) {
	user = os.Getenv("OPENVPN_USER")
	if user == "" {
		// Retro-compatibility
		user = os.Getenv("USER")
		if user != "" {
			r.onRetroActive("USER", "OPENVPN_USER")
		}
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	return strings.ReplaceAll(user, " ", "")
}

func (r *Reader) readOpenVPNPassword() (password string) {
	password = os.Getenv("OPENVPN_PASSWORD")
	if password != "" {
		return password
	}

	// Retro-compatibility
	password = os.Getenv("PASSWORD")
	if password != "" {
		r.onRetroActive("PASSWORD", "OPENVPN_PASSWORD")
	}
	return password
}

func readBase64OrNil(envKey string) (valueOrNil *string, err error) {
	value := os.Getenv(envKey)
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	decoded, err := decodeBase64(value)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

func (r *Reader) readPIAEncryptionPreset() (presetPtr *string) {
	preset := strings.ToLower(os.Getenv("PIA_ENCRYPTION"))
	if preset != "" {
		return &preset
	}

	// Retro-compatibility
	preset = strings.ToLower(os.Getenv("ENCRYPTION"))
	if preset != "" {
		r.onRetroActive("ENCRYPTION", "PIA_ENCRYPTION")
		return &preset
	}

	return nil
}
