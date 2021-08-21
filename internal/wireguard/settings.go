package wireguard

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

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
	Endpoint *net.UDPAddr
	// Addresses assigned to the client.
	Addresses []*net.IPNet
	// FirewallMark to be used in routing tables and IP rules.
	// It defaults to 51820 if left to 0.
	FirewallMark int
	// RulePriority is the priority for the rule created with the
	// FirewallMark.
	RulePriority int
}

func (s *Settings) SetDefaults() {
	if s.InterfaceName == "" {
		const defaultInterfaceName = "wg0"
		s.InterfaceName = defaultInterfaceName
	}

	if s.Endpoint != nil && s.Endpoint.Port == 0 {
		const defaultPort = 51820
		s.Endpoint.Port = defaultPort
	}

	if s.FirewallMark == 0 {
		const defaultFirewallMark = 51820
		s.FirewallMark = defaultFirewallMark
	}
}

var (
	ErrInterfaceNameInvalid = errors.New("invalid interface name")
	ErrPrivateKeyMissing    = errors.New("private key is missing")
	ErrPrivateKeyInvalid    = errors.New("cannot parse private key")
	ErrPublicKeyMissing     = errors.New("public key is missing")
	ErrPublicKeyInvalid     = errors.New("cannot parse public key")
	ErrPreSharedKeyInvalid  = errors.New("cannot parse pre-shared key")
	ErrEndpointMissing      = errors.New("endpoint is missing")
	ErrEndpointIPMissing    = errors.New("endpoint IP is missing")
	ErrEndpointPortMissing  = errors.New("endpoint port is missing")
	ErrAddressMissing       = errors.New("interface address is missing")
	ErrAddressNil           = errors.New("interface address is nil")
	ErrAddressIPMissing     = errors.New("interface address IP is missing")
	ErrAddressMaskMissing   = errors.New("interface address mask is missing")
	ErrFirewallMarkMissing  = errors.New("firewall mark is missing")
)

var interfaceNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Settings) Check() (err error) {
	if !interfaceNameRegexp.MatchString(s.InterfaceName) {
		return fmt.Errorf("%w: %s", ErrInterfaceNameInvalid, s.InterfaceName)
	}

	if s.PrivateKey == "" {
		return ErrPrivateKeyMissing
	} else if _, err := wgtypes.ParseKey(s.PrivateKey); err != nil {
		return ErrPrivateKeyInvalid
	}

	if s.PublicKey == "" {
		return ErrPublicKeyMissing
	} else if _, err := wgtypes.ParseKey(s.PublicKey); err != nil {
		return fmt.Errorf("%w: %s", ErrPublicKeyInvalid, s.PublicKey)
	}

	if s.PreSharedKey != "" {
		if _, err := wgtypes.ParseKey(s.PreSharedKey); err != nil {
			return ErrPreSharedKeyInvalid
		}
	}

	switch {
	case s.Endpoint == nil:
		return ErrEndpointMissing
	case s.Endpoint.IP == nil:
		return ErrEndpointIPMissing
	case s.Endpoint.Port == 0:
		return ErrEndpointPortMissing
	}

	if len(s.Addresses) == 0 {
		return ErrAddressMissing
	}
	for i, addr := range s.Addresses {
		switch {
		case addr == nil:
			return fmt.Errorf("%w: for address %d of %d",
				ErrAddressNil, i+1, len(s.Addresses))
		case addr.IP == nil:
			return fmt.Errorf("%w: for address %d of %d",
				ErrAddressIPMissing, i+1, len(s.Addresses))
		case addr.Mask == nil:
			return fmt.Errorf("%w: for address %d of %d",
				ErrAddressMaskMissing, i+1, len(s.Addresses))
		}
	}

	if s.FirewallMark == 0 {
		return ErrFirewallMarkMissing
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
	if s.Endpoint != nil {
		endpointStr = s.Endpoint.String()
	}
	lines = append(lines, fieldPrefix+"Endpoint: "+endpointStr)

	if s.FirewallMark != 0 {
		lines = append(lines, fieldPrefix+"Firewall mark: "+fmt.Sprint(s.FirewallMark))
	}

	if s.RulePriority != 0 {
		lines = append(lines, fieldPrefix+"Rule priority: "+fmt.Sprint(s.RulePriority))
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

	return lines
}
