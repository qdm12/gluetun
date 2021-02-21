package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	PrivatevpnCertificate        = "MIIErTCCA5WgAwIBAgIJAPp3HmtYGCIOMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYDVQQGEwJTRTELMAkGA1UECBMCQ0ExEjAQBgNVBAcTCVN0b2NraG9sbTETMBEGA1UEChMKUHJpdmF0ZVZQTjEWMBQGA1UEAxMNUHJpdmF0ZVZQTiBDQTETMBEGA1UEKRMKUHJpdmF0ZVZQTjEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBwcml2YXR2cG4uc2UwHhcNMTcwNTI0MjAxNTM3WhcNMjcwNTIyMjAxNTM3WjCBlTELMAkGA1UEBhMCU0UxCzAJBgNVBAgTAkNBMRIwEAYDVQQHEwlTdG9ja2hvbG0xEzARBgNVBAoTClByaXZhdGVWUE4xFjAUBgNVBAMTDVByaXZhdGVWUE4gQ0ExEzARBgNVBCkTClByaXZhdGVWUE4xIzAhBgkqhkiG9w0BCQEWFHN1cHBvcnRAcHJpdmF0dnBuLnNlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwjqTWbKk85WN8nd1TaBgBnBHceQWosp8mMHr4xWMTLagWRcq2Modfy7RPnBo9kyn5j/ZZwL/21gLWJbxidurGyZZdEV9Wb5KQl3DUNxa19kwAbkkEchdES61e99MjmQlWq4vGPXAHjEuDxOZ906AXglCyAvQoXcYW0mNm9yybWllVp1aBrCaZQrNYr7eoFvolqJXdQQ3FFsTBCYa5bHJcKQLBfsiqdJ/BAxhNkQtcmWNSgLy16qoxQpCsxNCxAcYnasuL4rwOP+RazBkJTPXA/2neCJC5rt+sXR9CSfiXdJGwMpYso5m31ZEd7JL2+is0FeAZ6ETrKMnEZMsTpTkdwIDAQABo4H9MIH6MB0GA1UdDgQWBBRCkBlC94zCY6VNncMnK36JxT7bazCBygYDVR0jBIHCMIG/gBRCkBlC94zCY6VNncMnK36JxT7ba6GBm6SBmDCBlTELMAkGA1UEBhMCU0UxCzAJBgNVBAgTAkNBMRIwEAYDVQQHEwlTdG9ja2hvbG0xEzARBgNVBAoTClByaXZhdGVWUE4xFjAUBgNVBAMTDVByaXZhdGVWUE4gQ0ExEzARBgNVBCkTClByaXZhdGVWUE4xIzAhBgkqhkiG9w0BCQEWFHN1cHBvcnRAcHJpdmF0dnBuLnNlggkA+ncea1gYIg4wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAayugvExKDHar7t1zyYn99Vt1NMf46J8x4Dt9TNjBml5mR9nKvWmreMUuuOhLaO8Da466KGdXeDFNLcBYZd/J2iTawE6/3fmrML9H2sa+k/+E4uU5nQ84ZGOwCinCkMalVjM8EZ0/H2RZvLAVUnvPuUz2JfJhmiRkbeE75fVuqpAm9qdE+/7lg3oICYzxa6BJPxT+Imdjy3Q/FWdsXqX6aallhohPAZlMZgZL4eXECnV8rAfzyjOJggkMDZQt3Flc0Y4iDMfzrEhSOWMkNFBFwjK0F/dnhsX+fPX6GGRpUZgZcCt/hWvypqc05/SnrdKM/vV/jV/yZe0NVzY7S8Ur5g=="
	PrivatevpnOpenvpnStaticKeyV1 = "a49082f082ca89d6a6bb4ecc7c047c6d428a1d3c8254a95206d38a61d7fbe65984214cd7d56eacc5a60803bffd677fa7294d4bfe555036339312de2dfb1335bd9d5fd94b04bba3a15fc5192aeb02fb6d8dd2ca831fad7509be5eefa8d1eaa689dc586c831a23b589c512662652ecf1bb3a4a673816aba434a04f6857b8c2f8bb265bfe48a7b8112539729d2f7d9734a720e1035188118c73fef1824d0237d5579ca382d703b4bb252acaedc753b12199f00154d3769efbcf85ef5ad6ee755cbeaa944cb98e7654286df54c793a8443f5363078e3da548ba0beed079df633283cefb256f6a4bcfc4ab2c4affc24955c1864d5458e84a7c210d0d186269e55dcf6"
)

