package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

// PortForwarding contains settings for port forwarding.
type PortForwarding struct {
	// Enabled is true if port forwarding should be activated.
	// It cannot be nil for the internal state.
	Enabled *bool
	// Provider is set to specify which custom port forwarding code
	// should be used. This is especially necessary for the custom
	// provider using Wireguard for a provider where Wireguard is not
	// natively supported but custom port forwading code is available.
	// It defaults to the empty string, meaning the current provider
	// should be the one used for port forwarding.
	// It cannot be nil for the internal state.
	Provider *string
	// Filepath is the port forwarding status file path
	// to use. It can be the empty string to indicate not
	// to write to a file. It cannot be nil for the
	// internal state
	Filepath *string
}

func (p PortForwarding) validate(vpnProvider string) (err error) {
	if !*p.Enabled {
		return nil
	}

	// Validate current provider or custom provider specified
	providerSelected := vpnProvider
	if *p.Provider != "" {
		providerSelected = *p.Provider
	}
	validProviders := []string{providers.PrivateInternetAccess}
	if err = validate.IsOneOf(providerSelected, validProviders...); err != nil {
		return fmt.Errorf("%w: %w", ErrPortForwardingEnabled, err)
	}

	// Validate Filepath
	if *p.Filepath != "" { // optional
		_, err := filepath.Abs(*p.Filepath)
		if err != nil {
			return fmt.Errorf("filepath is not valid: %w", err)
		}
	}

	return nil
}

func (p *PortForwarding) copy() (copied PortForwarding) {
	return PortForwarding{
		Enabled:  gosettings.CopyPointer(p.Enabled),
		Provider: gosettings.CopyPointer(p.Provider),
		Filepath: gosettings.CopyPointer(p.Filepath),
	}
}

func (p *PortForwarding) mergeWith(other PortForwarding) {
	p.Enabled = gosettings.MergeWithPointer(p.Enabled, other.Enabled)
	p.Provider = gosettings.MergeWithPointer(p.Provider, other.Provider)
	p.Filepath = gosettings.MergeWithPointer(p.Filepath, other.Filepath)
}

func (p *PortForwarding) overrideWith(other PortForwarding) {
	p.Enabled = gosettings.OverrideWithPointer(p.Enabled, other.Enabled)
	p.Provider = gosettings.OverrideWithPointer(p.Provider, other.Provider)
	p.Filepath = gosettings.OverrideWithPointer(p.Filepath, other.Filepath)
}

func (p *PortForwarding) setDefaults() {
	p.Enabled = gosettings.DefaultPointer(p.Enabled, false)
	p.Provider = gosettings.DefaultPointer(p.Provider, "")
	p.Filepath = gosettings.DefaultPointer(p.Filepath, "/tmp/gluetun/forwarded_port")
}

func (p PortForwarding) String() string {
	return p.toLinesNode().String()
}

func (p PortForwarding) toLinesNode() (node *gotree.Node) {
	if !*p.Enabled {
		return nil
	}

	node = gotree.New("Automatic port forwarding settings:")
	if *p.Provider == "" {
		node.Appendf("Use port forwarding code for current provider")
	} else {
		node.Appendf("Use code for provider: %s", *p.Provider)
	}

	filepath := *p.Filepath
	if filepath == "" {
		filepath = "[not set]"
	}
	node.Appendf("Forwarded port file path: %s", filepath)

	return node
}
