package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	FastestvpnCertificate        = "MIIFQjCCAyqgAwIBAgIIUfxepT+rr8owDQYJKoZIhvcNAQEMBQAwPzELMAkGA1UEBhMCS1kxEzARBgNVBAoTCkZhc3Rlc3RWUE4xGzAZBgNVBAMTEkZhc3Rlc3RWUE4gUm9vdCBDQTAeFw0xNzA5MTYwMDAxNDZaFw0yNzA5MTQwMDAxNDZaMD8xCzAJBgNVBAYTAktZMRMwEQYDVQQKEwpGYXN0ZXN0VlBOMRswGQYDVQQDExJGYXN0ZXN0VlBOIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC1Xj+WfPTozFynFqc+c3CVrggIllaXEl5bY5VgFynXkqCTM6lSrfC4pNjGXUbqWe6RnGJbM4/6kUn+lQDjFSQV1rzP2eDS8+r5+X2WXh4AoeNRUWhvSG+HiHD/B2EFK+Nd5BRSdUjpKWAtsCmT2bBt7nT0jN1OdeNrLJeyF8siAqv/oQzKznF9aIe/N01b2M8ZOFTzoXi2fZAckgGWui8NB/lzkVIJqSkAPRL8qiJLuRCPVOX1PFD8vV//R8/QumtfbcYBMo6vCk2HmWdrh5OQHPxb3KJtbtG+Z1j8x6HGEAe17djYepBiRMyCEQvYgfD6tvFylc4IquhqE9yaP60PJod5TxpWnRQ6HIGSeBm+S+rYSMalTZ8+pUqOOA+IQCYpfpx6EKIJL/VsW2C7cXdvudxDhXPI5lR/QidCb9Ohq3WkfxXaYwzrngdg2avmNqId9R4KESuM9GoHW0dszfyBCh5wYfeaffMElfDam3B92NUwyhZwtIiv623WVXY9PPz+EDjSJsIAu2Vi1vdJyA4nD4k9Lwmx/1zTc/UaYVLsiBqL2WdfvFTeoWoV+dNxQXSEPhB8gwi8x4O4lZW0cwVy/6fa8KMY8gZbcbSTr7U5bRERfW8l+jY+mYKQ/M/ccgpxaHiw1/+4LWfbJQ7VhJJrTyN0C36FQzY1URkSXg+53wIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBBjAdBgNVHQ4EFgQUmVEL4x6xdCqiqu2OBLs27EA8xGYwDQYJKoZIhvcNAQEMBQADggIBABCpITvO1+R4T9v2+onHiFxU5JjtCZ0zkXqRCMp/Z0UIYbeo1p07pZCPAUjBfGPCkAaR++OiG9sysALdJf8Y6HQKcyuAcWUqQnaIhoZ2JcAP7EKq7uCqsMhcYZD/j3O/3RPtSW5UOx6ItDU+Ua0t9Edho9whNw0VQXmo1JjYoP3FzPjuKoDWTSO1q5eYlZfwcTcs55O2shNkFafPg/6cCm5j6v9nyHrM3sk4LjkrBPUXVx2m/aoz219t8O9Ha9/CdMKXsPO/8gTUzpgnzSgPnGnBmi5xr1nspVN8X4E2f3D+DKqBim3YgslD68NcuFQvJ0/BxZzWVbrr+QXoyzaiCgXuogpIDc2bB6oRXqFnHNz36d4QJmJdWdSaijiS/peQ6EOPgOZ1GuObLWlDCBZLNeQ+N6QaiJxVO4XUj/s22i1IRtwdz84TRHrbWiIpEymsqmb/Ep5r4xV5d6+791axclfOTH7tQrY/SbPtTJI4OEgNekI8YfadQifpelF82MsFFEZuaQn0lj+fvLGtE/zKh3OdLTxRc5TAgBB+0T81+JQosygNr2aFFG0hxar1eyw/gLeG8H+7Ie50pyPvXO4OgB6Key8rSExpilQXlvAT1qX0qS3/K1i/9QkSE9ftIPT6vtwLV2sVQzfyanI4IZgWC6ryhvNLsRn0NFnQclor0+aq"
	FastestvpnOpenvpnStaticKeyV1 = "697fe793b32cb5091d30f2326d5d124a9412e93d0a44ef7361395d76528fcbfc82c3859dccea70a93cfa8fae409709bff75f844cf5ff0c237f426d0c20969233db0e706edb6bdf195ec3dc11b3f76bc807a77e74662d9a800c8cd1144ebb67b7f0d3f1281d1baf522bfe03b7c3f963b1364fc0769400e413b61ca7b43ab19fac9e0f77e41efd4bda7fd77b1de2d7d7855cbbe3e620cecceac72c21a825b243e651f44d90e290e09c3ad650de8fca99c858bc7caad584bc69b11e5c9fd9381c69c505ec487a65912c672d83ed0113b5a74ddfbd3ab33b3683cec593557520a72c4d6cce46111f56f3396cc3ce7183edce553c68ea0796cf6c4375fad00aaa2a42"
)