func PrivatevpnCountryChoices() (choices []string) {
	servers := PrivatevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeChoicesUnique(choices)
}

func PrivatevpnCityChoices() (choices []string) {
	servers := PrivatevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeChoicesUnique(choices)
}

func PrivatevpnHostnameChoices() (choices []string) {
	servers := PrivatevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeChoicesUnique(choices)
}

//nolint:lll
// PrivatevpnServers returns a slice of all the server information for Privatevpn.
func PrivatevpnServers() []models.PrivatevpnServer {
	return []models.PrivatevpnServer{
		{Country: "Argentina", City: "Buenos Aires", Hostname: "ar-bue.pvdata.host", IPs: []net.IP{{181, 119, 160, 59}}},
		{Country: "Australia", City: "Melbourne", Hostname: "au-mel.pvdata.host", IPs: []net.IP{{103, 231, 88, 203}}},
		{Country: "Australia", City: "Sydney", Hostname: "au-syd.pvdata.host", IPs: []net.IP{{143, 244, 63, 96}}},
		{Country: "Austria", City: "Wien", Hostname: "at-wie.pvdata.host", IPs: []net.IP{{185, 9, 19, 91}}},
		{Country: "Belgium", City: "Brussels", Hostname: "be-bru.pvdata.host", IPs: []net.IP{{185, 104, 186, 211}}},
		{Country: "Brazil", City: "Sao Paulo", Hostname: "br-sao.pvdata.host", IPs: []net.IP{{45, 162, 230, 59}}},
		{Country: "Bulgaria", City: "Sofia", Hostname: "bg-sof.pvdata.host", IPs: []net.IP{{185, 94, 192, 163}}},
		{Country: "Canada", City: "Montreal", Hostname: "ca-mon.pvdata.host", IPs: []net.IP{{37, 120, 237, 163}, {87, 101, 92, 131}}},
		{Country: "Canada", City: "Toronto", Hostname: "ca-tor.pvdata.host", IPs: []net.IP{{45, 148, 7, 3}, {45, 148, 7, 6}, {45, 148, 7, 8}}},
		{Country: "Canada", City: "Vancouver", Hostname: "ca-van.pvdata.host", IPs: []net.IP{{74, 3, 160, 19}}},
		{Country: "Chile", City: "Santiago", Hostname: "cl-san.pvdata.host", IPs: []net.IP{{216, 241, 14, 227}}},
		{Country: "Costa Rica", City: "San Jose", Hostname: "cr-san.pvdata.host", IPs: []net.IP{{190, 10, 8, 218}}},
		{Country: "Croatia", City: "Zagreb", Hostname: "hr-zag.pvdata.host", IPs: []net.IP{{85, 10, 56, 127}}},
		{Country: "Cyprus", City: "Nicosia", Hostname: "cy-nic.pvdata.host", IPs: []net.IP{{185, 173, 226, 47}}},
		{Country: "Czech Republic", City: "Prague", Hostname: "cz-pra.pvdata.host", IPs: []net.IP{{185, 156, 174, 179}}},
		{Country: "Denmark", City: "Copenhagen", Hostname: "dk-cop.pvdata.host", IPs: []net.IP{{62, 115, 255, 188}, {62, 115, 255, 189}}},
		{Country: "France", City: "Paris", Hostname: "fr-par.pvdata.host", IPs: []net.IP{{80, 239, 199, 102}, {80, 239, 199, 103}, {80, 239, 199, 104}, {80, 239, 199, 105}}},
		{Country: "Germany", City: "Frankfurt", Hostname: "de-fra.pvdata.host", IPs: []net.IP{{193, 180, 119, 130}, {193, 180, 119, 131}}},
		{Country: "Germany", City: "Nuremberg", Hostname: "de-nur.pvdata.host", IPs: []net.IP{{185, 89, 36, 3}}},
		{Country: "Greece", City: "Athens", Hostname: "gr-ath.pvdata.host", IPs: []net.IP{{154, 57, 3, 33}}},
		{Country: "Hong Kong", City: "Hong Kong", Hostname: "hk-hon.pvdata.host", IPs: []net.IP{{84, 17, 37, 58}}},
		{Country: "Hungary", City: "Budapest", Hostname: "hu-bud.pvdata.host", IPs: []net.IP{{185, 104, 187, 67}}},
		{Country: "Iceland", City: "Reykjavik", Hostname: "is-rey.pvdata.host", IPs: []net.IP{{82, 221, 113, 210}}},
		{Country: "Indonesia", City: "Jakarta", Hostname: "id-jak.pvdata.host", IPs: []net.IP{{23, 248, 170, 136}}},
		{Country: "Ireland", City: "Dublin", Hostname: "ie-dub.pvdata.host", IPs: []net.IP{{217, 138, 222, 67}}},
		{Country: "Isle of Man", City: "Ballasalla", Hostname: "im-bal.pvdata.host", IPs: []net.IP{{81, 27, 96, 89}}},
		{Country: "Italy", City: "Milan", Hostname: "it-mil.pvdata.host", IPs: []net.IP{{217, 212, 240, 90}, {217, 212, 240, 91}, {217, 212, 240, 92}, {217, 212, 240, 93}}},
		{Country: "Japan", City: "Tokyo", Hostname: "jp-tok.pvdata.host", IPs: []net.IP{{89, 187, 160, 154}}},
		{Country: "Korea", City: "Seoul", Hostname: "kr-seo.pvdata.host", IPs: []net.IP{{92, 223, 73, 37}}},
		{Country: "Latvia", City: "Riga", Hostname: "lv-rig.pvdata.host", IPs: []net.IP{{80, 233, 134, 165}}},
		{Country: "Lithuania", City: "Siauliai", Hostname: "lt-sia.pvdata.host", IPs: []net.IP{{5, 199, 171, 93}}},
		{Country: "Luxembourg", City: "Steinsel", Hostname: "lu-ste.pvdata.host", IPs: []net.IP{{94, 242, 250, 71}}},
		{Country: "Malaysia", City: "Kuala Lumpur", Hostname: "my-kua.pvdata.host", IPs: []net.IP{{128, 1, 160, 184}}},
		{Country: "Malta", City: "Qormi", Hostname: "mt-qor.pvdata.host", IPs: []net.IP{{130, 185, 255, 25}}},
		{Country: "Mexico", City: "Mexico City", Hostname: "mx-mex.pvdata.host", IPs: []net.IP{{190, 60, 16, 28}}},
		{Country: "Moldova", City: "Chisinau", Hostname: "md-chi.pvdata.host", IPs: []net.IP{{178, 17, 172, 99}}},
		{Country: "Netherlands", City: "Amsterdam", Hostname: "nl-ams.pvdata.host", IPs: []net.IP{{193, 180, 119, 194}, {193, 180, 119, 195}, {193, 180, 119, 196}, {193, 180, 119, 197}}},
		{Country: "New Zealand", City: "Auckland", Hostname: "nz-auc.pvdata.host", IPs: []net.IP{{45, 252, 191, 34}}},
		{Country: "Norway", City: "Oslo", Hostname: "no-osl.pvdata.host", IPs: []net.IP{{91, 205, 186, 26}}},
		{Country: "Panama", City: "Panama City", Hostname: "pa-pan.pvdata.host", IPs: []net.IP{{200, 110, 155, 235}}},
		{Country: "Peru", City: "Lima", Hostname: "pe-lim.pvdata.host", IPs: []net.IP{{170, 0, 81, 107}}},
		{Country: "Philippines", City: "Manila", Hostname: "ph-man.pvdata.host", IPs: []net.IP{{128, 1, 209, 12}}},
		{Country: "Portugal", City: "Lisbon", Hostname: "pt-lis.pvdata.host", IPs: []net.IP{{130, 185, 85, 107}}},
		{Country: "Romania", City: "Bukarest", Hostname: "ro-buk.pvdata.host", IPs: []net.IP{{89, 40, 181, 203}}},
		{Country: "Russian Federation", City: "Krasnoyarsk", Hostname: "ru-kra.pvdata.host", IPs: []net.IP{{92, 223, 87, 11}}},
		{Country: "Russian Federation", City: "Moscow", Hostname: "ru-mos.pvdata.host", IPs: []net.IP{{92, 223, 103, 138}}},
		{Country: "Russian Federation", City: "St Petersburg", Hostname: "ru-pet.pvdata.host", IPs: []net.IP{{95, 213, 148, 99}}},
		{Country: "Serbia", City: "Belgrade", Hostname: "rs-bel.pvdata.host", IPs: []net.IP{{141, 98, 103, 166}}},
		{Country: "Slovakia", City: "Bratislava", Hostname: "sg-sin.pvdata.host", IPs: []net.IP{{143, 244, 33, 81}}},
		{Country: "Spain", City: "Madrid", Hostname: "es-mad.pvdata.host", IPs: []net.IP{{217, 212, 244, 92}, {217, 212, 244, 93}}},
		{Country: "Sweden", City: "Gothenburg", Hostname: "se-got.pvdata.host", IPs: []net.IP{{193, 187, 91, 19}}},
		{Country: "Sweden", City: "Kista", Hostname: "se-kis.pvdata.host", IPs: []net.IP{{193, 187, 88, 216}, {193, 187, 88, 217}, {193, 187, 88, 218}, {193, 187, 88, 219}, {193, 187, 88, 220}, {193, 187, 88, 221}, {193, 187, 88, 222}}},
		{Country: "Sweden", City: "Stockholm", Hostname: "se-sto.pvdata.host", IPs: []net.IP{{193, 180, 119, 2}, {193, 180, 119, 6}, {193, 180, 119, 7}}},
		{Country: "Switzerland", City: "Zurich", Hostname: "ch-zur.pvdata.host", IPs: []net.IP{{217, 212, 245, 92}, {217, 212, 245, 93}}},
		{Country: "Taiwan", City: "Taipei", Hostname: "tw-tai.pvdata.host", IPs: []net.IP{{2, 58, 241, 51}}},
		{Country: "Thailand", City: "Bangkok", Hostname: "th-ban.pvdata.host", IPs: []net.IP{{103, 27, 203, 234}}},
		{Country: "Turkey", City: "Istanbul", Hostname: "tr-ist.pvdata.host", IPs: []net.IP{{92, 38, 180, 28}}},
		{Country: "Ukraine", City: "Kiev", Hostname: "ua-kie.pvdata.host", IPs: []net.IP{{192, 121, 68, 131}}},
		{Country: "Ukraine", City: "Nikolaev", Hostname: "ua-nik.pvdata.host", IPs: []net.IP{{194, 54, 83, 21}}},
		{Country: "United Arab Emirates", City: "Dubai", Hostname: "ae-dub.pvdata.host", IPs: []net.IP{{45, 9, 249, 59}}},
		{Country: "United Kingdom", City: "London", Hostname: "uk-lon.pvdata.host", IPs: []net.IP{{193, 180, 119, 66}, {193, 180, 119, 67}, {193, 180, 119, 68}, {193, 180, 119, 69}, {193, 180, 119, 70}}},
		{Country: "United Kingdom", City: "London", Hostname: "uk-lon2.pvdata.host", IPs: []net.IP{{185, 41, 242, 67}}},
		{Country: "United Kingdom", City: "London", Hostname: "uk-lon7.pvdata.host", IPs: []net.IP{{185, 125, 204, 179}}},
		{Country: "United Kingdom", City: "Manchester", Hostname: "uk-man.pvdata.host", IPs: []net.IP{{185, 206, 227, 181}}},
		{Country: "United States", City: "Buffalo", Hostname: "us-buf.pvdata.host", IPs: []net.IP{{172, 245, 13, 115}, {192, 210, 199, 35}}},
		{Country: "United States", City: "Chicago", Hostname: "us-chi.pvdata.host", IPs: []net.IP{{185, 93, 1, 114}}},
		{Country: "United States", City: "Dallas", Hostname: "us-dal.pvdata.host", IPs: []net.IP{{89, 187, 164, 97}}},
		{Country: "United States", City: "Las Vegas", Hostname: "us-las.pvdata.host", IPs: []net.IP{{82, 102, 30, 19}}},
		{Country: "United States", City: "Los Angeles", Hostname: "us-los.pvdata.host", IPs: []net.IP{{89, 187, 185, 78}, {185, 152, 67, 132}}},
		{Country: "United States", City: "Miami", Hostname: "us-mia.pvdata.host", IPs: []net.IP{{195, 181, 163, 139}}},
		{Country: "United States", City: "New York", Hostname: "us-nyc.pvdata.host", IPs: []net.IP{{45, 130, 86, 3}, {45, 130, 86, 5}, {45, 130, 86, 8}, {45, 130, 86, 10}, {45, 130, 86, 12}}},
		{Country: "United States", City: "Phoenix", Hostname: "us-pho.pvdata.host", IPs: []net.IP{{82, 102, 30, 131}}},
		{Country: "Vietnam", City: "Ho Chi Minh City", Hostname: "vn-hoc.pvdata.host", IPs: []net.IP{{210, 2, 64, 5}}},
	}
}
