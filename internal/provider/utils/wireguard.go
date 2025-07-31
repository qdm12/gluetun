package utils

import (
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/wireguard"
)

func BuildWireguardSettings(connection models.Connection,
	userSettings settings.Wireguard, ipv6Supported bool,
) (settings wireguard.Settings) {
	// If a custom config file is provided, parse and use it
	if userSettings.CustomConfigFile != nil && *userSettings.CustomConfigFile != "" {
		conf, err := wireguard.ParseConfFile(*userSettings.CustomConfigFile)
		if err == nil {
			settings.PrivateKey = conf.Interface["PrivateKey"]
			settings.PublicKey = conf.Peer["PublicKey"]
			if userSettings.Interface != nil && *userSettings.Interface != "" {
				settings.InterfaceName = *userSettings.Interface
			} else {
				settings.InterfaceName = "wg0"
			}
			settings.IPv6 = &ipv6Supported
			const rulePriority = 101
			settings.RulePriority = rulePriority
			// Parse Endpoint
			if endpoint, ok := conf.Peer["Endpoint"]; ok {
				if addrPort, err := netip.ParseAddrPort(endpoint); err == nil {
					settings.Endpoint = addrPort
				}
			}
			// Parse Addresses
			if address, ok := conf.Interface["Address"]; ok {
				parts := strings.Split(address, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if prefix, err := netip.ParsePrefix(part); err == nil {
						settings.Addresses = append(settings.Addresses, prefix)
					}
				}
			}
			// Parse AllowedIPs
			if allowed, ok := conf.Peer["AllowedIPs"]; ok {
				parts := strings.Split(allowed, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if prefix, err := netip.ParsePrefix(part); err == nil {
						settings.AllowedIPs = append(settings.AllowedIPs, prefix)
					}
				}
			}
			return settings
		}
	}

	if userSettings.PrivateKey != nil {
		settings.PrivateKey = *userSettings.PrivateKey
	}
	settings.PublicKey = connection.PubKey
	if userSettings.PreSharedKey != nil {
		settings.PreSharedKey = *userSettings.PreSharedKey
	}
	if userSettings.Interface != nil {
		settings.InterfaceName = *userSettings.Interface
	}
	settings.IPv6 = &ipv6Supported

	const rulePriority = 101 // 100 is to receive external connections
	settings.RulePriority = rulePriority

	settings.Endpoint = netip.AddrPortFrom(connection.IP, connection.Port)

	settings.Addresses = make([]netip.Prefix, 0, len(userSettings.Addresses))
	for _, address := range userSettings.Addresses {
		if !ipv6Supported && address.Addr().Is6() {
			continue
		}
		addressCopy := netip.PrefixFrom(address.Addr(), address.Bits())
		settings.Addresses = append(settings.Addresses, addressCopy)
	}

	settings.AllowedIPs = make([]netip.Prefix, 0, len(userSettings.AllowedIPs))
	for _, allowedIP := range userSettings.AllowedIPs {
		if !ipv6Supported && allowedIP.Addr().Is6() {
			continue
		}
		settings.AllowedIPs = append(settings.AllowedIPs, allowedIP)
	}

	return settings
}
