package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	TorguardCertificate        = "MIIDMTCCAhmgAwIBAgIJAKnGGJK6qLqSMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCVRHLVZQTi1DQTAgFw0xOTA1MjExNDIzMTFaGA8yMDU5MDUxMTE0MjMxMVowFDESMBAGA1UEAwwJVEctVlBOLUNBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlv0UgPD3xVAvhhP6q1HCmeAWbH+9HPkyQ2P6qM5oHY5dntjmq8YT48FZGHWv7+s9O47v6Bv7rEc4UwQx15cc2LByivX2JwmE8JACvNfwEnZXYAPq9WU3ZgRrAGvA09ItuLqK2fQ4A7h8bFhmyxCbSzP1sSIT/zJY6ebuh5rDQSMJRMaoI0t1zorEZ7PlEmh+o0w5GPs0D0vY50UcnEzB4GOdWC9pJREwEqppWYLN7RRdG8JyIqmA59mhARCnQFUo38HWic4trxFe71jtD7YInNV7ShQtg0S0sXo36Rqfz72Jo08qqI70dNs5DN1aGNkQ/tRK9DhL5DLmTkaCw7mEFQIDAQABo4GDMIGAMB0GA1UdDgQWBBR7DcymXBp6u/jAaZOPUjUhEyhXfjBEBgNVHSMEPTA7gBR7DcymXBp6u/jAaZOPUjUhEyhXfqEYpBYwFDESMBAGA1UEAwwJVEctVlBOLUNBggkAqcYYkrqoupIwDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAQYwDQYJKoZIhvcNAQELBQADggEBAE79ngbdSlP7IBbfnJ+2Ju7vqt9/GyhcsYtjibp6gsMUxKlD8HuvlSGj5kNO5wiwN7XXqsjYtJfdhmzzVbXksi8Fnbnfa8GhFl4IAjLJ5cxaWOxjr6wx2AhIs+BVVARjaU7iTK91RXJnl6u7UDHTkQylBTl7wgpMeG6GjhaHfcOL1t7D2w8x23cTO+p+n53P3cBq+9TiAUORdzXJvbCxlPMDSDArsgBjC57W7dtdnZo7gTfQG77JTDFBeSwPwLF7PjBB4S6rzU/4fcYwy83XKP6zDn9tgUJDnpFb/7jJ/PbNkK4BWYJp3XytOtt66v9SEKw+v/fJ+VkjU16vE/9Q3h4="
	TorguardOpenvpnStaticKeyV1 = "770e8de5fc56e0248cc7b5aab56be80d0e19cbf003c1b3ed68efbaf08613c3a1a019dac6a4b84f13a6198f73229ffc21fa512394e288f82aa2cf0180f01fb3eb1a71e00a077a20f6d7a83633f5b4f47f27e30617eaf8485dd8c722a8606d56b3c183f65da5d3c9001a8cbdb96c793d936251098b24fe52a6dd2472e98cfccbc466e63520d63ade7a0eacc36208c3142a1068236a52142fbb7b3ed83d785e12a28261bccfb3bcb62a8d2f6d18f5df5f3652e59c5627d8d9c8f7877c4d7b08e19a5c363556ba68d392be78b75152dd55ba0f74d45089e84f77f4492d886524ea6c82b9f4dd83d46528d4f5c3b51cfeaf2838d938bd0597c426b0e440434f2c451f"
)

