package settings

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

// DNS contains settings to configure DNS.
type DNS struct {
	// ServerAddress is the DNS server to use inside
	// the Go program and for the system.
	// It defaults to '127.0.0.1' to be used with the
	// DoT server. It cannot be nil in the internal
	// state.
	ServerAddress net.IP
	// KeepNameserver is true if the Docker DNS server
	// found in /etc/resolv.conf should be kept.
	// Note settings this to true will go around the
	// DoT server blocking.
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
		return fmt.Errorf("failed validating DoT settings: %w", err)
	}

	return nil
}

func (d *DNS) Copy() (copied DNS) {
	return DNS{
		ServerAddress:  helpers.CopyIP(d.ServerAddress),
		KeepNameserver: helpers.CopyBoolPtr(d.KeepNameserver),
		DoT:            d.DoT.copy(),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (d *DNS) mergeWith(other DNS) {
	d.ServerAddress = helpers.MergeWithIP(d.ServerAddress, other.ServerAddress)
	d.KeepNameserver = helpers.MergeWithBool(d.KeepNameserver, other.KeepNameserver)
	d.DoT.mergeWith(other.DoT)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DNS) overrideWith(other DNS) {
	d.ServerAddress = helpers.OverrideWithIP(d.ServerAddress, other.ServerAddress)
	d.KeepNameserver = helpers.OverrideWithBool(d.KeepNameserver, other.KeepNameserver)
	d.DoT.overrideWith(other.DoT)
}

func (d *DNS) setDefaults() {
	localhost := net.IPv4(127, 0, 0, 1) //nolint:gomnd
	d.ServerAddress = helpers.DefaultIP(d.ServerAddress, localhost)
	d.KeepNameserver = helpers.DefaultBool(d.KeepNameserver, false)
	d.DoT.setDefaults()
}

func (d DNS) String() string {
	return d.toLinesNode().String()
}

func (d DNS) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS settings:")
	node.Appendf("DNS server address to use: %s", d.ServerAddress)
	node.Appendf("Keep existing nameserver(s): %s", helpers.BoolPtrToYesNo(d.KeepNameserver))
	node.AppendNode(d.DoT.toLinesNode())
	return node
}
