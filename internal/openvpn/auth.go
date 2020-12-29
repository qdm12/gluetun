package openvpn

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions.
func (c *configurator) WriteAuthFile(user, password string, uid, gid int) error {
	const filepath = string(constants.OpenVPNAuthConf)
	file, err := c.os.OpenFile(filepath, os.O_RDONLY, 0)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		file, err = c.os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0400)
		if err != nil {
			return err
		}
		_, err = file.WriteString(user + "\n" + password)
		if err != nil {
			_ = file.Close()
			return err
		}
		err = file.Chown(uid, gid)
		if err != nil {
			_ = file.Close()
			return err
		}
		return file.Close()
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 1 && lines[0] == user && lines[1] == password {
		return nil
	}

	c.logger.Info("username and password changed in %s", constants.OpenVPNAuthConf)
	file, err = c.os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, 0400)
	if err != nil {
		return err
	}
	_, err = file.WriteString(user + "\n" + password)
	if err != nil {
		_ = file.Close()
		return err
	}
	err = file.Chown(uid, gid)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
