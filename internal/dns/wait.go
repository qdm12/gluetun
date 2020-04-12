package dns

import (
	"fmt"
	"time"
)

func (c *configurator) WaitForUnbound() (err error) {
	const maxTries = 10
	const hostToResolve = "github.com"
	for try := 1; try <= maxTries; try++ {
		_, err := c.lookupIP(hostToResolve)
		if err == nil {
			return nil
		}
		c.logger.Warn("could not resolve %s (try %d of %d): %s", hostToResolve, try, maxTries, err)
		time.Sleep(maxTries * 50 * time.Millisecond)
	}
	return fmt.Errorf("Unbound does not seem to be working after %d tries", maxTries)
}
