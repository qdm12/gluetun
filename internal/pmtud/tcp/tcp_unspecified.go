//go:build !linux && !windows

package tcp

func setMark(fd, excludeMark int) error {
	panic("not implemented")
}

func setMTUDiscovery(fd int) error {
	panic("not implemented")
}
