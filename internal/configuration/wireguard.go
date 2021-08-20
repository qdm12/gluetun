package configuration

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/params"
)

// Wireguard contains settings to configure the Wireguard client.
type Wireguard struct {
	PrivateKey   string     `json:"privatekey"`
	PreSharedKey string     `json:"presharedkey"`
	Address      *net.IPNet `json:"address"`
	Interface    string     `json:"interface"`
}

func (settings *Wireguard) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Wireguard) lines() (lines []string) {
	lines = append(lines, lastIndent+"Wireguard:")

	lines = append(lines, indent+lastIndent+"Network interface: "+settings.Interface)

	if settings.PrivateKey != "" {
		lines = append(lines, indent+lastIndent+"Private key is set")
	}

	if settings.PreSharedKey != "" {
		lines = append(lines, indent+lastIndent+"Pre-shared key is set")
	}

	if settings.Address != nil {
		lines = append(lines, indent+lastIndent+"Address: "+settings.Address.String())
	}

	return lines
}

func (settings *Wireguard) read(r reader) (err error) {
	settings.PrivateKey, err = r.env.Get("WIREGUARD_PRIVATE_KEY",
		params.CaseSensitiveValue(), params.Unset(), params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_PRIVATE_KEY: %w", err)
	}

	settings.PreSharedKey, err = r.env.Get("WIREGUARD_PRESHARED_KEY",
		params.CaseSensitiveValue(), params.Unset())
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_PRESHARED_KEY: %w", err)
	}

	addressString, err := r.env.Get("WIREGUARD_ADDRESS", params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_ADDRESS: %w", err)
	}
	ip, ipNet, err := net.ParseCIDR(addressString)
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_ADDRESS: %w", err)
	}
	ipNet.IP = ip
	settings.Address = ipNet

	settings.Interface, err = r.env.Get("WIREGUARD_INTERFACE", params.Default("wg0"))
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_INTERFACE: %w", err)
	}

	return nil
}
