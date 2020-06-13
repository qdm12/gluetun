package params

import (
	"fmt"
	"net"
	"strings"

	libparams "github.com/qdm12/golibs/params"
)

// GetFirewall obtains if the firewall should be enabled from the environment variable FIREWALL
func (r *reader) GetFirewall() (enabled bool, err error) {
	return r.envParams.GetOnOff("FIREWALL", libparams.Default("on"))
}

// GetExtraSubnets obtains the CIDR subnets from the comma separated list of the
// environment variable EXTRA_SUBNETS
func (r *reader) GetExtraSubnets() (extraSubnets []net.IPNet, err error) {
	s, err := r.envParams.GetEnv("EXTRA_SUBNETS")
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}
	subnets := strings.Split(s, ",")
	for _, subnet := range subnets {
		_, cidr, err := net.ParseCIDR(subnet)
		if err != nil {
			return nil, fmt.Errorf("could not parse subnet %q from environment variable with key EXTRA_SUBNETS: %w", subnet, err)
		} else if cidr == nil {
			return nil, fmt.Errorf("parsing subnet %q resulted in a nil CIDR", subnet)
		}
		extraSubnets = append(extraSubnets, *cidr)
	}
	return extraSubnets, nil
}
