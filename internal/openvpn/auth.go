package openvpn

import (
	"io"
	"os"
	"strings"
)

type AuthWriter interface {
	WriteAuthFile(user, password string, puid, pgid int) error
}

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions.
func (c *configurator) WriteAuthFile(user, password string, puid, pgid int) error {
	file, err := os.Open(c.authFilePath)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		file, err = os.OpenFile(c.authFilePath, os.O_WRONLY|os.O_CREATE, 0400)
		if err != nil {
			return err
		}
		_, err = file.WriteString(user + "\n" + password)
		if err != nil {
			_ = file.Close()
			return err
		}
		err = file.Chown(puid, pgid)
		if err != nil {
			_ = file.Close()
			return err
		}
		return file.Close()
	}

	data, err := io.ReadAll(file)
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

	c.logger.Info("username and password changed in " + c.authFilePath)
	file, err = os.OpenFile(c.authFilePath, os.O_TRUNC|os.O_WRONLY, 0400)
	if err != nil {
		return err
	}
	_, err = file.WriteString(user + "\n" + password)
	if err != nil {
		_ = file.Close()
		return err
	}
	err = file.Chown(puid, pgid)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
