package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	IvpnCA                 = "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"
	IvpnOpenvpnStaticKeyV1 = "ac470c93ff9f5602a8aab37dee84a52814d10f20490ad23c47d5d82120c1bf859e93d0696b455d4a1b8d55d40c2685c41ca1d0aef29a3efd27274c4ef09020a3978fe45784b335da6df2d12db97bbb838416515f2a96f04715fd28949c6fe296a925cfada3f8b8928ed7fc963c1563272f5cf46e5e1d9c845d7703ca881497b7e6564a9d1dea9358adffd435295479f47d5298fabf5359613ff5992cb57ff081a04dfb81a26513a6b44a9b5490ad265f8a02384832a59cc3e075ad545461060b7bcab49bac815163cb80983dd51d5b1fd76170ffd904d8291071e96efc3fb777856c717b148d08a510f5687b8a8285dcffe737b98916dd15ef6235dee4266d3b"
)

func IvpnCountryChoices() (choices []string) {
	servers := IvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func IvpnCityChoices() (choices []string) {
	servers := IvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func IvpnHostnameChoices() (choices []string) {
	servers := IvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

//nolint:lll
// IvpnServers returns a slice of all the server information for Ivpn.
func IvpnServers() []models.IvpnServer {
	return []models.IvpnServer{
		{Country: "Australia", City: "", Hostname: "au-nsw.gw.ivpn.net", IPs: []net.IP{{46, 102, 153, 242}}},
		{Country: "Austria", City: "", Hostname: "at.gw.ivpn.net", IPs: []net.IP{{185, 244, 212, 66}}},
		{Country: "Belgium", City: "", Hostname: "be.gw.ivpn.net", IPs: []net.IP{{194, 187, 251, 10}}},
		{Country: "Brazil", City: "", Hostname: "br.gw.ivpn.net", IPs: []net.IP{{45, 162, 229, 130}}},
		{Country: "Canada", City: "Montreal", Hostname: "ca-qc.gw.ivpn.net", IPs: []net.IP{{87, 101, 92, 26}}},
		{Country: "Canada", City: "Toronto", Hostname: "ca.gw.ivpn.net", IPs: []net.IP{{104, 254, 90, 178}}},
		{Country: "Czech Republic", City: "", Hostname: "cz.gw.ivpn.net", IPs: []net.IP{{195, 181, 160, 167}}},
		{Country: "Denmark", City: "", Hostname: "dk.gw.ivpn.net", IPs: []net.IP{{185, 245, 84, 226}}},
		{Country: "Finland", City: "", Hostname: "fi.gw.ivpn.net", IPs: []net.IP{{185, 112, 82, 12}}},
		{Country: "France", City: "", Hostname: "fr.gw.ivpn.net", IPs: []net.IP{{185, 246, 211, 179}}},
		{Country: "Germany", City: "", Hostname: "de.gw.ivpn.net", IPs: []net.IP{{178, 162, 211, 114}}},
		{Country: "Hong Kong", City: "", Hostname: "hk.gw.ivpn.net", IPs: []net.IP{{209, 58, 188, 13}}},
		{Country: "Hungary", City: "", Hostname: "hu.gw.ivpn.net", IPs: []net.IP{{185, 189, 114, 186}}},
		{Country: "Iceland", City: "", Hostname: "is.gw.ivpn.net", IPs: []net.IP{{82, 221, 107, 178}}},
		{Country: "Israel", City: "", Hostname: "il.gw.ivpn.net", IPs: []net.IP{{185, 191, 207, 194}}},
		{Country: "Italy", City: "", Hostname: "it.gw.ivpn.net", IPs: []net.IP{{158, 58, 172, 73}}},
		{Country: "Japan", City: "", Hostname: "jp.gw.ivpn.net", IPs: []net.IP{{91, 207, 174, 234}}},
		{Country: "Luxembourg", City: "", Hostname: "lu.gw.ivpn.net", IPs: []net.IP{{92, 223, 89, 53}}},
		{Country: "Netherlands", City: "", Hostname: "nl.gw.ivpn.net", IPs: []net.IP{{95, 211, 172, 95}}},
		{Country: "Norway", City: "", Hostname: "no.gw.ivpn.net", IPs: []net.IP{{194, 242, 10, 150}}},
		{Country: "Poland", City: "", Hostname: "pl.gw.ivpn.net", IPs: []net.IP{{185, 246, 208, 86}}},
		{Country: "Portugal", City: "", Hostname: "pt.gw.ivpn.net", IPs: []net.IP{{94, 46, 175, 112}}},
		{Country: "Romania", City: "", Hostname: "ro.gw.ivpn.net", IPs: []net.IP{{37, 120, 206, 50}}},
		{Country: "Serbia", City: "", Hostname: "rs.gw.ivpn.net", IPs: []net.IP{{141, 98, 103, 250}}},
		{Country: "Singapore", City: "", Hostname: "sg.gw.ivpn.net", IPs: []net.IP{{185, 128, 24, 186}}},
		{Country: "Slovakia", City: "", Hostname: "sk.gw.ivpn.net", IPs: []net.IP{{185, 245, 85, 250}}},
		{Country: "Sweden", City: "", Hostname: "se.gw.ivpn.net", IPs: []net.IP{{80, 67, 10, 138}}},
		{Country: "Switzerland", City: "", Hostname: "ch.gw.ivpn.net", IPs: []net.IP{{185, 212, 170, 138}}},
		{Country: "USA", City: "Atlanta", Hostname: "us-ga.gw.ivpn.net", IPs: []net.IP{{104, 129, 24, 146}}},
		{Country: "USA", City: "Chicago", Hostname: "us-il.gw.ivpn.net", IPs: []net.IP{{72, 11, 137, 146}}},
		{Country: "USA", City: "Dallas", Hostname: "us-tx.gw.ivpn.net", IPs: []net.IP{{96, 44, 189, 194}}},
		{Country: "USA", City: "Las Vegas", Hostname: "us-nv.gw.ivpn.net", IPs: []net.IP{{185, 242, 5, 34}}},
		{Country: "USA", City: "Los Angeles", Hostname: "us-ca.gw.ivpn.net", IPs: []net.IP{{69, 12, 80, 146}}},
		{Country: "USA", City: "Miami", Hostname: "us-fl.gw.ivpn.net", IPs: []net.IP{{173, 44, 49, 90}}},
		{Country: "USA", City: "New Jersey", Hostname: "us-nj.gw.ivpn.net", IPs: []net.IP{{23, 226, 128, 18}}},
		{Country: "USA", City: "New York", Hostname: "us-ny.gw.ivpn.net", IPs: []net.IP{{64, 120, 44, 114}}},
		{Country: "USA", City: "Phoenix", Hostname: "us-az.gw.ivpn.net", IPs: []net.IP{{193, 37, 254, 130}}},
		{Country: "USA", City: "Salt Lake City", Hostname: "us-ut.gw.ivpn.net", IPs: []net.IP{{198, 105, 216, 28}}},
		{Country: "USA", City: "Seattle", Hostname: "us-wa.gw.ivpn.net", IPs: []net.IP{{23, 19, 87, 209}}},
		{Country: "USA", City: "Washington", Hostname: "us-dc.gw.ivpn.net", IPs: []net.IP{{207, 244, 108, 207}}},
		{Country: "Ukraine", City: "", Hostname: "ua.gw.ivpn.net", IPs: []net.IP{{193, 203, 48, 54}}},
		{Country: "United Kingdom", City: "London", Hostname: "gb.gw.ivpn.net", IPs: []net.IP{{185, 59, 221, 133}, {185, 59, 221, 88}}},
		{Country: "United Kingdom", City: "Manchester", Hostname: "gb-man.gw.ivpn.net", IPs: []net.IP{{89, 238, 141, 228}}},
	}
}
