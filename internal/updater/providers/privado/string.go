package privado

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.PrivadoServer) (s string) {
	s = "func PrivadoServers() []models.PrivadoServer {\n"
	s += "	return []models.PrivadoServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
