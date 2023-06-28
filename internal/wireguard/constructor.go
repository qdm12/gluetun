package wireguard

type Wireguard struct {
	logger   Logger
	settings Settings
	netlink  NetLinker
	routing  Routing
}

func New(settings Settings, netlink NetLinker,
	routing Routing, logger Logger,
) (w *Wireguard, err error) {
	settings.SetDefaults()
	if err := settings.Check(); err != nil {
		return nil, err
	}

	return &Wireguard{
		logger:   logger,
		settings: settings,
		netlink:  netlink,
		routing:  routing,
	}, nil
}
