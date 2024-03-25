package settings

import (
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
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
		Enabled: gosettings.CopyPointer(v.Enabled),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (v *Version) overrideWith(other Version) {
	v.Enabled = gosettings.OverrideWithPointer(v.Enabled, other.Enabled)
}

func (v *Version) setDefaults() {
	v.Enabled = gosettings.DefaultPointer(v.Enabled, true)
}

func (v Version) String() string {
	return v.toLinesNode().String()
}

func (v Version) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Version settings:")

	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(v.Enabled))

	return node
}

func (v *Version) read(r *reader.Reader) (err error) {
	v.Enabled, err = r.BoolPtr("VERSION_INFORMATION")
	if err != nil {
		return err
	}

	return nil
}
