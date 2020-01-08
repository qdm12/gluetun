package settings

import (
	"net"
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Firewall contains settings to customize the firewall operation
type Firewall struct {
	AllowedSubnets []*net.IPNet
}

func (f *Firewall) String() string {
	var allowedSubnets []string
	for _, net := range f.AllowedSubnets {
		allowedSubnets = append(allowedSubnets, net.String())
	}
	settingsList := []string{
		"Allowed subnets: " + strings.Join(allowedSubnets, ", "),
	}
	return "Firewall settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetFirewallSettings obtains firewall settings from environment variables using the params package.
func GetFirewallSettings(envParams libparams.EnvParams) (settings Firewall, err error) {
	settings.AllowedSubnets, err = params.GetExtraSubnets(envParams)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
