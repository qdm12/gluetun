package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	VPNUnlimitedCertificateAuthority = "MIIEjjCCA/egAwIBAgIJAKsVbHBdakXuMA0GCSqGSIb3DQEBBQUAMIHgMQswCQYDVQQGEwJVUzELMAkGA1UECBMCTlkxETAPBgNVBAcTCE5ldyBZb3JrMR8wHQYDVQQKExZTaW1wbGV4IFNvbHV0aW9ucyBJbmMuMRYwFAYDVQQLEw1WcG4gVW5saW1pdGVkMSMwIQYDVQQDExpzZXJ2ZXIudnBudW5saW1pdGVkYXBwLmNvbTEjMCEGA1UEKRMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xLjAsBgkqhkiG9w0BCQEWH3N1cHBvcnRAc2ltcGxleHNvbHV0aW9uc2luYy5jb20wHhcNMTMxMjE2MTM1OTQ0WhcNMjMxMjE0MTM1OTQ0WjCB4DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMREwDwYDVQQHEwhOZXcgWW9yazEfMB0GA1UEChMWU2ltcGxleCBTb2x1dGlvbnMgSW5jLjEWMBQGA1UECxMNVnBuIFVubGltaXRlZDEjMCEGA1UEAxMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xIzAhBgNVBCkTGnNlcnZlci52cG51bmxpbWl0ZWRhcHAuY29tMS4wLAYJKoZIhvcNAQkBFh9zdXBwb3J0QHNpbXBsZXhzb2x1dGlvbnNpbmMuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDADUzS8QWGvdhLFKsEzAiq5+b0ukKjBza0k6/dmCeYTvCVqGKg/h1IAtQdVVLAkmEp0zvGH7PuOhXm7zZrCouBr/RiG4tEcoRHwc6AJmowkYERlY7+xGx3OuNrD00QceNTsan0bn78jwt0zhFNmHdoTtFjgK3oqmQYSAtbEVWYgwIDAQABo4IBTDCCAUgwHQYDVR0OBBYEFKClmYP+tMNgWagUJCCHjtaui2k+MIIBFwYDVR0jBIIBDjCCAQqAFKClmYP+tMNgWagUJCCHjtaui2k+oYHmpIHjMIHgMQswCQYDVQQGEwJVUzELMAkGA1UECBMCTlkxETAPBgNVBAcTCE5ldyBZb3JrMR8wHQYDVQQKExZTaW1wbGV4IFNvbHV0aW9ucyBJbmMuMRYwFAYDVQQLEw1WcG4gVW5saW1pdGVkMSMwIQYDVQQDExpzZXJ2ZXIudnBudW5saW1pdGVkYXBwLmNvbTEjMCEGA1UEKRMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xLjAsBgkqhkiG9w0BCQEWH3N1cHBvcnRAc2ltcGxleHNvbHV0aW9uc2luYy5jb22CCQCrFWxwXWpF7jAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBBQUAA4GBALkWhfw7SSV7it+ZYZmT+cQbExjlYgQ40zae2J2CqIYACRcfsDHvh7Q+fiwSocevv2NE0dWi6WB2H6SiudYjvDvubAX/QbzfBxtbxCCoAIlfPCm8xOnWFN7TUJAzWwHJkKgEnu29GZHu2x8J+7VeDbKH5RTMHHe8FkSxh6Zz/BMN"
	VPNUnlimitedCertificate          = "MIIFXTCCBMagAwIBAgIBATANBgkqhkiG9w0BAQUFADCB4DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMREwDwYDVQQHEwhOZXcgWW9yazEfMB0GA1UEChMWU2ltcGxleCBTb2x1dGlvbnMgSW5jLjEWMBQGA1UECxMNVnBuIFVubGltaXRlZDEjMCEGA1UEAxMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xIzAhBgNVBCkTGnNlcnZlci52cG51bmxpbWl0ZWRhcHAuY29tMS4wLAYJKoZIhvcNAQkBFh9zdXBwb3J0QHNpbXBsZXhzb2x1dGlvbnNpbmMuY29tMB4XDTIxMDYwMTE3NDU1M1oXDTMxMDUzMDE3NTA1M1owgeUxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJOWTERMA8GA1UEBxMITmV3IFlvcmsxHzAdBgNVBAoTFlNpbXBsZXggU29sdXRpb25zIEluYy4xIzAhBgNVBAsTGml0LW1pbC52cG51bmxpbWl0ZWRhcHAuY29tMS0wKwYDVQQDEyRLUzItN2M4YTg5MzNiN2VmOTZiOTQ4Yjg3YjZhN2I2MDNlZTgxLTArBgNVBCkTJEVCMjI5MjY3LUExRTItNDg3Ni1BNTE5LTZGNzZEQURFNDRGNTESMBAGCSqGSIb3DQEJARYDWFhYMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzzUL+7RnimM6gXHILVNlZypg0/HRNAAkMSmEvrQ+/EI4gyAuDJLEIDT9Pr+wBQv7Bkh8xyJikZ6oClbjEfY+JImUvjuJaeeaSWJfPzBLDtW4Kb0N4l4I5OlmZrmgMYT/8TnXZ8c7T7q9MSwR3cwVetlEc8KZh5DNEcioyEdP7/jQ5tOTcpQdopMu+PyDBeiuKiQbS+O0kyIw2yaBxMIEZaXvyr5DuPPjHa4zcUAFPWfLt4ea2/HhOUE46X/oMUqzErM0nkqYFKyCYdhBFDtr5b3BF58BQ79Lyfd2hUsqqHNHaPCx1w9spE6UBdUbpVCFg62w1xd/DgJ3LrjutOSOxwIDAQABo4IBmjCCAZYwCQYDVR0TBAIwADAtBglghkgBhvhCAQ0EIBYeRWFzeS1SU0EgR2VuZXJhdGVkIENlcnRpZmljYXRlMB0GA1UdDgQWBBSCiJtbYMNuHtRV6XbBYqf37sbPlDCCARcGA1UdIwSCAQ4wggEKgBSgpZmD/rTDYFmoFCQgh47WrotpPqGB5qSB4zCB4DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMREwDwYDVQQHEwhOZXcgWW9yazEfMB0GA1UEChMWU2ltcGxleCBTb2x1dGlvbnMgSW5jLjEWMBQGA1UECxMNVnBuIFVubGltaXRlZDEjMCEGA1UEAxMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xIzAhBgNVBCkTGnNlcnZlci52cG51bmxpbWl0ZWRhcHAuY29tMS4wLAYJKoZIhvcNAQkBFh9zdXBwb3J0QHNpbXBsZXhzb2x1dGlvbnNpbmMuY29tggkAqxVscF1qRe4wEwYDVR0lBAwwCgYIKwYBBQUHAwIwCwYDVR0PBAQDAgeAMA0GCSqGSIb3DQEBBQUAA4GBAFpRuj2vyLjDkUGnS6vmfCNZapEJHjLs8yDnqXXNPmyEyOsn1ivvW9uoWqOlA0VThkwoMb8+jAZCSwHavfoULBRPAQzddy5WCuG1OkhCs2eK/JUfIKEx8/XSf7hFJoEbBwRA4U7napmtaWGeVnV4SssPVIFGWuW8KM5Jv1hC9jys"
)

