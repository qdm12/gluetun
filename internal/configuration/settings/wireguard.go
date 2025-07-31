package settings

import (
	"fmt"
	"net/netip"
	"regexp"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Wireguard contains settings to configure the Wireguard client.
type Wireguard struct {
	// CustomConfigFile is the path to a custom Wireguard
	// configuration file to use. It is optional.
	// If it is used, settings such as the private key,
	// public key and addresses are ignored.
	CustomConfigFile *string `json:"custom_config_file"`
	// PrivateKey is the Wireguard client peer private key.
	// It is mandatory if CustomConfigFile is not set.
	PrivateKey *string `json:"private_key"`
	// PublicKey is the Wireguard server peer public key.
	// It is mandatory if CustomConfigFile is not set.
	PublicKey *string `json:"public_key"`
	// PreSharedKey is the Wireguard pre-shared key.
	// It is optional.
	PreSharedKey *string `json:"preshared_key"`
	// EndpointIP is the server endpoint IP address.
	// It is mandatory if CustomConfigFile is not set.
	EndpointIP netip.Addr `json:"endpoint_ip"`
	// EndpointPort is the server endpoint port.
	// It is mandatory if CustomConfigFile is not set.
	EndpointPort *uint16 `json:"endpoint_port"`
	// Addresses are the client tunnel IP addresses.
	// It is mandatory if CustomConfigFile is not set.
	Addresses []netip.Prefix `json:"addresses"`
	// AllowedIPs are the allowed IP addresses for the server peer.
	// It defaults to 0.0.0.0/0 and ::/0.
	// It is optional.
	AllowedIPs []netip.Prefix `json:"allowed_ips"`
	// PersistentKeepalive is the persistent keepalive interval.
	// It is optional.
	PersistentKeepalive *time.Duration `json:"persistent_keepalive"`
	// Interface is the name of the Wireguard interface to create.
	// It defaults to wg0.
	// It is optional.
	Interface *string `json:"interface"`
}

var regexpInterfaceName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// Validate validates Wireguard settings.
// It should only be ran if the VPN type chosen is Wireguard.
func (w Wireguard) validate(vpnProvider string, ipv6Supported bool) (err error) {
	if !helpers.IsOneOf(vpnProvider,
		providers.Airvpn,
		providers.Custom,
		providers.Fastestvpn,
		providers.Ivpn,
		providers.Mullvad,
		providers.Nordvpn,
		providers.Protonvpn,
		providers.Surfshark,
		providers.Windscribe,
	) {
		// do not validate for VPN provider not supporting Wireguard
		return nil
	}

	if w.CustomConfigFile != nil && *w.CustomConfigFile != "" {
		// skip validation if a custom config file is used
		fmt.Printf("DEBUG: Skipping validation, using custom config file: %s\n", *w.CustomConfigFile)
		return nil
	}

	fmt.Printf("DEBUG: Not using custom config, proceeding with validation. CustomConfigFile: %v\n", w.CustomConfigFile)

	// Validate PrivateKey
	if w.PrivateKey == nil || *w.PrivateKey == "" {
		return fmt.Errorf("wireguard private key is not set")
	}
	_, err = wgtypes.ParseKey(*w.PrivateKey)
	if err != nil {
		err = fmt.Errorf("private key is not valid: %w", err)
		if vpnProvider == providers.Nordvpn &&
			err.Error() == "wgtypes: incorrect key size: 48" {
			err = fmt.Errorf("%w - you might be using your access token instead of the Wireguard private key", err)
		}
		return err
	}

	// Validate PublicKey
	if w.PublicKey == nil || *w.PublicKey == "" {
		return fmt.Errorf("wireguard public key is not set")
	}
	_, err = wgtypes.ParseKey(*w.PublicKey)
	if err != nil {
		return fmt.Errorf("public key is not valid: %w", err)
	}

	if vpnProvider == providers.Airvpn {
		if w.PreSharedKey == nil || *w.PreSharedKey == "" {
			return fmt.Errorf("%w", ErrWireguardPreSharedKeyNotSet)
		}
	}

	// Validate PreSharedKey
	if w.PreSharedKey != nil && *w.PreSharedKey != "" { // Note: this is optional
		_, err = wgtypes.ParseKey(*w.PreSharedKey)
		if err != nil {
			return fmt.Errorf("pre-shared key is not valid: %w", err)
		}
	}

	// Validate EndpointIP
	if !w.EndpointIP.IsValid() {
		return fmt.Errorf("wireguard endpoint IP is not set")
	}

	// Validate Addresses
	if len(w.Addresses) == 0 {
		return fmt.Errorf("%w", ErrWireguardInterfaceAddressNotSet)
	}
	for i, ipNet := range w.Addresses {
		if !ipNet.IsValid() {
			return fmt.Errorf("%w: for address at index %d",
				ErrWireguardInterfaceAddressNotSet, i)
		}

		if !ipv6Supported && ipNet.Addr().Is6() {
			return fmt.Errorf("%w: address %s",
				ErrWireguardInterfaceAddressIPv6, ipNet.String())
		}
	}

	// Validate AllowedIPs
	if len(w.AllowedIPs) == 0 {
		return fmt.Errorf("%w", ErrWireguardAllowedIPsNotSet)
	}
	for i, allowedIP := range w.AllowedIPs {
		if !allowedIP.IsValid() {
			return fmt.Errorf("%w: for allowed ip %d of %d",
				ErrWireguardAllowedIPNotSet, i+1, len(w.AllowedIPs))
		}
	}

	// Validate interface
	if w.Interface == nil || !regexpInterfaceName.MatchString(*w.Interface) {
		return fmt.Errorf("%w: '%v' does not match regex '%s'",
			ErrWireguardInterfaceNotValid, w.Interface, regexpInterfaceName)
	}

	return nil
}

func (w *Wireguard) copy() (copied Wireguard) {
	return Wireguard{
		CustomConfigFile:      gosettings.CopyPointer(w.CustomConfigFile),
		PrivateKey:            gosettings.CopyPointer(w.PrivateKey),
		PublicKey:             gosettings.CopyPointer(w.PublicKey),
		PreSharedKey:          gosettings.CopyPointer(w.PreSharedKey),
		EndpointIP:            w.EndpointIP,
		EndpointPort:          gosettings.CopyPointer(w.EndpointPort),
		Addresses:             gosettings.CopySlice(w.Addresses),
		AllowedIPs:            gosettings.CopySlice(w.AllowedIPs),
		PersistentKeepalive:   gosettings.CopyPointer(w.PersistentKeepalive),
		Interface:             gosettings.CopyPointer(w.Interface),
	}
}

func (w *Wireguard) mergeWith(other Wireguard) {
	w.CustomConfigFile = gosettings.OverrideWithPointer(w.CustomConfigFile, other.CustomConfigFile)
	w.PrivateKey = gosettings.OverrideWithPointer(w.PrivateKey, other.PrivateKey)
	w.PublicKey = gosettings.OverrideWithPointer(w.PublicKey, other.PublicKey)
	w.PreSharedKey = gosettings.OverrideWithPointer(w.PreSharedKey, other.PreSharedKey)
	if other.EndpointIP.IsValid() {
		w.EndpointIP = other.EndpointIP
	}
	w.EndpointPort = gosettings.OverrideWithPointer(w.EndpointPort, other.EndpointPort)
	w.Addresses = gosettings.OverrideWithSlice(w.Addresses, other.Addresses)
	w.AllowedIPs = gosettings.OverrideWithSlice(w.AllowedIPs, other.AllowedIPs)
	w.PersistentKeepalive = gosettings.OverrideWithPointer(w.PersistentKeepalive, other.PersistentKeepalive)
	w.Interface = gosettings.OverrideWithPointer(w.Interface, other.Interface)
}

func (w *Wireguard) overrideWith(other Wireguard) {
	w.CustomConfigFile = gosettings.OverrideWithPointer(w.CustomConfigFile, other.CustomConfigFile)
	w.PrivateKey = gosettings.OverrideWithPointer(w.PrivateKey, other.PrivateKey)
	w.PublicKey = gosettings.OverrideWithPointer(w.PublicKey, other.PublicKey)
	w.PreSharedKey = gosettings.OverrideWithPointer(w.PreSharedKey, other.PreSharedKey)
	if other.EndpointIP.IsValid() {
		w.EndpointIP = other.EndpointIP
	}
	w.EndpointPort = gosettings.OverrideWithPointer(w.EndpointPort, other.EndpointPort)
	w.Addresses = gosettings.OverrideWithSlice(w.Addresses, other.Addresses)
	w.AllowedIPs = gosettings.OverrideWithSlice(w.AllowedIPs, other.AllowedIPs)
	w.PersistentKeepalive = gosettings.OverrideWithPointer(w.PersistentKeepalive, other.PersistentKeepalive)
	w.Interface = gosettings.OverrideWithPointer(w.Interface, other.Interface)
}

func (w *Wireguard) setDefaults(vpnProvider string) {
	w.CustomConfigFile = gosettings.DefaultPointer(w.CustomConfigFile, "")
	w.PrivateKey = gosettings.DefaultPointer(w.PrivateKey, "")
	w.PublicKey = gosettings.DefaultPointer(w.PublicKey, "")
	w.PreSharedKey = gosettings.DefaultPointer(w.PreSharedKey, "")
	w.AllowedIPs = gosettings.DefaultSlice(w.AllowedIPs,
		[]netip.Prefix{
			netip.MustParsePrefix("0.0.0.0/0"),
			netip.MustParsePrefix("::/0"),
		})
	w.Interface = gosettings.DefaultPointer(w.Interface, "wg0")
}

func (w Wireguard) String() string {
	return w.toLinesNode().String()
}

func (w Wireguard) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Wireguard settings:")

	if w.CustomConfigFile != nil && *w.CustomConfigFile != "" {
		node.Append(fmt.Sprintf("Custom configuration file: %s", *w.CustomConfigFile))
		return node
	}

	if w.PrivateKey != nil && *w.PrivateKey != "" {
		s := gosettings.ObfuscateKey(*w.PrivateKey)
		node.Append(fmt.Sprintf("Private key: %s", s))
	}

	if w.PublicKey != nil && *w.PublicKey != "" {
		node.Append(fmt.Sprintf("Public key: %s", *w.PublicKey))
	}

	if w.PreSharedKey != nil && *w.PreSharedKey != "" {
		s := gosettings.ObfuscateKey(*w.PreSharedKey)
		node.Append(fmt.Sprintf("Pre-shared key: %s", s))
	}

	if w.EndpointIP.IsValid() {
		node.Append(fmt.Sprintf("Endpoint IP address: %s", w.EndpointIP))
	}

	if w.EndpointPort != nil {
		node.Append(fmt.Sprintf("Endpoint port: %d", *w.EndpointPort))
	}

	if len(w.Addresses) > 0 {
		node.Append(fmt.Sprintf("Tunnel addresses: %s", commaJoin(w.Addresses)))
	}

	if len(w.AllowedIPs) > 0 {
		node.Append(fmt.Sprintf("Allowed IPs: %s", commaJoin(w.AllowedIPs)))
	}

	if w.PersistentKeepalive != nil {
		node.Append(fmt.Sprintf("Persistent keepalive: %s", *w.PersistentKeepalive))
	}

	if w.Interface != nil {
		node.Append(fmt.Sprintf("Interface: %s", *w.Interface))
	}

	return node
}

