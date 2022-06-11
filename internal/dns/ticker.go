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
	if period := *settings.DoT.UpdatePeriod; period > 0 {
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

			status := l.GetStatus()
			if status == constants.Running {
				if err := l.updateFiles(ctx); err != nil {
					l.statusManager.SetStatus(constants.Crashed)
					l.logger.Error(err.Error())
					l.logger.Warn("skipping Unbound restart due to failed files update")
					continue
				}
			}

			_, _ = l.statusManager.ApplyStatus(ctx, constants.Stopped)
			_, _ = l.statusManager.ApplyStatus(ctx, constants.Running)

			settings := l.GetSettings()
			timer.Reset(*settings.DoT.UpdatePeriod)
		case <-l.updateTicker:
			if !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			settings := l.GetSettings()
			newUpdatePeriod := *settings.DoT.UpdatePeriod
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
