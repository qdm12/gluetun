package settings

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

// PublicIP contains settings for port forwarding.
type PublicIP struct {
	// Period is the period to get the public IP address.
	// It can be set to 0 to disable periodic checking.
	// It cannot be nil for the internal state.
	// TODO change to value and add enabled field
	Period *time.Duration
	// IPFilepath is the public IP address status file path
	// to use. It can be the empty string to indicate not
	// to write to a file. It cannot be nil for the
	// internal state
	IPFilepath *string
}

func (p PublicIP) validate() (err error) {
	const minPeriod = 5 * time.Second
	if *p.Period < minPeriod {
		return fmt.Errorf("%w: %s must be at least %s",
			ErrPublicIPPeriodTooShort, p.Period, minPeriod)
	}

	if *p.IPFilepath != "" { // optional
		_, err := filepath.Abs(*p.IPFilepath)
		if err != nil {
			return fmt.Errorf("filepath is not valid: %w", err)
		}
	}

	return nil
}

func (p *PublicIP) copy() (copied PublicIP) {
	return PublicIP{
		Period:     helpers.CopyDurationPtr(p.Period),
		IPFilepath: helpers.CopyStringPtr(p.IPFilepath),
	}
}

func (p *PublicIP) mergeWith(other PublicIP) {
	p.Period = helpers.MergeWithDurationPtr(p.Period, other.Period)
	p.IPFilepath = helpers.MergeWithStringPtr(p.IPFilepath, other.IPFilepath)
}

func (p *PublicIP) overrideWith(other PublicIP) {
	p.Period = helpers.OverrideWithDurationPtr(p.Period, other.Period)
	p.IPFilepath = helpers.OverrideWithStringPtr(p.IPFilepath, other.IPFilepath)
}

func (p *PublicIP) setDefaults() {
	const defaultPeriod = 12 * time.Hour
	p.Period = helpers.DefaultDurationPtr(p.Period, defaultPeriod)
	p.IPFilepath = helpers.DefaultStringPtr(p.IPFilepath, "/tmp/gluetun/ip")
}

func (p PublicIP) String() string {
	return p.toLinesNode().String()
}

func (p PublicIP) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Public IP settings:")

	if *p.Period == 0 {
		node.Appendf("Enabled: no")
		return node
	}

	updatePeriod := "disabled"
	if *p.Period > 0 {
		updatePeriod = "every " + p.Period.String()
	}
	node.Appendf("Fetching: %s", updatePeriod)

	if *p.IPFilepath != "" {
		node.Appendf("IP file path: %s", *p.IPFilepath)
	}

	return node
}
