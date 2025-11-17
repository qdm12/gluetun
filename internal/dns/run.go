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
		// Upper scope variables for the DNS over TLS server only
		// Their values are to be used if DOT=off
		var runError <-chan error

		settings := l.GetSettings()
		for !*settings.KeepNameserver && *settings.DoT.Enabled {
			var err error
			runError, err = l.setupServer(ctx)
			if err == nil {
				l.backoffTime = defaultBackoffTime
				l.logger.Info("ready")
				break
			}

			l.signalOrSetStatus(constants.Crashed)

			if ctx.Err() != nil {
				return
			}

			if !errors.Is(err, errUpdateBlockLists) {
				const fallback = true
				l.useUnencryptedDNS(fallback)
			}
			l.logAndWait(ctx, err)
			settings = l.GetSettings()
		}
		l.signalOrSetStatus(constants.Running)

		settings = l.GetSettings()
		if !*settings.KeepNameserver && !*settings.DoT.Enabled {
			const fallback = false
			l.useUnencryptedDNS(fallback)
		}

		l.userTrigger = false

		exitLoop := l.runWait(ctx, runError)
		if exitLoop {
			return
		}
	}
}

func (l *Loop) runWait(ctx context.Context, runError <-chan error) (exitLoop bool) {
	for {
		select {
		case <-ctx.Done():
			if !*l.GetSettings().KeepNameserver {
				l.stopServer()
				// TODO revert OS and Go nameserver when exiting
			}
			return true
		case <-l.stop:
			l.userTrigger = true
			l.logger.Info("stopping")
			if !*l.GetSettings().KeepNameserver {
				const fallback = false
				l.useUnencryptedDNS(fallback)
				l.stopServer()
			}
			l.stopped <- struct{}{}
		case <-l.start:
			l.userTrigger = true
			l.logger.Info("starting")
			return false
		case err := <-runError: // unexpected error
			l.statusManager.SetStatus(constants.Crashed)
			const fallback = true
			l.useUnencryptedDNS(fallback)
			l.logAndWait(ctx, err)
			return false
		}
	}
}

func (l *Loop) stopServer() {
	stopErr := l.server.Stop()
	if stopErr != nil {
		l.logger.Error("stopping DoT server: " + stopErr.Error())
	}
}
