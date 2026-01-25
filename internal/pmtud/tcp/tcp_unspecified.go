//go:build !linux && !windows

package tcp

func setMTUDiscovery(fd int) error {
	panic("not implemented")
}
