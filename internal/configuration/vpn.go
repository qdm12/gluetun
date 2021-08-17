package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

type VPN struct {
	Type      string    `json:"type"`
	OpenVPN   OpenVPN   `json:"openvpn"`
	Wireguard Wireguard `json:"wireguard"`
	Provider  Provider  `json:"provider"`
}

func (settings *VPN) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *VPN) lines() (lines []string) {
	lines = append(lines, lastIndent+"VPN:")

	lines = append(lines, indent+lastIndent+"Type: "+settings.Type)

	var vpnLines []string
	switch settings.Type {
	case constants.OpenVPN:
		vpnLines = settings.OpenVPN.lines()
	case constants.Wireguard:
		vpnLines = settings.Wireguard.lines()
	}
	for _, line := range vpnLines {
		lines = append(lines, indent+line)
	}

	for _, line := range settings.Provider.lines() {
		lines = append(lines, indent+line)
	}

	return lines
}

var (
	errReadProviderSettings  = errors.New("cannot read provider settings")
	errReadOpenVPNSettings   = errors.New("cannot read OpenVPN settings")
	errReadWireguardSettings = errors.New("cannot read Wireguard settings")
)

func (settings *VPN) read(r reader) (err error) {
	vpnType, err := r.env.Inside("VPN_TYPE",
		[]string{constants.OpenVPN, constants.Wireguard},
		params.Default(constants.OpenVPN))
	if err != nil {
		return fmt.Errorf("environment variable VPN_TYPE: %w", err)
	}
	settings.Type = vpnType

	if !settings.isOpenVPNCustomConfig(r.env) {
		if err := settings.Provider.read(r, vpnType); err != nil {
			return fmt.Errorf("%w: %s", errReadProviderSettings, err)
		}
	}

	switch settings.Type {
	case constants.OpenVPN:
		err = settings.OpenVPN.read(r, settings.Provider.Name)
		if err != nil {
			return fmt.Errorf("%w: %s", errReadOpenVPNSettings, err)
		}
	case constants.Wireguard:
		err = settings.Wireguard.read(r)
		if err != nil {
			return fmt.Errorf("%w: %s", errReadWireguardSettings, err)
		}
	}

	return nil
}

func (settings VPN) isOpenVPNCustomConfig(env params.Env) (ok bool) {
	if settings.Type != constants.OpenVPN {
		return false
	}
	s, err := env.Get("OPENVPN_CUSTOM_CONFIG")
	return err == nil && s != ""
}
