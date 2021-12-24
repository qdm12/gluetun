package settings

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

// PublicIP contains settings for port forwarding.
type PublicIP struct {
	// Period is the period to get the public IP address.
	// It can be set to 0 to disable periodic checking.
	// It cannot be nil for the internal state.
	Period *time.Duration `json:"period,omitempty"`
	// IPFilepath is the public IP address status file path
	// to use. It can be the empty string to indicate not
	// to write to a file. It cannot be nil for the
	// internal state
	IPFilepath *string `json:"ip_filepath,omitempty"`
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
			return fmt.Errorf("%w: %s", ErrPublicIPFilepathNotValid, err)
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
	p.Period = helpers.MergeWithDuration(p.Period, other.Period)
	p.IPFilepath = helpers.MergeWithStringPtr(p.IPFilepath, other.IPFilepath)
}

func (p *PublicIP) overrideWith(other PublicIP) {
	p.Period = helpers.OverrideWithDuration(p.Period, other.Period)
	p.IPFilepath = helpers.OverrideWithStringPtr(p.IPFilepath, other.IPFilepath)
}

func (p *PublicIP) setDefaults() {
	const defaultPeriod = 12 * time.Hour
	p.Period = helpers.DefaultDuration(p.Period, defaultPeriod)
	p.IPFilepath = helpers.DefaultStringPtr(p.IPFilepath, "/tmp/gluetun/ip")
}
