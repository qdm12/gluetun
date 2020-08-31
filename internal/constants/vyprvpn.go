package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

const (
	VyprvpnCertificate = "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"
)

func VyprvpnRegionChoices() (choices []string) {
	servers := VyprvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return choices
}

func VyprvpnServers() []models.VyprvpnServer {
	return []models.VyprvpnServer{
		{Region: "Algeria", IPs: []net.IP{{209, 99, 75, 20}}},
		{Region: "Argentina", IPs: []net.IP{{209, 99, 109, 19}}},
		{Region: "Australia Melbourne", IPs: []net.IP{{209, 99, 117, 19}}},
		{Region: "Australia Perth", IPs: []net.IP{{209, 99, 1, 19}}},
		{Region: "Australia Sydney", IPs: []net.IP{{209, 99, 117, 18}}},
		{Region: "Austria", IPs: []net.IP{{128, 90, 96, 18}}},
		{Region: "Bahrain", IPs: []net.IP{{209, 99, 115, 19}}},
		{Region: "Belgium", IPs: []net.IP{{128, 90, 96, 20}}},
		{Region: "Brazil", IPs: []net.IP{{209, 99, 109, 20}}},
		{Region: "Bulgaria", IPs: []net.IP{{128, 90, 96, 22}}},
		{Region: "Canada", IPs: []net.IP{{209, 99, 21, 18}}},
		{Region: "Columbia", IPs: []net.IP{{209, 99, 109, 21}}},
		{Region: "Costa Rica", IPs: []net.IP{{209, 99, 109, 22}}},
		{Region: "Czech Republic", IPs: []net.IP{{128, 90, 96, 24}}},
		{Region: "Denmark", IPs: []net.IP{{128, 90, 96, 28}}},
		{Region: "Dubai", IPs: []net.IP{{128, 90, 45, 104}}},
		{Region: "Egypt", IPs: []net.IP{{128, 90, 228, 43}}},
		{Region: "El Salvador", IPs: []net.IP{{209, 99, 61, 20}}},
		{Region: "Finland", IPs: []net.IP{{128, 90, 96, 32}}},
		{Region: "France", IPs: []net.IP{{128, 90, 96, 34}}},
		{Region: "Germany", IPs: []net.IP{{128, 90, 96, 26}}},
		{Region: "Greece", IPs: []net.IP{{128, 90, 228, 59}}},
		{Region: "Hong Kong", IPs: []net.IP{{128, 90, 227, 18}}},
		{Region: "Iceland", IPs: []net.IP{{209, 99, 22, 20}}},
		{Region: "India", IPs: []net.IP{{209, 99, 115, 20}}},
		{Region: "Indonesia", IPs: []net.IP{{209, 99, 1, 20}}},
		{Region: "Ireland", IPs: []net.IP{{209, 99, 22, 19}}},
		{Region: "Israel", IPs: []net.IP{{128, 90, 228, 20}}},
		{Region: "Italy", IPs: []net.IP{{128, 90, 96, 36}}},
		{Region: "Japan", IPs: []net.IP{{209, 99, 113, 18}}},
		{Region: "Latvia", IPs: []net.IP{{128, 90, 96, 44}}},
		{Region: "Liechtenstein", IPs: []net.IP{{128, 90, 96, 38}}},
		{Region: "Lithuania", IPs: []net.IP{{128, 90, 96, 40}}},
		{Region: "Luxembourg", IPs: []net.IP{{128, 90, 96, 42}}},
		{Region: "Macao", IPs: []net.IP{{128, 90, 227, 36}}},
		{Region: "Malaysia", IPs: []net.IP{{209, 99, 1, 21}}},
		{Region: "Maldives", IPs: []net.IP{{209, 99, 1, 26}}},
		{Region: "Marshall Islands", IPs: []net.IP{{209, 99, 1, 25}}},
		{Region: "Mexico", IPs: []net.IP{{209, 99, 61, 19}}},
		{Region: "Netherlands", IPs: []net.IP{{128, 90, 96, 16}}},
		{Region: "New Zealand", IPs: []net.IP{{209, 99, 117, 20}}},
		{Region: "Norway", IPs: []net.IP{{128, 90, 96, 46}}},
		{Region: "Pakistan", IPs: []net.IP{{128, 90, 228, 67}}},
		{Region: "Panama", IPs: []net.IP{{209, 99, 109, 23}}},
		{Region: "Philippines", IPs: []net.IP{{209, 99, 1, 22}}},
		{Region: "Poland", IPs: []net.IP{{128, 90, 96, 48}}},
		{Region: "Portugal", IPs: []net.IP{{128, 90, 96, 50}}},
		{Region: "Qatar", IPs: []net.IP{{209, 99, 115, 21}}},
		{Region: "Romania", IPs: []net.IP{{128, 90, 96, 52}}},
		{Region: "Russia", IPs: []net.IP{{128, 90, 96, 54}}},
		{Region: "Saudi Arabia", IPs: []net.IP{{209, 99, 115, 22}}},
		{Region: "Singapore", IPs: []net.IP{{209, 99, 1, 18}}},
		{Region: "Slovakia", IPs: []net.IP{{128, 90, 96, 60}}},
		{Region: "Slovenia", IPs: []net.IP{{128, 90, 96, 58}}},
		{Region: "South Korea", IPs: []net.IP{{209, 99, 113, 19}}},
		{Region: "Spain", IPs: []net.IP{{128, 90, 96, 30}}},
		{Region: "Sweden", IPs: []net.IP{{128, 90, 96, 56}}},
		{Region: "Switzerland", IPs: []net.IP{{209, 99, 60, 18}}},
		{Region: "Taiwan", IPs: []net.IP{{128, 90, 227, 27}}},
		{Region: "Thailand", IPs: []net.IP{{209, 99, 1, 23}}},
		{Region: "Turkey", IPs: []net.IP{{128, 90, 96, 62}}},
		{Region: "USA Austin", IPs: []net.IP{{209, 99, 61, 18}}},
		{Region: "USA Chicago", IPs: []net.IP{{209, 99, 93, 18}}},
		{Region: "USA Los Angeles", IPs: []net.IP{{209, 99, 67, 18}}},
		{Region: "USA Miami", IPs: []net.IP{{209, 99, 109, 18}}},
		{Region: "USA New York", IPs: []net.IP{{209, 99, 63, 18}}},
		{Region: "USA San Francisco", IPs: []net.IP{{209, 99, 95, 18}}},
		{Region: "USA Seattle", IPs: []net.IP{{209, 99, 94, 18}}},
		{Region: "USA Washington", IPs: []net.IP{{209, 99, 62, 18}}},
		{Region: "USA Washington DC", IPs: []net.IP{{209, 99, 62, 18}}},
		{Region: "Ukraine", IPs: []net.IP{{128, 90, 96, 64}}},
		{Region: "United Kingdom", IPs: []net.IP{{209, 99, 22, 18}}},
		{Region: "Uruguay", IPs: []net.IP{{209, 99, 61, 21}}},
		{Region: "Vietnam", IPs: []net.IP{{209, 99, 1, 24}}},
	}
}
