package settings

import (
	"fmt"
	"net/netip"
	"regexp"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gotree"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Wireguard contains settings to configure the Wireguard client.
type Wireguard struct {
	// PrivateKey is the Wireguard client peer private key.
	// It cannot be nil in the internal state.
	PrivateKey *string
	// PreSharedKey is the Wireguard pre-shared key.
	// It can be the empty string to indicate there
	// is no pre-shared key.
	// It cannot be nil in the internal state.
	PreSharedKey *string
	// Addresses are the Wireguard interface addresses.
	Addresses []netip.Prefix
	// Interface is the name of the Wireguard interface
	// to create. It cannot be the empty string in the
	// internal state.
	Interface string
	// Implementation is the Wireguard implementation to use.
	// It can be "auto", "userspace" or "kernelspace".
	// It defaults to "auto" and cannot be the empty string
	// in the internal state.
	Implementation string
}

var regexpInterfaceName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// Validate validates Wireguard settings.
// It should only be ran if the VPN type chosen is Wireguard.
func (w Wireguard) validate(vpnProvider string, ipv6Supported bool) (err error) {
	if !helpers.IsOneOf(vpnProvider,
		providers.Custom,
		providers.Ivpn,
		providers.Mullvad,
		providers.Surfshark,
		providers.Windscribe,
	) {
		// do not validate for VPN provider not supporting Wireguard
		return nil
	}

	// Validate PrivateKey
	if *w.PrivateKey == "" {
		return fmt.Errorf("%w", ErrWireguardPrivateKeyNotSet)
	}
	_, err = wgtypes.ParseKey(*w.PrivateKey)
	if err != nil {
		return fmt.Errorf("private key is not valid: %w", err)
	}

	if vpnProvider == providers.Airvpn {
		if *w.PreSharedKey == "" {
			return fmt.Errorf("%w", ErrWireguardPreSharedKeyNotSet)
		}
	}

	// Validate PreSharedKey
	if *w.PreSharedKey != "" { // Note: this is optional
		_, err = wgtypes.ParseKey(*w.PreSharedKey)
		if err != nil {
			return fmt.Errorf("pre-shared key is not valid: %w", err)
		}
	}

	// Validate Addresses
	if len(w.Addresses) == 0 {
		return fmt.Errorf("%w", ErrWireguardInterfaceAddressNotSet)
	}
	for i, ipNet := range w.Addresses {
		if !ipNet.IsValid() {
			return fmt.Errorf("%w: for address at index %d: %s",
				ErrWireguardInterfaceAddressNotSet, i, ipNet.String())
		}

		if !ipv6Supported && ipNet.Addr().Is6() {
			return fmt.Errorf("%w: address %s",
				ErrWireguardInterfaceAddressIPv6, ipNet)
		}
	}

	// Validate interface
	if !regexpInterfaceName.MatchString(w.Interface) {
		return fmt.Errorf("%w: '%s' does not match regex '%s'",
			ErrWireguardInterfaceNotValid, w.Interface, regexpInterfaceName)
	}

	validImplementations := []string{"auto", "userspace", "kernelspace"}
	if !helpers.IsOneOf(w.Implementation, validImplementations...) {
		return fmt.Errorf("%w: %s must be one of %s", ErrWireguardImplementationNotValid,
			w.Implementation, helpers.ChoicesOrString(validImplementations))
	}

	return nil
}

func (w *Wireguard) copy() (copied Wireguard) {
	return Wireguard{
		PrivateKey:     helpers.CopyStringPtr(w.PrivateKey),
		PreSharedKey:   helpers.CopyStringPtr(w.PreSharedKey),
		Addresses:      helpers.CopyNetipPrefixesSlice(w.Addresses),
		Interface:      w.Interface,
		Implementation: w.Implementation,
	}
}

func (w *Wireguard) mergeWith(other Wireguard) {
	w.PrivateKey = helpers.MergeWithStringPtr(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = helpers.MergeWithStringPtr(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = helpers.MergeNetipPrefixesSlices(w.Addresses, other.Addresses)
	w.Interface = helpers.MergeWithString(w.Interface, other.Interface)
	w.Implementation = helpers.MergeWithString(w.Implementation, other.Implementation)
}

func (w *Wireguard) overrideWith(other Wireguard) {
	w.PrivateKey = helpers.OverrideWithStringPtr(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = helpers.OverrideWithStringPtr(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = helpers.OverrideWithNetipPrefixesSlice(w.Addresses, other.Addresses)
	w.Interface = helpers.OverrideWithString(w.Interface, other.Interface)
	w.Implementation = helpers.OverrideWithString(w.Implementation, other.Implementation)
}

func (w *Wireguard) setDefaults() {
	w.PrivateKey = helpers.DefaultStringPtr(w.PrivateKey, "")
	w.PreSharedKey = helpers.DefaultStringPtr(w.PreSharedKey, "")
	w.Interface = helpers.DefaultString(w.Interface, "wg0")
	w.Implementation = helpers.DefaultString(w.Implementation, "auto")
}

func (w Wireguard) String() string {
	return w.toLinesNode().String()
}

func (w Wireguard) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Wireguard settings:")

	if *w.PrivateKey != "" {
		s := helpers.ObfuscateWireguardKey(*w.PrivateKey)
		node.Appendf("Private key: %s", s)
	}

	if *w.PreSharedKey != "" {
		s := helpers.ObfuscateWireguardKey(*w.PreSharedKey)
		node.Appendf("Pre-shared key: %s", s)
	}

	addressesNode := node.Appendf("Interface addresses:")
	for _, address := range w.Addresses {
		addressesNode.Appendf(address.String())
	}

	node.Appendf("Network interface: %s", w.Interface)

	if w.Implementation != "auto" {
		node.Appendf("Implementation: %s", w.Implementation)
	}

	return node
}
