package settings

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

// Firewall contains settings to customize the firewall operation.
type Firewall struct {
	VPNInputPorts   []uint16     `json:"vpn_input_ports,omitempty"`
	InputPorts      []uint16     `json:"input_ports,omitempty"`
	OutboundSubnets []*net.IPNet `json:"outbount_subnets,omitempty"`
	Enabled         *bool        `json:"enabled,omitempty"`
	Debug           *bool        `json:"debug,omitempty"`
}

func (f Firewall) validate() (err error) {
	if hasZeroPort(f.VPNInputPorts) {
		return fmt.Errorf("VPN input ports: %w", ErrFirewallZeroPort)
	}

	if hasZeroPort(f.InputPorts) {
		return fmt.Errorf("input ports: %w", ErrFirewallZeroPort)
	}

	return nil
}

func hasZeroPort(ports []uint16) (has bool) {
	for _, port := range ports {
		if port == 0 {
			return true
		}
	}
	return false
}

func (f *Firewall) copy() (copied Firewall) {
	return Firewall{
		VPNInputPorts:   helpers.CopyUint16Slice(f.VPNInputPorts),
		InputPorts:      helpers.CopyUint16Slice(f.InputPorts),
		OutboundSubnets: helpers.CopyIPNetSlice(f.OutboundSubnets),
		Enabled:         helpers.CopyBoolPtr(f.Enabled),
		Debug:           helpers.CopyBoolPtr(f.Debug),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
// It merges values of slices together, even if they
// are set in the receiver settings.
func (f *Firewall) mergeWith(other Firewall) {
	f.VPNInputPorts = helpers.MergeUint16Slices(f.VPNInputPorts, other.VPNInputPorts)
	f.InputPorts = helpers.MergeUint16Slices(f.InputPorts, other.InputPorts)
	f.OutboundSubnets = helpers.MergeIPNetsSlices(f.OutboundSubnets, other.OutboundSubnets)
	f.Enabled = helpers.MergeWithBool(f.Enabled, other.Enabled)
	f.Debug = helpers.MergeWithBool(f.Debug, other.Debug)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (f *Firewall) overrideWith(other Firewall) {
	f.VPNInputPorts = helpers.OverrideWithUint16Slice(f.VPNInputPorts, other.VPNInputPorts)
	f.InputPorts = helpers.OverrideWithUint16Slice(f.InputPorts, other.InputPorts)
	f.OutboundSubnets = helpers.OverrideWithIPNetsSlice(f.OutboundSubnets, other.OutboundSubnets)
	f.Enabled = helpers.OverrideWithBool(f.Enabled, other.Enabled)
	f.Debug = helpers.OverrideWithBool(f.Debug, other.Debug)
}

func (f *Firewall) setDefaults() {
	f.Enabled = helpers.DefaultBool(f.Enabled, true)
	f.Debug = helpers.DefaultBool(f.Debug, false)
}
