package healthcheck

import "github.com/go-ping/ping"

//go:generate mockgen -destination=pinger_mock_test.go -package healthcheck . Pinger

type Pinger interface {
	Run() error
	Stop()
}

func newPinger(addrToPing string) (pinger *ping.Pinger) {
	const count = 1
	pinger = ping.New(addrToPing)
	pinger.Count = count
	pinger.SetPrivileged(true) // see https://github.com/go-ping/ping#linux
	return pinger
}
