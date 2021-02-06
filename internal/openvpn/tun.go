package openvpn

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/unix"
)

// CheckTUN checks the tunnel device is present and accessible.
func (c *configurator) CheckTUN() error {
	c.logger.Info("checking for device %s", constants.TunnelDevice)
	f, err := c.os.OpenFile(constants.TunnelDevice, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("TUN device is not available: %w", err)
	}
	if err := f.Close(); err != nil {
		c.logger.Warn("Could not close TUN device file: %s", err)
	}
	return nil
}

func (c *configurator) CreateTUN() error {
	c.logger.Info("creating %s", constants.TunnelDevice)
	if err := c.os.MkdirAll("/dev/net", 0751); err != nil {
		return err
	}

	const (
		major = 10
		minor = 200
	)
	dev := c.unix.Mkdev(major, minor)
	if err := c.unix.Mknod(constants.TunnelDevice, unix.S_IFCHR, int(dev)); err != nil {
		return err
	}

	file, err := c.os.OpenFile(constants.TunnelDevice, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	const readWriteAllPerms os.FileMode = 0666
	if err := file.Chmod(readWriteAllPerms); err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
