package wireguard

import (
	"errors"
	"fmt"
	"net"
	"regexp"

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
}

func (s *Settings) SetDefaults() {
	if s.InterfaceName == "" {
		const defaultInterfaceName = "wg0"
		s.InterfaceName = defaultInterfaceName
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
