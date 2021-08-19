package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

type VPN struct {
	Type     string   `json:"type"`
	OpenVPN  OpenVPN  `json:"openvpn"`
	Provider Provider `json:"provider"`
}

func (settings *VPN) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *VPN) lines() (lines []string) {
	lines = append(lines, lastIndent+"VPN:")

	lines = append(lines, indent+lastIndent+"Type: "+settings.Type)

	for _, line := range settings.OpenVPN.lines() {
		lines = append(lines, indent+line)
	}

	for _, line := range settings.Provider.lines() {
		lines = append(lines, indent+line)
	}

	return lines
}

var (
	errReadProviderSettings = errors.New("cannot read provider settings")
	errReadOpenVPNSettings  = errors.New("cannot read OpenVPN settings")
)

func (settings *VPN) read(r reader) (err error) {
	vpnType, err := r.env.Inside("VPN_TYPE",
		[]string{constants.OpenVPN}, params.Default(constants.OpenVPN))
	if err != nil {
		return fmt.Errorf("environment variable VPN_TYPE: %w", err)
	}
	settings.Type = vpnType

	if !settings.isOpenVPNCustomConfig(r.env) {
		if err := settings.Provider.read(r, vpnType); err != nil {
			return fmt.Errorf("%w: %s", errReadProviderSettings, err)
		}
	}

	err = settings.OpenVPN.read(r, settings.Provider.Name)
	if err != nil {
		return fmt.Errorf("%w: %s", errReadOpenVPNSettings, err)
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
