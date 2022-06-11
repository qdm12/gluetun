package publicip

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) RunRestartTicker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	timer := time.NewTimer(time.Hour)
	timer.Stop() // 1 hour, cannot be a race condition
	timerIsStopped := true
	if period := *l.state.GetSettings().Period; period > 0 {
		timerIsStopped = false
		timer.Reset(period)
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
			_, _ = l.statusManager.ApplyStatus(ctx, constants.Running)
			timer.Reset(*l.state.GetSettings().Period)
		case <-l.updateTicker:
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			period := *l.state.GetSettings().Period
			if period == 0 {
				continue
			}
			var waited time.Duration
			if lastTick.UnixNano() > 0 {
				waited = l.timeNow().Sub(lastTick)
			}
			leftToWait := period - waited
			timer.Reset(leftToWait)
			timerIsStopped = false
		}
	}
}
