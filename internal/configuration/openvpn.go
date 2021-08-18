package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

// OpenVPN contains settings to configure the OpenVPN client.
type OpenVPN struct {
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Verbosity int      `json:"verbosity"`
	Flags     []string `json:"flags"`
	MSSFix    uint16   `json:"mssfix"`
	Root      bool     `json:"run_as_root"`
	Cipher    string   `json:"cipher"`
	Auth      string   `json:"auth"`
	Config    string   `json:"custom_config"`
	Version   string   `json:"version"`
	ClientCrt string   `json:"-"`                 // Cyberghost
	ClientKey string   `json:"-"`                 // Cyberghost, VPNUnlimited
	EncPreset string   `json:"encryption_preset"` // PIA
	IPv6      bool     `json:"ipv6"`              // Mullvad
	ProcUser  string   `json:"procuser"`          // Process username
}

func (settings *OpenVPN) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *OpenVPN) lines() (lines []string) {
	lines = append(lines, lastIndent+"OpenVPN:")

	lines = append(lines, indent+lastIndent+"Version: "+settings.Version)

	lines = append(lines, indent+lastIndent+"Verbosity level: "+strconv.Itoa(settings.Verbosity))

	if len(settings.Flags) > 0 {
		lines = append(lines, indent+lastIndent+"Flags: "+strings.Join(settings.Flags, " "))
	}

	if settings.Root {
		lines = append(lines, indent+lastIndent+"Run as root: enabled")
	}

	if len(settings.Cipher) > 0 {
		lines = append(lines, indent+lastIndent+"Custom cipher: "+settings.Cipher)
	}
	if len(settings.Auth) > 0 {
		lines = append(lines, indent+lastIndent+"Custom auth algorithm: "+settings.Auth)
	}

	if len(settings.Config) > 0 {
		lines = append(lines, indent+lastIndent+"Custom configuration: "+settings.Config)
	}

	if settings.ClientKey != "" {
		lines = append(lines, indent+lastIndent+"Client key is set")
	}

	if settings.ClientCrt != "" {
		lines = append(lines, indent+lastIndent+"Client certificate is set")
	}

	if settings.IPv6 {
		lines = append(lines, indent+lastIndent+"IPv6: enabled")
	}

	if settings.EncPreset != "" { // PIA only
		lines = append(lines, indent+lastIndent+"Encryption preset: "+settings.EncPreset)
	}

	return lines
}

func (settings *OpenVPN) read(r reader, serviceProvider string) (err error) {
	settings.Config, err = r.env.Get("OPENVPN_CUSTOM_CONFIG", params.CaseSensitiveValue())
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_CUSTOM_CONFIG: %w", err)
	}

	credentialsRequired := settings.Config == "" && serviceProvider != constants.VPNUnlimited

	settings.User, err = r.getFromEnvOrSecretFile("OPENVPN_USER", credentialsRequired, []string{"USER"})
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_USER: %w", err)
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")

	if serviceProvider == constants.Mullvad {
		settings.Password = "m"
	} else {
		settings.Password, err = r.getFromEnvOrSecretFile("OPENVPN_PASSWORD", credentialsRequired, []string{"PASSWORD"})
		if err != nil {
			return err
		}
	}

	settings.Version, err = r.env.Inside("OPENVPN_VERSION",
		[]string{constants.Openvpn24, constants.Openvpn25}, params.Default(constants.Openvpn25))
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_VERSION: %w", err)
	}

	settings.Verbosity, err = r.env.IntRange("OPENVPN_VERBOSITY", 0, 6, params.Default("1")) //nolint:gomnd
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_VERBOSITY: %w", err)
	}

	settings.Flags = []string{}
	flagsStr, err := r.env.Get("OPENVPN_FLAGS")
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_FLAGS: %w", err)
	}
	if flagsStr != "" {
		settings.Flags = strings.Fields(flagsStr)
	}

	settings.Root, err = r.env.YesNo("OPENVPN_ROOT", params.Default("yes"))
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_ROOT: %w", err)
	}

	settings.Cipher, err = r.env.Get("OPENVPN_CIPHER")
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_CIPHER: %w", err)
	}

	settings.Auth, err = r.env.Get("OPENVPN_AUTH")
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_AUTH: %w", err)
	}

	const maxMSSFix = 10000
	mssFix, err := r.env.IntRange("OPENVPN_MSSFIX", 0, maxMSSFix, params.Default("0"))
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_MSSFIX: %w", err)
	}
	settings.MSSFix = uint16(mssFix)

	settings.IPv6, err = r.env.OnOff("OPENVPN_IPV6", params.Default("off"))
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_IPV6: %w", err)
	}

	settings.EncPreset, err = getPIAEncryptionPreset(r)
	if err != nil {
		return err
	}

	switch serviceProvider {
	case constants.Cyberghost:
		err = settings.readCyberghost(r)
	case constants.VPNUnlimited:
		err = settings.readVPNUnlimited(r)
	}
	if err != nil {
		return err
	}

	return nil
}

func readProtocol(env params.Env) (tcp bool, err error) {
	protocol, err := env.Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, params.Default(constants.UDP))
	if err != nil {
		return false, fmt.Errorf("environment variable PROTOCOL: %w", err)
	}
	return protocol == constants.TCP, nil
}
