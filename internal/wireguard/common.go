package wireguard

import (
	"net"
)

type tunDevice interface {
	Close() error
	Name() (string, error)
}

type bind interface {
	Close() error
}

type userspaceDevice interface {
	Close()
	Wait() chan struct{}
	IpcHandle(net.Conn)
	IpcSet(string) error
}

type userSpaceBackend struct {
	createTun    func(string, int) (tunDevice, error)
	createBind   func() bind
	createDevice func(tunDevice, bind, Logger) userspaceDevice
	preStart     func(userspaceDevice, Settings) error
}
