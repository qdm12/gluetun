package portforward

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/service"
)

type Loop struct {
	// State
	settings service.Settings
	service  Service
	// Fixed injected objets
	client      *http.Client
	portAllower PortAllower
	logger      Logger
	// Fixed parameters
	uid, gid int
	// Internal channels and locks
	// runCtx is used to detect when the loop has exited
	// when performing an update
	runCtx    context.Context //nolint:containedctx
	runCancel context.CancelFunc
	updateCh  chan<- service.Settings
	runDone   <-chan struct{}
}

func NewLoop(settings settings.PortForwarding,
	client *http.Client, portAllower PortAllower,
	logger Logger, uid, gid int) *Loop {
	return &Loop{
		settings: service.Settings{
			Settings: settings,
		},
		client:      client,
		portAllower: portAllower,
		logger:      logger,
		uid:         uid,
		gid:         gid,
	}
}

func (l *Loop) Start(_ context.Context) (runError <-chan error, _ error) {
	l.runCtx, l.runCancel = context.WithCancel(context.Background())
	runDone := make(chan struct{})
	l.runDone = runDone

	updateCh := make(chan service.Settings)
	l.updateCh = updateCh
	runErrorCh := make(chan error)

	go l.run(l.runCtx, runDone, runErrorCh, updateCh)

	return runErrorCh, nil
}

func (l *Loop) run(runCtx context.Context, runDone chan<- struct{},
	runErrorCh chan<- error, updateCh <-chan service.Settings) {
	defer close(runDone)

	var serviceRunError <-chan error
	for {
		var partialUpdate service.Settings
		select {
		case <-runCtx.Done():
			// Stop call takes care of stopping the service
			return
		case partialUpdate = <-updateCh: // first and subsequent start trigger
		case err := <-serviceRunError:
			runErrorCh <- err
			return
		}

		firstRun := l.service == nil
		if !firstRun {
			err := l.service.Stop()
			if err != nil {
				runErrorCh <- fmt.Errorf("stopping previous service: %w", err)
				return
			}
		}

		err := l.settings.UpdateWith(partialUpdate)
		if err != nil {
			runErrorCh <- fmt.Errorf("updating settings: %w", err)
			return
		}

		l.service = service.New(l.settings, l.client,
			l.portAllower, l.logger, l.uid, l.gid)

		serviceRunError, err = l.service.Start(runCtx)
		if err != nil {
			if runCtx.Err() == nil { // crashed but NOT stopped
				runErrorCh <- fmt.Errorf("starting new service: %w", err)
			}
			return
		}
	}
}

func (l *Loop) Update(partialUpdate service.Settings) {
	select {
	case l.updateCh <- partialUpdate:
	case <-l.runCtx.Done():
		// loop has been stopped, no update can be done
	}
}

func (l *Loop) Stop() (err error) {
	l.runCancel()
	<-l.runDone

	if l.service != nil {
		return l.service.Stop()
	}
	return nil
}

func (l *Loop) GetSettings() (settings service.Settings) {
	return l.settings
}

func (l *Loop) GetPortForwarded() (port uint16) {
	if l.service == nil {
		return 0
	}
	return l.service.GetPortForwarded()
}
