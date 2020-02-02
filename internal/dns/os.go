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
	s := strings.TrimSuffix(string(data), "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	found := false
	for i := range lines {
		if strings.HasPrefix(lines[i], "nameserver ") {
			lines[i] = "nameserver 127.0.0.1"
			found = true
		}
	}
	if !found {
		lines = append(lines, "nameserver 127.0.0.1")
	}
	data = []byte(strings.Join(lines, "\n"))
	return c.fileManager.WriteToFile(string(constants.ResolvConf), data)
}
