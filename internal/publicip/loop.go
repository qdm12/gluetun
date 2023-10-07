package publicip

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
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
	runTrigger    chan<- struct{}
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
	runTrigger := make(chan struct{})
	l.runTrigger = runTrigger
	updateTrigger := make(chan settings.PublicIP)
	l.updateTrigger = updateTrigger
	updatedResult := make(chan error)
	l.updatedResult = updatedResult

	go l.run(l.runCtx, runDone, runTrigger, updateTrigger, updatedResult)

	return nil, nil //nolint:nilnil
}

func (l *Loop) run(runCtx context.Context, runDone chan<- struct{},
	runTrigger <-chan struct{}, updateTrigger <-chan settings.PublicIP,
	updatedResult chan<- error) {
	defer close(runDone)

	timer := time.NewTimer(time.Hour)
	defer timer.Stop()
	_ = timer.Stop()
	timerIsReadyToReset := true
	lastFetch := time.Unix(0, 0)

	for {
		select {
		case <-runCtx.Done():
			return
		case <-runTrigger:
		case <-timer.C:
			timerIsReadyToReset = true
		case partialUpdate := <-updateTrigger:
			var err error
			timerIsReadyToReset, err = l.update(partialUpdate, lastFetch, timer, timerIsReadyToReset)
			updatedResult <- err
			continue
		}

		result, err := l.fetchIPData(runCtx)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		lastFetch = l.timeNow()
		timerIsReadyToReset = l.updateTimer(*l.settings.Period, lastFetch, timer, timerIsReadyToReset)

		if errors.Is(err, ipinfo.ErrTooManyRequests) {
			continue
		}

		message := "Public IP address is " + result.IP.String()
		message += " (" + result.Country + ", " + result.Region + ", " + result.City + ")"
		l.logger.Info(message)

		l.ipDataMutex.Lock()
		l.ipData = result.ToPublicIPModel()
		l.ipDataMutex.Unlock()

		filepath := *l.settings.IPFilepath
		err = persistPublicIP(filepath, result.IP.String(), l.puid, l.pgid)
		if err != nil { // non critical error, which can be fixed with settings updates.
			l.logger.Error(err.Error())
		}
	}
}

func (l *Loop) fetchIPData(ctx context.Context) (result ipinfo.Response, err error) {
	// keep retrying since settings updates won't change the
	// behavior of the following code.
	const defaultBackoffTime = 5 * time.Second
	backoffTime := defaultBackoffTime
	for {
		result, err = l.fetcher.FetchInfo(ctx, netip.Addr{})
		switch {
		case err == nil:
			return result, nil
		case ctx.Err() != nil:
			return result, err
		case errors.Is(err, ipinfo.ErrTooManyRequests):
			l.logger.Warn(err.Error() + "; not retrying.")
			return result, err
		}

		l.logger.Error(fmt.Sprintf("%s - retrying in %s", err, backoffTime))
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(backoffTime):
		}
		const backoffTimeMultipler = 2
		backoffTime *= backoffTimeMultipler
	}
}

func (l *Loop) StartSingleRun() {
	l.runTrigger <- struct{}{}
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
