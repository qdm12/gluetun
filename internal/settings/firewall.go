package settings

import (
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Firewall contains settings to customize the firewall operation
type Firewall struct {
	AllowedSubnets []net.IPNet
}

func (f *Firewall) String() string {
	allowedSubnets := make([]string, len(f.AllowedSubnets))
	for i := range f.AllowedSubnets {
		allowedSubnets[i] = f.AllowedSubnets[i].String()
	}
	settingsList := []string{
		"Firewall settings:",
		"Allowed subnets: " + strings.Join(allowedSubnets, ", "),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetFirewallSettings obtains firewall settings from environment variables using the params package.
func GetFirewallSettings(paramsReader params.Reader) (settings Firewall, err error) {
	settings.AllowedSubnets, err = paramsReader.GetExtraSubnets()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
