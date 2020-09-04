package settings

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// Firewall contains settings to customize the firewall operation
type Firewall struct {
	AllowedSubnets []net.IPNet
	VPNInputPorts  []uint16
	Enabled        bool
	Debug          bool
}

func (f *Firewall) String() string {
	allowedSubnets := make([]string, len(f.AllowedSubnets))
	for i := range f.AllowedSubnets {
		allowedSubnets[i] = f.AllowedSubnets[i].String()
	}
	if !f.Enabled {
		return "Firewall settings: disabled"
	}
	vpnInputPorts := make([]string, len(f.VPNInputPorts))
	for i, port := range f.VPNInputPorts {
		vpnInputPorts[i] = fmt.Sprintf("%d", port)
	}

	settingsList := []string{
		"Firewall settings:",
		"Allowed subnets: " + strings.Join(allowedSubnets, ", "),
		"VPN input ports: " + strings.Join(vpnInputPorts, ", "),
	}
	if f.Debug {
		settingsList = append(settingsList, "Debug: on")
	}
	return strings.Join(settingsList, "\n |--")
}

// GetFirewallSettings obtains firewall settings from environment variables using the params package.
func GetFirewallSettings(paramsReader params.Reader) (settings Firewall, err error) {
	settings.AllowedSubnets, err = paramsReader.GetExtraSubnets()
	if err != nil {
		return settings, err
	}
	settings.VPNInputPorts, err = paramsReader.GetVPNInputPorts()
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
