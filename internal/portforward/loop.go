package portforward

import (
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/portforward/state"
)

type Loop struct {
	statusManager *loopstate.State
	state         *state.State
	// Fixed parameters
	puid int
	pgid int
	// Objects
	client      *http.Client
	portAllower PortAllower
	logger      Logger
	// Internal channels and locks
	start       chan struct{}
	running     chan models.LoopStatus
	stop        chan struct{}
	stopped     chan struct{}
	startMu     sync.Mutex
	backoffTime time.Duration
	userTrigger bool
}

const defaultBackoffTime = 5 * time.Second

func NewLoop(settings settings.PortForwarding,
	client *http.Client, portAllower PortAllower,
	logger Logger, puid, pgid int) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, settings)

	return &Loop{
		statusManager: statusManager,
		state:         state,
		puid:          puid,
		pgid:          pgid,
		// Objects
		client:      client,
		portAllower: portAllower,
		logger:      logger,
		start:       start,
		running:     running,
		stop:        stop,
		stopped:     stopped,
		userTrigger: true,
		backoffTime: defaultBackoffTime,
	}
}
