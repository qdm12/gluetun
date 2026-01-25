//go:build !linux && !windows && !darwin

package tcp

func setIPv4HeaderIncluded(fd int) error {
	panic("not implemented")
}
