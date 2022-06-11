package tun

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

var (
	ErrTUNInfo    = errors.New("cannot get syscall stat info of TUN file")
	ErrTUNBadRdev = errors.New("TUN file has an unexpected rdev")
)

// Check checks the tunnel device specified by path is present and accessible.
func (t *Tun) Check(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("TUN device is not available: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat TUN file: %w", err)
	}

	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return ErrTUNInfo
	}

	const expectedRdev = 2760 // corresponds to major 10 and minor 200
	if sys.Rdev != expectedRdev {
		return fmt.Errorf("%w: %d instead of expected %d",
			ErrTUNBadRdev, sys.Rdev, expectedRdev)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("cannot close TUN device: %w", err)
	}

	return nil
}
