package wireguard

import "github.com/qdm12/gluetun/internal/netlink"

var _ Wireguarder = (*Wireguard)(nil)

type Wireguarder interface {
	Runner
	Runner
}

type Wireguard struct {
	logger   Logger
	settings Settings
	netlink  netlink.NetLinker
}

func New(settings Settings, netlink NetLinker,
	logger Logger) (w *Wireguard, err error) {
	settings.SetDefaults()
	if err := settings.Check(); err != nil {
		return nil, err
	}

	return &Wireguard{
		logger:   logger,
		settings: settings,
		netlink:  netlink,
	}, nil
}
