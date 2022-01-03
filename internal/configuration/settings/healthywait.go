package settings

import (
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

type HealthyWait struct {
	// Initial is the initial duration to wait for the program
	// to be healthy before taking action.
	// It cannot be nil in the internal state.
	Initial *time.Duration `json:"initial,omitempty"`
	// Addition is the duration to add to the Initial duration
	// after Initial has expired to wait longer for the program
	// to be healthy.
	// It cannot be nil in the internal state.
	Addition *time.Duration `json:"addition,omitempty"`
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
	h.Initial = helpers.MergeWithDuration(h.Initial, other.Initial)
	h.Addition = helpers.MergeWithDuration(h.Addition, other.Addition)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *HealthyWait) overrideWith(other HealthyWait) {
	h.Initial = helpers.OverrideWithDuration(h.Initial, other.Initial)
	h.Addition = helpers.OverrideWithDuration(h.Addition, other.Addition)
}

func (h *HealthyWait) setDefaults() {
	const initialDurationDefault = 6 * time.Second
	const additionDurationDefault = 5 * time.Second
	h.Initial = helpers.DefaultDuration(h.Initial, initialDurationDefault)
	h.Addition = helpers.DefaultDuration(h.Addition, additionDurationDefault)
}