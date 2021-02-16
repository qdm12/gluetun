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
	return choices
}

func TorguardCityChoices() (choices []string) {
	servers := TorguardServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return choices
}

func TorguardHostnamesChoices() (choices []string) {
	servers := TorguardServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return choices
}

//nolint:lll
// TorguardServers returns a slice of all the server information for Torguard.
func TorguardServers() []models.TorguardServer {
	return []models.TorguardServer{
		{Country: "Australia", City: "Sydney", Hostname: "au.torguardvpnaccess.com", IP: net.IP{168, 1, 210, 188}},
		{Country: "Austria", City: "", Hostname: "aus.torguardvpnaccess.com", IP: net.IP{37, 120, 155, 18}},
		{Country: "Belarus", City: "", Hostname: "bl.torguardvpnaccess.com", IP: net.IP{95, 47, 99, 12}},
		{Country: "Belgium", City: "", Hostname: "bg.torguardvpnaccess.com", IP: net.IP{89, 249, 73, 250}},
		{Country: "Brazil", City: "", Hostname: "br.torguardvpnaccess.com", IP: net.IP{169, 57, 191, 3}},
		{Country: "Bulgaria", City: "", Hostname: "bul.torguardvpnaccess.com", IP: net.IP{82, 102, 23, 186}},
		{Country: "Canada", City: "Toronto", Hostname: "ca.torguardvpnaccess.com", IP: net.IP{184, 75, 209, 250}},
		{Country: "Canada", City: "Vancouver", Hostname: "vanc.ca.west.torguardvpnaccess.com", IP: net.IP{107, 181, 189, 48}},
		{Country: "Chile", City: "", Hostname: "chil.torguardvpnaccess.com", IP: net.IP{37, 235, 52, 19}},
		{Country: "Cyprus", City: "", Hostname: "cp.torguardvpnaccess.com", IP: net.IP{185, 173, 226, 49}},
		{Country: "Czech", City: "", Hostname: "czech.torguardvpnaccess.com", IP: net.IP{185, 189, 115, 103}},
		{Country: "Denmark", City: "", Hostname: "den.torguardvpnaccess.com", IP: net.IP{37, 120, 131, 29}},
		{Country: "Finland", City: "", Hostname: "fin.torguardvpnaccess.com", IP: net.IP{91, 233, 116, 228}},
		{Country: "France", City: "", Hostname: "fr.torguardvpnaccess.com", IP: net.IP{93, 177, 75, 90}},
		{Country: "Germany", City: "", Hostname: "gr.torguardvpnaccess.com", IP: net.IP{93, 177, 73, 90}},
		{Country: "Greece", City: "", Hostname: "gre.torguardvpnaccess.com", IP: net.IP{45, 92, 33, 10}},
		{Country: "Hong", City: "Kong", Hostname: "hk.torguardvpnaccess.com", IP: net.IP{45, 133, 181, 158}},
		{Country: "Hungary", City: "", Hostname: "hg.torguardvpnaccess.com", IP: net.IP{37, 120, 144, 106}},
		{Country: "Iceland", City: "", Hostname: "ice.torguardvpnaccess.com", IP: net.IP{82, 221, 111, 11}},
		{Country: "India", City: "Bangalore", Hostname: "in.torguardvpnaccess.com", IP: net.IP{139, 59, 62, 109}},
		{Country: "Ireland", City: "", Hostname: "ire.torguardvpnaccess.com", IP: net.IP{81, 17, 246, 108}},
		{Country: "Israel", City: "", Hostname: "isr.torguardvpnaccess.com", IP: net.IP{91, 223, 106, 201}},
		{Country: "Italy", City: "", Hostname: "it.torguardvpnaccess.com", IP: net.IP{192, 145, 127, 190}},
		{Country: "Japan", City: "", Hostname: "jp.torguardvpnaccess.com", IP: net.IP{91, 207, 174, 50}},
		{Country: "Latvia", City: "", Hostname: "lv.torguardvpnaccess.com", IP: net.IP{109, 248, 149, 167}},
		{Country: "Luxembourg", City: "", Hostname: "lux.torguardvpnaccess.com", IP: net.IP{94, 242, 238, 66}},
		{Country: "Mexico", City: "", Hostname: "mx.torguardvpnaccess.com", IP: net.IP{45, 133, 180, 34}},
		{Country: "Moldova", City: "", Hostname: "md.torguardvpnaccess.com", IP: net.IP{178, 175, 128, 74}},
		{Country: "Netherlands", City: "", Hostname: "nl.torguardvpnaccess.com", IP: net.IP{88, 202, 177, 181}},
		{Country: "New", City: "Zealand", Hostname: "nz.torguardvpnaccess.com", IP: net.IP{103, 108, 94, 58}},
		{Country: "Norway", City: "", Hostname: "no.torguardvpnaccess.com", IP: net.IP{185, 125, 169, 31}},
		{Country: "Poland", City: "", Hostname: "pl.torguardvpnaccess.com", IP: net.IP{37, 120, 156, 194}},
		{Country: "Portugal", City: "", Hostname: "por.torguardvpnaccess.com", IP: net.IP{94, 46, 179, 75}},
		{Country: "Romania", City: "", Hostname: "ro.torguardvpnaccess.com", IP: net.IP{93, 120, 27, 162}},
		{Country: "Singapore", City: "", Hostname: "singp.torguardvpnaccess.com", IP: net.IP{206, 189, 43, 152}},
		{Country: "Slovakia", City: "", Hostname: "slk.torguardvpnaccess.com", IP: net.IP{46, 29, 2, 113}},
		{Country: "South", City: "Korea", Hostname: "sk.torguardvpnaccess.com", IP: net.IP{169, 56, 83, 216}},
		{Country: "Spain", City: "", Hostname: "sp.torguardvpnaccess.com", IP: net.IP{192, 145, 124, 242}},
		{Country: "Sweden", City: "", Hostname: "swe.torguardvpnaccess.com", IP: net.IP{37, 120, 153, 72}},
		{Country: "Switzerland", City: "", Hostname: "swiss.torguardvpnaccess.com", IP: net.IP{195, 206, 105, 37}},
		{Country: "Taiwan", City: "", Hostname: "tw.torguardvpnaccess.com", IP: net.IP{61, 216, 159, 176}},
		{Country: "Thailand", City: "", Hostname: "thai.torguardvpnaccess.com", IP: net.IP{202, 129, 16, 104}},
		{Country: "UAE", City: "", Hostname: "uae.secureconnect.me", IP: net.IP{45, 9, 250, 10}},
		{Country: "UK", City: "London", Hostname: "uk.torguardvpnaccess.com", IP: net.IP{109, 123, 118, 13}},
		{Country: "USA", City: "Atlanta", Hostname: "atl.east.usa.torguardvpnaccess.com", IP: net.IP{104, 223, 95, 50}},
		{Country: "USA", City: "Chicago", Hostname: "chi.central.usa.torguardvpnaccess.com", IP: net.IP{167, 160, 172, 106}},
		{Country: "USA", City: "Dallas", Hostname: "dal.central.usa.torguardvpnaccess.com", IP: net.IP{96, 44, 145, 26}},
		{Country: "USA", City: "Las Vegas", Hostname: "lv.west.usa.torguardvpnaccess.com", IP: net.IP{76, 164, 203, 130}},
		{Country: "USA", City: "Los Angeles", Hostname: "la.west.usa.torguardvpnaccess.com", IP: net.IP{67, 215, 236, 58}},
		{Country: "USA", City: "Miami", Hostname: "fl.east.usa.torguardvpnaccess.com", IP: net.IP{96, 47, 226, 42}},
		{Country: "USA", City: "New Jersey", Hostname: "nj.east.usa.torguardvpnaccess.com", IP: net.IP{23, 226, 128, 146}},
		{Country: "USA", City: "New York", Hostname: "ny.east.usa.torguardvpnaccess.com", IP: net.IP{209, 95, 50, 116}},
		{Country: "USA", City: "San Francisco", Hostname: "sf.west.usa.torguardvpnaccess.com", IP: net.IP{206, 189, 218, 238}},
		{Country: "USA", City: "Seattle", Hostname: "sa.west.usa.torguardvpnaccess.com", IP: net.IP{199, 229, 250, 38}},
	}
}
