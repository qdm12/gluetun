package os

import (
	"io"
	nativeos "os"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . File

type File interface {
	io.ReadWriteCloser
	WriteString(s string) (int, error)
	Chown(uid, gid int) error
	Chmod(mode nativeos.FileMode) error
}
