package dns

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) Start(ctx context.Context, verbosityDetailsLevel uint8) (stdout io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("starting unbound")
	args := []string{"-d", "-c", string(constants.UnboundConf)}
	if verbosityDetailsLevel > 0 {
		args = append(args, "-"+strings.Repeat("v", int(verbosityDetailsLevel)))
	}
	// Only logs to stderr
	_, stdout, waitFn, err = c.commander.Start(ctx, "unbound", args...)
	return stdout, waitFn, err
}

func (c *configurator) Version(ctx context.Context) (version string, err error) {
	output, err := c.commander.Run(ctx, "unbound", "-V")
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
