package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readFirewall() (firewall settings.Firewall, err error) {
	firewall.VPNInputPorts, err = s.env.CSVUint16("FIREWALL_VPN_INPUT_PORTS")
	if err != nil {
		return firewall, err
	}

	firewall.InputPorts, err = s.env.CSVUint16("FIREWALL_INPUT_PORTS")
	if err != nil {
		return firewall, err
	}

	firewall.OutboundSubnets, err = s.env.CSVNetipPrefixes("FIREWALL_OUTBOUND_SUBNETS",
		env.RetroKeys("EXTRA_SUBNETS"))
	if err != nil {
		return firewall, err
	}

	firewall.Enabled, err = s.env.BoolPtr("FIREWALL")
	if err != nil {
		return firewall, err
	}

	firewall.Debug, err = s.env.BoolPtr("FIREWALL_DEBUG")
	if err != nil {
		return firewall, err
	}

	return firewall, nil
}
