package settings

import (
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

type HealthyWait struct {
	// Initial is the initial duration to wait for the program
	// to be healthy before taking action.
	// It cannot be nil in the internal state.
	Initial *time.Duration
	// Addition is the duration to add to the Initial duration
	// after Initial has expired to wait longer for the program
	// to be healthy.
	// It cannot be nil in the internal state.
	Addition *time.Duration
}

func (h HealthyWait) validate() (err error) {
	return nil
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *HealthyWait) copy() (copied HealthyWait) {
	return HealthyWait{
		Initial:  helpers.CopyDurationPtr(h.Initial),
		Addition: helpers.CopyDurationPtr(h.Addition),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *HealthyWait) mergeWith(other HealthyWait) {
	h.Initial = helpers.MergeWithDurationPtr(h.Initial, other.Initial)
	h.Addition = helpers.MergeWithDurationPtr(h.Addition, other.Addition)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *HealthyWait) overrideWith(other HealthyWait) {
	h.Initial = helpers.OverrideWithDurationPtr(h.Initial, other.Initial)
	h.Addition = helpers.OverrideWithDurationPtr(h.Addition, other.Addition)
}

func (h *HealthyWait) setDefaults() {
	const initialDurationDefault = 6 * time.Second
	const additionDurationDefault = 5 * time.Second
	h.Initial = helpers.DefaultDurationPtr(h.Initial, initialDurationDefault)
	h.Addition = helpers.DefaultDurationPtr(h.Addition, additionDurationDefault)
}

func (h HealthyWait) String() string {
	return h.toLinesNode("Health").String()
}

func (h HealthyWait) toLinesNode(kind string) (node *gotree.Node) {
	node = gotree.New(kind + " wait durations:")
	node.Appendf("Initial duration: %s", *h.Initial)
	node.Appendf("Additional duration: %s", *h.Addition)
	return node
}
