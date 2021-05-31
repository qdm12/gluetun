package ipvanish

import "github.com/qdm12/gluetun/internal/models"

func Stringify(servers []models.IpvanishServer) (s string) {
	s = "func IpvanishServers() []models.IpvanishServer {\n"
	s += "	return []models.IpvanishServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
