package params

import (
	"fmt"
	"net"
	"strings"
)

// GetOutboundSubnets obtains the CIDR subnets from the comma separated list of the
// environment variable FIREWALL_OUTBOUND_SUBNETS.
func (r *reader) GetOutboundSubnets() (outboundSubnets []net.IPNet, err error) {
	const key = "FIREWALL_OUTBOUND_SUBNETS"
	s, err := r.envParams.GetEnv(key)
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