func TorguardCountryChoices() (choices []string) {
	servers := TorguardServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func TorguardCityChoices() (choices []string) {
	servers := TorguardServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func TorguardHostnameChoices() (choices []string) {
	servers := TorguardServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

//nolint:lll
// TorguardServers returns a slice of all the server information for Torguard.
func TorguardServers() []models.TorguardServer {
	return []models.TorguardServer{
		{Country: "Australia", City: "Sydney", Hostname: "au.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{193, 56, 253, 50}, {217, 138, 205, 106}, {93, 115, 35, 106}, {93, 115, 35, 122}, {217, 138, 205, 218}, {93, 115, 35, 98}, {217, 138, 205, 98}, {193, 56, 253, 146}, {217, 138, 205, 114}, {93, 115, 35, 130}, {217, 138, 205, 194}, {93, 115, 35, 114}, {93, 115, 35, 154}, {193, 56, 253, 98}, {93, 115, 35, 138}, {193, 56, 253, 82}, {193, 56, 253, 66}, {193, 56, 253, 18}, {193, 56, 253, 34}, {93, 115, 35, 146}, {217, 138, 205, 202}, {193, 56, 253, 114}, {217, 138, 205, 210}, {193, 56, 253, 130}}},
		{Country: "Austria", City: "", Hostname: "aus.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{37, 120, 155, 26}, {37, 120, 155, 10}, {37, 120, 155, 18}, {37, 120, 155, 2}, {37, 120, 155, 34}}},
		{Country: "Belarus", City: "", Hostname: "bl.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{95, 47, 99, 12}}},
		{Country: "Belgium", City: "", Hostname: "bg.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 232, 21, 210}, {89, 249, 73, 162}, {89, 249, 73, 130}, {89, 249, 73, 170}, {185, 104, 186, 2}, {185, 232, 21, 34}, {185, 232, 21, 242}, {194, 187, 251, 34}, {185, 232, 21, 250}, {89, 249, 73, 250}, {185, 232, 21, 42}}},
		{Country: "Brazil", City: "", Hostname: "br.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 133, 180, 162}, {45, 133, 180, 154}, {45, 133, 180, 138}, {45, 133, 180, 130}, {45, 133, 180, 146}}},
		{Country: "Bulgaria", City: "", Hostname: "bul.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{82, 102, 23, 170}, {82, 102, 23, 202}, {82, 102, 23, 194}, {82, 102, 23, 210}, {82, 102, 23, 186}, {82, 102, 23, 178}, {82, 102, 23, 218}}},
		{Country: "Canada", City: "Toronto", Hostname: "ca.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{184, 75, 209, 250}}},
		{Country: "Canada", City: "Vancouver", Hostname: "vanc.ca.west.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{107, 181, 189, 41}, {107, 181, 189, 37}, {107, 181, 189, 43}, {107, 181, 189, 36}, {107, 181, 189, 38}, {107, 181, 189, 34}, {107, 181, 189, 42}, {107, 181, 189, 45}, {107, 181, 189, 35}, {107, 181, 189, 40}, {107, 181, 189, 48}, {107, 181, 189, 44}, {107, 181, 189, 46}, {107, 181, 189, 47}, {107, 181, 189, 39}}},
		{Country: "Chile", City: "", Hostname: "chil.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{37, 235, 52, 71}, {193, 235, 146, 104}, {37, 235, 52, 42}, {37, 235, 52, 19}}},
		{Country: "Cyprus", City: "", Hostname: "cp.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 173, 226, 49}}},
		{Country: "Czech", City: "", Hostname: "czech.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 189, 115, 98}, {185, 189, 115, 103}, {185, 189, 115, 108}, {185, 189, 115, 113}, {185, 189, 115, 118}}},
		{Country: "Denmark", City: "", Hostname: "den.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 12, 221, 26}, {45, 12, 221, 2}, {2, 58, 46, 146}, {185, 245, 84, 74}, {2, 58, 46, 138}, {2, 58, 46, 170}, {45, 12, 221, 42}, {45, 12, 221, 34}, {45, 12, 221, 18}, {2, 58, 46, 178}, {2, 58, 46, 162}, {2, 58, 46, 154}, {2, 58, 46, 186}, {45, 12, 221, 10}}},
		{Country: "Finland", City: "", Hostname: "fin.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{91, 132, 197, 192}, {91, 132, 197, 188}, {91, 132, 197, 186}}},
		{Country: "France", City: "", Hostname: "fr.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{93, 177, 75, 58}, {93, 177, 75, 66}, {93, 177, 75, 122}, {93, 177, 75, 90}, {93, 177, 75, 106}, {93, 177, 75, 82}, {93, 177, 75, 2}, {93, 177, 75, 138}, {93, 177, 75, 146}, {93, 177, 75, 26}, {93, 177, 75, 18}, {93, 177, 75, 34}, {93, 177, 75, 10}, {93, 177, 75, 210}, {37, 120, 158, 138}, {93, 177, 75, 162}, {93, 177, 75, 114}, {93, 177, 75, 202}, {93, 177, 75, 130}, {93, 177, 75, 98}, {93, 177, 75, 154}, {93, 177, 75, 42}, {93, 177, 75, 74}, {93, 177, 75, 50}, {93, 177, 75, 218}}},
		{Country: "Germany", City: "", Hostname: "gr.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{93, 177, 73, 90}}},
		{Country: "Greece", City: "", Hostname: "gre.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 92, 33, 10}, {45, 92, 33, 18}, {45, 92, 33, 2}}},
		{Country: "Hong", City: "Kong", Hostname: "hk.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 133, 181, 158}}},
		{Country: "Hungary", City: "", Hostname: "hg.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{37, 120, 144, 98}, {37, 120, 144, 106}}},
		{Country: "Iceland", City: "", Hostname: "ice.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 133, 192, 250}, {45, 133, 192, 218}, {45, 133, 192, 214}, {45, 133, 192, 226}, {45, 133, 192, 230}, {45, 133, 192, 210}, {45, 133, 192, 234}, {45, 133, 192, 254}}},
		{Country: "India", City: "Bangalore", Hostname: "in.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{43, 251, 180, 6}, {103, 78, 121, 182}, {103, 78, 121, 142}, {43, 251, 180, 26}, {43, 251, 180, 10}, {172, 107, 172, 6}, {43, 251, 180, 2}}},
		{Country: "Ireland", City: "", Hostname: "ire.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{77, 81, 139, 90}, {77, 81, 139, 66}, {77, 81, 139, 82}, {77, 81, 139, 74}, {77, 81, 139, 58}}},
		{Country: "Israel", City: "", Hostname: "isr.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{91, 223, 106, 201}}},
		{Country: "Italy", City: "", Hostname: "it.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 9, 251, 182}, {192, 145, 127, 190}, {185, 128, 27, 106}, {45, 9, 251, 14}, {45, 9, 251, 178}}},
		{Country: "Japan", City: "", Hostname: "jp.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{91, 207, 174, 50}}},
		{Country: "Latvia", City: "", Hostname: "lv.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{109, 248, 149, 167}}},
		{Country: "Luxembourg", City: "", Hostname: "lux.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{5, 253, 204, 74}, {5, 253, 204, 90}, {5, 253, 204, 66}, {5, 253, 204, 58}, {5, 253, 204, 82}}},
		{Country: "Mexico", City: "", Hostname: "mx.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{45, 133, 180, 26}, {45, 133, 180, 34}, {45, 133, 180, 2}, {45, 133, 180, 18}, {45, 133, 180, 10}}},
		{Country: "Moldova", City: "", Hostname: "md.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{178, 175, 131, 114}, {178, 175, 131, 106}}},
		{Country: "Netherlands", City: "", Hostname: "nl.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{88, 202, 177, 181}}},
		{Country: "New", City: "Zealand", Hostname: "nz.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{103, 108, 94, 58}}},
		{Country: "Norway", City: "", Hostname: "no.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 181, 61, 38}, {185, 125, 169, 32}, {185, 181, 61, 40}, {185, 125, 169, 28}, {185, 181, 61, 37}, {185, 181, 61, 39}, {185, 125, 169, 30}, {185, 125, 169, 31}, {185, 125, 169, 29}, {185, 125, 169, 24}, {185, 125, 169, 26}, {185, 125, 169, 27}, {185, 125, 169, 23}, {185, 125, 169, 25}, {185, 125, 168, 247}, {185, 125, 168, 248}, {185, 181, 61, 36}}},
		{Country: "Poland", City: "", Hostname: "pl.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{37, 120, 156, 194}}},
		{Country: "Portugal", City: "", Hostname: "por.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{94, 46, 179, 75}}},
		{Country: "Romania", City: "", Hostname: "ro.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 45, 14, 250}, {194, 59, 248, 202}, {31, 14, 252, 146}, {31, 14, 252, 18}, {89, 46, 103, 106}, {194, 59, 248, 210}, {89, 46, 103, 2}, {93, 120, 27, 162}, {31, 14, 252, 178}, {185, 45, 14, 122}, {31, 14, 252, 90}, {185, 45, 15, 106}, {89, 40, 71, 106}}},
		{Country: "Singapore", City: "", Hostname: "singp.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{91, 245, 253, 134}, {185, 200, 117, 142}, {91, 245, 253, 138}, {92, 119, 178, 22}, {82, 102, 25, 2}, {185, 200, 116, 250}, {185, 200, 117, 138}, {82, 102, 25, 226}, {92, 119, 178, 26}, {185, 200, 117, 186}}},
		{Country: "Slovakia", City: "", Hostname: "slk.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{193, 37, 255, 162}, {193, 37, 255, 130}, {193, 37, 255, 122}, {193, 37, 255, 146}, {193, 37, 255, 138}}},
		{Country: "South", City: "Korea", Hostname: "sk.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{169, 56, 83, 216}}},
		{Country: "Spain", City: "", Hostname: "sp.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{192, 145, 124, 242}}},
		{Country: "Sweden", City: "", Hostname: "swe.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{37, 120, 153, 82}, {37, 120, 153, 72}, {37, 120, 153, 32}, {37, 120, 153, 97}, {37, 120, 153, 102}, {37, 120, 153, 92}, {37, 120, 153, 77}, {37, 120, 153, 42}, {37, 120, 153, 52}, {37, 120, 153, 47}, {37, 120, 153, 22}, {37, 120, 153, 57}, {37, 120, 153, 107}, {37, 120, 153, 67}, {37, 120, 153, 7}, {37, 120, 153, 37}, {37, 120, 153, 2}, {37, 120, 153, 12}, {37, 120, 153, 87}, {37, 120, 153, 17}, {37, 120, 153, 62}, {37, 120, 153, 27}}},
		{Country: "Switzerland", City: "", Hostname: "swiss.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{195, 206, 105, 17}, {195, 206, 105, 22}, {195, 206, 105, 32}, {195, 206, 105, 52}, {195, 206, 105, 7}, {185, 236, 201, 210}, {195, 206, 105, 42}, {217, 138, 203, 242}, {195, 206, 105, 12}, {195, 206, 105, 37}, {185, 236, 201, 215}, {195, 206, 105, 27}, {195, 206, 105, 57}, {195, 206, 105, 47}, {195, 206, 105, 2}, {217, 138, 203, 154}}},
		{Country: "Taiwan", City: "", Hostname: "tw.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{61, 216, 159, 176}}},
		{Country: "Thailand", City: "", Hostname: "thai.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{202, 129, 16, 106}, {202, 129, 16, 104}, {202, 129, 16, 42}, {202, 129, 16, 140}, {202, 129, 16, 143}}},
		{Country: "UAE", City: "", Hostname: "uae.secureconnect.me", TCP: true, UDP: true, IPs: []net.IP{{45, 9, 249, 158}, {45, 9, 249, 238}, {45, 9, 250, 10}}},
		{Country: "UK", City: "London", Hostname: "uk.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{109, 123, 118, 13}}},
		{Country: "USA", City: "Atlanta", Hostname: "atl.east.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{104, 223, 95, 50}}},
		{Country: "USA", City: "Chicago", Hostname: "chi.central.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{167, 160, 172, 106}}},
		{Country: "USA", City: "Dallas", Hostname: "dal.central.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{96, 44, 145, 26}}},
		{Country: "USA", City: "Las Vegas", Hostname: "lv.west.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{185, 242, 5, 82}, {45, 89, 173, 106}, {37, 120, 147, 250}, {185, 242, 5, 74}, {185, 242, 5, 178}, {45, 89, 173, 114}, {45, 89, 173, 98}, {185, 242, 5, 66}, {185, 242, 5, 186}, {37, 120, 147, 234}, {185, 242, 5, 170}, {37, 120, 147, 242}, {185, 242, 5, 90}, {45, 89, 173, 162}, {185, 242, 5, 162}, {45, 89, 173, 122}}},
		{Country: "USA", City: "Los Angeles", Hostname: "la.west.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{67, 215, 236, 58}}},
		{Country: "USA", City: "Miami", Hostname: "fl.east.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{96, 47, 226, 42}}},
		{Country: "USA", City: "New Jersey", Hostname: "nj.east.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{23, 226, 128, 146}}},
		{Country: "USA", City: "New York", Hostname: "ny.east.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{67, 213, 221, 12}, {67, 213, 221, 26}, {67, 213, 221, 6}, {67, 213, 221, 7}, {67, 213, 221, 15}, {67, 213, 221, 5}, {67, 213, 221, 22}, {67, 213, 221, 24}, {67, 213, 221, 13}, {67, 213, 221, 23}, {67, 213, 221, 9}, {67, 213, 221, 16}, {67, 213, 221, 19}, {67, 213, 221, 25}, {67, 213, 221, 3}, {67, 213, 221, 10}, {67, 213, 221, 20}, {67, 213, 221, 18}, {67, 213, 221, 17}, {67, 213, 221, 21}, {67, 213, 221, 14}, {67, 213, 221, 8}, {67, 213, 221, 4}, {67, 213, 221, 11}, {67, 213, 221, 27}}},
		{Country: "USA", City: "San Francisco", Hostname: "sf.west.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{206, 189, 218, 238}, {206, 189, 64, 126}, {206, 189, 218, 114}, {206, 189, 169, 41}, {206, 189, 64, 78}, {167, 99, 109, 166}, {206, 189, 214, 52}, {206, 189, 208, 52}, {206, 189, 218, 112}, {206, 189, 208, 113}, {206, 189, 214, 46}}},
		{Country: "USA", City: "Seattle", Hostname: "sa.west.usa.torguardvpnaccess.com", TCP: true, UDP: true, IPs: []net.IP{{199, 229, 250, 38}}},
	}
}
