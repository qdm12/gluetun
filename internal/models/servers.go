package models

type AllServers struct {
	Version      uint16              `json:"version"` // used for migration of the top level scheme
	Cyberghost   CyberghostServers   `json:"cyberghost"`
	Fastestvpn   FastestvpnServers   `json:"fastestvpn"`
	HideMyAss    HideMyAssServers    `json:"hidemyass"`
	Ipvanish     IpvanishServers     `json:"ipvanish"`
	Ivpn         IvpnServers         `json:"ivpn"`
	Mullvad      MullvadServers      `json:"mullvad"`
	Nordvpn      NordvpnServers      `json:"nordvpn"`
	Privado      PrivadoServers      `json:"privado"`
	Pia          PiaServers          `json:"pia"`
	Privatevpn   PrivatevpnServers   `json:"privatevpn"`
	Protonvpn    ProtonvpnServers    `json:"protonvpn"`
	Purevpn      PurevpnServers      `json:"purevpn"`
	Surfshark    SurfsharkServers    `json:"surfshark"`
	Torguard     TorguardServers     `json:"torguard"`
	VPNUnlimited VPNUnlimitedServers `json:"vpnunlimited"`
	Vyprvpn      VyprvpnServers      `json:"vyprvpn"`
	Wevpn        WevpnServers        `json:"wevpn"`
	Windscribe   WindscribeServers   `json:"windscribe"`
}

func (a *AllServers) Count() int {
	return len(a.Cyberghost.Servers) +
		len(a.Fastestvpn.Servers) +
		len(a.HideMyAss.Servers) +
		len(a.Ipvanish.Servers) +
		len(a.Ivpn.Servers) +
		len(a.Mullvad.Servers) +
		len(a.Nordvpn.Servers) +
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

type CyberghostServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []CyberghostServer `json:"servers"`
}
type FastestvpnServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []FastestvpnServer `json:"servers"`
}
type HideMyAssServers struct {
	Version   uint16            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Servers   []HideMyAssServer `json:"servers"`
}
type IpvanishServers struct {
	Version   uint16           `json:"version"`
	Timestamp int64            `json:"timestamp"`
	Servers   []IpvanishServer `json:"servers"`
}
type IvpnServers struct {
	Version   uint16       `json:"version"`
	Timestamp int64        `json:"timestamp"`
	Servers   []IvpnServer `json:"servers"`
}
type MullvadServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []MullvadServer `json:"servers"`
}
type NordvpnServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []NordvpnServer `json:"servers"`
}
type PrivadoServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []PrivadoServer `json:"servers"`
}
type PiaServers struct {
	Version   uint16      `json:"version"`
	Timestamp int64       `json:"timestamp"`
	Servers   []PIAServer `json:"servers"`
}
type PrivatevpnServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []PrivatevpnServer `json:"servers"`
}
type ProtonvpnServers struct {
	Version   uint16            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Servers   []ProtonvpnServer `json:"servers"`
}
type PurevpnServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []PurevpnServer `json:"servers"`
}
type SurfsharkServers struct {
	Version   uint16            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Servers   []SurfsharkServer `json:"servers"`
}
type TorguardServers struct {
	Version   uint16           `json:"version"`
	Timestamp int64            `json:"timestamp"`
	Servers   []TorguardServer `json:"servers"`
}
type VPNUnlimitedServers struct {
	Version   uint16               `json:"version"`
	Timestamp int64                `json:"timestamp"`
	Servers   []VPNUnlimitedServer `json:"servers"`
}
type VyprvpnServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []VyprvpnServer `json:"servers"`
}
type WevpnServers struct {
	Version   uint16        `json:"version"`
	Timestamp int64         `json:"timestamp"`
	Servers   []WevpnServer `json:"servers"`
}
type WindscribeServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []WindscribeServer `json:"servers"`
}
