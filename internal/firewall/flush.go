package firewall

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/firewall/iptables"
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
			if errors.Is(err, iptables.ErrKernelModuleMissing) {
				c.logger.Debugf("falling back to killing connections for one second because marking packets failed: %s", err)
				return c.rejectOutputTrafficTemporarily(ctx)
			}
			return fmt.Errorf("accepting only new output public traffic: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("flushing conntrack: %w", err)
	}
}

func (c *Config) rejectOutputTrafficTemporarily(ctx context.Context) error {
	remove := false
	err := c.impl.RejectOutputPublicTraffic(ctx, remove)
	if err != nil {
		return fmt.Errorf("rejecting only new output public traffic: %w", err)
	}
	timer := time.NewTimer(time.Second)
	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	}
	remove = true
	// Use [context.Background] to make sure this is removed, even if the context
	// passed to this function is canceled.
	err = c.impl.RejectOutputPublicTraffic(context.Background(), remove)
	if err != nil {
		return fmt.Errorf("reverting rejecting only new output public traffic: %w", err)
	}
	return nil
}
