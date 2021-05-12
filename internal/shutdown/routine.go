package shutdown

import (
	"context"
	"fmt"
	"time"
)

type routine struct {
	name    string
	cancel  context.CancelFunc
	done    <-chan struct{}
	timeout time.Duration
}

func (r *routine) shutdown(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	r.cancel()

	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("for routine %q: %w", r.name, ctx.Err())
	}
}
