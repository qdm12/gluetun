package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

// Version contains settings to configure the version
// information fetcher.
type Version struct {
	// Enabled is true if the version information should
	// be fetched from Github.
	Enabled *bool
}

func (u Version) validate() (err error) {
	return nil
}

func (u *Version) copy() (copied Version) {
	return Version{
		Enabled: helpers.CopyBoolPtr(u.Enabled),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (u *Version) mergeWith(other Version) {
	u.Enabled = helpers.MergeWithBool(u.Enabled, other.Enabled)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (u *Version) overrideWith(other Version) {
	u.Enabled = helpers.OverrideWithBool(u.Enabled, other.Enabled)
}

func (u *Version) setDefaults() {
	u.Enabled = helpers.DefaultBool(u.Enabled, true)
}
