package settings

import (
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/dns/pkg/blacklist"
	"github.com/qdm12/dns/pkg/unbound"
)

// DoT contains settings to configure the DoT server.
type DoT struct {
	// Enabled is true if the DoT server should be running
	// and used. It defaults to true, and cannot be nil
	// in the internal state.
	Enabled *bool `json:"enabled"`
	// UpdatePeriod is the period to update DNS block
	// lists and cryptographic files for DNSSEC validation.
	// It can be set to 0 to disable the update.
	// It defaults to TODO and cannot be nil in
	// the internal state.
	UpdatePeriod *time.Duration
	// Unbound contains settings to configure Unbound.
	Unbound        unbound.Settings
	BlacklistBuild blacklist.BuilderSettings
}

var (
	ErrDoTUpdatePeriodTooShort = errors.New("update period is too short")
)

func (d DoT) validate() (err error) {
	const minUpdatePeriod = 30 * time.Second
	if *d.UpdatePeriod < minUpdatePeriod {
		return fmt.Errorf("%w: %s must be bigger than %s",
			ErrDoTUpdatePeriodTooShort, *d.UpdatePeriod, minUpdatePeriod)
	}

	return nil
}

func validateUnbound(settings unbound.Settings) (err error) {

}

func (d *DoT) copy() (copied DoT) {
	// TODO
	return DoT{}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (d *DoT) mergeWith(other DoT) {
	// TODO
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DoT) overrideWith(other DoT) {
	// TODO
}

func (d *DoT) setDefaults() {
	// TODO
}
