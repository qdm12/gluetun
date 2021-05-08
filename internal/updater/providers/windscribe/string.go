package windscribe

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.WindscribeServer) (s string) {
	s = "func WindscribeServers() []models.WindscribeServer {\n"
	s += "	return []models.WindscribeServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
