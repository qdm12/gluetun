package settings

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
)

// PortForwarding contains settings for port forwarding.
type PortForwarding struct {
	// Enabled is true if port forwarding should be activated.
	// It cannot be nil for the internal state.
	Enabled *bool `json:"enabled,omitempty"`
	// Filepath is the port forwarding status file path
	// to use. It can be the empty string to indicate not
	// to write to a file. It cannot be nil for the
	// internal state
	Filepath *string `json:"filepath,omitempty"`
}

func (p PortForwarding) validate(vpnProvider string) (err error) {
	if !*p.Enabled {
		return nil
	}

	// Validate Enabled
	validProviders := []string{constants.PrivateInternetAccess}
	if !helpers.IsOneOf(vpnProvider, validProviders...) {
		return fmt.Errorf("%w: for provider %s, it is only available for %s",
			ErrPortForwardingEnabled, vpnProvider, strings.Join(validProviders, ", "))
	}

	// Validate Filepath
	if *p.Filepath != "" { // optional
		_, err := filepath.Abs(*p.Filepath)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrPortForwardingFilepathNotValid, err)
		}
	}

	return nil
}

func (p *PortForwarding) copy() (copied PortForwarding) {
	return PortForwarding{
		Enabled:  helpers.CopyBoolPtr(p.Enabled),
		Filepath: helpers.CopyStringPtr(p.Filepath),
	}
}

func (p *PortForwarding) mergeWith(other PortForwarding) {
	p.Enabled = helpers.MergeWithBool(p.Enabled, other.Enabled)
	p.Filepath = helpers.MergeWithStringPtr(p.Filepath, other.Filepath)
}

func (p *PortForwarding) overrideWith(other PortForwarding) {
	p.Enabled = helpers.OverrideWithBool(p.Enabled, other.Enabled)
	p.Filepath = helpers.OverrideWithStringPtr(p.Filepath, other.Filepath)
}

func (p *PortForwarding) setDefaults() {
	p.Filepath = helpers.DefaultStringPtr(p.Filepath, "/tmp/gluetun/forwarded_port")
}
