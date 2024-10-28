//go:build linux || darwin

package tun

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Create creates a TUN device at the path specified.
func (t *Tun) Create(path string) (err error) {
	parentDir := filepath.Dir(path)
	err = os.MkdirAll(parentDir, 0o751) //nolint:mnd
	if err != nil {
		return err
	}

	const (
		major = 10
		minor = 200
	)
	dev := unix.Mkdev(major, minor)
	if dev > math.MaxInt {
		panic("dev is too high")
	}
	err = unix.Mknod(path, unix.S_IFCHR, int(dev))
	if err != nil {
		return fmt.Errorf("creating TUN device file node: %w", err)
	}

	fd, err := unix.Open(path, 0, 0)
	if err != nil {
		if err.Error() == "operation not permitted" {
			err = fmt.Errorf("%w (did you specify --device /dev/net/tun to your container command?)", err)
		}
		return fmt.Errorf("unix opening TUN device file: %w", err)
	}

	const nonBlocking = true
	err = unix.SetNonblock(fd, nonBlocking)
	if err != nil {
		return fmt.Errorf("setting non block to TUN device file descriptor: %w", err)
	}

	return nil
}
