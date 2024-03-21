package firewall

import (
	"fmt"
	"net/netip"
)

type Logger interface {
	Debug(s string)
	Info(s string)
	Error(s string)
}

func (c *Config) logIgnoredSubnetFamily(subnet netip.Prefix) {
	c.logger.Info(fmt.Sprintf("ignoring subnet %s which has "+
		"no default route matching its family", subnet))
}