func (w *Wireguard) read(r *reader.Reader) (err error) {
	w.CustomConfigFile = r.Get("WIREGUARD_CUSTOM_CONFIG_FILE")
	w.PrivateKey = r.Get("WIREGUARD_PRIVATE_KEY", reader.ForceLowercase(false))
	w.PublicKey = r.Get("WIREGUARD_PUBLIC_KEY", reader.ForceLowercase(false))
	w.PreSharedKey = r.Get("WIREGUARD_PRESHARED_KEY", reader.ForceLowercase(false))

	w.EndpointIP, err = r.NetipAddr("WIREGUARD_ENDPOINT_IP")
	if err != nil {
		return err
	}

	w.EndpointPort, err = r.Uint16Ptr("WIREGUARD_ENDPOINT_PORT")
	if err != nil {
		return err
	}

	// Parse comma-separated addresses for Addresses and AllowedIPs
	addressesStr := r.Get("WIREGUARD_ADDRESSES")
	if addressesStr != nil && *addressesStr != "" {
		addressList := strings.Split(*addressesStr, ",")
		w.Addresses = make([]netip.Prefix, 0, len(addressList))
		for _, addr := range addressList {
			addr = strings.TrimSpace(addr)
			if addr == "" {
				continue
			}
			prefix, err := netip.ParsePrefix(addr)
			if err != nil {
				return fmt.Errorf("invalid address prefix: %w", err)
			}
			w.Addresses = append(w.Addresses, prefix)
		}
	}

	allowedIPsStr := r.Get("WIREGUARD_ALLOWED_IPS")
	if allowedIPsStr != nil && *allowedIPsStr != "" {
		allowedList := strings.Split(*allowedIPsStr, ",")
		w.AllowedIPs = make([]netip.Prefix, 0, len(allowedList))
		for _, ip := range allowedList {
			ip = strings.TrimSpace(ip)
			if ip == "" {
				continue
			}
			prefix, err := netip.ParsePrefix(ip)
			if err != nil {
				return fmt.Errorf("invalid allowed ip prefix: %w", err)
			}
			w.AllowedIPs = append(w.AllowedIPs, prefix)
		}
	}

	w.PersistentKeepalive, err = r.DurationPtr("WIREGUARD_PERSISTENT_KEEPALIVE_INTERVAL")
	if err != nil {
		return err
	}

	w.Interface = r.Get("WIREGUARD_INTERFACE")

	return nil
}

func commaJoin[T any](slice []T) string {
	ss := make([]string, len(slice))
	for i, value := range slice {
		ss[i] = fmt.Sprint(value)
	}
	return strings.Join(ss, ",")
}
