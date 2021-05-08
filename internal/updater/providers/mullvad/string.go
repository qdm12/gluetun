package mullvad

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.MullvadServer) (s string) {
	s = "func MullvadServers() []models.MullvadServer {\n"
	s += "	return []models.MullvadServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
