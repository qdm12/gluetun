package publicip

import (
	"context"
	"errors"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		getCtx, getCancel := context.WithCancel(ctx)
		defer getCancel()

		resultCh := make(chan models.PublicIP)
		errorCh := make(chan error)
		go func() {
			result, err := l.fetcher.FetchInfo(getCtx, nil)
			if err != nil {
				if getCtx.Err() == nil {
					errorCh <- err
				}
				return
			}
			resultCh <- result.ToPublicIPModel()
		}()

		if l.userTrigger {
			l.userTrigger = false
			l.running <- constants.Running
		} else { // crash
			l.backoffTime = defaultBackoffTime
			l.statusManager.SetStatus(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				getCancel()
				close(errorCh)
				filepath := *l.state.GetSettings().IPFilepath
				l.logger.Info("Removing ip file " + filepath)
				if err := os.Remove(filepath); err != nil {
					l.logger.Error(err.Error())
				}
				return
			case <-l.start:
				l.userTrigger = true
				getCancel()
				stayHere = false
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				getCancel()
				<-errorCh
				l.stopped <- struct{}{}
			case result := <-resultCh:
				getCancel()

				message := "Public IP address is " + result.IP.String()
				message += " (" + result.Country + ", " + result.Region + ", " + result.City + ")"
				l.logger.Info(message)

				l.state.SetData(result)

				filepath := *l.state.GetSettings().IPFilepath
				err := persistPublicIP(filepath, result.IP.String(), l.puid, l.pgid)
				if err != nil {
					l.logger.Error(err.Error())
				}
				l.statusManager.SetStatus(constants.Completed)
			case err := <-errorCh:
				if errors.Is(err, ipinfo.ErrTooManyRequests) {
					l.logger.Warn(err.Error())
					l.statusManager.SetStatus(constants.Crashed)
					break
				}
				getCancel()
				close(resultCh)
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
		close(errorCh)
	}
}
