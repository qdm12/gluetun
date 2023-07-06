package wireguard

import (
	"errors"
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Settings struct {
	// Interface name for the Wireguard interface.
	// It defaults to wg0 if unset.
	InterfaceName string
	// Private key in base 64 format
	PrivateKey string
	// Public key in base 64 format
	PublicKey string
	// Pre shared key in base 64 format
	PreSharedKey string
	// Wireguard server endpoint to connect to.
	Endpoint netip.AddrPort
	// Addresses assigned to the client.
	// Note IPv6 addresses are ignored if IPv6 is not supported.
	Addresses []netip.Prefix
	// AllowedIPs is the IP networks to be routed through
	// the Wireguard interface.
	// Note IPv6 addresses are ignored if IPv6 is not supported.
	AllowedIPs []netip.Prefix
	// FirewallMark to be used in routing tables and IP rules.
	// It defaults to 51820 if left to 0.
	FirewallMark int
	// Maximum Transmission Unit (MTU) setting for the network interface.
	// It defaults to device.DefaultMTU from wireguard-go which is 1420
	MTU uint16
	// RulePriority is the priority for the rule created with the
	// FirewallMark.
	RulePriority int
	// IPv6 can bet set to true if IPv6 should be handled.
	// It defaults to false if left unset.
	IPv6 *bool
	// Implementation is the implementation to use.
	// It can be auto, kernelspace or userspace, and defaults to auto.
	Implementation string
}

func (s *Settings) SetDefaults() {
	if s.InterfaceName == "" {
		const defaultInterfaceName = "wg0"
		s.InterfaceName = defaultInterfaceName
	}

	if s.Endpoint.IsValid() && s.Endpoint.Port() == 0 {
		const defaultPort = 51820
		s.Endpoint = netip.AddrPortFrom(s.Endpoint.Addr(), defaultPort)
	}

	if s.FirewallMark == 0 {
		const defaultFirewallMark = 51820
		s.FirewallMark = defaultFirewallMark
	}

	if s.MTU == 0 {
		s.MTU = device.DefaultMTU
	}

	if s.IPv6 == nil {
		ipv6 := false // this should be injected from host
		s.IPv6 = &ipv6
	}

	if len(s.AllowedIPs) == 0 {
		s.AllowedIPs = append(s.AllowedIPs, allIPv4())
		if *s.IPv6 {
			s.AllowedIPs = append(s.AllowedIPs, allIPv6())
		}
	}

	if s.Implementation == "" {
		const defaultImplementation = "auto"
		s.Implementation = defaultImplementation
	}
}

var (
	ErrInterfaceNameInvalid    = errors.New("invalid interface name")
	ErrPrivateKeyMissing       = errors.New("private key is missing")
	ErrPrivateKeyInvalid       = errors.New("cannot parse private key")
	ErrPublicKeyMissing        = errors.New("public key is missing")
	ErrPublicKeyInvalid        = errors.New("cannot parse public key")
	ErrPreSharedKeyInvalid     = errors.New("cannot parse pre-shared key")
	ErrEndpointAddrMissing     = errors.New("endpoint address is missing")
	ErrEndpointPortMissing     = errors.New("endpoint port is missing")
	ErrAddressMissing          = errors.New("interface address is missing")
	ErrAddressNotValid         = errors.New("interface address is not valid")
	ErrAllowedIPsMissing       = errors.New("allowed IPs are missing")
	ErrAllowedIPNotValid       = errors.New("allowed IP is not valid")
	ErrAllowedIPv6NotSupported = errors.New("allowed IPv6 address not supported")
	ErrFirewallMarkMissing     = errors.New("firewall mark is missing")
	ErrMTUMissing              = errors.New("MTU is missing")
	ErrImplementationInvalid   = errors.New("invalid implementation")
)

var interfaceNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Settings) Check() (err error) {
	if !interfaceNameRegexp.MatchString(s.InterfaceName) {
		return fmt.Errorf("%w: %s", ErrInterfaceNameInvalid, s.InterfaceName)
	}

	if s.PrivateKey == "" {
		return fmt.Errorf("%w", ErrPrivateKeyMissing)
	} else if _, err := wgtypes.ParseKey(s.PrivateKey); err != nil {
		return fmt.Errorf("%w", ErrPrivateKeyInvalid)
	}

	if s.PublicKey == "" {
		return fmt.Errorf("%w", ErrPublicKeyMissing)
	} else if _, err := wgtypes.ParseKey(s.PublicKey); err != nil {
		return fmt.Errorf("%w: %s", ErrPublicKeyInvalid, s.PublicKey)
	}

	if s.PreSharedKey != "" {
		if _, err := wgtypes.ParseKey(s.PreSharedKey); err != nil {
			return fmt.Errorf("%w", ErrPreSharedKeyInvalid)
		}
	}

	switch {
	case !s.Endpoint.Addr().IsValid():
		return fmt.Errorf("%w", ErrEndpointAddrMissing)
	case s.Endpoint.Port() == 0:
		return fmt.Errorf("%w", ErrEndpointPortMissing)
	}

	if len(s.Addresses) == 0 {
		return fmt.Errorf("%w", ErrAddressMissing)
	}
	for i, addr := range s.Addresses {
		if !addr.IsValid() {
			return fmt.Errorf("%w: for address %d of %d",
				ErrAddressNotValid, i+1, len(s.Addresses))
		}
	}

	if len(s.AllowedIPs) == 0 {
		return fmt.Errorf("%w", ErrAllowedIPsMissing)
	}
	for i, allowedIP := range s.AllowedIPs {
		switch {
		case !allowedIP.IsValid():
			return fmt.Errorf("%w: for allowed IP %d of %d",
				ErrAllowedIPNotValid, i+1, len(s.AllowedIPs))
		case allowedIP.Addr().Is6() && !*s.IPv6:
			return fmt.Errorf("%w: for allowed IP %s",
				ErrAllowedIPv6NotSupported, allowedIP)
		}
	}

	if s.FirewallMark == 0 {
		return fmt.Errorf("%w", ErrFirewallMarkMissing)
	}

	if s.MTU == 0 {
		return fmt.Errorf("%w", ErrMTUMissing)
	}

	switch s.Implementation {
	case "auto", "kernelspace", "userspace":
	default:
		return fmt.Errorf("%w: %s", ErrImplementationInvalid, s.Implementation)
	}

	return nil
}

