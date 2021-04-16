package configuration

import (
	"errors"
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
	MSSFix    uint16   `json:"mssfix"`
	Root      bool     `json:"run_as_root"`
	Cipher    string   `json:"cipher"`
	Auth      string   `json:"auth"`
	Provider  Provider `json:"provider"`
	Config    string   `json:"custom_config"`
}

func (settings *OpenVPN) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *OpenVPN) lines() (lines []string) {
	lines = append(lines, lastIndent+"OpenVPN:")

	lines = append(lines, indent+lastIndent+"Verbosity level: "+strconv.Itoa(settings.Verbosity))

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

	lines = append(lines, indent+lastIndent+"Provider:")
	for _, line := range settings.Provider.lines() {
		lines = append(lines, indent+indent+line)
	}

	return lines
}

var (
	ErrInvalidVPNProvider = errors.New("invalid VPN provider")
)

func (settings *OpenVPN) read(r reader) (err error) {
	vpnsp, err := r.env.Inside("VPNSP", []string{
		"cyberghost", "fastestvpn", "hidemyass", "mullvad", "nordvpn",
		"privado", "pia", "private internet access", "privatevpn", "protonvpn",
		"purevpn", "surfshark", "torguard", "vyprvpn", "windscribe"},
		params.Default("private internet access"))
	if err != nil {
		return err
	}
	if vpnsp == "pia" { // retro compatibility
		vpnsp = "private internet access"
	}

	settings.Provider.Name = vpnsp

	settings.Config, err = r.env.Get("OPENVPN_CUSTOM_CONFIG", params.CaseSensitiveValue())
	if err != nil {
		return err
	}

	credentialsRequired := len(settings.Config) == 0

	settings.User, err = r.getFromEnvOrSecretFile("OPENVPN_USER", credentialsRequired, []string{"USER"})
	if err != nil {
		return err
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")

	if settings.Provider.Name == constants.Mullvad {
		settings.Password = "m"
	} else {
		settings.Password, err = r.getFromEnvOrSecretFile("OPENVPN_PASSWORD", credentialsRequired, []string{"PASSWORD"})
		if err != nil {
			return err
		}
	}

	settings.Verbosity, err = r.env.IntRange("OPENVPN_VERBOSITY", 0, 6, params.Default("1"))
	if err != nil {
		return err
	}

	settings.Root, err = r.env.YesNo("OPENVPN_ROOT", params.Default("yes"))
	if err != nil {
		return err
	}

	settings.Cipher, err = r.env.Get("OPENVPN_CIPHER")
	if err != nil {
		return err
	}

	settings.Auth, err = r.env.Get("OPENVPN_AUTH")
	if err != nil {
		return err
	}

	mssFix, err := r.env.IntRange("OPENVPN_MSSFIX", 0, 10000, params.Default("0"))
	if err != nil {
		return err
	}
	settings.MSSFix = uint16(mssFix)

	var readProvider func(r reader) error
	switch settings.Provider.Name {
	case constants.Cyberghost:
		readProvider = settings.Provider.readCyberghost
	case constants.Fastestvpn:
		readProvider = settings.Provider.readFastestvpn
	case constants.HideMyAss:
		readProvider = settings.Provider.readHideMyAss
	case constants.Mullvad:
		readProvider = settings.Provider.readMullvad
	case constants.Nordvpn:
		readProvider = settings.Provider.readNordvpn
	case constants.Privado:
		readProvider = settings.Provider.readPrivado
	case constants.PrivateInternetAccess:
		readProvider = settings.Provider.readPrivateInternetAccess
	case constants.Privatevpn:
		readProvider = settings.Provider.readPrivatevpn
	case constants.Protonvpn:
		readProvider = settings.Provider.readProtonvpn
	case constants.Purevpn:
		readProvider = settings.Provider.readPurevpn
	case constants.Surfshark:
		readProvider = settings.Provider.readSurfshark
	case constants.Torguard:
		readProvider = settings.Provider.readTorguard
	case constants.Vyprvpn:
		readProvider = settings.Provider.readVyprvpn
	case constants.Windscribe:
		readProvider = settings.Provider.readWindscribe
	default:
		return fmt.Errorf("%w: %s", ErrInvalidVPNProvider, settings.Provider.Name)
	}

	return readProvider(r)
}
