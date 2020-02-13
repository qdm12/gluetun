package openvpn

import (
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) Start() (stdout io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("%s: starting openvpn", logPrefix)
	stdout, _, waitFn, err = c.commander.Start("openvpn", "--config", string(constants.OpenVPNConf))
	return stdout, waitFn, err
}

func (c *configurator) Version() (string, error) {
	output, err := c.commander.Run("openvpn", "--version")
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}
	firstLine := strings.Split(output, "\n")[0]
	words := strings.Fields(firstLine)
	if len(words) < 2 {
		return "", fmt.Errorf("openvpn --version: first line is too short: %q", firstLine)
	}
	return words[1], nil
}
