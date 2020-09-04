package openvpn

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
	"golang.org/x/sys/unix"
)

// CheckTUN checks the tunnel device is present and accessible
func (c *configurator) CheckTUN() error {
	c.logger.Info("checking for device %s", constants.TunnelDevice)
	f, err := c.openFile(string(constants.TunnelDevice), os.O_RDWR, 0)
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
	if err := c.fileManager.CreateDir("/dev/net"); err != nil {
		return err
	}
	dev := c.mkDev(10, 200)
	if err := c.mkNod(string(constants.TunnelDevice), unix.S_IFCHR, int(dev)); err != nil {
		return err
	}
	if err := c.fileManager.SetUserPermissions(string(constants.TunnelDevice), 0666); err != nil {
		return err
	}
	return nil
}
