package tun

import "golang.org/x/sys/unix"

var _ Interface = (*Tun)(nil)

type Interface interface {
	Checker
	Creator
}

type Tun struct {
	mknod func(path string, mode uint32, dev int) (err error)
}

func New() *Tun {
	return &Tun{
		mknod: unix.Mknod,
	}
}
