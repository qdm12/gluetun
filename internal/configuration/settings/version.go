package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

// Version contains settings to configure the version
// information fetcher.
type Version struct {
	// Enabled is true if the version information should
	// be fetched from Github.
	Enabled *bool
}

func (v Version) validate() (err error) {
	return nil
}

func (v *Version) copy() (copied Version) {
	return Version{
		Enabled: helpers.CopyBoolPtr(v.Enabled),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (v *Version) mergeWith(other Version) {
	v.Enabled = helpers.MergeWithBool(v.Enabled, other.Enabled)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (v *Version) overrideWith(other Version) {
	v.Enabled = helpers.OverrideWithBool(v.Enabled, other.Enabled)
}

func (v *Version) setDefaults() {
	v.Enabled = helpers.DefaultBool(v.Enabled, true)
}

func (v Version) String() string {
	return v.toLinesNode().String()
}

func (v Version) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Version settings:")

	node.Appendf("Enabled: %s", helpers.BoolPtrToYesNo(v.Enabled))

	return node
}
