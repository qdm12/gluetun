package dns

import (
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) Start(verbosityDetailsLevel uint8) (stdout io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("%s: starting unbound", logPrefix)
	args := []string{"-d", "-c", string(constants.UnboundConf)}
	if verbosityDetailsLevel > 0 {
		args = append(args, "-"+strings.Repeat("v", int(verbosityDetailsLevel)))
	}
	// Only logs to stderr
	_, stdout, waitFn, err = c.commander.Start("unbound", args...)
	return stdout, waitFn, err
}

func (c *configurator) Version() (version string, err error) {
	output, err := c.commander.Run("unbound", "-V")
	if err != nil {
		return "", fmt.Errorf("unbound version: %w", err)
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Version ") {
			words := strings.Fields(line)
			if len(words) < 2 {
				continue
			}
			version = words[1]
		}
	}
	if version == "" {
		return "", fmt.Errorf("unbound version was not found in %q", output)
	}
	return version, nil
}
