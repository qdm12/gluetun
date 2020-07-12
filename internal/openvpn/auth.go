package openvpn

import (
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions
func (c *configurator) WriteAuthFile(user, password string, uid, gid int) error {
	exists, err := c.fileManager.FileExists(string(constants.OpenVPNAuthConf))
	if err != nil {
		return err
	} else if exists {
		data, err := c.fileManager.ReadFile(string(constants.OpenVPNAuthConf))
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		if len(lines) > 1 && lines[0] == user && lines[1] == password {
			return nil
		}
		c.logger.Info("username and password changed", constants.OpenVPNAuthConf)
	}
	return c.fileManager.WriteLinesToFile(
		string(constants.OpenVPNAuthConf),
		[]string{user, password},
		files.Ownership(uid, gid),
		files.Permissions(0400))
}
