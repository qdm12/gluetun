// Package unix defines interfaces to interact with Unix related objects.
// Its primary use is to be used in tests.
package unix

import sysunix "golang.org/x/sys/unix"

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Unix

type Unix interface {
	Mkdev(major uint32, minor uint32) uint64
	Mknod(path string, mode uint32, dev int) (err error)
}

func New() Unix {
	return &unix{}
}

type unix struct{}

func (u *unix) Mkdev(major uint32, minor uint32) uint64 {
	return sysunix.Mkdev(major, minor)
}

func (u *unix) Mknod(path string, mode uint32, dev int) (err error) {
	return sysunix.Mknod(path, mode, dev)
}
