package shutdown

import (
	"context"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Wave interface {
	Add(name string, timeout time.Duration) (
		ctx context.Context, done chan struct{})
	size() int
	shutdown(ctx context.Context, logger logging.Logger) (incomplete int)
}

type wave struct {
	name     string
	routines []routine
}

func NewWave(name string) Wave {
	return &wave{
		name: name,
	}
}

func (w *wave) Add(name string, timeout time.Duration) (ctx context.Context, done chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	done = make(chan struct{})
	routine := routine{
		name:    name,
		cancel:  cancel,
		done:    done,
		timeout: timeout,
	}
	w.routines = append(w.routines, routine)
	return ctx, done
}

func (w *wave) size() int { return len(w.routines) }

func (w *wave) shutdown(ctx context.Context, logger logging.Logger) (incomplete int) {
	completed := make(chan bool)

	for _, r := range w.routines {
		go func(r routine) {
			if err := r.shutdown(ctx); err != nil {
				logger.Warn(w.name + " routines: " + err.Error() + " ⚠️")
				completed <- false
			} else {
				logger.Info(w.name + " routines: " + r.name + " terminated ✔️")
				completed <- err == nil
			}
		}(r)
	}

	for range w.routines {
		c := <-completed
		if !c {
			incomplete++
		}
	}

	return incomplete
}
