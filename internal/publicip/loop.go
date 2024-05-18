package publicip

import (
	"context"
	"fmt"
	"net/netip"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type Loop struct {
	// State
	settings      settings.PublicIP
	settingsMutex sync.RWMutex
	ipData        models.PublicIP
	ipDataMutex   sync.RWMutex
	// Fixed injected objets
	fetcher Fetcher
	logger  Logger
	// Fixed parameters
	puid int
	pgid int
	// Internal channels and locks
	// runCtx is used to detect when the loop has exited
	// when performing an update
	runCtx        context.Context //nolint:containedctx
	runCancel     context.CancelFunc
	runTrigger    chan<- context.Context
	runResult     <-chan error
	updateTrigger chan<- settings.PublicIP
	updatedResult <-chan error
	runDone       <-chan struct{}
	// Mock functions
	timeNow func() time.Time
}

func NewLoop(fetcher Fetcher, logger Logger,
	settings settings.PublicIP, puid, pgid int) *Loop {
	return &Loop{
		settings: settings,
		fetcher:  fetcher,
		logger:   logger,
		puid:     puid,
		pgid:     pgid,
		timeNow:  time.Now,
	}
}

func (l *Loop) String() string {
	return "public ip loop"
}

func (l *Loop) Start(_ context.Context) (_ <-chan error, err error) {
	l.runCtx, l.runCancel = context.WithCancel(context.Background())
	runDone := make(chan struct{})
	l.runDone = runDone
	runTrigger := make(chan context.Context)
	l.runTrigger = runTrigger
	runResult := make(chan error)
	l.runResult = runResult
	updateTrigger := make(chan settings.PublicIP)
	l.updateTrigger = updateTrigger
	updatedResult := make(chan error)
	l.updatedResult = updatedResult

	go l.run(l.runCtx, runDone, runTrigger, runResult, updateTrigger, updatedResult)

	return nil, nil //nolint:nilnil
}

func (l *Loop) run(runCtx context.Context, runDone chan<- struct{},
	runTrigger <-chan context.Context, runResult chan<- error,
	updateTrigger <-chan settings.PublicIP, updatedResult chan<- error) {
	defer close(runDone)

	timer := time.NewTimer(time.Hour)
	defer timer.Stop()
	_ = timer.Stop()
	timerIsReadyToReset := true
	lastFetch := time.Unix(0, 0)

	for {
		singleRunCtx := runCtx
		var singleRunResult chan<- error
		select {
		case <-runCtx.Done():
			return
		case singleRunCtx = <-runTrigger:
			// Note singleRunCtx is canceled if runCtx is canceled.
			singleRunResult = runResult
		case <-timer.C:
			timerIsReadyToReset = true
		case partialUpdate := <-updateTrigger:
			var err error
			timerIsReadyToReset, err = l.update(partialUpdate, lastFetch, timer, timerIsReadyToReset)
			updatedResult <- err
			continue
		}

		lastFetch = l.timeNow()
		timerIsReadyToReset = l.updateTimer(*l.settings.Period, lastFetch, timer, timerIsReadyToReset)

		result, err := l.fetcher.FetchInfo(singleRunCtx, netip.Addr{})
		if err != nil {
			err = fmt.Errorf("fetching information: %w", err)
			if singleRunResult != nil {
				singleRunResult <- err
			} else {
				l.logger.Error(err.Error())
			}
			continue
		}

		message := "Public IP address is " + result.IP.String()
		message += " (" + result.Country + ", " + result.Region + ", " + result.City + ")"
		l.logger.Info(message)

		l.ipDataMutex.Lock()
		l.ipData = result
		l.ipDataMutex.Unlock()

		filepath := *l.settings.IPFilepath
		err = persistPublicIP(filepath, result.IP.String(), l.puid, l.pgid)
		if err != nil {
			err = fmt.Errorf("persisting public ip address: %w", err)
		}

		if singleRunResult != nil {
			singleRunResult <- err
		} else if err != nil {
			l.logger.Error(err.Error())
		}
	}
}

func (l *Loop) RunOnce(ctx context.Context) (err error) {
	singleRunCtx, singleRunCancel := context.WithCancel(ctx)
	select {
	case l.runTrigger <- singleRunCtx:
	case <-ctx.Done(): // in case writing to run trigger is blocking
		singleRunCancel()
		return ctx.Err()
	}

	select {
	case err = <-l.runResult:
		singleRunCancel()
		return err
	case <-l.runCtx.Done():
		singleRunCancel()
		<-l.runResult
		return l.runCtx.Err()
	}
}

func (l *Loop) UpdateWith(partialUpdate settings.PublicIP) (err error) {
	select {
	case l.updateTrigger <- partialUpdate:
		select {
		case err = <-l.updatedResult:
			return err
		case <-l.runCtx.Done():
			return l.runCtx.Err()
		}
	case <-l.runCtx.Done():
		// loop has been stopped, no update can be done
		return l.runCtx.Err()
	}
}

func (l *Loop) Stop() (err error) {
	l.runCancel()
	<-l.runDone
	return l.ClearData()
}
