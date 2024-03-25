package settings

import (
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// DoT contains settings to configure the DoT server.
type DoT struct {
	// Enabled is true if the DoT server should be running
	// and used. It defaults to true, and cannot be nil
	// in the internal state.
	Enabled *bool
	// UpdatePeriod is the period to update DNS block
	// lists and cryptographic files for DNSSEC validation.
	// It can be set to 0 to disable the update.
	// It defaults to 24h and cannot be nil in
	// the internal state.
	UpdatePeriod *time.Duration
	// Unbound contains settings to configure Unbound.
	Unbound Unbound
	// Blacklist contains settings to configure the filter
	// block lists.
	Blacklist DNSBlacklist
}

var (
	ErrDoTUpdatePeriodTooShort = errors.New("update period is too short")
)

func (d DoT) validate() (err error) {
	const minUpdatePeriod = 30 * time.Second
	if *d.UpdatePeriod != 0 && *d.UpdatePeriod < minUpdatePeriod {
		return fmt.Errorf("%w: %s must be bigger than %s",
			ErrDoTUpdatePeriodTooShort, *d.UpdatePeriod, minUpdatePeriod)
	}

	err = d.Unbound.validate()
	if err != nil {
		return err
	}

	err = d.Blacklist.validate()
	if err != nil {
		return err
	}

	return nil
}

func (d *DoT) copy() (copied DoT) {
	return DoT{
		Enabled:      gosettings.CopyPointer(d.Enabled),
		UpdatePeriod: gosettings.CopyPointer(d.UpdatePeriod),
		Unbound:      d.Unbound.copy(),
		Blacklist:    d.Blacklist.copy(),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DoT) overrideWith(other DoT) {
	d.Enabled = gosettings.OverrideWithPointer(d.Enabled, other.Enabled)
	d.UpdatePeriod = gosettings.OverrideWithPointer(d.UpdatePeriod, other.UpdatePeriod)
	d.Unbound.overrideWith(other.Unbound)
	d.Blacklist.overrideWith(other.Blacklist)
}

func (d *DoT) setDefaults() {
	d.Enabled = gosettings.DefaultPointer(d.Enabled, true)
	const defaultUpdatePeriod = 24 * time.Hour
	d.UpdatePeriod = gosettings.DefaultPointer(d.UpdatePeriod, defaultUpdatePeriod)
	d.Unbound.setDefaults()
	d.Blacklist.setDefaults()
}

func (d DoT) String() string {
	return d.toLinesNode().String()
}

func (d DoT) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS over TLS settings:")

	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(d.Enabled))
	if !*d.Enabled {
		return node
	}

	update := "disabled" //nolint:goconst
	if *d.UpdatePeriod > 0 {
		update = "every " + d.UpdatePeriod.String()
	}
	node.Appendf("Update period: %s", update)

	node.AppendNode(d.Unbound.toLinesNode())
	node.AppendNode(d.Blacklist.toLinesNode())

	return node
}

func (d *DoT) read(reader *reader.Reader) (err error) {
	d.Enabled, err = reader.BoolPtr("DOT")
	if err != nil {
		return err
	}

	d.UpdatePeriod, err = reader.DurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return err
	}

	err = d.Unbound.read(reader)
	if err != nil {
		return err
	}

	err = d.Blacklist.read(reader)
	if err != nil {
		return err
	}

	return nil
}
