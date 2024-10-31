package settings

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireguardSelection struct {
	EndpointIP   netip.Addr `json:"endpoint_ip"`
	EndpointPort *uint16    `json:"endpoint_port"`
	PublicKey    string     `json:"public_key"`
	UseIPv6      bool       `json:"use_ipv6"`
}

func (w WireguardSelection) validate(vpnProvider string) (err error) {
	// Validate IPv6 usage setting
	if !w.UseIPv6 && w.EndpointIP.Is6() {
		return fmt.Errorf("%w: IPv6 address is disabled by configuration", ErrWireguardEndpointIPNotSet)
	}

	// Validate EndpointIP
	switch vpnProvider {
	case providers.Airvpn, providers.Fastestvpn, providers.Ivpn,
		providers.Mullvad, providers.Nordvpn, providers.Protonvpn,
		providers.Surfshark, providers.Windscribe:
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
	case providers.Fastestvpn, providers.Nordvpn,
		providers.Protonvpn, providers.Surfshark:
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
	case providers.Fastestvpn, providers.Ivpn, providers.Mullvad,
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

func (w *WireguardSelection) read(r *reader.Reader) (err error) {
	w.EndpointIP, err = r.NetipAddr("WIREGUARD_ENDPOINT_IP", reader.RetroKeys("VPN_ENDPOINT_IP"))
	if err != nil {
		return err
	}

	w.EndpointPort, err = r.Uint16Ptr("WIREGUARD_ENDPOINT_PORT", reader.RetroKeys("VPN_ENDPOINT_PORT"))
	if err != nil {
		return err
	}

	w.PublicKey = r.String("WIREGUARD_PUBLIC_KEY", reader.ForceLowercase(false))

	// Read the IPv6 usage setting
	w.UseIPv6, err = r.Bool("WIREGUARD_USE_IPV6_SERVER", reader.RetroKeys("VPN_USE_IPV6_SERVER"))
	if err != nil {
		return err
	}

	return nil
}
