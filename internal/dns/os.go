package dns

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) SetLocalNameserver() error {
	data, err := c.fileManager.ReadFile(string(constants.ResolvConf))
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		if strings.HasPrefix(lines[i], "nameserver ") {
			lines[i] = "nameserver 127.0.0.1"
		}
	}
	return c.fileManager.WriteLinesToFile(string(constants.ResolvConf), lines)
}
