package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"github.com/qdm12/ss-server/pkg/tcpudp"
)

// Shadowsocks contains settings to configure the Shadowsocks server.
type Shadowsocks struct {
	// Enabled is true if the server should be running.
	// It defaults to false, and cannot be nil in the internal state.
	Enabled *bool
	// Settings are settings for the TCP+UDP server.
	tcpudp.Settings
}

func (s Shadowsocks) validate() (err error) {
	return s.Settings.Validate()
}

func (s *Shadowsocks) copy() (copied Shadowsocks) {
	return Shadowsocks{
		Enabled:  helpers.CopyBoolPtr(s.Enabled),
		Settings: s.Settings.Copy(),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (s *Shadowsocks) mergeWith(other Shadowsocks) {
	s.Enabled = helpers.MergeWithBool(s.Enabled, other.Enabled)
	s.Settings.MergeWith(other.Settings)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (s *Shadowsocks) overrideWith(other Shadowsocks) {
	s.Enabled = helpers.OverrideWithBool(s.Enabled, other.Enabled)
	s.Settings.OverrideWith(other.Settings)
}

func (s *Shadowsocks) setDefaults() {
	s.Enabled = helpers.DefaultBool(s.Enabled, false)
	s.Settings.SetDefaults()
}

func (s Shadowsocks) String() string {
	return s.toLinesNode().String()
}

func (s Shadowsocks) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Shadowsocks server settings:")

	node.Appendf("Enabled: %s", helpers.BoolPtrToYesNo(s.Enabled))
	if !*s.Enabled {
		return node
	}

	// TODO have ToLinesNode in qdm12/ss-server
	node.Appendf("Listening address: %s", s.Address)
	node.Appendf("Cipher: %s", s.CipherName)
	node.Appendf("Password: %s", helpers.ObfuscatePassword(*s.Password))
	node.Appendf("Log addresses: %s", helpers.BoolPtrToYesNo(s.LogAddresses))

	return node
}
