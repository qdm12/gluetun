package cyberghost

import "github.com/qdm12/gluetun/internal/models"

// Stringify converts servers to code string format.
func Stringify(servers []models.CyberghostServer) (s string) {
	s = "func CyberghostServers() []models.CyberghostServer {\n"
	s += "	return []models.CyberghostServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
