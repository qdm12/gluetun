package vpnunlimited

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.VPNUnlimitedServer) (s string) {
	s = "func VPNUnlimitedServers() []models.VPNUnlimitedServer {\n"
	s += "	return []models.VPNUnlimitedServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
