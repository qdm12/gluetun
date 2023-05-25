package settings

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// Firewall contains settings to customize the firewall operation.
type Firewall struct {
	VPNInputPorts   []uint16
	InputPorts      []uint16
	OutboundSubnets []netip.Prefix
	Enabled         *bool
	Debug           *bool
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
		VPNInputPorts:   gosettings.CopySlice(f.VPNInputPorts),
		InputPorts:      gosettings.CopySlice(f.InputPorts),
		OutboundSubnets: gosettings.CopySlice(f.OutboundSubnets),
		Enabled:         gosettings.CopyPointer(f.Enabled),
		Debug:           gosettings.CopyPointer(f.Debug),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
// It merges values of slices together, even if they
// are set in the receiver settings.
func (f *Firewall) mergeWith(other Firewall) {
	f.VPNInputPorts = gosettings.MergeWithSlice(f.VPNInputPorts, other.VPNInputPorts)
	f.InputPorts = gosettings.MergeWithSlice(f.InputPorts, other.InputPorts)
	f.OutboundSubnets = gosettings.MergeWithSlice(f.OutboundSubnets, other.OutboundSubnets)
	f.Enabled = gosettings.MergeWithPointer(f.Enabled, other.Enabled)
	f.Debug = gosettings.MergeWithPointer(f.Debug, other.Debug)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (f *Firewall) overrideWith(other Firewall) {
	f.VPNInputPorts = gosettings.OverrideWithSlice(f.VPNInputPorts, other.VPNInputPorts)
	f.InputPorts = gosettings.OverrideWithSlice(f.InputPorts, other.InputPorts)
	f.OutboundSubnets = gosettings.OverrideWithSlice(f.OutboundSubnets, other.OutboundSubnets)
	f.Enabled = gosettings.OverrideWithPointer(f.Enabled, other.Enabled)
	f.Debug = gosettings.OverrideWithPointer(f.Debug, other.Debug)
}

func (f *Firewall) setDefaults() {
	f.Enabled = gosettings.DefaultPointer(f.Enabled, true)
	f.Debug = gosettings.DefaultPointer(f.Debug, false)
}

func (f Firewall) String() string {
	return f.toLinesNode().String()
}

func (f Firewall) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Firewall settings:")

	node.Appendf("Enabled: %s", helpers.BoolPtrToYesNo(f.Enabled))
	if !*f.Enabled {
		return node
	}

	if *f.Debug {
		node.Appendf("Debug mode: on")
	}

	if len(f.VPNInputPorts) > 0 {
		vpnInputPortsNode := node.Appendf("VPN input ports:")
		for _, port := range f.VPNInputPorts {
			vpnInputPortsNode.Appendf("%d", port)
		}
	}

	if len(f.InputPorts) > 0 {
		inputPortsNode := node.Appendf("Input ports:")
		for _, port := range f.InputPorts {
			inputPortsNode.Appendf("%d", port)
		}
	}

	if len(f.OutboundSubnets) > 0 {
		outboundSubnets := node.Appendf("Outbound subnets:")
		for _, subnet := range f.OutboundSubnets {
			subnet := subnet
			outboundSubnets.Appendf("%s", &subnet)
		}
	}

	return node
}
