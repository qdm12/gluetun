package provider

import (
	"context"
	"time"

	"github.com/qdm12/golibs/logging"
)

func tryUntilSuccessful(ctx context.Context, logger logging.Logger, fn func() error) {
	const retryPeriod = 10 * time.Second
	for {
		err := fn()
		if err == nil {
			break
		}
		logger.Error(err)
		logger.Info("Trying again in %s", retryPeriod)
		timer := time.NewTimer(retryPeriod)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}
}
