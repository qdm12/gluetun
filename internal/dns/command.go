package dns

import (
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) Start() (stdout io.ReadCloser, err error) {
	stdout, _, err = c.commander.Start("unbound", "-d", "-c", string(constants.UnboundConf))
	return stdout, err
}

func (c *configurator) Version() (version string, err error) {
	output, err := c.commander.Run("unbound", "-V")
	if err != nil {
		return "", err
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
		return "", fmt.Errorf("unbound -h: version was not found in %q", output)
	}
	return version, nil
}
