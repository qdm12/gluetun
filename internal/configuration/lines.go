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

func boolToEnabled(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func boolToYes(yes bool) string {
	if yes {
		return "yes"
	}
	return "no"
}

func boolToOn(on bool) string {
	if on {
		return "on"
	}
	return "off"
}
