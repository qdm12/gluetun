package settings

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// DNS contains settings to configure DNS.
type DNS struct {
	// ServerAddress is the DNS server to use inside
	// the Go program and for the system.
	// It defaults to '127.0.0.1' to be used with the
	// DoT server. It cannot be the zero value in the internal
	// state.
	ServerAddress netip.Addr
	// KeepNameserver is true if the existing DNS server
	// found in /etc/resolv.conf should be used
	// Note setting this to true will likely DNS traffic
	// outside the VPN tunnel since it would go through
	// the local DNS server of your Docker/Kubernetes
	// configuration, which is likely not going through the tunnel.
	// This will also disable the DNS over TLS server and the
	// `ServerAddress` field will be ignored.
	// It defaults to false and cannot be nil in the
	// internal state.
	KeepNameserver *bool
	// DOT contains settings to configure the DoT
	// server.
	DoT DoT
}

func (d DNS) validate() (err error) {
	err = d.DoT.validate()
	if err != nil {
		return fmt.Errorf("validating DoT settings: %w", err)
	}

	return nil
}

func (d *DNS) Copy() (copied DNS) {
	return DNS{
		ServerAddress:  d.ServerAddress,
		KeepNameserver: gosettings.CopyPointer(d.KeepNameserver),
		DoT:            d.DoT.copy(),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (d *DNS) mergeWith(other DNS) {
	d.ServerAddress = gosettings.MergeWithValidator(d.ServerAddress, other.ServerAddress)
	d.KeepNameserver = gosettings.MergeWithPointer(d.KeepNameserver, other.KeepNameserver)
	d.DoT.mergeWith(other.DoT)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DNS) overrideWith(other DNS) {
	d.ServerAddress = gosettings.OverrideWithValidator(d.ServerAddress, other.ServerAddress)
	d.KeepNameserver = gosettings.OverrideWithPointer(d.KeepNameserver, other.KeepNameserver)
	d.DoT.overrideWith(other.DoT)
}

func (d *DNS) setDefaults() {
	localhost := netip.AddrFrom4([4]byte{127, 0, 0, 1})
	d.ServerAddress = gosettings.DefaultValidator(d.ServerAddress, localhost)
	d.KeepNameserver = gosettings.DefaultPointer(d.KeepNameserver, false)
	d.DoT.setDefaults()
}

func (d DNS) String() string {
	return d.toLinesNode().String()
}

func (d DNS) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS settings:")
	node.Appendf("Keep existing nameserver(s): %s", gosettings.BoolToYesNo(d.KeepNameserver))
	if *d.KeepNameserver {
		return node
	}
	node.Appendf("DNS server address to use: %s", d.ServerAddress)
	node.AppendNode(d.DoT.toLinesNode())
	return node
}
