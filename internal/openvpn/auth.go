package openvpn

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// WriteAuthFile writes the OpenVPN auth file to disk with the right permissions.
func (c *Configurator) WriteAuthFile(user, password string) error {
	content := strings.Join([]string{user, password}, "\n")
	return writeIfDifferent(c.authFilePath, content, c.puid, c.pgid)
}

// WriteAskPassFile writes the OpenVPN askpass file to disk with the right permissions.
func (c *Configurator) WriteAskPassFile(passphrase string) error {
	return writeIfDifferent(c.askPassPath, passphrase, c.puid, c.pgid)
}

func writeIfDifferent(path, content string, puid, pgid int) (err error) {
	fileStat, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("obtaining file information: %w", err)
	}

	const perm = os.FileMode(0400)
	var writeData, setChown bool
	if os.IsNotExist(err) {
		writeData = true
		setChown = true
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
		writeData = string(data) != content
		setChown = fileStat.Mode().Perm() != perm
	}

	if writeData {
		err = ioutil.WriteFile(path, []byte(content), perm)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}

	if setChown {
		err = os.Chown(path, puid, pgid)
		if err != nil {
			return fmt.Errorf("setting file permissions: %w", err)
		}
	}

	return nil
}
