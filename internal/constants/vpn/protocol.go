package vpn

const (
	OpenVPN   = "openvpn"
	Wireguard = "wireguard"
	Both      = "openvpn+wireguard"
)

func IsWireguard(s string) bool {
	return s == Wireguard || s == Both
}

func IsOpenVPN(s string) bool {
	return s == OpenVPN || s == Both
}
