package pinger

import (
	"github.com/mehrdadrad/ping"
)

func Ping() (bool, error) {
	p, err := ping.New("ipv6.test-ipv6.com")
	p.SetForceV6()
	if err != nil {
		return false, err
	}
	p.SetCount(4)
	_, err = p.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
