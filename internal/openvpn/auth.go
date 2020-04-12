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
		c.logger.Info("%s already exists", constants.OpenVPNAuthConf)
		return nil
	}
	c.logger.Info("writing auth file %s", constants.OpenVPNAuthConf)
	return c.fileManager.WriteLinesToFile(
		string(constants.OpenVPNAuthConf),
		[]string{user, password},
		files.Ownership(uid, gid),
		files.Permissions(0400))
}
