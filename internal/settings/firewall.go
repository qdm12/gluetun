package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// Firewall contains settings to customize the firewall operation.
type Firewall struct {
	VPNInputPorts []uint16
	InputPorts    []uint16
	Enabled       bool
	Debug         bool
}

func (f *Firewall) String() string {
	if !f.Enabled {
		return "Firewall settings: disabled"
	}
	vpnInputPorts := make([]string, len(f.VPNInputPorts))
	for i, port := range f.VPNInputPorts {
		vpnInputPorts[i] = fmt.Sprintf("%d", port)
	}
	inputPorts := make([]string, len(f.InputPorts))
	for i, port := range f.InputPorts {
		inputPorts[i] = fmt.Sprintf("%d", port)
	}

	settingsList := []string{
		"Firewall settings:",
		"VPN input ports: " + strings.Join(vpnInputPorts, ", "),
		"Input ports: " + strings.Join(inputPorts, ", "),
	}
	if f.Debug {
		settingsList = append(settingsList, "Debug: on")
	}
	return strings.Join(settingsList, "\n |--")
}

// GetFirewallSettings obtains firewall settings from environment variables using the params package.
func GetFirewallSettings(paramsReader params.Reader) (settings Firewall, err error) {
	settings.VPNInputPorts, err = paramsReader.GetVPNInputPorts()
	if err != nil {
		return settings, err
	}
	settings.InputPorts, err = paramsReader.GetInputPorts()
	if err != nil {
		return settings, err
	}
	settings.Enabled, err = paramsReader.GetFirewall()
	if err != nil {
		return settings, err
	}
	settings.Debug, err = paramsReader.GetFirewallDebug()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
