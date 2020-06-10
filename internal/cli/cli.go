package cli

import (
	"fmt"
	"net"
)

func HealthCheck() error {
	ips, err := net.LookupIP("github.com")
	if err != nil {
		return fmt.Errorf("cannot resolve github.com (%s)", err)
	} else if len(ips) == 0 {
		return fmt.Errorf("resolved no IP addresses for github.com")
	}
	return nil
}
