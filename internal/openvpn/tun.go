package openvpn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/unix"
)

// CheckTUN checks the tunnel device is present and accessible.
func (c *configurator) CheckTUN() error {
	c.logger.Info("checking for device " + c.tunDevPath)
	f, err := os.OpenFile(c.tunDevPath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("TUN device is not available: %w", err)
	}
	if err := f.Close(); err != nil {
		c.logger.Warn("Could not close TUN device file: %s", err)
	}
	return nil
}

func (c *configurator) CreateTUN() error {
	c.logger.Info("creating " + c.tunDevPath)

	parentDir := filepath.Dir(c.tunDevPath)
	if err := os.MkdirAll(parentDir, 0751); err != nil { //nolint:gomnd
		return err
	}

	const (
		major = 10
		minor = 200
	)
	dev := c.unix.Mkdev(major, minor)
	if err := c.unix.Mknod(c.tunDevPath, unix.S_IFCHR, int(dev)); err != nil {
		return err
	}

	const readWriteAllPerms os.FileMode = 0666
	file, err := os.OpenFile(c.tunDevPath, os.O_WRONLY, readWriteAllPerms)
	if err != nil {
		return err
	}

	return file.Close()
}
