package protonvpn

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.ProtonvpnServer) (s string) {
	s = "func ProtonvpnServers() []models.ProtonvpnServer {\n"
	s += "	return []models.ProtonvpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
