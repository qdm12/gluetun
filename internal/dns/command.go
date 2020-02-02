package dns

import (
	"fmt"
	"io"
	"strings"
)

func (c *configurator) Start() (stdout io.ReadCloser, err error) {
	stdout, _, err = c.commander.Start("unbound")
	return stdout, err
}

func (c *configurator) Version() (version string, err error) {
	output, err := c.commander.Run("unbound", "-h")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Version ") {
			words := strings.Split(line, " ")
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
