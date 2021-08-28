package configuration

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/params"
)

// Wireguard contains settings to configure the Wireguard client.
type Wireguard struct {
	PrivateKey   string       `json:"privatekey"`
	PreSharedKey string       `json:"presharedkey"`
	Addresses    []*net.IPNet `json:"addresses"`
	Interface    string       `json:"interface"`
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

	if len(settings.Addresses) > 0 {
		lines = append(lines, indent+lastIndent+"Addresses: ")
		for _, address := range settings.Addresses {
			lines = append(lines, indent+indent+lastIndent+address.String())
		}
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

	err = settings.readAddresses(r.env)
	if err != nil {
		return err
	}

	settings.Interface, err = r.env.Get("WIREGUARD_INTERFACE", params.Default("wg0"))
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_INTERFACE: %w", err)
	}

	return nil
}

func (settings *Wireguard) readAddresses(env params.Interface) (err error) {
	addressStrings, err := env.CSV("WIREGUARD_ADDRESS", params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_ADDRESS: %w", err)
	}

	for _, addressString := range addressStrings {
		ip, ipNet, err := net.ParseCIDR(addressString)
		if err != nil {
			return fmt.Errorf("environment variable WIREGUARD_ADDRESS: address %s: %w", addressString, err)
		}
		ipNet.IP = ip
		settings.Addresses = append(settings.Addresses, ipNet)
	}

	return nil
}
