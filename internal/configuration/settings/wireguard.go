package settings

import (
	"fmt"
	"net/netip"
	"regexp"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
	wireguarddevice "golang.zx2c4.com/wireguard/device"
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
	// Maximum Transmission Unit (MTU) of the Wireguard interface.
	// It cannot be zero in the internal state, and defaults to
	// the wireguard-go MTU default of 1420.
	MTU uint16
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
		providers.Airvpn,
		providers.Custom,
		providers.Ivpn,
		providers.Mullvad,
		providers.Nordvpn,
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
	if err := validate.IsOneOf(w.Implementation, validImplementations...); err != nil {
		return fmt.Errorf("%w: %w", ErrWireguardImplementationNotValid, err)
	}

	return nil
}

func (w *Wireguard) copy() (copied Wireguard) {
	return Wireguard{
		PrivateKey:     gosettings.CopyPointer(w.PrivateKey),
		PreSharedKey:   gosettings.CopyPointer(w.PreSharedKey),
		Addresses:      gosettings.CopySlice(w.Addresses),
		Interface:      w.Interface,
		MTU:            w.MTU,
		Implementation: w.Implementation,
	}
}

func (w *Wireguard) mergeWith(other Wireguard) {
	w.PrivateKey = gosettings.MergeWithPointer(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = gosettings.MergeWithPointer(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = gosettings.MergeWithSlice(w.Addresses, other.Addresses)
	w.Interface = gosettings.MergeWithString(w.Interface, other.Interface)
	w.MTU = gosettings.MergeWithNumber(w.MTU, other.MTU)
	w.Implementation = gosettings.MergeWithString(w.Implementation, other.Implementation)
}

func (w *Wireguard) overrideWith(other Wireguard) {
	w.PrivateKey = gosettings.OverrideWithPointer(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = gosettings.OverrideWithPointer(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = gosettings.OverrideWithSlice(w.Addresses, other.Addresses)
	w.Interface = gosettings.OverrideWithString(w.Interface, other.Interface)
	w.MTU = gosettings.OverrideWithNumber(w.MTU, other.MTU)
	w.Implementation = gosettings.OverrideWithString(w.Implementation, other.Implementation)
}

func (w *Wireguard) setDefaults(vpnProvider string) {
	w.PrivateKey = gosettings.DefaultPointer(w.PrivateKey, "")
	w.PreSharedKey = gosettings.DefaultPointer(w.PreSharedKey, "")
	if vpnProvider == providers.Nordvpn {
		defaultNordVPNAddress := netip.AddrFrom4([4]byte{10, 5, 0, 2})
		defaultNordVPNPrefix := netip.PrefixFrom(defaultNordVPNAddress, defaultNordVPNAddress.BitLen())
		w.Addresses = gosettings.DefaultSlice(w.Addresses, []netip.Prefix{defaultNordVPNPrefix})
	}
	w.Interface = gosettings.DefaultString(w.Interface, "wg0")
	w.MTU = gosettings.DefaultNumber(w.MTU, wireguarddevice.DefaultMTU)
	w.Implementation = gosettings.DefaultString(w.Implementation, "auto")
}

func (w Wireguard) String() string {
	return w.toLinesNode().String()
}

func (w Wireguard) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Wireguard settings:")

	if *w.PrivateKey != "" {
		s := gosettings.ObfuscateKey(*w.PrivateKey)
		node.Appendf("Private key: %s", s)
	}

	if *w.PreSharedKey != "" {
		s := gosettings.ObfuscateKey(*w.PreSharedKey)
		node.Appendf("Pre-shared key: %s", s)
	}

	addressesNode := node.Appendf("Interface addresses:")
	for _, address := range w.Addresses {
		addressesNode.Appendf(address.String())
	}

	interfaceNode := node.Appendf("Network interface: %s", w.Interface)
	interfaceNode.Appendf("MTU: %d", w.MTU)

	if w.Implementation != "auto" {
		node.Appendf("Implementation: %s", w.Implementation)
	}

	return node
}
