package settings

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// Firewall contains settings to customize the firewall operation.
type Firewall struct {
	VPNInputPorts   []uint16
	InputPorts      []uint16
	OutboundSubnets []net.IPNet
	Enabled         bool
	Debug           bool
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
	outboundSubnets := make([]string, len(f.OutboundSubnets))
	for i := range f.OutboundSubnets {
		outboundSubnets[i] = f.OutboundSubnets[i].String()
	}

	settingsList := []string{
		"Firewall settings:",
		"VPN input ports: " + strings.Join(vpnInputPorts, ", "),
		"Input ports: " + strings.Join(inputPorts, ", "),
		"Outbound subnets: " + strings.Join(outboundSubnets, ", "),
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
	settings.OutboundSubnets, err = paramsReader.GetOutboundSubnets()
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
