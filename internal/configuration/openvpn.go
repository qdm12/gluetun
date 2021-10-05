package configuration

import (
	"errors"
	"fmt"
	"regexp"
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
	Ciphers   []string `json:"ciphers"`
	Auth      string   `json:"auth"`
	ConfFile  string   `json:"conf_file"`
	Version   string   `json:"version"`
	ClientCrt string   `json:"-"`                 // Cyberghost
	ClientKey string   `json:"-"`                 // Cyberghost, VPNUnlimited
	EncPreset string   `json:"encryption_preset"` // PIA
	IPv6      bool     `json:"ipv6"`              // Mullvad
	ProcUser  string   `json:"procuser"`          // Process username
	Interface string   `json:"interface"`
}

func (settings *OpenVPN) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *OpenVPN) lines() (lines []string) {
	lines = append(lines, lastIndent+"OpenVPN:")

	lines = append(lines, indent+lastIndent+"Version: "+settings.Version)

	lines = append(lines, indent+lastIndent+"Verbosity level: "+strconv.Itoa(settings.Verbosity))

	lines = append(lines, indent+lastIndent+"Network interface: "+settings.Interface)

	if len(settings.Flags) > 0 {
		lines = append(lines, indent+lastIndent+"Flags: "+strings.Join(settings.Flags, " "))
	}

	if settings.Root {
		lines = append(lines, indent+lastIndent+"Run as root: enabled")
	}

	if len(settings.Ciphers) > 0 {
		lines = append(lines, indent+lastIndent+"Custom ciphers: "+commaJoin(settings.Ciphers))
	}
	if len(settings.Auth) > 0 {
		lines = append(lines, indent+lastIndent+"Custom auth algorithm: "+settings.Auth)
	}

	if settings.ConfFile != "" {
		lines = append(lines, indent+lastIndent+"Configuration file: "+settings.ConfFile)
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
	credentialsRequired := false
	switch serviceProvider {
	case constants.Custom:
	case constants.VPNUnlimited:
	default:
		credentialsRequired = true
	}

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

	settings.Root, err = r.env.YesNo("OPENVPN_ROOT", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_ROOT: %w", err)
	}

	settings.Ciphers, err = r.env.CSV("OPENVPN_CIPHER")
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

	settings.Interface, err = readInterface(r.env)
	if err != nil {
		return err
	}

	switch serviceProvider {
	case constants.Custom:
		err = settings.readCustom(r) // read OPENVPN_CUSTOM_CONFIG
	case constants.Cyberghost:
		err = settings.readCyberghost(r)
	case constants.PrivateInternetAccess:
		settings.EncPreset, err = getPIAEncryptionPreset(r)
	case constants.VPNUnlimited:
		err = settings.readVPNUnlimited(r)
	case constants.Wevpn:
		err = settings.readWevpn(r)
	}
	if err != nil {
		return err
	}

	return nil
}

func readOpenVPNProtocol(r reader) (tcp bool, err error) {
	protocol, err := r.env.Inside("OPENVPN_PROTOCOL", []string{constants.TCP, constants.UDP},
		params.Default(constants.UDP), params.RetroKeys([]string{"PROTOCOL"}, r.onRetroActive))
	if err != nil {
		return false, fmt.Errorf("environment variable OPENVPN_PROTOCOL: %w", err)
	}
	return protocol == constants.TCP, nil
}

const openvpnIntfRegexString = `^.*[0-9]$`

var openvpnIntfRegexp = regexp.MustCompile(openvpnIntfRegexString)
var errInterfaceNameNotValid = errors.New("interface name is not valid")

func readInterface(env params.Interface) (intf string, err error) {
	intf, err = env.Get("OPENVPN_INTERFACE", params.Default("tun0"))
	if err != nil {
		return "", fmt.Errorf("environment variable OPENVPN_INTERFACE: %w", err)
	}

	if !openvpnIntfRegexp.MatchString(intf) {
		return "", fmt.Errorf("%w: does not match regex %s: %s",
			errInterfaceNameNotValid, openvpnIntfRegexString, intf)
	}

	return intf, nil
}
