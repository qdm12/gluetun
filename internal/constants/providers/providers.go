package providers

const (
	// Custom is the VPN provider name for custom
	// VPN configurations.
	Custom                = "custom"
	Cyberghost            = "cyberghost"
	Example               = "example"
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
	SlickVPN              = "slickvpn"
	Surfshark             = "surfshark"
	Torguard              = "torguard"
	VPNUnlimited          = "vpn unlimited"
	Vyprvpn               = "vyprvpn"
	Wevpn                 = "wevpn"
	Windscribe            = "windscribe"
)

// All returns all the providers except the custom provider.
func All() []string {
	return []string{
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
		SlickVPN,
		Surfshark,
		Torguard,
		VPNUnlimited,
		Vyprvpn,
		Wevpn,
		Windscribe,
	}
}

func AllWithCustom() []string {
	allProviders := All()
	allProvidersWithCustom := make([]string, len(allProviders)+1)
	copy(allProvidersWithCustom, allProviders)
	allProvidersWithCustom[len(allProvidersWithCustom)-1] = Custom
	return allProvidersWithCustom
}
