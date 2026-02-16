//go:build !linux && !windows && !darwin

package ip

func SetIPv6HeaderIncluded(fd int) error {
	panic("not implemented")
}
