package wireguard

type Wireguarder interface {
	Runner
	Starter
}

type Wireguard struct {
	logger   Logger
	settings Settings
}

func New(settings Settings, logger Logger) (w *Wireguard, err error) {
	settings.SetDefaults()
	if err := settings.Check(); err != nil {
		return nil, err
	}

	return &Wireguard{
		logger:   logger,
		settings: settings,
	}, nil
}