package models

type AllServers struct {
	Version    uint16            `json:"version"`
	Cyberghost CyberghostServers `json:"cyberghost"`
	HideMyAss  HideMyAssServers  `json:"hidemyass"`
	Mullvad    MullvadServers    `json:"mullvad"`
	Nordvpn    NordvpnServers    `json:"nordvpn"`
	Privado    PrivadoServers    `json:"privado"`
	Pia        PiaServers        `json:"pia"`
	Purevpn    PurevpnServers    `json:"purevpn"`
	Surfshark  SurfsharkServers  `json:"surfshark"`
	Torguard   TorguardServers   `json:"torguard"`
	Vyprvpn    VyprvpnServers    `json:"vyprvpn"`
	Windscribe WindscribeServers `json:"windscribe"`
}

func (a *AllServers) Count() int {
	return len(a.Cyberghost.Servers) +
		len(a.HideMyAss.Servers) +
		len(a.Mullvad.Servers) +
		len(a.Nordvpn.Servers) +
		len(a.Privado.Servers) +
		len(a.Pia.Servers) +
		len(a.Purevpn.Servers) +
		len(a.Surfshark.Servers) +
		len(a.Torguard.Servers) +
		len(a.Vyprvpn.Servers) +
		len(a.Windscribe.Servers)
}

type CyberghostServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []CyberghostServer `json:"servers"`
}
type HideMyAssServers struct {
	Version   uint16            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Servers   []HideMyAssServer `json:"servers"`
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
type VyprvpnServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []VyprvpnServer `json:"servers"`
}
type WindscribeServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []WindscribeServer `json:"servers"`
}
