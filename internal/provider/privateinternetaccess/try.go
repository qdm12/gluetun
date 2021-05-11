package privateinternetaccess

import (
	"context"
	"time"

	"github.com/qdm12/golibs/logging"
)

func tryUntilSuccessful(ctx context.Context, logger logging.Logger, fn func() error) {
	const initialRetryPeriod = 5 * time.Second
	retryPeriod := initialRetryPeriod
	for {
		err := fn()
		if err == nil {
			break
		}
		logger.Error(err)
		logger.Info("Trying again in " + retryPeriod.String())
		timer := time.NewTimer(retryPeriod)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
		retryPeriod *= 2
	}
}
