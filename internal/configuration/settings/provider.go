package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
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

func (p Provider) validate(vpnType string, allServers models.AllServers) (err error) {
	// Validate Name
	var validNames []string
	if vpnType == constants.OpenVPN {
		validNames = constants.AllProviders()
		validNames = append(validNames, "pia") // Retro-compatibility
	} else { // Wireguard
		validNames = []string{
			constants.Custom,
			constants.Ivpn,
			constants.Mullvad,
			constants.Windscribe,
		}
	}
	if !helpers.IsOneOf(*p.Name, validNames...) {
		return fmt.Errorf("%w: %q can only be one of %s",
			ErrVPNProviderNameNotValid, *p.Name, helpers.ChoicesOrString(validNames))
	}

	err = p.ServerSelection.validate(*p.Name, allServers)
	if err != nil {
		return fmt.Errorf("server selection settings validation failed: %w", err)
	}

	err = p.PortForwarding.validate(*p.Name)
	if err != nil {
		return fmt.Errorf("port forwarding settings validation failed: %w", err)
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
	p.Name = helpers.DefaultStringPtr(p.Name, constants.PrivateInternetAccess)
	p.ServerSelection.setDefaults(*p.Name)
	p.PortForwarding.setDefaults()
}
