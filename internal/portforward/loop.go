package portforward

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/service"
)

type Loop struct {
	// State
	settings      service.Settings
	settingsMutex sync.RWMutex
	service       Service
	// Fixed injected objets
	client      *http.Client
	portAllower PortAllower
	logger      Logger
	// Fixed parameters
	uid, gid int
	// Internal channels and locks
	// runCtx is used to detect when the loop has exited
	// when performing an update
	runCtx        context.Context //nolint:containedctx
	runCancel     context.CancelFunc
	updatedSignal chan<- struct{}
	runDone       <-chan struct{}
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

	updatedSignal := make(chan struct{})
	l.updatedSignal = updatedSignal
	runErrorCh := make(chan error)

	go l.run(l.runCtx, runDone, runErrorCh, updatedSignal)

	return runErrorCh, nil
}

func (l *Loop) run(runCtx context.Context, runDone chan<- struct{},
	runErrorCh chan<- error, updatedSignal <-chan struct{}) {
	defer close(runDone)

	var serviceRunError <-chan error
	for {
		select {
		case <-runCtx.Done():
			// Stop call takes care of stopping the service
			return
		case <-updatedSignal: // first and subsequent start trigger
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

		l.settingsMutex.RLock()
		l.service = service.New(l.settings, l.client,
			l.portAllower, l.logger, l.uid, l.gid)
		l.settingsMutex.RUnlock()

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

func (l *Loop) UpdateWith(partialUpdate service.Settings) (err error) {
	l.settingsMutex.Lock()
	err = l.settings.UpdateWith(partialUpdate)
	l.settingsMutex.Unlock()
	if err != nil {
		return err
	}

	select {
	case l.updatedSignal <- struct{}{}:
		// Settings are validated and if the service fails to start
		// or crashes at runtime, the loop will stop and signal its
		// parent goroutine. Settings validation should be the only
		// error feedback for the caller of `Update`.
		return nil
	case <-l.runCtx.Done():
		// loop has been stopped, no update can be done
		return l.runCtx.Err()
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
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *Loop) GetPortForwarded() (port uint16) {
	if l.service == nil {
		return 0
	}
	return l.service.GetPortForwarded()
}
