package utils

import (
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func filterServers(servers []models.Server,
	selection settings.ServerSelection) (filtered []models.Server) {
	for _, server := range servers {
		if filterServer(server, selection) {
			continue
		}

		filtered = append(filtered, server)
	}

	return filtered
}

func filterServer(server models.Server,
	selection settings.ServerSelection) (filtered bool) {
	// Note each condition is split to make sure
	// we have full testing coverage.
	if server.VPN != selection.VPN {
		return true
	}

	if filterByProtocol(selection, server.TCP, server.UDP) {
		return true
	}

	if *selection.MultiHopOnly && !server.MultiHop {
		return true
	}

	if *selection.FreeOnly && !server.Free {
		return true
	}

	if *selection.PremiumOnly && !server.Premium {
		return true
	}

	if *selection.StreamOnly && !server.Stream {
		return true
	}

	if *selection.OwnedOnly && !server.Owned {
		return true
	}

	if filterByPossibilities(server.Country, selection.Countries) {
		return true
	}

	if filterByPossibilities(server.Region, selection.Regions) {
		return true
	}

	if filterByPossibilities(server.City, selection.Cities) {
		return true
	}

	if filterByPossibilities(server.ISP, selection.ISPs) {
		return true
	}

	if filterByPossibilitiesUint16(server.Number, selection.Numbers) {
		return true
	}

	if filterByPossibilities(server.ServerName, selection.Names) {
		return true
	}

	if filterByPossibilities(server.Hostname, selection.Hostnames) {
		return true
	}

	// TODO filter port forward server for PIA

	return false
}

func filterByPossibilities(value string, possibilities []string) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}
	for _, possibility := range possibilities {
		if strings.EqualFold(value, possibility) {
			return false
		}
	}
	return true
}

// TODO merge with filterByPossibilities with generics in Go 1.18.
func filterByPossibilitiesUint16(value uint16, possibilities []uint16) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}
	for _, possibility := range possibilities {
		if value == possibility {
			return false
		}
	}
	return true
}
