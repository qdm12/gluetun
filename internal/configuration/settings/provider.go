package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gotree"
)

// Provider contains settings specific to a VPN provider.
type Provider struct {
	// Name is the VPN service provider name.
	// It cannot be nil in the internal state.
	Name *string
	// ServerSelection is the settings to
	// select the VPN server.
	ServerSelection ServerSelection
	// PortForwarding is the settings about port forwarding.
	PortForwarding PortForwarding
}

// TODO v4 remove pointer for receiver (because of Surfshark).
func (p *Provider) validate(vpnType string, storage Storage) (err error) {
	// Validate Name
	var validNames []string
	if vpnType == vpn.OpenVPN {
		validNames = providers.AllWithCustom()
		validNames = append(validNames, "pia") // Retro-compatibility
	} else { // Wireguard
		validNames = []string{
			providers.Airvpn,
			providers.Custom,
			providers.Ivpn,
			providers.Mullvad,
			providers.Surfshark,
			providers.Windscribe,
		}
	}
	if !helpers.IsOneOf(*p.Name, validNames...) {
		return fmt.Errorf("%w for Wireguard: %q can only be one of %s",
			ErrVPNProviderNameNotValid, *p.Name, helpers.ChoicesOrString(validNames))
	}

	err = p.ServerSelection.validate(*p.Name, storage)
	if err != nil {
		return fmt.Errorf("server selection: %w", err)
	}

	err = p.PortForwarding.validate(*p.Name)
	if err != nil {
		return fmt.Errorf("port forwarding: %w", err)
	}

	return nil
}

func (p *Provider) copy() (copied Provider) {
	return Provider{
		Name:            helpers.CopyStringPtr(p.Name),
		ServerSelection: p.ServerSelection.copy(),
		PortForwarding:  p.PortForwarding.copy(),
	}
}

func (p *Provider) mergeWith(other Provider) {
	p.Name = helpers.MergeWithStringPtr(p.Name, other.Name)
	p.ServerSelection.mergeWith(other.ServerSelection)
	p.PortForwarding.mergeWith(other.PortForwarding)
}

func (p *Provider) overrideWith(other Provider) {
	p.Name = helpers.OverrideWithStringPtr(p.Name, other.Name)
	p.ServerSelection.overrideWith(other.ServerSelection)
	p.PortForwarding.overrideWith(other.PortForwarding)
}

func (p *Provider) setDefaults() {
	p.Name = helpers.DefaultStringPtr(p.Name, providers.PrivateInternetAccess)
	p.ServerSelection.setDefaults(*p.Name)
	p.PortForwarding.setDefaults()
}

func (p Provider) String() string {
	return p.toLinesNode().String()
}

func (p Provider) toLinesNode() (node *gotree.Node) {
	node = gotree.New("VPN provider settings:")
	node.Appendf("Name: %s", *p.Name)
	node.AppendNode(p.ServerSelection.toLinesNode())
	node.AppendNode(p.PortForwarding.toLinesNode())
	return node
}
