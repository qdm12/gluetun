//go:build !linux && !windows && !darwin

package tcp

func setIPv6HeaderIncluded(fd int) error {
	panic("not implemented")
}
