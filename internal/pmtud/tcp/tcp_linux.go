package tcp

import "golang.org/x/sys/unix"

// setMark sets a mark on each packets sent through this socket.
// This is used in conjunction with iptables to block outgoing kernel automated
// RST packets, since the kernel is not aware of us handling the connection manually.
// For example:
// iptables -A OUTPUT -p tcp --tcp-flags RST RST -m mark ! --mark 123 -j DROP
//
//nolint:dupword
func setMark(fd, excludeMark int) error {
	return unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_MARK, excludeMark)
}

func setMTUDiscovery(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_MTU_DISCOVER, unix.IP_PMTUDISC_PROBE)
}
