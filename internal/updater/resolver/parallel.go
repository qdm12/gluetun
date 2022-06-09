package resolver

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type Parallel struct {
	repeatResolver *Repeat
	settings       ParallelSettings
}

func NewParallelResolver(settings ParallelSettings) *Parallel {
	return &Parallel{
		repeatResolver: NewRepeat(settings.Repeat),
		settings:       settings,
	}
}

type ParallelSettings struct {
	Repeat    RepeatSettings
	FailEarly bool
	// Maximum ratio of the hosts failing DNS resolution
	// divided by the total number of hosts requested.
	// This value is between 0 and 1. Note this is only
	// applicable if FailEarly is not set to true.
	MaxFailRatio float64
}

type parallelResult struct {
	host string
	IPs  []net.IP
}

var (
	ErrMinFound     = errors.New("not enough hosts found")
	ErrMaxFailRatio = errors.New("maximum failure ratio reached")
)

func (pr *Parallel) Resolve(ctx context.Context, hosts []string, minToFind int) (
	hostToIPs map[string][]net.IP, warnings []string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan parallelResult)
	defer close(results)
	errors := make(chan error)
	defer close(errors)

	for _, host := range hosts {
		go pr.resolveAsync(ctx, host, results, errors)
	}

	hostToIPs = make(map[string][]net.IP, len(hosts))
	maxFails := int(pr.settings.MaxFailRatio * float64(len(hosts)))

	for range hosts {
		select {
		case newErr := <-errors:
			if pr.settings.FailEarly {
				if err == nil {
					// only set the error to the first error encountered
					// and not the context canceled errors coming after.
					err = newErr
					cancel()
				}
				break
			}

			// do not add warnings coming from the call to cancel()
			if len(warnings) < maxFails {
				warnings = append(warnings, newErr.Error())
			}

			if len(warnings) == maxFails {
				cancel() // cancel only once when we reach maxFails
			}
		case result := <-results:
			hostToIPs[result.host] = result.IPs
		}
	}

	if err != nil { // fail early
		return nil, warnings, err
	}

	if len(hostToIPs) < minToFind {
		return nil, warnings,
			fmt.Errorf("%w: found %d hosts but expected at least %d",
				ErrMinFound, len(hostToIPs), minToFind)
	}

	failureRatio := float64(len(warnings)) / float64(len(hosts))
	if failureRatio > pr.settings.MaxFailRatio {
		return hostToIPs, warnings,
			fmt.Errorf("%w: %.2f failure ratio reached", ErrMaxFailRatio, failureRatio)
	}

	return hostToIPs, warnings, nil
}

func (pr *Parallel) resolveAsync(ctx context.Context, host string,
	results chan<- parallelResult, errors chan<- error) {
	IPs, err := pr.repeatResolver.Resolve(ctx, host)
	if err != nil {
		errors <- err
		return
	}
	results <- parallelResult{
		host: host,
		IPs:  IPs,
	}
}
