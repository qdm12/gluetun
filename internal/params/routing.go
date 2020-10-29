package params

import (
	"fmt"
	"net"
	"strings"

	libparams "github.com/qdm12/golibs/params"
)

// GetOutboundSubnets obtains the CIDR subnets from the comma separated list of the
// environment variable FIREWALL_OUTBOUND_SUBNETS.
func (r *reader) GetOutboundSubnets() (outboundSubnets []net.IPNet, err error) {
	const key = "FIREWALL_OUTBOUND_SUBNETS"
	retroOption := libparams.RetroKeys(
		[]string{"EXTRA_SUBNETS"},
		r.onRetroActive,
	)
	s, err := r.envParams.GetEnv(key, retroOption)
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}
	subnets := strings.Split(s, ",")
	for _, subnet := range subnets {
		_, cidr, err := net.ParseCIDR(subnet)
		if err != nil {
			return nil, fmt.Errorf("cannot parse outbound subnet %q from environment variable with key %s: %w", subnet, key, err)
		} else if cidr == nil {
			return nil, fmt.Errorf("cannot parse outbound subnet %q from environment variable with key %s: subnet is nil",
				subnet, key)
		}
		outboundSubnets = append(outboundSubnets, *cidr)
	}
	return outboundSubnets, nil
}
