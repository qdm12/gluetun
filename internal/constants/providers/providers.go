package providers

const (
	// Custom is the VPN provider name for custom
	// VPN configurations.
	Custom                = "custom"
	Cyberghost            = "cyberghost"
	Expressvpn            = "expressvpn"
	Fastestvpn            = "fastestvpn"
	HideMyAss             = "hidemyass"
	Ipvanish              = "ipvanish"
	Ivpn                  = "ivpn"
	Mullvad               = "mullvad"
	Nordvpn               = "nordvpn"
	Perfectprivacy        = "perfect privacy"
	Privado               = "privado"
	PrivateInternetAccess = "private internet access"
	Privatevpn            = "privatevpn"
	Protonvpn             = "protonvpn"
	Purevpn               = "purevpn"
	Surfshark             = "surfshark"
	Torguard              = "torguard"
	VPNUnlimited          = "vpn unlimited"
	Vyprvpn               = "vyprvpn"
	Wevpn                 = "wevpn"
	Windscribe            = "windscribe"
)

func All() []string {
	return []string{
		Custom,
		Cyberghost,
		Expressvpn,
		Fastestvpn,
		HideMyAss,
		Ipvanish,
		Ivpn,
		Mullvad,
		Nordvpn,
		Perfectprivacy,
		Privado,
		PrivateInternetAccess,
		Privatevpn,
		Protonvpn,
		Purevpn,
		Surfshark,
		Torguard,
		VPNUnlimited,
		Vyprvpn,
		Wevpn,
		Windscribe,
	}
}
