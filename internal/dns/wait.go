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
		c.logger.Warn("could not resolve %s (try %d of %d)", hostToResolve, try, maxTries)
		time.Sleep(time.Duration(maxTries * 50 * time.Millisecond))
	}
	return fmt.Errorf("Unbound does not seem to be working after %d tries", maxTries)
}
