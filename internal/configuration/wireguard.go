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
	CustomPort   uint16     `json:"port"`
}

func (settings *Wireguard) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Wireguard) lines() (lines []string) {
	lines = append(lines, lastIndent+"Wireguard:")

	if settings.PrivateKey != "" {
		lines = append(lines, indent+lastIndent+"Private key is set")
	}

	if settings.PreSharedKey != "" {
		lines = append(lines, indent+lastIndent+"Pre-shared key is set")
	}

	if settings.Address != nil {
		lines = append(lines, indent+lastIndent+"Address: "+settings.Address.String())
	}

	if settings.CustomPort != 0 {
		lines = append(lines, indent+lastIndent+"Custom port: "+fmt.Sprint(settings.CustomPort))
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

	settings.CustomPort, err = readPortOrZero(r.env, "WIREGUARD_PORT")
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_PORT: %w", err)
	}

	return nil
}
