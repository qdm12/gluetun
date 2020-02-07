package openvpn

import (
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions
func (c *configurator) WriteAuthFile(user, password string, uid, gid int) error {
	authExists, err := c.fileManager.FileExists(string(constants.OpenVPNAuthConf))
	if err != nil {
		return err
	} else if authExists { // in case of container stop/start
		c.logger.Info("%s: %s already exists", logPrefix, constants.OpenVPNAuthConf)
		return nil
	}
	c.logger.Info("%s: writing auth file %s", logPrefix, constants.OpenVPNAuthConf)
	return c.fileManager.WriteLinesToFile(
		string(constants.OpenVPNAuthConf),
		[]string{user, password},
		files.FileOwnership(uid, gid),
		files.FilePermissions(0400))
}