func VPNUnlimitedCountryChoices() (choices []string) {
	servers := VPNUnlimitedServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func VPNUnlimitedCityChoices() (choices []string) {
	servers := VPNUnlimitedServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func VPNUnlimitedHostnameChoices() (choices []string) {
	servers := VPNUnlimitedServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

//nolint:lll
// VPNUnlimitedServers returns a slice of all the server information for VPNUnlimited.
func VPNUnlimitedServers() []models.VPNUnlimitedServer {
	return []models.VPNUnlimitedServer{
		{Country: "Argentina", City: "", Hostname: "ar.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{131, 255, 4, 187}, {131, 255, 4, 235}}},
		{Country: "Australia", City: "Sydney", Hostname: "au-syd.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{139, 99, 131, 38}, {139, 99, 130, 220}}},
		{Country: "Austria", City: "", Hostname: "at.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{50, 7, 115, 19}, {185, 210, 219, 194}, {185, 210, 219, 198}}},
		{Country: "Belarus", City: "", Hostname: "by.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 66, 68, 23}, {185, 66, 71, 108}, {185, 66, 71, 115}, {185, 66, 68, 84}, {185, 66, 71, 22}, {185, 66, 71, 122}, {185, 66, 71, 8}, {185, 66, 68, 18}}},
		{Country: "Belgium", City: "", Hostname: "be.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{37, 120, 143, 178}, {82, 102, 19, 110}}},
		{Country: "Bosnia and Herzegovina", City: "", Hostname: "ba.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 164, 35, 37}}},
		{Country: "Brazil", City: "", Hostname: "br.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{66, 90, 70, 7}, {177, 67, 82, 222}, {177, 67, 82, 218}, {177, 67, 82, 228}, {177, 67, 82, 221}, {177, 67, 82, 219}, {177, 67, 82, 210}, {177, 67, 82, 223}}},
		{Country: "Bulgaria", City: "", Hostname: "bg.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{37, 120, 152, 26}}},
		{Country: "Canada", City: "", Hostname: "ca.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{192, 99, 7, 170}, {192, 99, 6, 164}, {192, 99, 14, 158}, {192, 99, 37, 199}, {192, 99, 6, 166}}},
		{Country: "Canada", City: "Toronto", Hostname: "ca-tr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{184, 75, 213, 154}, {162, 253, 131, 106}, {162, 219, 176, 163}, {104, 254, 90, 58}, {104, 254, 92, 42}, {184, 75, 213, 194}, {204, 187, 100, 82}, {104, 254, 90, 34}}},
		{Country: "Canada", City: "Vancouver", Hostname: "ca-vn.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{162, 221, 202, 138}, {162, 221, 202, 134}, {162, 221, 202, 17}}},
		{Country: "Costa Rica", City: "", Hostname: "cr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{138, 59, 19, 8}}},
		{Country: "Croatia", City: "", Hostname: "hr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{85, 10, 56, 3}, {85, 10, 56, 100}, {85, 10, 51, 3}, {85, 10, 56, 4}}},
		{Country: "Cyprus", City: "", Hostname: "cy.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{194, 30, 136, 123}, {194, 30, 136, 102}}},
		{Country: "Czech Republic", City: "", Hostname: "cz.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 216, 35, 46}}},
		{Country: "Denmark", City: "", Hostname: "dk.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{37, 120, 194, 174}, {89, 45, 7, 90}}},
		{Country: "Estonia", City: "", Hostname: "ee.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 155, 96, 132}}},
		{Country: "Finland", City: "", Hostname: "fi.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 112, 82, 26}, {91, 233, 116, 44}}},
		{Country: "France", City: "", Hostname: "fr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{195, 154, 189, 85}, {195, 154, 204, 36}, {195, 154, 211, 84}, {62, 210, 38, 83}, {195, 154, 169, 66}, {195, 154, 166, 20}, {62, 210, 204, 161}, {62, 210, 205, 16}, {62, 210, 139, 124}, {195, 154, 221, 54}, {195, 154, 199, 175}, {195, 154, 189, 212}, {195, 154, 180, 96}, {195, 154, 199, 155}, {195, 154, 209, 149}, {195, 154, 222, 168}}},
		{Country: "France", City: "Roubaix", Hostname: "fr-rbx.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{151, 80, 27, 199}, {51, 255, 71, 16}, {147, 135, 137, 107}}},
		{Country: "Germany", City: "", Hostname: "de.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{84, 16, 236, 168}, {178, 162, 205, 77}, {178, 162, 205, 115}, {178, 162, 205, 82}, {84, 16, 238, 220}}},
		{Country: "Germany", City: "Düsseldorf", Hostname: "de-dus.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{146, 0, 42, 77}, {146, 0, 32, 121}, {85, 14, 243, 42}}},
		{Country: "Greece", City: "", Hostname: "gr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{154, 57, 3, 36}}},
		{Country: "Hungary", City: "", Hostname: "hu.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{91, 219, 238, 174}}},
		{Country: "Iceland", City: "", Hostname: "is.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{37, 235, 49, 75}, {151, 236, 24, 114}, {37, 235, 49, 42}, {37, 235, 49, 244}, {37, 235, 49, 70}, {37, 235, 49, 15}, {151, 236, 24, 142}}},
		{Country: "India", City: "", Hostname: "in.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{103, 26, 204, 46}, {103, 26, 204, 84}}},
		{Country: "India", City: "Karnataka", Hostname: "in-ka.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{139, 59, 67, 224}, {139, 59, 65, 136}, {139, 59, 24, 197}, {139, 59, 24, 198}, {139, 59, 24, 199}, {139, 59, 24, 191}, {139, 59, 16, 78}}},
		{Country: "Ireland", City: "Dublin", Hostname: "ie-dub.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{188, 241, 178, 118}}},
		{Country: "Isle of Man", City: "", Hostname: "im.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{192, 71, 211, 15}, {37, 235, 55, 62}, {37, 235, 55, 68}, {37, 235, 55, 14}, {37, 235, 55, 45}}},
		{Country: "Israel", City: "", Hostname: "il.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{193, 182, 144, 170}, {193, 182, 144, 151}, {193, 182, 144, 23}}},
		{Country: "Italy", City: "Milan", Hostname: "it-mil.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{91, 193, 5, 50}}},
		{Country: "Japan", City: "", Hostname: "jp.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{172, 104, 67, 80}, {172, 104, 75, 121}, {161, 202, 97, 109}, {172, 104, 71, 35}, {172, 104, 126, 64}, {172, 104, 64, 213}, {85, 208, 110, 122}, {172, 104, 78, 67}, {172, 104, 83, 34}, {172, 104, 99, 172}, {161, 202, 97, 101}, {172, 104, 83, 13}, {172, 104, 87, 163}, {172, 104, 99, 14}}},
		{Country: "Korea", City: "", Hostname: "kr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{160, 202, 162, 60}, {160, 202, 163, 122}}},
		{Country: "Kuala Lumpur", City: "", Hostname: "mys.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{111, 90, 141, 34}, {111, 90, 158, 159}}},
		{Country: "Latvia", City: "", Hostname: "lv.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{195, 123, 213, 172}, {195, 123, 209, 168}}},
		{Country: "Libya", City: "", Hostname: "ly.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{41, 208, 71, 39}}},
		{Country: "Lithuania", City: "", Hostname: "lt.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 64, 105, 93}, {185, 64, 104, 142}, {185, 25, 48, 81}}},
		{Country: "Luxembourg", City: "", Hostname: "lu.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{94, 242, 246, 45}, {94, 242, 246, 46}, {94, 242, 246, 77}, {94, 242, 246, 38}, {94, 242, 246, 47}}},
		{Country: "Mexico", City: "", Hostname: "mx.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{169, 57, 35, 104}}},
		{Country: "Moldova", City: "", Hostname: "md.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 163, 46, 141}}},
		{Country: "Netherlands", City: "", Hostname: "nl.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{190, 2, 132, 115}, {190, 2, 132, 16}, {190, 2, 132, 52}}},
		{Country: "New Zealand", City: "", Hostname: "nz.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{103, 75, 119, 109}, {103, 75, 119, 107}, {103, 75, 119, 105}, {103, 75, 119, 102}, {103, 75, 119, 103}, {103, 75, 119, 104}, {103, 75, 119, 106}, {103, 75, 119, 108}, {103, 75, 119, 101}, {103, 75, 119, 11}}},
		{Country: "Norway", City: "", Hostname: "no.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 90, 61, 21}, {185, 90, 61, 74}}},
		{Country: "Oman", City: "", Hostname: "om.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 226, 124, 110}}},
		{Country: "Poland", City: "", Hostname: "pl.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{37, 120, 156, 234}}},
		{Country: "Portugal", City: "", Hostname: "pt.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 31, 159, 110}}},
		{Country: "Romania", City: "", Hostname: "ro.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 144, 83, 11}, {93, 115, 92, 207}, {185, 144, 83, 13}, {77, 81, 98, 70}, {93, 115, 92, 208}}},
		{Country: "Serbia", City: "", Hostname: "rs.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{152, 89, 160, 142}}},
		{Country: "Singapore", City: "", Hostname: "sg-free.vpnunlimitedapp.com", Free: true, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{178, 128, 48, 177}, {188, 166, 177, 236}, {206, 189, 80, 158}, {178, 128, 117, 139}}},
		{Country: "Singapore", City: "", Hostname: "sg.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{50, 7, 60, 52}, {117, 20, 41, 9}, {117, 20, 41, 10}, {139, 99, 62, 61}}},
		{Country: "Slovakia", City: "", Hostname: "sk.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{46, 29, 2, 131}}},
		{Country: "Slovenia", City: "", Hostname: "si.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{192, 71, 244, 38}, {192, 71, 244, 28}}},
		{Country: "South Africa", City: "", Hostname: "za.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{129, 232, 130, 187}, {129, 232, 219, 196}, {129, 232, 133, 41}, {129, 232, 129, 157}, {129, 232, 219, 195}, {129, 232, 134, 122}, {129, 232, 130, 186}, {129, 232, 130, 179}}},
		{Country: "Spain", City: "", Hostname: "es.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{194, 99, 104, 58}, {195, 206, 107, 134}}},
		{Country: "Sweden", City: "", Hostname: "se.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{91, 132, 138, 174}}},
		{Country: "Switzerland", City: "", Hostname: "ch.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{152, 89, 162, 90}, {152, 89, 162, 86}}},
		{Country: "Thailand", City: "", Hostname: "th.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{103, 159, 3, 121}}},
		{Country: "Turkey", City: "", Hostname: "tr.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 17, 115, 62}, {185, 73, 202, 218}}},
		{Country: "United Arab Emirates", City: "", Hostname: "ae.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{45, 9, 249, 201}, {45, 9, 249, 209}, {45, 9, 250, 147}, {45, 9, 250, 249}, {45, 9, 250, 138}, {45, 9, 249, 211}, {45, 9, 250, 144}, {45, 9, 249, 216}, {45, 9, 250, 145}, {45, 9, 250, 143}, {45, 9, 250, 134}, {45, 9, 250, 146}, {45, 9, 249, 202}}},
		{Country: "United Kingdom", City: "", Hostname: "uk.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{176, 227, 198, 122}, {77, 245, 65, 2}, {109, 73, 77, 18}, {5, 152, 213, 186}, {88, 150, 224, 74}, {88, 150, 180, 82}}},
		{Country: "United Kingdom", City: "London", Hostname: "uk-cv.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{5, 101, 169, 146}, {5, 101, 143, 66}, {178, 159, 10, 78}}},
		{Country: "United Kingdom", City: "London", Hostname: "uk-lon.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{78, 110, 160, 6}, {50, 7, 114, 220}, {5, 101, 136, 154}, {23, 106, 33, 89}}},
		{Country: "United States", City: "", Hostname: "us-stream.vpnunlimitedapp.com", Free: false, Stream: true, TCP: false, UDP: true, IPs: []net.IP{{167, 99, 96, 113}, {198, 199, 114, 155}, {192, 241, 194, 208}, {165, 227, 49, 171}, {159, 89, 128, 183}, {138, 197, 203, 142}, {64, 227, 111, 143}, {143, 110, 225, 79}, {164, 90, 144, 44}, {198, 199, 103, 243}, {198, 199, 96, 52}, {128, 199, 9, 51}, {143, 110, 156, 9}, {178, 128, 79, 75}, {198, 199, 97, 247}, {198, 199, 105, 87}, {198, 199, 103, 88}, {206, 189, 211, 205}, {206, 189, 165, 44}, {161, 35, 236, 242}}},
		{Country: "United States", City: "", Hostname: "us.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{199, 115, 117, 81}, {199, 115, 115, 139}, {207, 244, 72, 212}, {199, 115, 115, 160}, {199, 115, 117, 73}}},
		{Country: "United States", City: "Chicago", Hostname: "us-chi.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{66, 23, 205, 226}, {108, 62, 202, 211}, {173, 234, 41, 130}, {174, 34, 184, 130}}},
		{Country: "United States", City: "Dallas", Hostname: "us-dal.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{172, 241, 115, 99}, {172, 241, 113, 19}, {172, 241, 112, 92}}},
		{Country: "United States", City: "Denver", Hostname: "us-den.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{23, 237, 26, 131}, {23, 237, 26, 67}}},
		{Country: "United States", City: "Houston", Hostname: "us-hou.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{162, 218, 228, 194}, {162, 218, 228, 196}, {66, 187, 75, 122}, {162, 218, 229, 106}}},
		{Country: "United States", City: "Las Vegas", Hostname: "us-lv.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{185, 242, 5, 18}, {185, 242, 5, 22}}},
		{Country: "United States", City: "Los Angeles", Hostname: "us-la.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{64, 31, 33, 154}, {192, 184, 48, 10}, {23, 83, 37, 209}, {23, 83, 37, 213}, {45, 136, 131, 40}, {69, 162, 99, 70}}},
		{Country: "United States", City: "Miami", Hostname: "us-mia.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{186, 233, 185, 29}, {186, 233, 185, 79}, {186, 233, 185, 80}, {186, 233, 184, 187}, {186, 233, 184, 188}, {186, 233, 184, 37}, {186, 233, 184, 31}}},
		{Country: "United States", City: "New York", Hostname: "us-ny-free.vpnunlimitedapp.com", Free: true, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{64, 94, 215, 66}, {64, 94, 215, 162}, {64, 94, 215, 170}}},
		{Country: "United States", City: "New York", Hostname: "us-ny.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{64, 42, 178, 202}, {23, 105, 134, 162}, {64, 42, 178, 226}, {23, 237, 58, 112}, {23, 108, 31, 122}}},
		{Country: "United States", City: "Saint Louis", Hostname: "us-sl.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{69, 64, 58, 110}, {69, 64, 58, 255}, {209, 239, 123, 77}, {209, 239, 123, 107}}},
		{Country: "United States", City: "Salt Lake City", Hostname: "us-slc.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{209, 95, 53, 223}, {209, 95, 53, 225}}},
		{Country: "United States", City: "San Francisco", Hostname: "us-sf.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{209, 58, 139, 34}, {209, 58, 135, 72}, {209, 58, 130, 210}, {209, 58, 135, 106}, {209, 58, 135, 108}, {209, 58, 135, 74}, {209, 58, 139, 35}, {209, 58, 137, 94}, {172, 241, 251, 164}, {209, 58, 135, 120}}},
		{Country: "United States", City: "Seattle", Hostname: "us-sea.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{216, 244, 82, 50}, {23, 81, 209, 137}, {216, 244, 82, 210}}},
		{Country: "Vietnam", City: "", Hostname: "vn.vpnunlimitedapp.com", Free: false, Stream: false, TCP: false, UDP: true, IPs: []net.IP{{103, 9, 78, 84}, {146, 196, 67, 7}}},
	}
}
