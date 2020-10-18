package provider

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type timeNowFunc func() time.Time

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

func pickRandomConnection(connections []models.OpenVPNConnection, source rand.Source) models.OpenVPNConnection {
	return connections[rand.New(source).Intn(len(connections))] //nolint:gosec
}

func filterByPossibilities(value string, possibilities []string) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}
	for _, possibility := range possibilities {
		if strings.EqualFold(value, possibility) {
			return false
		}
	}
	return true
}

func commaJoin(slice []string) string {
	return strings.Join(slice, ",")
}
