package dns

import (
	"fmt"
	"time"
)

func (c *configurator) WaitForUnbound() (err error) {
	const hostToResolve = "github.com"
	waitDurations := [...]time.Duration{
		300 * time.Millisecond,
		100 * time.Millisecond,
		300 * time.Millisecond,
		500 * time.Millisecond,
		time.Second,
		2 * time.Second,
	}
	maxTries := len(waitDurations)
	for i, waitDuration := range waitDurations {
		time.Sleep(waitDuration)
		_, err := c.lookupIP(hostToResolve)
		if err == nil {
			return nil
		}
		c.logger.Warn("could not resolve %s (try %d of %d): %s", hostToResolve, i+1, maxTries, err)
	}
	return fmt.Errorf("Unbound does not seem to be working after %d tries", maxTries)
}
