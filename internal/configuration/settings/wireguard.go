package settings

import (
	"fmt"
	"net"
	"regexp"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Wireguard contains settings to configure the Wireguard client.
type Wireguard struct {
	// PrivateKey is the Wireguard client peer private key.
	// It cannot be nil in the internal state.
	PrivateKey *string `json:"privatekey,omitempty"`
	// PreSharedKey is the Wireguard pre-shared key.
	// It can be the empty string to indicate there
	// is no pre-shared key.
	// It cannot be nil in the internal state.
	PreSharedKey *string `json:"presharedkey,omitempty"`
	// Addresses are the Wireguard interface addresses.
	Addresses []net.IPNet `json:"addresses,omitempty"`
	// Interface is the name of the Wireguard interface
	// to create. It cannot be the empty string in the
	// internal state.
	Interface string `json:"interface,omitempty"`
}

var regexpInterfaceName = regexp.MustCompile(`^[a-z_0-9]+$`)

// Validate validates Wireguard settings.
// It should only be ran if the VPN type chosen is Wireguard.
func (w Wireguard) validate(vpnProvider string) (err error) {
	if !helpers.IsOneOf(vpnProvider,
		constants.Custom,
		constants.Ivpn,
		constants.Mullvad,
		constants.Windscribe,
	) {
		// do not validate for VPN provider not supporting Wireguard
		return nil
	}

	// Validate PrivateKey
	if *w.PrivateKey == "" {
		return ErrWireguardPrivateKeyNotSet
	}
	_, err = wgtypes.ParseKey(*w.PrivateKey)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrWireguardPrivateKeyNotValid, err)
	}

	// Validate PreSharedKey
	if *w.PreSharedKey != "" { // Note: this is optional
		_, err = wgtypes.ParseKey(*w.PreSharedKey)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrWireguardPreSharedKeyNotValid, err)
		}
	}

	// Validate Addresses
	if len(w.Addresses) == 0 {
		return ErrWireguardInterfaceAddressNotSet
	}
	for i, ipNet := range w.Addresses {
		if ipNet.IP == nil || ipNet.Mask == nil {
			return fmt.Errorf("%w: for address at index %d: %s",
				ErrWireguardInterfaceAddressNotSet, i, ipNet.String())
		}
	}

	// Validate interface
	if regexpInterfaceName.MatchString(w.Interface) {
		return fmt.Errorf("%w: '%s' does not match regex '%s'",
			ErrWireguardInterfaceNotValid, w.Interface, regexpInterfaceName)
	}

	return nil
}

func (w *Wireguard) copy() (copied Wireguard) {
	return Wireguard{
		PrivateKey:   helpers.CopyStringPtr(w.PrivateKey),
		PreSharedKey: helpers.CopyStringPtr(w.PreSharedKey),
		Addresses:    helpers.CopyIPNetSlice(w.Addresses),
		Interface:    w.Interface,
	}
}

func (w *Wireguard) mergeWith(other Wireguard) {
	w.PrivateKey = helpers.MergeWithStringPtr(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = helpers.MergeWithStringPtr(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = helpers.MergeIPNetsSlices(w.Addresses, other.Addresses)
	w.Interface = helpers.MergeWithString(w.Interface, other.Interface)
}

func (w *Wireguard) overrideWith(other Wireguard) {
	w.PrivateKey = helpers.OverrideWithStringPtr(w.PrivateKey, other.PrivateKey)
	w.PreSharedKey = helpers.OverrideWithStringPtr(w.PreSharedKey, other.PreSharedKey)
	w.Addresses = helpers.OverrideWithIPNetsSlice(w.Addresses, other.Addresses)
	w.Interface = helpers.OverrideWithString(w.Interface, other.Interface)
}

func (w *Wireguard) setDefaults() {
	w.PrivateKey = helpers.DefaultStringPtr(w.PrivateKey, "")
	w.PreSharedKey = helpers.DefaultStringPtr(w.PreSharedKey, "")
	w.Interface = helpers.DefaultString(w.Interface, "wg0")
}
