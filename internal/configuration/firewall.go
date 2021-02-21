package configuration

import (
	"net"
	"strings"

	"github.com/qdm12/golibs/params"
)

// Firewall contains settings to customize the firewall operation.
type Firewall struct {
	VPNInputPorts   []uint16
	InputPorts      []uint16
	OutboundSubnets []net.IPNet
	Enabled         bool
	Debug           bool
}

func (settings *Firewall) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Firewall) lines() (lines []string) {
	if !settings.Enabled {
		lines = append(lines, lastIndent+"Firewall: disabled ⚠️")
		return lines
	}

	lines = append(lines, lastIndent+"Firewall:")

	if settings.Debug {
		lines = append(lines, indent+lastIndent+"Debug: on")
	}

	if len(settings.VPNInputPorts) > 0 {
		lines = append(lines, indent+lastIndent+"VPN input ports: "+
			strings.Join(uint16sToStrings(settings.VPNInputPorts), ", "))
	}

	if len(settings.InputPorts) > 0 {
		lines = append(lines, indent+lastIndent+"Input ports: "+
			strings.Join(uint16sToStrings(settings.InputPorts), ", "))
	}

	if len(settings.OutboundSubnets) > 0 {
		lines = append(lines, indent+lastIndent+"Outbound subnets: "+
			strings.Join(ipNetsToStrings(settings.OutboundSubnets), ", "))
	}

	return lines
}

func (settings *Firewall) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("FIREWALL", params.Default("on"))
	if err != nil {
		return err
	}

	settings.Debug, err = r.env.OnOff("FIREWALL_DEBUG", params.Default("off"))
	if err != nil {
		return err
	}

	if err := settings.readVPNInputPorts(r.env); err != nil {
		return err
	}

	if err := settings.readInputPorts(r.env); err != nil {
		return err
	}

	if err := settings.readOutboundSubnets(r); err != nil {
		return err
	}

	return nil
}

func (settings *Firewall) readVPNInputPorts(env params.Env) (err error) {
	settings.VPNInputPorts, err = readCSVPorts(env, "FIREWALL_VPN_INPUT_PORTS")
	return err
}

func (settings *Firewall) readInputPorts(env params.Env) (err error) {
	settings.InputPorts, err = readCSVPorts(env, "FIREWALL_INPUT_PORTS")
	return err
}

func (settings *Firewall) readOutboundSubnets(r reader) (err error) {
	retroOption := params.RetroKeys([]string{"EXTRA_SUBNETS"}, r.onRetroActive)
	settings.OutboundSubnets, err = readCSVIPNets(r.env, "FIREWALL_OUTBOUND_SUBNETS", retroOption)
	return err
}
