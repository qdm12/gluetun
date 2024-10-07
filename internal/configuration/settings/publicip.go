package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/publicip/api"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// PublicIP contains settings for port forwarding.
type PublicIP struct {
	// Enabled is set to true to fetch the public ip address
	// information on VPN connection. It defaults to true.
	Enabled *bool
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
		Enabled:    gosettings.CopyPointer(p.Enabled),
		IPFilepath: gosettings.CopyPointer(p.IPFilepath),
		API:        p.API,
		APIToken:   gosettings.CopyPointer(p.APIToken),
	}
}

func (p *PublicIP) overrideWith(other PublicIP) {
	p.Enabled = gosettings.OverrideWithPointer(p.Enabled, other.Enabled)
	p.IPFilepath = gosettings.OverrideWithPointer(p.IPFilepath, other.IPFilepath)
	p.API = gosettings.OverrideWithComparable(p.API, other.API)
	p.APIToken = gosettings.OverrideWithPointer(p.APIToken, other.APIToken)
}

func (p *PublicIP) setDefaults() {
	p.Enabled = gosettings.DefaultPointer(p.Enabled, true)
	p.IPFilepath = gosettings.DefaultPointer(p.IPFilepath, "/tmp/gluetun/ip")
	p.API = gosettings.DefaultComparable(p.API, "ipinfo")
	p.APIToken = gosettings.DefaultPointer(p.APIToken, "")
}

func (p PublicIP) String() string {
	return p.toLinesNode().String()
}

func (p PublicIP) toLinesNode() (node *gotree.Node) {
	if !*p.Enabled {
		return gotree.New("Public IP settings: disabled")
	}

	node = gotree.New("Public IP settings:")

	if *p.IPFilepath != "" {
		node.Appendf("IP file path: %s", *p.IPFilepath)
	}

	node.Appendf("Public IP data API: %s", p.API)

	if *p.APIToken != "" {
		node.Appendf("API token: %s", gosettings.ObfuscateKey(*p.APIToken))
	}

	return node
}

func (p *PublicIP) read(r *reader.Reader, warner Warner) (err error) {
	p.Enabled, err = readPublicIPEnabled(r, warner)
	if err != nil {
		return err
	}

	p.IPFilepath = r.Get("PUBLICIP_FILE",
		reader.ForceLowercase(false), reader.RetroKeys("IP_STATUS_FILE"))
	p.API = r.String("PUBLICIP_API")
	p.APIToken = r.Get("PUBLICIP_API_TOKEN")
	return nil
}

func readPublicIPEnabled(r *reader.Reader, warner Warner) (
	enabled *bool, err error) {
	periodPtr, err := r.DurationPtr("PUBLICIP_PERIOD") // Retro-compatibility
	if err != nil {
		return nil, err
	} else if periodPtr == nil {
		return r.BoolPtr("PUBLICIP_ENABLED")
	}

	if *periodPtr == 0 {
		warner.Warn("please replace PUBLICIP_PERIOD=0 with PUBLICIP_ENABLED=no")
		return ptrTo(false), nil
	}

	warner.Warn("PUBLICIP_PERIOD is no longer used. " +
		"It is assumed from its non-zero value you want PUBLICIP_ENABLED=yes. " +
		"Please migrate to use PUBLICIP_ENABLED only in the future.")
	return ptrTo(true), nil
}
