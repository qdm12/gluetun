package surfshark

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.SurfsharkServer) (s string) {
	s = "func SurfsharkServers() []models.SurfsharkServer {\n"
	s += "	return []models.SurfsharkServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