func (s Settings) String() string {
	lines := s.ToLines(ToLinesSettings{})
	return strings.Join(lines, "\n")
}

type ToLinesSettings struct {
	// Indent defaults to 4 spaces "    ".
	Indent *string
	// FieldPrefix defaults to "├── ".
	FieldPrefix *string
	// LastFieldPrefix defaults to "└── ".
	LastFieldPrefix *string
}

func (settings *ToLinesSettings) setDefaults() {
	toStringPtr := func(s string) *string { return &s }
	if settings.Indent == nil {
		settings.Indent = toStringPtr("    ")
	}
	if settings.FieldPrefix == nil {
		settings.FieldPrefix = toStringPtr("├── ")
	}
	if settings.LastFieldPrefix == nil {
		settings.LastFieldPrefix = toStringPtr("└── ")
	}
}

// ToLines serializes the settings to a slice of strings for display.
func (s Settings) ToLines(settings ToLinesSettings) (lines []string) {
	settings.setDefaults()

	indent := *settings.Indent
	fieldPrefix := *settings.FieldPrefix
	lastFieldPrefix := *settings.LastFieldPrefix

	lines = append(lines, fieldPrefix+"Interface name: "+s.InterfaceName)
	const (
		set    = "set"
		notSet = "not set"
	)

	isSet := notSet
	if s.PrivateKey != "" {
		isSet = set
	}
	lines = append(lines, fieldPrefix+"Private key: "+isSet)

	if s.PublicKey != "" {
		lines = append(lines, fieldPrefix+"PublicKey: "+s.PublicKey)
	}

	isSet = notSet
	if s.PreSharedKey != "" {
		isSet = set
	}
	lines = append(lines, fieldPrefix+"Pre shared key: "+isSet)

	endpointStr := notSet
	if s.Endpoint.Addr().IsValid() {
		endpointStr = s.Endpoint.String()
	}
	lines = append(lines, fieldPrefix+"Endpoint: "+endpointStr)

	ipv6Status := "disabled"
	if *s.IPv6 {
		ipv6Status = "enabled"
	}
	lines = append(lines, fieldPrefix+"IPv6: "+ipv6Status)

	if s.FirewallMark != 0 {
		lines = append(lines, fieldPrefix+"Firewall mark: "+fmt.Sprint(s.FirewallMark))
	}

	if s.MTU != 0 {
		lines = append(lines, fieldPrefix+"MTU: "+fmt.Sprint(s.MTU))
	}

	if s.RulePriority != 0 {
		lines = append(lines, fieldPrefix+"Rule priority: "+fmt.Sprint(s.RulePriority))
	}

	if s.Implementation != "auto" {
		lines = append(lines, fieldPrefix+"Implementation: "+s.Implementation)
	}

	if len(s.Addresses) == 0 {
		lines = append(lines, lastFieldPrefix+"Addresses: "+notSet)
	} else {
		lines = append(lines, lastFieldPrefix+"Addresses:")
		for i, address := range s.Addresses {
			prefix := fieldPrefix
			if i == len(s.Addresses)-1 {
				prefix = lastFieldPrefix
			}
			lines = append(lines, indent+prefix+address.String())
		}
	}

	if len(s.AllowedIPs) > 0 {
		lines = append(lines, fieldPrefix+"Allowed IPs:")
		for i, allowedIP := range s.AllowedIPs {
			prefix := fieldPrefix
			if i == len(s.AllowedIPs)-1 {
				prefix = lastFieldPrefix
			}
			lines = append(lines, indent+prefix+allowedIP.String())
		}
	}

	return lines
}
