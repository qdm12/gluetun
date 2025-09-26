//go:build !linux && !windows

package pmtud

// setDontFragment for platforms other than Linux and Windows
// is not implemented, so we just return assuming the don't
// fragment flag is set on IP packets.
func setDontFragment(fd uintptr) (err error) {
	return nil
}
