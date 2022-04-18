package models

type AllServers struct {
	Version        uint16  `json:"version"` // used for migration of the top level scheme
	Cyberghost     Servers `json:"cyberghost"`
	Expressvpn     Servers `json:"expressvpn"`
	Fastestvpn     Servers `json:"fastestvpn"`
	HideMyAss      Servers `json:"hidemyass"`
	Ipvanish       Servers `json:"ipvanish"`
	Ivpn           Servers `json:"ivpn"`
	Mullvad        Servers `json:"mullvad"`
	Perfectprivacy Servers `json:"perfectprivacy"`
	Nordvpn        Servers `json:"nordvpn"`
	Privado        Servers `json:"privado"`
	Pia            Servers `json:"pia"`
	Privatevpn     Servers `json:"privatevpn"`
	Protonvpn      Servers `json:"protonvpn"`
	Purevpn        Servers `json:"purevpn"`
	Surfshark      Servers `json:"surfshark"`
	Torguard       Servers `json:"torguard"`
	VPNUnlimited   Servers `json:"vpnunlimited"`
	Vyprvpn        Servers `json:"vyprvpn"`
	Wevpn          Servers `json:"wevpn"`
	Windscribe     Servers `json:"windscribe"`
}

func (a *AllServers) Count() int {
	return len(a.Cyberghost.Servers) +
		len(a.Expressvpn.Servers) +
		len(a.Fastestvpn.Servers) +
		len(a.HideMyAss.Servers) +
		len(a.Ipvanish.Servers) +
		len(a.Ivpn.Servers) +
		len(a.Mullvad.Servers) +
		len(a.Nordvpn.Servers) +
		len(a.Perfectprivacy.Servers) +
		len(a.Privado.Servers) +
		len(a.Pia.Servers) +
		len(a.Privatevpn.Servers) +
		len(a.Protonvpn.Servers) +
		len(a.Purevpn.Servers) +
		len(a.Surfshark.Servers) +
		len(a.Torguard.Servers) +
		len(a.VPNUnlimited.Servers) +
		len(a.Vyprvpn.Servers) +
		len(a.Wevpn.Servers) +
		len(a.Windscribe.Servers)
}

func (a *AllServers) SetDefaults() {
	a.Cyberghost.SetDefaults()
	a.Expressvpn.SetDefaults()
	a.Fastestvpn.SetDefaults()
	a.HideMyAss.SetDefaults()
	a.Ipvanish.SetDefaults()
	a.Ivpn.SetDefaults()
	a.Mullvad.SetDefaults()
	a.Nordvpn.SetDefaults()
	a.Perfectprivacy.SetDefaults()
	a.Privado.SetDefaults()
	a.Pia.SetDefaults()
	a.Privatevpn.SetDefaults()
	a.Protonvpn.SetDefaults()
	a.Purevpn.SetDefaults()
	a.Surfshark.SetDefaults()
	a.Torguard.SetDefaults()
	a.VPNUnlimited.SetDefaults()
	a.Vyprvpn.SetDefaults()
	a.Wevpn.SetDefaults()
	a.Windscribe.SetDefaults()
}

type Servers struct {
	Version   uint16   `json:"version"`
	Timestamp int64    `json:"timestamp"`
	Servers   []Server `json:"servers"`
}

func (s *Servers) SetDefaults() {
	for i := range s.Servers {
		s.Servers[i].setDefaults()
	}
}
