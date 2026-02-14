//go:build !linux && !windows && !darwin

package ip

func SetIPv4HeaderIncluded(fd int) error {
	panic("not implemented")
}
