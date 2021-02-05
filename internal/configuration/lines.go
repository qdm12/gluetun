package configuration

import (
	"net"
	"strconv"
)

func uint16sToStrings(uint16s []uint16) (strings []string) {
	strings = make([]string, len(uint16s))
	for i := range uint16s {
		strings[i] = strconv.Itoa(int(uint16s[i]))
	}
	return strings
}

func ipNetsToStrings(ipNets []net.IPNet) (strings []string) {
	strings = make([]string, len(ipNets))
	for i := range ipNets {
		strings[i] = ipNets[i].String()
	}
	return strings
}
