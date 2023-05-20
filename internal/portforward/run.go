package portforward

import (
	"context"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start: // l.state.SetStartData called beforehand
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		pfCtx, pfCancel := context.WithCancel(ctx)

		portCh := make(chan uint16)
		errorCh := make(chan error)

		startData := l.state.GetStartData()

		go func(ctx context.Context, startData StartData) {
			port, err := startData.PortForwarder.PortForward(ctx, l.client, l.logger,
				startData.Gateway, startData.ServerName)
			if err != nil {
				errorCh <- err
				return
			}
			portCh <- port

			// Infinite loop
			err = startData.PortForwarder.KeepPortForward(ctx, port,
				startData.Gateway, startData.ServerName, l.logger)
			errorCh <- err
		}(pfCtx, startData)

		if l.userTrigger {
			l.userTrigger = false
			l.running <- constants.Running
		} else { // crash
			l.backoffTime = defaultBackoffTime
			l.statusManager.SetStatus(constants.Running)
		}

		stayHere := true
		stopped := false
		for stayHere {
			select {
			case <-ctx.Done():
				pfCancel()
				if stopped {
					return
				}
				<-errorCh
				close(errorCh)
				close(portCh)
				l.removePortForwardedFile()
				l.firewallBlockPort(ctx)
				l.state.SetPortForwarded(0)
				return
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				pfCancel()
				stayHere = false
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				pfCancel()
				<-errorCh
				l.removePortForwardedFile()
				l.firewallBlockPort(ctx)
				l.state.SetPortForwarded(0)
				l.stopped <- struct{}{}
				stopped = true
			case port := <-portCh:
				l.logger.Info("port forwarded is " + strconv.Itoa(int(port)))
				l.firewallBlockPort(ctx)
				l.state.SetPortForwarded(port)
				l.firewallAllowPort(ctx)
				l.writePortForwardedFile(port)
			case err := <-errorCh:
				pfCancel()
				close(errorCh)
				close(portCh)
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
		pfCancel() // for linting
	}
}
