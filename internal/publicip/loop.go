package publicip

import (
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip/state"
)

type Loop struct {
	statusManager *loopstate.State
	state         *state.State
	// Objects
	fetcher Fetcher
	logger  Logger
	// Fixed settings
	puid int
	pgid int
	// Internal channels and locks
	start        chan struct{}
	running      chan models.LoopStatus
	stop         chan struct{}
	stopped      chan struct{}
	updateTicker chan struct{}
	backoffTime  time.Duration
	userTrigger  bool
	// Mock functions
	timeNow func() time.Time
}

const defaultBackoffTime = 5 * time.Second

func NewLoop(fetcher Fetcher, logger Logger,
	settings settings.PublicIP, puid, pgid int) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})
	updateTicker := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, settings, updateTicker)

	return &Loop{
		statusManager: statusManager,
		state:         state,
		// Objects
		fetcher:      fetcher,
		logger:       logger,
		puid:         puid,
		pgid:         pgid,
		start:        start,
		running:      running,
		stop:         stop,
		stopped:      stopped,
		updateTicker: updateTicker,
		userTrigger:  true,
		backoffTime:  defaultBackoffTime,
		timeNow:      time.Now,
	}
}
