package dns

import (
	"context"
	"errors"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	if *l.GetSettings().KeepNameserver {
		l.logger.Warn("⚠️⚠️⚠️  keeping the default container nameservers, " +
			"this will likely leak DNS traffic outside the VPN " +
			"and go through your container network DNS outside the VPN tunnel!")
	} else {
		const fallback = false
		l.useUnencryptedDNS(fallback)
	}

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		// Upper scope variables for Unbound only
		// Their values are to be used if DOT=off
		waitError := make(chan error)
		unboundCancel := func() { waitError <- nil }
		closeStreams := func() {}

		settings := l.GetSettings()
		for !*settings.KeepNameserver && *settings.DoT.Enabled {
			var err error
			unboundCancel, waitError, closeStreams, err = l.setupUnbound(ctx)
			if err == nil {
				l.backoffTime = defaultBackoffTime
				l.logger.Info("ready")
				l.signalOrSetStatus(constants.Running)
				break
			}

			l.signalOrSetStatus(constants.Crashed)

			if ctx.Err() != nil {
				return
			}

			if !errors.Is(err, errUpdateFiles) {
				const fallback = true
				l.useUnencryptedDNS(fallback)
			}
			l.logAndWait(ctx, err)
		}

		settings = l.GetSettings()
		if !*settings.KeepNameserver && !*settings.DoT.Enabled {
			const fallback = false
			l.useUnencryptedDNS(fallback)
		}

		l.userTrigger = false

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				unboundCancel()
				<-waitError
				close(waitError)
				closeStreams()
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				const fallback = false
				l.useUnencryptedDNS(fallback)
				unboundCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				closeStreams()
				l.stopped <- struct{}{}
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				closeStreams()

				unboundCancel()
				l.statusManager.SetStatus(constants.Crashed)
				const fallback = true
				l.useUnencryptedDNS(fallback)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
	}
}
