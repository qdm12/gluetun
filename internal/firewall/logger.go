package firewall

import (
	"fmt"
	"net/netip"
)

func (c *Config) logIgnoredSubnetFamily(subnet netip.Prefix) {
	c.logger.Info(fmt.Sprintf("ignoring subnet %s which has "+
		"no default route matching its family", subnet))
}
