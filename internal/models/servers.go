package models

type AllServers struct {
	Version    uint16            `json:"version"`
	Cyberghost CyberghostServers `json:"cyberghost"`
	Mullvad    MullvadServers    `json:"mullvad"`
	Nordvpn    NordvpnServers    `json:"nordvpn"`
	PiaOld     PiaServers        `json:"piaOld"`
	Pia        PiaServers        `json:"pia"`
	Purevpn    PurevpnServers    `json:"purevpn"`
	Surfshark  SurfsharkServers  `json:"surfshark"`
	Vyprvpn    VyprvpnServers    `json:"vyprvpn"`
	Windscribe WindscribeServers `json:"windscribe"`
}

type CyberghostServers struct {
	Version   uint16             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Servers   []CyberghostServer `json:"servers"`
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
type PiaServers struct {
	Version   uint16      `json:"version"`
	Timestamp int64       `json:"timestamp"`
	Servers   []PIAServer `json:"servers"`
}
type PurevpnServers struct {
	Version   uint16          `json:"version"`
	Timestamp int64           `json:"timestamp"`
	Servers   []PurevpnServer `json:"purevpn"`
}
type SurfsharkServers struct {
	Version   uint16            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Servers   []SurfsharkServer `json:"servers"`
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
