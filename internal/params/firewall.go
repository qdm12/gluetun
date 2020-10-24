package params

import (
	"fmt"
	"strconv"
	"strings"

	libparams "github.com/qdm12/golibs/params"
)

// GetFirewall obtains if the firewall should be enabled from the environment variable FIREWALL.
func (r *reader) GetFirewall() (enabled bool, err error) {
	return r.envParams.GetOnOff("FIREWALL", libparams.Default("on"))
}

// GetAllowedVPNInputPorts obtains a list of input ports to allow from the
// VPN server side in the firewall, from the environment variable FIREWALL_VPN_INPUT_PORTS.
func (r *reader) GetVPNInputPorts() (ports []uint16, err error) {
	s, err := r.envParams.GetEnv("FIREWALL_VPN_INPUT_PORTS", libparams.Default(""))
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}
	portsStr := strings.Split(s, ",")
	ports = make([]uint16, len(portsStr))
	for i := range portsStr {
		portInt, err := strconv.Atoi(portsStr[i])
		if err != nil {
			return nil, fmt.Errorf("VPN input port %q is not valid (%s)", portInt, err)
		} else if portInt <= 0 || portInt > 65535 {
			return nil, fmt.Errorf("VPN input port %d must be between 1 and 65535", portInt)
		}
		ports[i] = uint16(portInt)
	}
	return ports, nil
}

// GetInputPorts obtains a list of input ports to allow through the
// default interface in the firewall, from the environment variable FIREWALL_INPUT_PORTS.
func (r *reader) GetInputPorts() (ports []uint16, err error) {
	s, err := r.envParams.GetEnv("FIREWALL_INPUT_PORTS", libparams.Default(""))
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}
	portsStr := strings.Split(s, ",")
	ports = make([]uint16, len(portsStr))
	for i := range portsStr {
		portInt, err := strconv.Atoi(portsStr[i])
		if err != nil {
			return nil, fmt.Errorf("Input port %q is not valid (%s)", portInt, err)
		} else if portInt <= 0 || portInt > 65535 {
			return nil, fmt.Errorf("Input port %d must be between 1 and 65535", portInt)
		}
		ports[i] = uint16(portInt)
	}
	return ports, nil
}

// GetFirewallDebug obtains if the firewall should run in debug verbose mode
// from the environment variable FIREWALL_DEBUG.
func (r *reader) GetFirewallDebug() (debug bool, err error) {
	return r.envParams.GetOnOff("FIREWALL_DEBUG", libparams.Default("off"))
}
