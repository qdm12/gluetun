package shutdown

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Order interface {
	Append(waves ...Wave)
	Shutdown(timeout time.Duration, logger logging.Logger) (err error)
}

type order struct {
	waves []Wave
	total int // for logging only
}

func NewOrder() Order {
	return &order{}
}

var ErrIncomplete = errors.New("one or more routines did not terminate gracefully")

func (o *order) Append(waves ...Wave) {
	o.waves = append(o.waves, waves...)
}

func (o *order) Shutdown(timeout time.Duration, logger logging.Logger) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	total := 0
	incomplete := 0

	for _, wave := range o.waves {
		total += wave.size()
		incomplete += wave.shutdown(ctx, logger)
	}

	if incomplete == 0 {
		return nil
	}

	return fmt.Errorf("%w: %d not terminated on %d routines",
		ErrIncomplete, incomplete, total)
}
