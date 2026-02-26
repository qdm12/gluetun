package dns

import (
	"context"

	"github.com/qdm12/dns/v2/pkg/nameserver"
	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	var err error
	l.localResolvers, err = nameserver.GetPrivateDNSServers()
	if err != nil {
		l.logger.Error("getting private DNS servers: " + err.Error())
		return
	}

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
		// Upper scope variables for the DNS forwarder server only
		// Their values are to be used if DOT=off
		var runError <-chan error

		settings := l.GetSettings()
		for !*settings.KeepNameserver && *settings.ServerEnabled {
			var err error
			runError, err = l.setupServer(ctx)
			if err == nil {
				l.backoffTime = defaultBackoffTime
				l.logger.Info("ready and using DNS server at address " + settings.ServerAddress.String())

				err = l.updateFiles(ctx, settings)
				if err != nil {
					l.logger.Warn("downloading block lists failed, skipping: " + err.Error())
				}
				break
			}

			l.signalOrSetStatus(constants.Crashed)

			if ctx.Err() != nil {
				return
			}
			l.logAndWait(ctx, err)
			settings = l.GetSettings()
		}
		l.signalOrSetStatus(constants.Running)

		settings = l.GetSettings()
		if !*settings.KeepNameserver && !*settings.ServerEnabled {
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
			settings := l.GetSettings()
			if !*settings.KeepNameserver && *settings.ServerEnabled {
				l.stopServer()
				// TODO revert OS and Go nameserver when exiting
			}
			return true
		case <-l.stop:
			l.userTrigger = true
			l.logger.Info("stopping")
			settings := l.GetSettings()
			if !*settings.KeepNameserver && *settings.ServerEnabled {
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
		l.logger.Error("stopping server: " + stopErr.Error())
	}
}
