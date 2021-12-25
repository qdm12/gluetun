package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/ss-server/pkg/tcpudp"
)

// Shadowsocks contains settings to configure the Shadowsocks server.
type Shadowsocks struct {
	// Enabled is true if the server should be running.
	// It defaults to false, and cannot be nil in the internal state.
	Enabled *bool `json:"enabled"`
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
