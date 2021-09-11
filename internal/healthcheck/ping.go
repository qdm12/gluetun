package healthcheck

import "github.com/go-ping/ping"

//go:generate mockgen -destination=pinger_mock_test.go -package healthcheck . Pinger

type Pinger interface {
	Run() error
	Stop()
}

func newPinger() (pinger *ping.Pinger) {
	const addrToPing = "1.1.1.1"
	const count = 1
	pinger = ping.New(addrToPing)
	pinger.Count = count
	pinger.SetPrivileged(true)
	return pinger
}
