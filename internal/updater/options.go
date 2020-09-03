package updater

type Options struct {
	PIA        bool
	PIAold     bool
	Mullvad    bool
	Vyprvpn    bool
	Surfshark  bool
	Nordvpn    bool
	Purevpn    bool
	File       bool // update JSON file (user side)
	Stdout     bool // update constants file (maintainer side)
	DNSAddress string
}
