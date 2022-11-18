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

	l.logger.Debug("starting port forwarding run loop")

	for ctx.Err() == nil {
		pfCtx, pfCancel := context.WithCancel(ctx)

		portCh := make(chan uint16)
		errorCh := make(chan error)

		startData := l.state.GetStartData()
		l.logger.Debug("start data obtained")

		go func(ctx context.Context, startData StartData) {
			port, err := startData.PortForwarder.PortForward(ctx, l.client, l.logger,
				startData.Gateway, startData.ServerName)
			if err != nil {
				errorCh <- err
				return
			}
			l.logger.Debug("port forwarded obtained")
			portCh <- port
			l.logger.Debug("port forwarded notified")

			// Infinite loop
			l.logger.Debug("keeping port forwarded...")
			err = startData.PortForwarder.KeepPortForward(ctx, l.client,
				port, startData.Gateway, startData.ServerName)
			l.logger.Debug("done keeping port forwarded...")
			errorCh <- err
			l.logger.Debug("notifying keeping port forward eventual error...")
		}(pfCtx, startData)

		if l.userTrigger {
			l.userTrigger = false
			l.running <- constants.Running
			l.logger.Debug("was user triggered, status changed to running")
		} else { // crash
			l.backoffTime = defaultBackoffTime
			l.statusManager.SetStatus(constants.Running)
			l.logger.Debug("succeeded after crash, status changed to running")
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
				l.logger.Debug("port forward start triggered from start channel")
				l.userTrigger = true
				l.logger.Info("starting")
				pfCancel()
				stayHere = false
			case <-l.stop:
				l.logger.Debug("port forward stop triggered from start channel")
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
				l.logger.Debug("port forward crashed, restarting it")
				pfCancel()
				close(errorCh)
				close(portCh)
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
		pfCancel() // for linting
		l.logger.Debug("going back to top of function")
	}
}
