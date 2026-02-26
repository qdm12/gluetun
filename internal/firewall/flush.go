package firewall

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
)

// Note remove is a no-op if conntrack netlink is supported by the kernel.
func (c *Config) flushExistingConnections(ctx context.Context) error {
	err := c.netlinker.FlushConntrack()
	switch {
	case err == nil:
		return nil
	case errors.Is(err, netlink.ErrConntrackNetlinkNotSupported):
		c.logger.Debugf("falling back to marking and filtering unmarked packets because flush conntrack failed: %s", err)
		err = c.impl.AcceptOutputPublicOnlyNewTraffic(ctx)
		if err != nil {
			return fmt.Errorf("accepting only new output public traffic: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("flushing conntrack: %w", err)
	}
}
