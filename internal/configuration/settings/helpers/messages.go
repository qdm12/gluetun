package helpers

import (
	"fmt"
	"strings"
)

func ChoicesOrString(choices []string) string {
	return strings.Join(
		choices[:len(choices)-1], ", ") +
		" or " + choices[len(choices)-1]
}

func PortChoicesOrString(ports []uint16) (s string) {
	switch len(ports) {
	case 0:
		return "there is no allowed port"
	case 1:
		return "allowed port is " + fmt.Sprint(ports[0])
	}

	s = "allowed ports are "
	portStrings := make([]string, len(ports))
	for i := range ports {
		portStrings[i] = fmt.Sprint(ports[i])
	}
	s += ChoicesOrString(portStrings)
	return s
}
