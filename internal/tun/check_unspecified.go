//go:build !linux && !darwin

package tun

func (t *Tun) Check(path string) error {
	panic("not implemented")
}
