package tinyproxy

import (
	"fmt"
	"io"
	"strings"
)

func (c *configurator) Start() (stdout io.ReadCloser, err error) {
	stdout, _, err = c.commander.Start("tinyproxy", "-d")
	return stdout, err
}

// Version obtains the version of the installed Tinyproxy server
func (c *configurator) Version() (string, error) {
	output, err := c.commander.Run("tinyproxy", "-v")
	if err != nil {
		return "", err
	}
	words := strings.Fields(output)
	if len(words) < 2 {
		return "", fmt.Errorf("tinyproxy -v: output is too short: %q", output)
	}
	return words[1], nil
}
