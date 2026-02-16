package dns

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) RunRestartTicker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	// Timer that acts as a ticker
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	timerIsStopped := true
	settings := l.GetSettings()
	if period := *settings.UpdatePeriod; period > 0 {
		timer.Reset(period)
		timerIsStopped = false
	}
	lastTick := time.Unix(0, 0)
	for {
		select {
		case <-ctx.Done():
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			lastTick = l.timeNow()
			settings := l.GetSettings()
			if l.GetStatus() == constants.Running {
				if err := l.updateFiles(ctx, settings); err != nil {
					l.logger.Warn("updating block lists failed, skipping: " + err.Error())
				}
			}
			timer.Reset(*settings.UpdatePeriod)
		case <-l.updateTicker:
			if !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			settings := l.GetSettings()
			newUpdatePeriod := *settings.UpdatePeriod
			if newUpdatePeriod == 0 {
				continue
			}
			var waited time.Duration
			if lastTick.UnixNano() != 0 {
				waited = l.timeSince(lastTick)
			}
			leftToWait := newUpdatePeriod - waited
			timer.Reset(leftToWait)
			timerIsStopped = false
		}
	}
}
