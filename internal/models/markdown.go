package models

import (
	"fmt"
	"strings"
)

func boolToMarkdown(b bool) string {
	if b {
		return "✅"
	}
	return "❎"
}

func markdownTableHeading(legendFields ...string) (markdown string) {
	return "| " + strings.Join(legendFields, " | ") + " |\n" +
		"|" + strings.Repeat(" --- |", len(legendFields)) + "\n"
}

func (s *CyberghostServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "Group", "Hostname")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s CyberghostServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` |", s.Region, s.Group, s.Hostname)
}

func (s *FastestvpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *FastestvpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | `%s` | %s | %s |",
		s.Country, s.Hostname, boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *HideMyAssServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "Region", "City", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *HideMyAssServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | `%s` | %s | %s |",
		s.Country, s.Region, s.City, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *IpvanishServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *IpvanishServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` | %s | %s |",
		s.Country, s.City, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *IvpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "ISP", "Hostname", "VPN", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *IvpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | `%s` | %s | %s | %s |",
		s.Country, s.City, s.ISP, s.Hostname, s.VPN,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *MullvadServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "ISP", "Owned",
		"Hostname", "VPN")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *MullvadServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | %s | `%s` | %s |",
		s.Country, s.City, s.ISP, boolToMarkdown(s.Owned),
		s.Hostname, s.VPN)
}

func (s *NordvpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *NordvpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | `%s` | %s | %s |",
		s.Region, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *PrivadoServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "Region", "City", "Hostname")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *PrivadoServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | `%s` |",
		s.Country, s.Region, s.City, s.Hostname)
}

func (s *PiaServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *PIAServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | `%s` | %s | %s |",
		s.Region, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *PrivatevpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "Hostname")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *PrivatevpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` |",
		s.Country, s.City, s.Hostname)
}

func (s *ProtonvpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "Region", "City", "Hostname", "Free tier")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *ProtonvpnServer) ToMarkdown() (markdown string) {
	isFree := strings.Contains(strings.ToLower(s.Name), "free")
	return fmt.Sprintf("| %s | %s | %s | `%s` | %s |",
		s.Country, s.Region, s.City, s.Hostname, boolToMarkdown(isFree))
}

func (s *PurevpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "Region", "City", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *PurevpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | `%s` | %s | %s |",
		s.Country, s.Region, s.City, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *SurfsharkServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "Country", "City", "Hostname", "Multi-hop", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *SurfsharkServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | %s | `%s` | %s | %s | %s |",
		s.Region, s.Country, s.City, s.Hostname, boolToMarkdown(s.MultiHop),
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *TorguardServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *TorguardServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` | %s | %s |",
		s.Country, s.City, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *VPNUnlimitedServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Country", "City", "Hostname", "Free tier", "Streaming", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *VPNUnlimitedServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` | %s | %s | %s | %s |",
		s.Country, s.City, s.Hostname,
		boolToMarkdown(s.Free), boolToMarkdown(s.Stream),
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *VyprvpnServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "Hostname", "TCP", "UDP")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *VyprvpnServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | `%s` | %s | %s |",
		s.Region, s.Hostname,
		boolToMarkdown(s.TCP), boolToMarkdown(s.UDP))
}

func (s *WindscribeServers) ToMarkdown() (markdown string) {
	markdown = markdownTableHeading("Region", "City", "Hostname", "VPN")
	for _, server := range s.Servers {
		markdown += server.ToMarkdown() + "\n"
	}
	return markdown
}

func (s *WindscribeServer) ToMarkdown() (markdown string) {
	return fmt.Sprintf("| %s | %s | `%s` | %s |",
		s.Region, s.City, s.Hostname, s.VPN)
}
