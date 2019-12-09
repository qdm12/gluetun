package params

import (
	"fmt"
	"net"
	"strings"

	libparams "github.com/qdm12/golibs/params"
)

// GetExtraSubnets obtains the CIDR subnets from the comma separated list of the
// environment variable EXTRA_SUBNETS
func GetExtraSubnets() (extraSubnets []*net.IPNet, err error) {
	s := libparams.GetEnv("EXTRA_SUBNETS", "")
	if len(s) == 0 {
		return nil, nil
	}
	subnets := strings.Split(s, ",")
	for _, subnet := range subnets {
		_, cidr, err := net.ParseCIDR(subnet)
		if err != nil {
			return nil, fmt.Errorf("could not parse subnet %q from environment variable with key EXTRA_SUBNETS: %w", subnet, err)
		}
		extraSubnets = append(extraSubnets, cidr)
	}
	return extraSubnets, nil
}
