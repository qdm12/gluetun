package firewall

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/firewall/iptables"
	"github.com/qdm12/gluetun/internal/netlink"
)

func (c *Config) flushExistingConnections(ctx context.Context) error {
	tries := []struct {
		name string
		f    func(ctx context.Context) error
	}{
		{name: "flushing conntrack", f: func(_ context.Context) error {
			return c.netlinker.FlushConntrack()
		}},
		{name: "marking and filtering unmarked packets", f: c.impl.AcceptOutputPublicOnlyNewTraffic},
		{name: "rejecting connections for one second", f: c.rejectOutputTrafficTemporarily},
		{name: "dropping connections for one second", f: c.dropOutputTrafficTemporarily},
	}
	errs := make([]error, 0, len(tries))
	for i, try := range tries {
		if i > 0 {
			c.logger.Debugf("falling back to %s because %s failed: %s", try.name, tries[i-1].name, errs[i-1])
		}
		err := try.f(ctx)
		if err == nil {
			return nil
		}
		err = fmt.Errorf("%s: %w", try.name, err)
		if !errors.Is(err, iptables.ErrKernelModuleMissing) && !errors.Is(err, netlink.ErrConntrackNetlinkNotSupported) {
			return err
		}
		errs = append(errs, err)
	}
	return fmt.Errorf("all tries failed: %v", errs) //nolint:err113
}

func (c *Config) rejectOutputTrafficTemporarily(ctx context.Context) error {
	return setupThenRevert(ctx, c.impl.RejectOutputPublicTraffic)
}

func (c *Config) dropOutputTrafficTemporarily(ctx context.Context) error {
	return setupThenRevert(ctx, c.impl.DropOutputPublicTraffic)
}

// setupThenRevert is a helper function to run a setup function that takes a remove boolean argument,
// and then run the same function with remove set to true after one second or when the context is canceled,
// whichever comes first.
func setupThenRevert(ctx context.Context, f func(ctx context.Context, remove bool) error) error {
	remove := false
	err := f(ctx, remove)
	if err != nil {
		return fmt.Errorf("setting up: %w", err)
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
	err = f(context.Background(), remove)
	if err != nil {
		return fmt.Errorf("reverting: %w", err)
	}
	return nil
}
