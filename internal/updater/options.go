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
	Stdout     bool // in order to update constants file (maintainer side)
	Verbose    bool
	DNSAddress string
}

func NewOptions(dnsAddress string) Options {
	return Options{
		Cyberghost: true,
		Mullvad:    true,
		Nordvpn:    true,
		PIA:        true,
		PIAold:     true,
		Purevpn:    true,
		Surfshark:  true,
		Vyprvpn:    true,
		Windscribe: true,
		Stdout:     false,
		Verbose:    false,
		DNSAddress: dnsAddress,
	}
}
