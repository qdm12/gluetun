package portforward

import (
	"context"
	"fmt"
	"net/http"

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
	runCancel context.CancelFunc
	updateCh  chan<- service.Settings
	runDone   <-chan struct{}
}

func NewLoop(client *http.Client, portAllower PortAllower,
	logger Logger, uid, gid int) *Loop {
	return &Loop{
		client:      client,
		portAllower: portAllower,
		logger:      logger,
		uid:         uid,
		gid:         gid,
	}
}

func (l *Loop) Start(_ context.Context) (runError <-chan error, _ error) {
	runCtx, runCancel := context.WithCancel(context.Background())
	l.runCancel = runCancel
	runDone := make(chan struct{})
	l.runDone = runDone

	updateCh := make(chan service.Settings)
	l.updateCh = updateCh
	runErrorCh := make(chan error)

	go l.run(runCtx, runDone, runErrorCh, updateCh)

	return runErrorCh, nil
}

func (l *Loop) run(runCtx context.Context, runDone chan<- struct{},
	runErrorCh chan<- error, updateCh <-chan service.Settings) {
	defer close(runDone)

	var serviceRunError <-chan error
	for {
		var update service.Settings
		select {
		case <-runCtx.Done():
			// Stop call takes care of stopping the service
			return
		case update = <-updateCh: // first and subsequent start trigger
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

		l.settings = update

		l.service = service.New(update, l.client,
			l.portAllower, l.logger, l.uid, l.gid)

		var err error
		serviceRunError, err = l.service.Start(runCtx)
		if err != nil {
			if runCtx.Err() == nil { // crashed but NOT stopped
				runErrorCh <- fmt.Errorf("starting new service: %w", err)
			}
			return
		}
	}
}

func (l *Loop) Update(settings service.Settings) {
	l.updateCh <- settings
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
