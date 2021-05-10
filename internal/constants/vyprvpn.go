package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	VyprvpnCertificate = "MIIGDjCCA/agAwIBAgIJAL2ON5xbane/MA0GCSqGSIb3DQEBDQUAMIGTMQswCQYDVQQGEwJDSDEQMA4GA1UECAwHTHVjZXJuZTEPMA0GA1UEBwwGTWVnZ2VuMRkwFwYDVQQKDBBHb2xkZW4gRnJvZyBHbWJIMSEwHwYDVQQDDBhHb2xkZW4gRnJvZyBHbWJIIFJvb3QgQ0ExIzAhBgkqhkiG9w0BCQEWFGFkbWluQGdvbGRlbmZyb2cuY29tMB4XDTE5MTAxNzIwMTQxMFoXDTM5MTAxMjIwMTQxMFowgZMxCzAJBgNVBAYTAkNIMRAwDgYDVQQIDAdMdWNlcm5lMQ8wDQYDVQQHDAZNZWdnZW4xGTAXBgNVBAoMEEdvbGRlbiBGcm9nIEdtYkgxITAfBgNVBAMMGEdvbGRlbiBGcm9nIEdtYkggUm9vdCBDQTEjMCEGCSqGSIb3DQEJARYUYWRtaW5AZ29sZGVuZnJvZy5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCtuddaZrpWZ+nUuJpG+ohTquO3XZtq6d4U0E2oiPeIiwm+WWLY49G+GNJb5aVrlrBojaykCAc2sU6NeUlpg3zuqrDqLcz7PAE4OdNiOdrLBF1o9ZHrcITDZN304eAY5nbyHx5V6x/QoDVCi4g+5OVTA+tZjpcl4wRIpgknWznO73IKCJ6YckpLn1BsFrVCb2ehHYZLg7Js58FzMySIxBmtkuPeHQXL61DFHh3cTFcMxqJjzh7EGsWRyXfbAaBGYnT+TZwzpLXXt8oBGpNXG8YBDrPdK0A+lzMnJ4nS0rgHDSRF0brx+QYk/6CgM510uFzB7zytw9UTD3/5TvKlCUmTGGgI84DbJ3DEvjxbgiQnJXCUZKKYSHwrK79Y4Qn+lXu4Bu0ZTCJBje0GUVMTPAvBCeDvzSe0iRcVSNMJVM68d4kD1PpSY/zWfCz5hiOjHWuXinaoZ0JJqRF8kGbJsbDlDYDtVvh/Cd4aWN6Q/2XLpszBsG5i8sdkS37nzkdlRwNEIZwsKfcXwdTOlDinR1LUG68LmzJAwfNE47xbrZUsdGGfG+HSPsrqFFiLGe7Y4e2+a7vGdSY9qR9PAzyx0ijCCrYzZDIsb2dwjLctUx6a3LNV8cpfhKX+s6tfMldGufPI7byHT1Ybf0NtMS1d1RjD6IbqedXQdCKtaw68kTX//wIDAQABo2MwYTAdBgNVHQ4EFgQU2EbQvBd1r/EADr2jCPMXsH7zEXEwHwYDVR0jBBgwFoAU2EbQvBd1r/EADr2jCPMXsH7zEXEwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwDQYJKoZIhvcNAQENBQADggIBAAViCPieIronV+9asjZyo5oSZSNWUkWRYdezjezsf49+fwT12iRgnkSEQeoj5caqcOfNm/eRpN4G7jhhCcxy9RGF+GurIlZ4v0mChZbx1jcxqr9/3/Z2TqvHALyWngBYDv6pv1iWcd9a4+QL9kj1Tlp8vUDIcHMtDQkEHnkhC+MnjyrdsdNE5wjlLljjFR2Qy5a6/kWwZ1JQVYof1J1EzY6mU7YLMHOdjfmeci5i0vg8+9kGMsc/7Wm69L1BeqpDB3ZEAgmOtda2jwOevJ4sABmRoSThFp4DeMcxb62HW1zZCCpgzWv/33+pZdPvnZHSz7RGoxH4Ln7eBf3oo2PMlu7wCsid3HUdgkRf2Og1RJIrFfEjb7jga1JbKX2Qo/FH3txzdUimKiDRv3ccFmEOqjndUG6hP+7/EsI43oCPYOvZR+u5GdOkhYrDGZlvjXeJ1CpQxTR/EX+Vt7F8YG+i2LkO7lhPLb+LzgPAxVPCcEMHruuUlE1BYxxzRMOW4X4kjHvJjZGISxa9lgTY3e0mnoQNQVBHKfzI2vGLwvcrFcCIrVxeEbj2dryfByyhZlrNPFbXyf7P4OSfk+fVh6Is1IF1wksfLY/6gWvcmXB8JwmKFDa9s5NfzXnzP3VMrNUWXN3G8Eee6qzKKTDsJ70OrgAx9j9a+dMLfe1vP5t6GQj5"
)

