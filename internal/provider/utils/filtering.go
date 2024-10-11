package utils

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func filterServers(servers []models.Server,
	selection settings.ServerSelection,
) (filtered []models.Server) {
	for _, server := range servers {
		if filterServer(server, selection) {
			continue
		}

		filtered = append(filtered, server)
	}

	return filtered
}

func filterServer(server models.Server,
	selection settings.ServerSelection,
) (filtered bool) {
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

	if *selection.PortForwardOnly && !server.PortForward {
		return true
	}

	if *selection.SecureCoreOnly && !server.SecureCore {
		return true
	}

	if *selection.TorOnly && !server.Tor {
		return true
	}

	if filterByPossibilities(server.Country, selection.Countries) {
		return true
	}

	if filterAnyByPossibilities(server.Categories, selection.Categories) {
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

	if filterByPossibilities(server.Number, selection.Numbers) {
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

func filterByPossibilities[T string | uint16](value T, possibilities []T) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}
	for _, possibility := range possibilities {
		if strings.EqualFold(fmt.Sprint(value), fmt.Sprint(possibility)) {
			return false
		}
	}
	return true
}

func filterAnyByPossibilities(values, possibilities []string) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}

	for _, value := range values {
		if !filterByPossibilities(value, possibilities) {
			return false // found a valid value
		}
	}

	return true
}
