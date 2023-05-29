//go:build !linux && !darwin

package tun

// Create creates a TUN device at the path specified.
func (t *Tun) Create(path string) error {
	panic("not implemented")
}