func FastestvpnCountriesChoices() (choices []string) {
	servers := FastestvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return choices
}

func FastestvpnHostnameChoices() (choices []string) {
	servers := FastestvpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return choices
}

// FastestvpnServers returns the list of all VPN servers for FastestVPN.
//nolint:lll
func FastestvpnServers() []models.FastestvpnServer {
	return []models.FastestvpnServer{
		{Country: "Australia", Hostname: "au-sd-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{139, 99, 149, 10}}},
		{Country: "Australia", Hostname: "au-sd-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{139, 99, 149, 10}}},
		{Country: "Australia", Hostname: "au2-sd-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{139, 99, 131, 126}}},
		{Country: "Australia", Hostname: "au2-sd-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{139, 99, 131, 126}}},
		{Country: "Austria", Hostname: "at.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{86, 107, 21, 146}}},
		{Country: "Belgium", Hostname: "bel1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{217, 138, 211, 67}}},
		{Country: "Belgium", Hostname: "bel2.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{217, 138, 211, 68}}},
		{Country: "Belgium", Hostname: "bel3.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{217, 138, 211, 69}}},
		{Country: "Brazil", Hostname: "br-jp-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{45, 179, 88, 31}}},
		{Country: "Brazil", Hostname: "br-jp-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{45, 179, 88, 31}}},
		{Country: "Bulgaria", Hostname: "bg.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{37, 46, 114, 46}}},
		{Country: "Canada", Hostname: "canada.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{158, 69, 26, 75}}},
		{Country: "Czechia", Hostname: "cz-pr-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{185, 216, 35, 218}}},
		{Country: "Czechia", Hostname: "cz-pr-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{185, 216, 35, 218}}},
		{Country: "Denmark", Hostname: "dk.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{185, 245, 84, 70}}},
		{Country: "Finland", Hostname: "fi-hs-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{194, 34, 132, 19}}},
		{Country: "Finland", Hostname: "fi-hs-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{194, 34, 132, 19}}},
		{Country: "France", Hostname: "fr-rb-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{37, 59, 172, 213}}},
		{Country: "France", Hostname: "fr-rb-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{37, 59, 172, 213}}},
		{Country: "Germany", Hostname: "de1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{83, 143, 245, 254}}},
		{Country: "Hong.Kong", Hostname: "hk-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{64, 120, 88, 115}}},
		{Country: "Hong.Kong", Hostname: "hk-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{64, 120, 88, 115}}},
		{Country: "India", Hostname: "in50.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{103, 104, 74, 32}}},
		{Country: "India-Stream", Hostname: "in-stream.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{103, 104, 74, 30}}},
		{Country: "Italy", Hostname: "it.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{37, 120, 207, 90}}},
		{Country: "Japan", Hostname: "jp-tk-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{202, 239, 38, 147}}},
		{Country: "Japan", Hostname: "jp-tk-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{202, 239, 38, 147}}},
		{Country: "Luxembourg", Hostname: "lux1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{94, 242, 195, 147}}},
		{Country: "Netherlands", Hostname: "nl.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{213, 5, 64, 22}}},
		{Country: "Netherlands", Hostname: "nl2.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{89, 46, 223, 251}}},
		{Country: "Netherlands", Hostname: "nl3.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{89, 46, 223, 252}}},
		{Country: "Norway", Hostname: "nr-ol-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{185, 90, 61, 20}}},
		{Country: "Norway", Hostname: "nr-ol-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{185, 90, 61, 20}}},
		{Country: "Poland", Hostname: "pl2.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{194, 15, 196, 117}}},
		{Country: "Portugal", Hostname: "pt.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{185, 90, 57, 146}}},
		{Country: "Romania", Hostname: "ro.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{91, 199, 50, 131}}},
		{Country: "Russia", Hostname: "russia.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{95, 213, 193, 52}}},
		{Country: "Serbia", Hostname: "rs.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{37, 46, 115, 246}}},
		{Country: "Singapore", Hostname: "sg-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{209, 58, 174, 195}}},
		{Country: "Singapore", Hostname: "sg-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{209, 58, 174, 195}}},
		{Country: "South.Korea", Hostname: "kr-so-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{103, 249, 31, 36}}},
		{Country: "South.Korea", Hostname: "kr-so-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{103, 249, 31, 36}}},
		{Country: "Spain", Hostname: "es-bl-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{193, 148, 19, 155}}},
		{Country: "Spain", Hostname: "es-bl-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{193, 148, 19, 155}}},
		{Country: "Sweden", Hostname: "se-st-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{128, 127, 104, 200}}},
		{Country: "Sweden", Hostname: "se-st-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{128, 127, 104, 201}}},
		{Country: "Sweden", Hostname: "se2.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{79, 142, 76, 142}}},
		{Country: "Switzerland", Hostname: "ch-zr-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{82, 102, 24, 254}}},
		{Country: "Switzerland", Hostname: "ch-zr-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{82, 102, 24, 254}}},
		{Country: "Turkey", Hostname: "tr-iz-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{185, 123, 102, 57}}},
		{Country: "Turkey", Hostname: "tr.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{185, 123, 102, 57}}},
		{Country: "UAE-Dubai", Hostname: "ue-db-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{45, 9, 249, 110}}},
		{Country: "UAE-Dubai", Hostname: "ue-db-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{45, 9, 249, 110}}},
		{Country: "UK", Hostname: "uk.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{5, 226, 139, 143}}},
		{Country: "UK", Hostname: "uk6.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{5, 226, 139, 148}}},
		{Country: "UK-Stream", Hostname: "uk-stream.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{195, 206, 169, 171}}},
		{Country: "US-Atlanta", Hostname: "us-at-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{23, 82, 10, 205}}},
		{Country: "US-Atlanta", Hostname: "us-at-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{23, 82, 10, 205}}},
		{Country: "US-Charlotte", Hostname: "us-cf-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{192, 154, 253, 6}}},
		{Country: "US-Charlotte", Hostname: "us-cf-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{192, 154, 253, 6}}},
		{Country: "US-Chicago", Hostname: "us-ch1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{174, 34, 154, 209}}},
		{Country: "US-Chicago", Hostname: "us-ch2.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{174, 34, 154, 207}}},
		{Country: "US-Dallas", Hostname: "us-dl-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{74, 63, 219, 202}}},
		{Country: "US-Dallas", Hostname: "us-dl-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{74, 63, 219, 202}}},
		{Country: "US-Denver", Hostname: "us-dv1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{173, 248, 157, 107}}},
		{Country: "US-Los.Angeles", Hostname: "us-la-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{64, 31, 35, 222}}},
		{Country: "US-Los.Angeles", Hostname: "us-la-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{64, 31, 35, 222}}},
		{Country: "US-Miami", Hostname: "us-mi-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{162, 255, 138, 231}}},
		{Country: "US-Miami", Hostname: "us-mi-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{162, 255, 138, 232}}},
		{Country: "US-Netflix", Hostname: "netflix.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{37, 59, 172, 215}}},
		{Country: "US-New.York", Hostname: "us-ny-ovtcp-01.jumptoserver.com", UDP: false, TCP: true, IPs: []net.IP{{38, 132, 102, 107}}},
		{Country: "US-New.York", Hostname: "us-ny-ovudp-01.jumptoserver.com", UDP: true, TCP: false, IPs: []net.IP{{38, 132, 102, 107}}},
		{Country: "US-Phoenix", Hostname: "us-ph1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{23, 83, 184, 71}}},
		{Country: "US-Seattle", Hostname: "us-se1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{23, 82, 33, 99}}},
		{Country: "US-St.Louis", Hostname: "us-st1.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{148, 72, 173, 28}}},
		{Country: "US-St.Louis", Hostname: "us-st3.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{148, 72, 173, 30}}},
		{Country: "US-St.Louis", Hostname: "us-st4.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{148, 72, 173, 31}}},
		{Country: "US-St.Louis", Hostname: "us-st5.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{148, 72, 173, 32}}},
		{Country: "US-Washington", Hostname: "us-wt.jumptoserver.com", UDP: true, TCP: true, IPs: []net.IP{{23, 82, 15, 90}}},
	}
}
