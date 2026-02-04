package pmtud

import (
	"syscall"
)

func setDontFragment(fd uintptr) (err error) {
	return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP,
		syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_PROBE)
}
