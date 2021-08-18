package tun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type Creator interface {
	Create(path string) error
}

var (
	ErrMknod = errors.New("cannot create TUN device file node")
)

// Create creates a TUN device at the path specified.
func (t *Tun) Create(path string) error {
	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0751); err != nil { //nolint:gomnd
		return err
	}

	const (
		major = 10
		minor = 200
	)
	dev := unix.Mkdev(major, minor)
	err := t.mknod(path, unix.S_IFCHR, int(dev))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrMknod, err)
	}

	return nil
}
