package os

import (
	nativeos "os"
)

type OpenFileFunc func(name string, flag int, perm FileMode) (File, error)
type MkdirAllFunc func(name string, perm nativeos.FileMode) error
type RemoveFunc func(name string) error
type ChownFunc func(name string, uid int, gid int) error
