package service

import (
	"fmt"
	"strings"
)

func portsToString(ports []uint16) (s string) {
	switch len(ports) {
	case 0:
		return "no port forwarded"
	case 1:
		return "port forwarded is " + fmt.Sprint(int(ports[0]))
	default:
		portStrings := make([]string, len(ports))
		for i, port := range ports {
			portStrings[i] = fmt.Sprint(int(port))
		}
		return "ports forwarded are " + strings.Join(portStrings[:len(portStrings)-1], ", ") +
			" and " + portStrings[len(portStrings)-1]
	}
}
