package loopstate

import (
	"sync"

	"github.com/qdm12/gluetun/internal/models"
)

var _ Manager = (*State)(nil)

type Manager interface {
	Locker
	Getter
	Setter
	Applier
}

func New(status models.LoopStatus,
	start chan<- struct{}, running <-chan models.LoopStatus,
	stop chan<- struct{}, stopped <-chan struct{}) *State {
	return &State{
		status:  status,
		start:   start,
		running: running,
		stop:    stop,
		stopped: stopped,
	}
}

type State struct {
	loopMu sync.RWMutex

	status   models.LoopStatus
	statusMu sync.RWMutex

	start   chan<- struct{}
	running <-chan models.LoopStatus
	stop    chan<- struct{}
	stopped <-chan struct{}
}
