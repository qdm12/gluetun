package hidemyass

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.HideMyAssServer) (s string) {
	s = "func HideMyAssServers() []models.HideMyAssServer {\n"
	s += "	return []models.HideMyAssServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
