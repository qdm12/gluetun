package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type VPN struct {
	// Type is the VPN type and can only be
	// 'openvpn' or 'wireguard'. It cannot be the
	// empty string in the internal state.
	Type      string    `json:"type"`
	Provider  Provider  `json:"provider"`
	OpenVPN   OpenVPN   `json:"openvpn"`
	Wireguard Wireguard `json:"wireguard"`
}

// TODO v4 remove pointer for receiver (because of Surfshark).
func (v *VPN) Validate(storage Storage, ipv6Supported bool) (err error) {
	// Validate Type
	validVPNTypes := []string{vpn.OpenVPN, vpn.Wireguard}
	if err = validate.IsOneOf(v.Type, validVPNTypes...); err != nil {
		return fmt.Errorf("%w: %w", ErrVPNTypeNotValid, err)
	}

	err = v.Provider.validate(v.Type, storage)
	if err != nil {
		return fmt.Errorf("provider settings: %w", err)
	}

	if v.Type == vpn.OpenVPN {
		err := v.OpenVPN.validate(*v.Provider.Name)
		if err != nil {
			return fmt.Errorf("OpenVPN settings: %w", err)
		}
	} else {
		err := v.Wireguard.validate(*v.Provider.Name, ipv6Supported)
		if err != nil {
			return fmt.Errorf("Wireguard settings: %w", err)
		}
	}

	return nil
}

func (v *VPN) Copy() (copied VPN) {
	return VPN{
		Type:      v.Type,
		Provider:  v.Provider.copy(),
		OpenVPN:   v.OpenVPN.copy(),
		Wireguard: v.Wireguard.copy(),
	}
}

func (v *VPN) mergeWith(other VPN) {
	v.Type = gosettings.MergeWithString(v.Type, other.Type)
	v.Provider.mergeWith(other.Provider)
	v.OpenVPN.mergeWith(other.OpenVPN)
	v.Wireguard.mergeWith(other.Wireguard)
}

func (v *VPN) OverrideWith(other VPN) {
	v.Type = gosettings.OverrideWithString(v.Type, other.Type)
	v.Provider.overrideWith(other.Provider)
	v.OpenVPN.overrideWith(other.OpenVPN)
	v.Wireguard.overrideWith(other.Wireguard)
}

func (v *VPN) setDefaults() {
	v.Type = gosettings.DefaultString(v.Type, vpn.OpenVPN)
	v.Provider.setDefaults()
	v.OpenVPN.setDefaults(*v.Provider.Name)
	v.Wireguard.setDefaults(*v.Provider.Name)
}

func (v VPN) String() string {
	return v.toLinesNode().String()
}

func (v VPN) toLinesNode() (node *gotree.Node) {
	node = gotree.New("VPN settings:")

	node.AppendNode(v.Provider.toLinesNode())

	if v.Type == vpn.OpenVPN {
		node.AppendNode(v.OpenVPN.toLinesNode())
	} else {
		node.AppendNode(v.Wireguard.toLinesNode())
	}

	return node
}
