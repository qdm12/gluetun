package torguard

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.TorguardServer) (s string) {
	s = "func TorguardServers() []models.TorguardServer {\n"
	s += "	return []models.TorguardServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
