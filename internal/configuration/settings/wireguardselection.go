package settings

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireguardSelection struct {
	// EndpointIP is the server endpoint IP address.
	// It is only used with VPN providers generating Wireguard
	// configurations specific to each server and user.
	// To indicate it should not be used, it should be set
	// to netip.IPv4Unspecified(). It can never be the zero value
	// in the internal state.
	EndpointIP netip.Addr `json:"endpoint_ip"`
	// EndpointPort is a the server port to use for the VPN server.
	// It is optional for VPN providers IVPN, Mullvad, Surfshark
	// and Windscribe, and compulsory for the others.
	// When optional, it can be set to 0 to indicate not use
	// a custom endpoint port. It cannot be nil in the internal
	// state.
	EndpointPort *uint16 `json:"endpoint_port"`
	// PublicKey is the server public key.
	// It is only used with VPN providers generating Wireguard
	// configurations specific to each server and user.
	PublicKey string `json:"public_key"`
}

// Validate validates WireguardSelection settings.
// It should only be ran if the VPN type chosen is Wireguard.
func (w WireguardSelection) validate(vpnProvider string) (err error) {
	// Validate EndpointIP
	switch vpnProvider {
	case providers.Airvpn, providers.Ivpn, providers.Mullvad,
		providers.Nordvpn, providers.Surfshark, providers.Windscribe:
		// endpoint IP addresses are baked in
	case providers.Custom:
		if !w.EndpointIP.IsValid() || w.EndpointIP.IsUnspecified() {
			return fmt.Errorf("%w", ErrWireguardEndpointIPNotSet)
		}
	default: // Providers not supporting Wireguard
	}

	// Validate EndpointPort
	switch vpnProvider {
	// EndpointPort is required
	case providers.Custom:
		if *w.EndpointPort == 0 {
			return fmt.Errorf("%w", ErrWireguardEndpointPortNotSet)
		}
	// EndpointPort cannot be set
	case providers.Surfshark, providers.Nordvpn:
		if *w.EndpointPort != 0 {
			return fmt.Errorf("%w", ErrWireguardEndpointPortSet)
		}
	case providers.Airvpn, providers.Ivpn, providers.Mullvad, providers.Windscribe:
		// EndpointPort is optional and can be 0
		if *w.EndpointPort == 0 {
			break // no custom endpoint port set
		}
		if vpnProvider == providers.Mullvad {
			break // no restriction on custom endpoint port value
		}
		var allowed []uint16
		switch vpnProvider {
		case providers.Airvpn:
			allowed = []uint16{1637, 47107}
		case providers.Ivpn:
			allowed = []uint16{2049, 2050, 53, 30587, 41893, 48574, 58237}
		case providers.Windscribe:
			allowed = []uint16{53, 80, 123, 443, 1194, 65142}
		}

		err = validate.IsOneOf(*w.EndpointPort, allowed...)
		if err == nil {
			break
		}
		return fmt.Errorf("%w: for VPN service provider %s: %w",
			ErrWireguardEndpointPortNotAllowed, vpnProvider, err)
	default: // Providers not supporting Wireguard
	}

	// Validate PublicKey
	switch vpnProvider {
	case providers.Ivpn, providers.Mullvad,
		providers.Surfshark, providers.Windscribe:
		// public keys are baked in
	case providers.Custom:
		if w.PublicKey == "" {
			return fmt.Errorf("%w", ErrWireguardPublicKeyNotSet)
		}
	default: // Providers not supporting Wireguard
	}
	if w.PublicKey != "" {
		_, err := wgtypes.ParseKey(w.PublicKey)
		if err != nil {
			return fmt.Errorf("%w: %s: %s",
				ErrWireguardPublicKeyNotValid, w.PublicKey, err)
		}
	}

	return nil
}

func (w *WireguardSelection) copy() (copied WireguardSelection) {
	return WireguardSelection{
		EndpointIP:   w.EndpointIP,
		EndpointPort: gosettings.CopyPointer(w.EndpointPort),
		PublicKey:    w.PublicKey,
	}
}

func (w *WireguardSelection) mergeWith(other WireguardSelection) {
	w.EndpointIP = gosettings.MergeWithValidator(w.EndpointIP, other.EndpointIP)
	w.EndpointPort = gosettings.MergeWithPointer(w.EndpointPort, other.EndpointPort)
	w.PublicKey = gosettings.MergeWithString(w.PublicKey, other.PublicKey)
}

func (w *WireguardSelection) overrideWith(other WireguardSelection) {
	w.EndpointIP = gosettings.OverrideWithValidator(w.EndpointIP, other.EndpointIP)
	w.EndpointPort = gosettings.OverrideWithPointer(w.EndpointPort, other.EndpointPort)
	w.PublicKey = gosettings.OverrideWithString(w.PublicKey, other.PublicKey)
}

func (w *WireguardSelection) setDefaults() {
	w.EndpointIP = gosettings.DefaultValidator(w.EndpointIP, netip.IPv4Unspecified())
	w.EndpointPort = gosettings.DefaultPointer(w.EndpointPort, 0)
}

func (w WireguardSelection) String() string {
	return w.toLinesNode().String()
}

func (w WireguardSelection) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Wireguard selection settings:")

	if !w.EndpointIP.IsUnspecified() {
		node.Appendf("Endpoint IP address: %s", w.EndpointIP)
	}

	if *w.EndpointPort != 0 {
		node.Appendf("Endpoint port: %d", *w.EndpointPort)
	}

	if w.PublicKey != "" {
		node.Appendf("Server public key: %s", w.PublicKey)
	}

	return node
}