func VyprvpnRegionChoices() (choices []string) {
	servers := VyprvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

//nolint:lll
func VyprvpnServers() []models.VyprvpnServer {
	return []models.VyprvpnServer{
		{Region: "Algeria", Hostname: "dz1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 75, 20}}},
		{Region: "Argentina", Hostname: "ar1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 19}}},
		{Region: "Australia Melbourne", Hostname: "au2.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 117, 19}}},
		{Region: "Australia Perth", Hostname: "au3.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 19}}},
		{Region: "Australia Sydney", Hostname: "au1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 117, 18}}},
		{Region: "Austria", Hostname: "at1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 18}}},
		{Region: "Bahrain", Hostname: "bh1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 115, 19}}},
		{Region: "Belgium", Hostname: "be1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 20}}},
		{Region: "Brazil", Hostname: "br1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 20}}},
		{Region: "Bulgaria", Hostname: "bg1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 22}}},
		{Region: "Canada", Hostname: "ca1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 21, 18}}},
		{Region: "Columbia", Hostname: "co1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 21}}},
		{Region: "Costa Rica", Hostname: "cr1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 22}}},
		{Region: "Czech Republic", Hostname: "cz1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 24}}},
		{Region: "Denmark", Hostname: "dk1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 28}}},
		{Region: "Dubai", Hostname: "ae1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 45, 104}}},
		{Region: "Egypt", Hostname: "eg1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 75, 21}}},
		{Region: "El Salvador", Hostname: "sv1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 61, 20}}},
		{Region: "Finland", Hostname: "fi1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 32}}},
		{Region: "France", Hostname: "fr1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 34}}},
		{Region: "Germany", Hostname: "de1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 26}}},
		{Region: "Greece", Hostname: "gr1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 75, 22}}},
		{Region: "Hong Kong", Hostname: "hk1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 227, 18}}},
		{Region: "Iceland", Hostname: "is1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 22, 20}}},
		{Region: "India", Hostname: "in1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 115, 20}}},
		{Region: "Indonesia", Hostname: "id1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 20}}},
		{Region: "Ireland", Hostname: "ie1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 22, 19}}},
		{Region: "Israel", Hostname: "il1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 75, 18}}},
		{Region: "Italy", Hostname: "it1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 36}}},
		{Region: "Japan", Hostname: "jp1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 113, 18}}},
		{Region: "Latvia", Hostname: "lv1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 44}}},
		{Region: "Liechtenstein", Hostname: "li1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 38}}},
		{Region: "Lithuania", Hostname: "lt1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 40}}},
		{Region: "Luxembourg", Hostname: "lu1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 42}}},
		{Region: "Macao", Hostname: "mo1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 227, 36}}},
		{Region: "Malaysia", Hostname: "my1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 21}}},
		{Region: "Maldives", Hostname: "mv1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 26}}},
		{Region: "Marshall Islands", Hostname: "mh1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 25}}},
		{Region: "Mexico", Hostname: "mx1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 61, 19}}},
		{Region: "Netherlands", Hostname: "eu1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 16}}},
		{Region: "New Zealand", Hostname: "nz1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 117, 20}}},
		{Region: "Norway", Hostname: "no1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 46}}},
		{Region: "Pakistan", Hostname: "pk1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 75, 23}}},
		{Region: "Panama", Hostname: "pa1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 23}}},
		{Region: "Philippines", Hostname: "ph1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 22}}},
		{Region: "Poland", Hostname: "pl1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 48}}},
		{Region: "Portugal", Hostname: "pt1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 50}}},
		{Region: "Qatar", Hostname: "qa1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 115, 21}}},
		{Region: "Romania", Hostname: "ro1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 52}}},
		{Region: "Russia", Hostname: "ru1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 54}}},
		{Region: "Saudi Arabia", Hostname: "sa1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 115, 22}}},
		{Region: "Singapore", Hostname: "sg1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 18}}},
		{Region: "Slovakia", Hostname: "sk1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 60}}},
		{Region: "Slovenia", Hostname: "si1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 58}}},
		{Region: "South Korea", Hostname: "kr1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 113, 19}}},
		{Region: "Spain", Hostname: "es1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 30}}},
		{Region: "Sweden", Hostname: "se1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 56}}},
		{Region: "Switzerland", Hostname: "ch1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 60, 18}}},
		{Region: "Taiwan", Hostname: "tw1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 227, 27}}},
		{Region: "Thailand", Hostname: "th1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 23}}},
		{Region: "Turkey", Hostname: "tr1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 62}}},
		{Region: "USA Austin", Hostname: "us3.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 61, 18}}},
		{Region: "USA Chicago", Hostname: "us6.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 93, 18}}},
		{Region: "USA Los Angeles", Hostname: "us1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 67, 18}}},
		{Region: "USA Miami", Hostname: "us4.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 109, 18}}},
		{Region: "USA New York", Hostname: "us5.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 63, 18}}},
		{Region: "USA San Francisco", Hostname: "us7.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 95, 18}}},
		{Region: "USA Seattle", Hostname: "us8.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 94, 18}}},
		{Region: "USA Washington", Hostname: "us2.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 62, 18}}},
		{Region: "Ukraine", Hostname: "ua1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{128, 90, 96, 64}}},
		{Region: "United Kingdom", Hostname: "uk1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 22, 18}}},
		{Region: "Uruguay", Hostname: "uy1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 61, 21}}},
		{Region: "Vietnam", Hostname: "vn1.vyprvpn.com", TCP: false, UDP: true, IPs: []net.IP{{209, 99, 1, 24}}},
	}
}
