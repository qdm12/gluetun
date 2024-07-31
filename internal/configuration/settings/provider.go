package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

// Provider contains settings specific to a VPN provider.
type Provider struct {
	// Name is the VPN service provider name.
	// It cannot be the empty string in the internal state.
	Name string `json:"name"`
	// ServerSelection is the settings to
	// select the VPN server.
	ServerSelection ServerSelection `json:"server_selection"`
	// PortForwarding is the settings about port forwarding.
	PortForwarding PortForwarding `json:"port_forwarding"`
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
			providers.Fastestvpn,
			providers.Ivpn,
			providers.Mullvad,
			providers.Nordvpn,
			providers.Surfshark,
			providers.Windscribe,
		}
	}
	if err = validate.IsOneOf(p.Name, validNames...); err != nil {
		return fmt.Errorf("%w for Wireguard: %w", ErrVPNProviderNameNotValid, err)
	}

	err = p.ServerSelection.validate(p.Name, storage)
	if err != nil {
		return fmt.Errorf("server selection: %w", err)
	}

	err = p.PortForwarding.Validate(p.Name)
	if err != nil {
		return fmt.Errorf("port forwarding: %w", err)
	}

	return nil
}

func (p *Provider) copy() (copied Provider) {
	return Provider{
		Name:            p.Name,
		ServerSelection: p.ServerSelection.copy(),
		PortForwarding:  p.PortForwarding.Copy(),
	}
}

func (p *Provider) overrideWith(other Provider) {
	p.Name = gosettings.OverrideWithComparable(p.Name, other.Name)
	p.ServerSelection.overrideWith(other.ServerSelection)
	p.PortForwarding.OverrideWith(other.PortForwarding)
}

func (p *Provider) setDefaults() {
	p.Name = gosettings.DefaultComparable(p.Name, providers.PrivateInternetAccess)
	p.ServerSelection.setDefaults(p.Name)
	p.PortForwarding.setDefaults()
}

func (p Provider) String() string {
	return p.toLinesNode().String()
}

func (p Provider) toLinesNode() (node *gotree.Node) {
	node = gotree.New("VPN provider settings:")
	node.Appendf("Name: %s", p.Name)
	node.AppendNode(p.ServerSelection.toLinesNode())
	node.AppendNode(p.PortForwarding.toLinesNode())
	return node
}

func (p *Provider) read(r *reader.Reader, vpnType string) (err error) {
	p.Name = readVPNServiceProvider(r, vpnType)

	err = p.ServerSelection.read(r, p.Name, vpnType)
	if err != nil {
		return fmt.Errorf("server selection: %w", err)
	}

	err = p.PortForwarding.read(r)
	if err != nil {
		return fmt.Errorf("port forwarding: %w", err)
	}

	return nil
}

func readVPNServiceProvider(r *reader.Reader, vpnType string) (vpnProvider string) {
	vpnProvider = r.String("VPN_SERVICE_PROVIDER", reader.RetroKeys("VPNSP"))
	if vpnProvider == "" {
		if vpnType != vpn.Wireguard && r.Get("OPENVPN_CUSTOM_CONFIG") != nil {
			// retro compatibility
			return providers.Custom
		}
		return ""
	}

	vpnProvider = strings.ToLower(vpnProvider)
	if vpnProvider == "pia" { // retro compatibility
		return providers.PrivateInternetAccess
	}

	return vpnProvider
}
