package publicip

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) update(partialUpdate settings.PublicIP,
	lastTick time.Time, timer *time.Timer, timerIsReadyToReset bool) (
	newTimerIsReadyToReset bool, err error) {
	newTimerIsReadyToReset = timerIsReadyToReset
	// No need to lock the mutex since it can only be written
	// in the code below in this goroutine.
	updatedSettings, err := l.settings.UpdateWith(partialUpdate)
	if err != nil {
		return newTimerIsReadyToReset, err
	}

	if *l.settings.Period != *updatedSettings.Period {
		newTimerIsReadyToReset = l.updateTimer(*updatedSettings.Period, lastTick,
			timer, timerIsReadyToReset)
	}

	if *l.settings.IPFilepath != *updatedSettings.IPFilepath {
		switch {
		case *l.settings.IPFilepath == "":
			err = persistPublicIP(*updatedSettings.IPFilepath,
				l.ipData.IP.String(), l.puid, l.pgid)
			if err != nil {
				return newTimerIsReadyToReset, fmt.Errorf("persisting ip data: %w", err)
			}
		case *updatedSettings.IPFilepath == "":
			err = os.Remove(*l.settings.IPFilepath)
			if err != nil {
				return newTimerIsReadyToReset, fmt.Errorf("removing ip data file path: %w", err)
			}
		default:
			err = os.Rename(*l.settings.IPFilepath, *updatedSettings.IPFilepath)
			if err != nil {
				return newTimerIsReadyToReset, fmt.Errorf("renaming ip data file path: %w", err)
			}
		}
	}

	l.settingsMutex.Lock()
	l.settings = updatedSettings
	l.settingsMutex.Unlock()

	return newTimerIsReadyToReset, nil
}

func (l *Loop) updateTimer(period time.Duration, lastFetch time.Time,
	timer *time.Timer, timerIsReadyToReset bool) (newTimerIsReadyToReset bool) {
	disableTimer := period == 0
	if disableTimer {
		if !timer.Stop() {
			<-timer.C
		}
		return true
	}

	if !timerIsReadyToReset {
		if !timer.Stop() {
			<-timer.C
		}
	}

	var waited time.Duration
	if lastFetch.UnixNano() > 0 {
		waited = l.timeNow().Sub(lastFetch)
	}
	leftToWait := period - waited
	if leftToWait <= 0 {
		leftToWait = time.Nanosecond
	}

	timer.Reset(leftToWait)
	return false
}
