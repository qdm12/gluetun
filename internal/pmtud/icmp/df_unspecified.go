//go:build !linux && !windows

package icmp

// setDontFragment for platforms other than Linux and Windows
// is not implemented, so we just return assuming the don't
// fragment flag is set on IP packets.
func setDontFragment(fd uintptr, ipv4 bool) (err error) {
	return nil
}
