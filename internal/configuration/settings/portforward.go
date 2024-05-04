package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

// PortForwarding contains settings for port forwarding.
type PortForwarding struct {
	// Enabled is true if port forwarding should be activated.
	// It cannot be nil for the internal state.
	Enabled *bool `json:"enabled"`
	// Provider is set to specify which custom port forwarding code
	// should be used. This is especially necessary for the custom
	// provider using Wireguard for a provider where Wireguard is not
	// natively supported but custom port forwarding code is available.
	// It defaults to the empty string, meaning the current provider
	// should be the one used for port forwarding. It can also be set
	// to "mock" to run a mock port forwarding for testing.
	// It cannot be nil for the internal state.
	Provider *string `json:"provider"`
	// Filepath is the port forwarding status file path
	// to use. It can be the empty string to indicate not
	// to write to a file. It cannot be nil for the
	// internal state
	Filepath *string `json:"status_file_path"`
	// ListeningPort is the port traffic would be redirected to from the
	// forwarded port. The redirection is disabled if it is set to 0, which
	// is its default as well.
	ListeningPort *uint16 `json:"listening_port"`
}

func (p PortForwarding) Validate(vpnProvider string) (err error) {
	if !*p.Enabled {
		return nil
	}

	// Validate current provider or custom provider specified
	providerSelected := vpnProvider
	if *p.Provider != "" {
		providerSelected = *p.Provider
	}
	validProviders := []string{
		providers.Mock,
		providers.PrivateInternetAccess,
		providers.Protonvpn,
	}
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

func (p *PortForwarding) Copy() (copied PortForwarding) {
	return PortForwarding{
		Enabled:       gosettings.CopyPointer(p.Enabled),
		Provider:      gosettings.CopyPointer(p.Provider),
		Filepath:      gosettings.CopyPointer(p.Filepath),
		ListeningPort: gosettings.CopyPointer(p.ListeningPort),
	}
}

func (p *PortForwarding) OverrideWith(other PortForwarding) {
	p.Enabled = gosettings.OverrideWithPointer(p.Enabled, other.Enabled)
	p.Provider = gosettings.OverrideWithPointer(p.Provider, other.Provider)
	p.Filepath = gosettings.OverrideWithPointer(p.Filepath, other.Filepath)
	p.ListeningPort = gosettings.OverrideWithPointer(p.ListeningPort, other.ListeningPort)
}

func (p *PortForwarding) setDefaults() {
	p.Enabled = gosettings.DefaultPointer(p.Enabled, false)
	p.Provider = gosettings.DefaultPointer(p.Provider, "")
	p.Filepath = gosettings.DefaultPointer(p.Filepath, "/tmp/gluetun/forwarded_port")
	p.ListeningPort = gosettings.DefaultPointer(p.ListeningPort, 0)
}

func (p PortForwarding) String() string {
	return p.toLinesNode().String()
}

func (p PortForwarding) toLinesNode() (node *gotree.Node) {
	if !*p.Enabled {
		return nil
	}

	node = gotree.New("Automatic port forwarding settings:")

	listeningPort := "disabled"
	if *p.ListeningPort != 0 {
		listeningPort = fmt.Sprintf("%d", *p.ListeningPort)
	}
	node.Appendf("Redirection listening port: %s", listeningPort)

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

func (p *PortForwarding) read(r *reader.Reader) (err error) {
	p.Enabled, err = r.BoolPtr("VPN_PORT_FORWARDING",
		reader.RetroKeys(
			"PORT_FORWARDING",
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING",
		))
	if err != nil {
		return err
	}

	p.Provider = r.Get("VPN_PORT_FORWARDING_PROVIDER")

	p.Filepath = r.Get("VPN_PORT_FORWARDING_STATUS_FILE",
		reader.ForceLowercase(false),
		reader.RetroKeys(
			"PORT_FORWARDING_STATUS_FILE",
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING_STATUS_FILE",
		))

	p.ListeningPort, err = r.Uint16Ptr("VPN_PORT_FORWARDING_LISTENING_PORT")
	if err != nil {
		return err
	}

	return nil
}
