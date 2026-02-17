package tcp

import "golang.org/x/sys/unix"

func setMTUDiscovery(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_MTU_DISCOVER, unix.IP_PMTUDISC_PROBE)
}
