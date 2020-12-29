package os

import nativeos "os"

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . OS

type OS interface {
	OpenFile(name string, flag int, perm FileMode) (File, error)
	MkdirAll(name string, perm FileMode) error
	Remove(name string) error
	Chown(name string, uid int, gid int) error
	Stat(name string) (nativeos.FileInfo, error)
}

func New() OS {
	return &os{}
}

type os struct{}

func (o *os) OpenFile(name string, flag int, perm FileMode) (File, error) {
	return nativeos.OpenFile(name, flag, nativeos.FileMode(perm))
}
func (o *os) MkdirAll(name string, perm FileMode) error {
	return nativeos.MkdirAll(name, nativeos.FileMode(perm))
}
func (o *os) Remove(name string) error {
	return nativeos.Remove(name)
}
func (o *os) Chown(name string, uid, gid int) error {
	return nativeos.Chown(name, uid, gid)
}
func (o *os) Stat(name string) (nativeos.FileInfo, error) {
	return nativeos.Stat(name)
}
