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
	cityHeader        = "City"
	countryHeader     = "Country"
	freeHeader        = "Free"
	hostnameHeader    = "Hostname"
	ispHeader         = "ISP"
	multiHopHeader    = "MultiHop"
	numberHeader      = "Number"
	ownedHeader       = "Owned"
	portForwardHeader = "Port forwarding"
	regionHeader      = "Region"
	streamHeader      = "Stream"
	tcpHeader         = "TCP"
	udpHeader         = "UDP"
	vpnHeader         = "VPN"
)

func (s *Server) ToMarkdown(headers ...string) (markdown string) {
	if len(headers) == 0 {
		return ""
	}

	fields := make([]string, len(headers))
	for i, header := range headers {
		switch header {
		case cityHeader:
			fields[i] = s.City
		case countryHeader:
			fields[i] = s.Country
		case freeHeader:
			fields[i] = boolToMarkdown(s.Free)
		case hostnameHeader:
			fields[i] = fmt.Sprintf("`%s`", s.Hostname)
		case ispHeader:
			fields[i] = s.ISP
		case multiHopHeader:
			fields[i] = boolToMarkdown(s.MultiHop)
		case numberHeader:
			fields[i] = fmt.Sprint(s.Number)
		case ownedHeader:
			fields[i] = boolToMarkdown(s.Owned)
		case portForwardHeader:
			fields[i] = boolToMarkdown(s.PortForward)
		case regionHeader:
			fields[i] = s.Region
		case streamHeader:
			fields[i] = boolToMarkdown(s.Stream)
		case tcpHeader:
			fields[i] = boolToMarkdown(s.TCP)
		case udpHeader:
			fields[i] = boolToMarkdown(s.UDP)
		case vpnHeader:
			fields[i] = s.VPN
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
		return []string{regionHeader, hostnameHeader, tcpHeader, udpHeader, portForwardHeader}
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
