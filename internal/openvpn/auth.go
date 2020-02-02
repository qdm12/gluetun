package openvpn

import (
	"os"
	libuser "os/user"
	"strconv"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions
func (c *configurator) WriteAuthFile(user, password string) error {
	authExists, err := c.fileManager.FileExists(string(constants.OpenVPNAuthConf))
	if err != nil {
		return err
	} else if authExists { // in case of container stop/start
		c.logger.Info("openvpn configurator: %s already exists", constants.OpenVPNAuthConf)
		return nil
	}
	c.logger.Info("openvpn configurator: writing auth file %s", constants.OpenVPNAuthConf)
	c.fileManager.WriteLinesToFile(string(constants.OpenVPNAuthConf), []string{user, password})
	userObject, err := libuser.Lookup("nonrootuser")
	if err != nil {
		return err
	}
	// Operations below are run as root
	uid, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		return err
	}
	if err := os.Chown(string(constants.OpenVPNAuthConf), uid, gid); err != nil {
		return err
	}
	if err := os.Chmod(string(constants.OpenVPNAuthConf), 0400); err != nil {
		return err
	}
	return nil
}
