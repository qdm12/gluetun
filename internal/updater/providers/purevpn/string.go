package purevpn

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.PurevpnServer) (s string) {
	s = "func PurevpnServers() []models.PurevpnServer {\n"
	s += "	return []models.PurevpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
