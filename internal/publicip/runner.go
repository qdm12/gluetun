package publicip

import (
	"context"
	"net"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
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

		ipCh := make(chan net.IP)
		errorCh := make(chan error)
		go func() {
			ip, err := l.fetcher.FetchPublicIP(getCtx)
			if err != nil {
				if getCtx.Err() == nil {
					errorCh <- err
				}
				return
			}
			ipCh <- ip
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
			case ip := <-ipCh:
				getCancel()

				message := "Public IP address is " + ip.String()
				result, err := Info(ctx, l.client, ip)
				if err != nil {
					l.logger.Warn(err.Error())
				} else {
					message += " (" + result.Country + ", " + result.Region + ", " + result.City + ")"
				}
				l.logger.Info(message)

				result.SetIP(ip)
				l.state.SetData(result)

				filepath := *l.state.GetSettings().IPFilepath
				err = persistPublicIP(filepath, ip.String(), l.puid, l.pgid)
				if err != nil {
					l.logger.Error(err.Error())
				}
				l.statusManager.SetStatus(constants.Completed)
			case err := <-errorCh:
				getCancel()
				close(ipCh)
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
		close(errorCh)
	}
}
