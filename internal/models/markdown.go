package models

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants/providers"
)

func boolToMarkdown(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

func markdownTableHeading(legendFields ...string) (markdown string) {
	return "| " + strings.Join(legendFields, " | ") + " |\n" +
		"|" + strings.Repeat(" --- |", len(legendFields))
}

const (
	vpnHeader         = "VPN"
	countryHeader     = "Country"
	regionHeader      = "Region"
	cityHeader        = "City"
	ispHeader         = "ISP"
	ownedHeader       = "Owned"
	numberHeader      = "Number"
	hostnameHeader    = "Hostname"
	tcpHeader         = "TCP"
	udpHeader         = "UDP"
	retroLocHeader    = "Retro region" // TODO
	multiHopHeader    = "MultiHop"
	freeHeader        = "Free"
	streamHeader      = "Stream"
	portForwardHeader = "Port forwarding"
)

func (s *Server) ToMarkdown(headers ...string) (markdown string) {
	if len(headers) == 0 {
		return ""
	}

	fields := make([]string, len(headers))
	for i, header := range headers {
		switch header {
		case vpnHeader:
			fields[i] = s.VPN
		case countryHeader:
			fields[i] = s.Country
		case regionHeader:
			fields[i] = s.Region
		case cityHeader:
			fields[i] = s.City
		case ispHeader:
			fields[i] = s.ISP
		case ownedHeader:
			fields[i] = boolToMarkdown(s.Owned)
		case numberHeader:
			fields[i] = fmt.Sprint(s.Number)
		case hostnameHeader:
			fields[i] = fmt.Sprintf("`%s`", s.Hostname)
		case tcpHeader:
			fields[i] = boolToMarkdown(s.TCP)
		case udpHeader:
			fields[i] = boolToMarkdown(s.UDP)
		case retroLocHeader:
			fields[i] = s.RetroLoc
		case multiHopHeader:
			fields[i] = boolToMarkdown(s.MultiHop)
		case freeHeader:
			fields[i] = boolToMarkdown(s.Free)
		case streamHeader:
			fields[i] = boolToMarkdown(s.Stream)
		case portForwardHeader:
			fields[i] = boolToMarkdown(s.PortForward)
		}
	}

	return "| " + strings.Join(fields, " | ") + " |"
}

func (s *Servers) ToMarkdown(vpnProvider string) (markdown string) {
	headers := getMarkdownHeaders(vpnProvider)

	legend := markdownTableHeading(headers...)

	entries := make([]string, len(s.Servers))
	for i, server := range s.Servers {
		entries[i] = server.ToMarkdown(headers...)
	}

	markdown = legend + "\n" +
		strings.Join(entries, "\n") + "\n"
	return markdown
}

func getMarkdownHeaders(vpnProvider string) (headers []string) {
	switch vpnProvider {
	case providers.Cyberghost:
		return []string{countryHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Expressvpn:
		return []string{countryHeader, cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Fastestvpn:
		return []string{countryHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.HideMyAss:
		return []string{countryHeader, regionHeader, cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Ipvanish:
		return []string{countryHeader, cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Ivpn:
		return []string{countryHeader, cityHeader, ispHeader, hostnameHeader, vpnHeader, tcpHeader, udpHeader}
	case providers.Mullvad:
		return []string{countryHeader, cityHeader, ispHeader, ownedHeader, hostnameHeader, vpnHeader}
	case providers.Nordvpn:
		return []string{countryHeader, regionHeader, cityHeader, hostnameHeader}
	case providers.Perfectprivacy:
		return []string{cityHeader, tcpHeader, udpHeader}
	case providers.Privado:
		return []string{countryHeader, regionHeader, cityHeader, hostnameHeader}
	case providers.PrivateInternetAccess:
		return []string{regionHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Privatevpn:
		return []string{countryHeader, cityHeader, hostnameHeader}
	case providers.Protonvpn:
		return []string{countryHeader, regionHeader, cityHeader, hostnameHeader, freeHeader}
	case providers.Purevpn:
		return []string{countryHeader, regionHeader, cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Surfshark:
		return []string{regionHeader, countryHeader, cityHeader, hostnameHeader, multiHopHeader, tcpHeader, udpHeader}
	case providers.Torguard:
		return []string{countryHeader, cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.VPNUnlimited:
		return []string{countryHeader, cityHeader, hostnameHeader, freeHeader, streamHeader, tcpHeader, udpHeader}
	case providers.Vyprvpn:
		return []string{regionHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Wevpn:
		return []string{cityHeader, hostnameHeader, tcpHeader, udpHeader}
	case providers.Windscribe:
		return []string{regionHeader, cityHeader, hostnameHeader, vpnHeader}
	default:
		return nil
	}
}
