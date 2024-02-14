package settings

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/qdm12/gluetun/internal/publicip/api"
	"github.com/qdm12/gosettings"
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
	// API is the API name to use to fetch public IP information.
	// It can be ipinfo or ip2location. It defaults to ipinfo.
	API string
	// APIToken is the token to use for the IP data service
	// such as ipinfo.io. It can be the empty string to
	// indicate not to use a token. It cannot be nil for the
	// internal state.
	APIToken *string
}

// UpdateWith deep copies the receiving settings, overrides the copy with
// fields set in the partialUpdate argument, validates the new settings
// and returns them if they are valid, or returns an error otherwise.
// In all cases, the receiving settings are unmodified.
func (p PublicIP) UpdateWith(partialUpdate PublicIP) (updatedSettings PublicIP, err error) {
	updatedSettings = p.copy()
	updatedSettings.overrideWith(partialUpdate)
	err = updatedSettings.validate()
	if err != nil {
		return updatedSettings, fmt.Errorf("validating updated settings: %w", err)
	}
	return updatedSettings, nil
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

	_, err = api.ParseProvider(p.API)
	if err != nil {
		return fmt.Errorf("API name: %w", err)
	}

	return nil
}

func (p *PublicIP) copy() (copied PublicIP) {
	return PublicIP{
		Period:     gosettings.CopyPointer(p.Period),
		IPFilepath: gosettings.CopyPointer(p.IPFilepath),
		API:        p.API,
		APIToken:   gosettings.CopyPointer(p.APIToken),
	}
}

func (p *PublicIP) mergeWith(other PublicIP) {
	p.Period = gosettings.MergeWithPointer(p.Period, other.Period)
	p.IPFilepath = gosettings.MergeWithPointer(p.IPFilepath, other.IPFilepath)
	p.API = gosettings.MergeWithString(p.API, other.API)
	p.APIToken = gosettings.MergeWithPointer(p.APIToken, other.APIToken)
}

func (p *PublicIP) overrideWith(other PublicIP) {
	p.Period = gosettings.OverrideWithPointer(p.Period, other.Period)
	p.IPFilepath = gosettings.OverrideWithPointer(p.IPFilepath, other.IPFilepath)
	p.API = gosettings.OverrideWithString(p.API, other.API)
	p.APIToken = gosettings.OverrideWithPointer(p.APIToken, other.APIToken)
}

func (p *PublicIP) setDefaults() {
	const defaultPeriod = 12 * time.Hour
	p.Period = gosettings.DefaultPointer(p.Period, defaultPeriod)
	p.IPFilepath = gosettings.DefaultPointer(p.IPFilepath, "/tmp/gluetun/ip")
	p.API = gosettings.DefaultString(p.API, "ipinfo")
	p.APIToken = gosettings.DefaultPointer(p.APIToken, "")
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

	node.Appendf("Public IP data API: %s", p.API)

	if *p.APIToken != "" {
		node.Appendf("API token: %s", gosettings.ObfuscateKey(*p.APIToken))
	}

	return node
}
