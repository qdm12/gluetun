package updater

type Options struct {
	Cyberghost bool
	Mullvad    bool
	Nordvpn    bool
	PIA        bool
	PIAold     bool
	Purevpn    bool
	Surfshark  bool
	Vyprvpn    bool
	Windscribe bool
	File       bool // update JSON file (user side)
	Stdout     bool // update constants file (maintainer side)
	DNSAddress string
}
