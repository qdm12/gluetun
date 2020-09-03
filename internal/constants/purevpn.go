package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

const (
	PurevpnCertificateAuthority = "MIIE6DCCA9CgAwIBAgIJAMjXFoeo5uSlMA0GCSqGSIb3DQEBCwUAMIGoMQswCQYDVQQGEwJISzEQMA4GA1UECBMHQ2VudHJhbDELMAkGA1UEBxMCSEsxGDAWBgNVBAoTD1NlY3VyZS1TZXJ2ZXJDQTELMAkGA1UECxMCSVQxGDAWBgNVBAMTD1NlY3VyZS1TZXJ2ZXJDQTEYMBYGA1UEKRMPU2VjdXJlLVNlcnZlckNBMR8wHQYJKoZIhvcNAQkBFhBtYWlsQGhvc3QuZG9tYWluMB4XDTE2MDExNTE1MzQwOVoXDTI2MDExMjE1MzQwOVowgagxCzAJBgNVBAYTAkhLMRAwDgYDVQQIEwdDZW50cmFsMQswCQYDVQQHEwJISzEYMBYGA1UEChMPU2VjdXJlLVNlcnZlckNBMQswCQYDVQQLEwJJVDEYMBYGA1UEAxMPU2VjdXJlLVNlcnZlckNBMRgwFgYDVQQpEw9TZWN1cmUtU2VydmVyQ0ExHzAdBgkqhkiG9w0BCQEWEG1haWxAaG9zdC5kb21haW4wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDluufhyLlyvXzPUL16kAWAdivl1roQv3QHbuRshyKacf/1Er1JqEbtW3Mx9Fvr/u27qU2W8lQI6DaJhU2BfijPe/KHkib55mvHzIVvoexxya26nk79F2c+d9PnuuMdThWQO3El5a/i2AASnM7T7piIBT2WRZW2i8RbfJaTT7G7LP7OpMKIV1qyBg/cWoO7cIWQW4jmzqrNryIkF0AzStLN1DxvnQZwgXBGv0CwuAkfQuNSLu0PQgPp0PhdukNZFllv5D29IhPr0Z+kwPtrAgPQo+lHlOBHBMUpDT4XChTPeAvMaUSBsqmonAE8UUHEabWrqYN/kWNHCNkYXMkiVmK1AgMBAAGjggERMIIBDTAdBgNVHQ4EFgQU456ijsFrYnzHBShLAPpOUqQ+Z2cwgd0GA1UdIwSB1TCB0oAU456ijsFrYnzHBShLAPpOUqQ+Z2ehga6kgaswgagxCzAJBgNVBAYTAkhLMRAwDgYDVQQIEwdDZW50cmFsMQswCQYDVQQHEwJISzEYMBYGA1UEChMPU2VjdXJlLVNlcnZlckNBMQswCQYDVQQLEwJJVDEYMBYGA1UEAxMPU2VjdXJlLVNlcnZlckNBMRgwFgYDVQQpEw9TZWN1cmUtU2VydmVyQ0ExHzAdBgkqhkiG9w0BCQEWEG1haWxAaG9zdC5kb21haW6CCQDI1xaHqObkpTAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCvga2HMwOtUxWH/inL2qk24KX2pxLg939JNhqoyNrUpbDHag5xPQYXUmUpKrNJZ0z+o/ZnNUPHydTSXE7Z7E45J0GDN5E7g4pakndKnDLSjp03NgGsCGW+cXnz6UBPM5FStFvGdDeModeSUyoS9fjk+mYROvmiy5EiVDP91sKGcPLR7Ym0M7zl2aaqV7bb98HmMoBOxpeZQinof67nKrCsgz/xjktWFgcmPl4/PQSsmqQD0fTtWxGuRX+FzwvF2OCMCAJgp1RqJNlk2g50/kBIoJVPPCfjDFeDU5zGaWGSQ9+z1L6/z7VXdjUiHL0ouOcHwbiS4ZjTr9nMn6WdAHU2"
	PurevpnCertificate          = "MIIEnzCCA4egAwIBAgIBAzANBgkqhkiG9w0BAQsFADCBqDELMAkGA1UEBhMCSEsxEDAOBgNVBAgTB0NlbnRyYWwxCzAJBgNVBAcTAkhLMRgwFgYDVQQKEw9TZWN1cmUtU2VydmVyQ0ExCzAJBgNVBAsTAklUMRgwFgYDVQQDEw9TZWN1cmUtU2VydmVyQ0ExGDAWBgNVBCkTD1NlY3VyZS1TZXJ2ZXJDQTEfMB0GCSqGSIb3DQEJARYQbWFpbEBob3N0LmRvbWFpbjAeFw0xNjAxMTUxNjE1MzhaFw0yNjAxMTIxNjE1MzhaMIGdMQswCQYDVQQGEwJISzEQMA4GA1UECBMHQ2VudHJhbDELMAkGA1UEBxMCSEsxFjAUBgNVBAoTDVNlY3VyZS1DbGllbnQxCzAJBgNVBAsTAklUMRYwFAYDVQQDEw1TZWN1cmUtQ2xpZW50MREwDwYDVQQpEwhjaGFuZ2VtZTEfMB0GCSqGSIb3DQEJARYQbWFpbEBob3N0LmRvbWFpbjCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAxsnyn4v6xxDPnuDaYS0b9M1N8nxgg7OBPBlK+FWRxdTQ8yxt5U5CZGm7riVp7fya2J2iPZIgmHQEv/KbxztsHAVlYSfYYlalrnhEL3bDP2tY+N43AwB1k5BrPq2s1pPLT2XG951drDKG4PUuFHUP1sHzW5oQlfVCmxgIMAP8OYkCAwEAAaOCAV8wggFbMAkGA1UdEwQCMAAwLQYJYIZIAYb4QgENBCAWHkVhc3ktUlNBIEdlbmVyYXRlZCBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQU9MwUnUDbQKKZKjoeieD2OD5NlAEwgd0GA1UdIwSB1TCB0oAU456ijsFrYnzHBShLAPpOUqQ+Z2ehga6kgaswgagxCzAJBgNVBAYTAkhLMRAwDgYDVQQIEwdDZW50cmFsMQswCQYDVQQHEwJISzEYMBYGA1UEChMPU2VjdXJlLVNlcnZlckNBMQswCQYDVQQLEwJJVDEYMBYGA1UEAxMPU2VjdXJlLVNlcnZlckNBMRgwFgYDVQQpEw9TZWN1cmUtU2VydmVyQ0ExHzAdBgkqhkiG9w0BCQEWEG1haWxAaG9zdC5kb21haW6CCQDI1xaHqObkpTATBgNVHSUEDDAKBggrBgEFBQcDAjALBgNVHQ8EBAMCB4AwDQYJKoZIhvcNAQELBQADggEBAFyFo2VUX/UFixsdPdK9/Yt6mkCWc+XS1xbapGXXb9U1d+h1iBCIV9odUHgNCXWpz1hR5Uu/OCzaZ0asLE4IFMZlQmJs8sMT0c1tfPPGW45vxbL0lhqnQ8PNcBH7huNK7VFjUh4szXRKmaQPaM4S91R3L4CaNfVeHfAg7mN2m9Zn5Gto1Q1/CFMGKu2hxwGEw5p+X1czBWEvg/O09ckx/ggkkI1NcZsNiYQ+6Pz8DdGGX3+05YwLZu94+O6iIMrzxl/il0eK83g3YPbsOrASARvw6w/8sOnJCK5eOacl21oww875KisnYdWjHB1FiI+VzQ1/gyoDsL5kPTJVuu2CoG8="
	PurevpnKey                  = "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMbJ8p+L+scQz57g2mEtG/TNTfJ8YIOzgTwZSvhVkcXU0PMsbeVOQmRpu64lae38mtidoj2SIJh0BL/ym8c7bBwFZWEn2GJWpa54RC92wz9rWPjeNwMAdZOQaz6trNaTy09lxvedXawyhuD1LhR1D9bB81uaEJX1QpsYCDAD/DmJAgMBAAECgYEAvTHbDupE5U0krUvHzBEIuHblptGlcfNYHoDcD3oxYR3pOGeiuElBexv+mgHVzcFLBrsQfJUlHLPfCWi3xmjRvDQcr7N7U1u7NIzazy/PpRBaKolMRiM1KMYi2DG0i4ZONwFT8bvNHOIrZzCLY54KDrqOn55OzC70WYjWh4t5evkCQQDkkzZUAeskBC9+JP/zLps8jhwfoLBWGw/zbC9ePDmX0N8MTZdcUpg6KUTf1wbkLUyVtIRjS2ao6qu1jWG6K0x3AkEA3qPWyaWQWCynhNDqu2U1cPb2kh5AJip+gqxO3emikAdajsSxeoyEC2AfyBITbeB1tvCUZH17J4i/0+OFTEQp/wJAb/zEOGJ8PzghwK8GC7JA8mk51DEZVAaMSRovFv9wxDXcoh191AjPdmdzzCuAv9iF1i8MUc3GbWoUWK39PIYsPwJAWh63sqfx5b8tj/WBDpnJKBDPfhYAoXJSA1L8GZeY1fQkE+ZKcPCwAmrGcpXeh3t0Krj3WDXyw+32uC5Apr5wwQJAPZwOOReaC4YNfBPZN9BdHvVjOYGGUffpI+X+hWpLRnQFJteAi+eqwyk0Oi0SkJB+a7jcerK2d7q7xhec5WHlng=="
	PurevpnOpenvpnStaticKeyV1   = "e30af995f56d07426d9ba1f824730521d4283db4b4d0cdda9c6e8759a3799dcb7939b6a5989160c9660de0f6125cbb1f585b41c074b2fe88ecfcf17eab9a33be1352379cdf74952b588fb161a93e13df9135b2b29038231e02d657a6225705e6868ccb0c384ed11614690a1894bfbeb274cebf1fe9c2329bdd5c8a40fe8820624d2ea7540cd79ab76892db51fc371a3ac5fc9573afecb3fffe3281e61d72e91579d9b03d8cbf7909b3aebf4d90850321ee6b7d0a7846d15c27d8290e031e951e19438a4654663cad975e138f5bc5af89c737ad822f27e19057731f41e1e254cc9c95b7175c622422cde9f1f2cfd3510add94498b4d7133d3729dd214a16b27fb"
)

func PurevpnRegionChoices() (choices []string) {
	servers := PurevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return choices
}

func PurevpnCountryChoices() (choices []string) {
	servers := PurevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return choices
}

func PurevpnCityChoices() (choices []string) {
	servers := PurevpnServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return choices
}

func PurevpnServers() []models.PurevpnServer {
	return []models.PurevpnServer{
		{Region: "Africa", Country: "Algeria", City: "Algiers", IPs: []net.IP{{172, 94, 64, 2}}},
		{Region: "Africa", Country: "Angola", City: "Benguela", IPs: []net.IP{{45, 115, 26, 2}}},
		{Region: "Africa", Country: "Cape Verde", City: "Praia", IPs: []net.IP{{45, 74, 25, 2}}},
		{Region: "Africa", Country: "Egypt", City: "Cairo", IPs: []net.IP{{192, 198, 120, 122}}},
		{Region: "Africa", Country: "Ethiopia", City: "Addis Ababa", IPs: []net.IP{{104, 250, 178, 4}}},
		{Region: "Africa", Country: "Ghana", City: "Accra", IPs: []net.IP{{196, 251, 67, 4}}},
		{Region: "Africa", Country: "Kenya", City: "Mombasa", IPs: []net.IP{{102, 135, 0, 2}}},
		{Region: "Africa", Country: "Madagascar", City: "Antananarivo", IPs: []net.IP{{206, 123, 156, 131}}},
		{Region: "Africa", Country: "Mauritania", City: "Nouakchott", IPs: []net.IP{{206, 123, 158, 63}}},
		{Region: "Africa", Country: "Mauritius", City: "Port Louis", IPs: []net.IP{{104, 250, 181, 4}}},
		{Region: "Africa", Country: "Morocco", City: "Rabat", IPs: []net.IP{{104, 243, 250, 126}}},
		{Region: "Africa", Country: "Niger", City: "Niamey", IPs: []net.IP{{206, 123, 157, 131}}},
		{Region: "Africa", Country: "Nigeria", City: "Suleja", IPs: []net.IP{{102, 165, 25, 38}}},
		{Region: "Africa", Country: "Senegal", City: "Dakar", IPs: []net.IP{{206, 123, 158, 131}}},
		{Region: "Africa", Country: "Seychelles", City: "Victoria", IPs: []net.IP{{172, 111, 128, 126}}},
		{Region: "Africa", Country: "South Africa", City: "Johannesburg", IPs: []net.IP{{45, 74, 45, 2}}},
		{Region: "Africa", Country: "Tanzania", City: "Dar Es Salaam", IPs: []net.IP{{102, 135, 0, 2}}},
		{Region: "Africa", Country: "Tunisia", City: "Tunis", IPs: []net.IP{{206, 123, 159, 4}}},
		{Region: "Asia", Country: "Afghanistan", City: "Kabul", IPs: []net.IP{{172, 111, 208, 2}}},
		{Region: "Asia", Country: "Armenia", City: "Singapore", IPs: []net.IP{{37, 120, 208, 147}}},
		{Region: "Asia", Country: "Azerbaijan", City: "Baku", IPs: []net.IP{{104, 250, 177, 4}}},
		{Region: "Asia", Country: "Bangladesh", City: "Dhaka", IPs: []net.IP{{206, 123, 154, 190}}},
		{Region: "Asia", Country: "Brunei Darussalam", City: "Bandar Seri Begawan", IPs: []net.IP{{36, 255, 98, 2}}},
		{Region: "Asia", Country: "Cambodia", City: "Phnom Penh", IPs: []net.IP{{104, 250, 176, 122}}},
		{Region: "Asia", Country: "India", City: "Chennai", IPs: []net.IP{{129, 227, 107, 242}}},
		{Region: "Asia", Country: "Indonesia", City: "Jakarta", IPs: []net.IP{{103, 55, 9, 2}}},
		{Region: "Asia", Country: "Japan", City: "Tokyo", IPs: []net.IP{{172, 94, 56, 2}}},
		{Region: "Asia", Country: "Kazakhstan", City: "Almaty", IPs: []net.IP{{206, 123, 152, 4}}},
		{Region: "Asia", Country: "Korea, South", City: "Seoul", IPs: []net.IP{{45, 115, 25, 1}}},
		{Region: "Asia", Country: "Kyrgyzstan", City: "Bishkek", IPs: []net.IP{{206, 123, 151, 131}}},
		{Region: "Asia", Country: "Laos", City: "Vientiane", IPs: []net.IP{{206, 123, 153, 4}}},
		{Region: "Asia", Country: "Macao", City: "Beyrouth", IPs: []net.IP{{104, 243, 240, 121}}},
		{Region: "Asia", Country: "Malaysia", City: "Johor Baharu", IPs: []net.IP{{103, 28, 90, 54}, {103, 28, 90, 55}, {103, 28, 90, 71}, {103, 28, 90, 72}, {103, 117, 20, 21}, {103, 117, 20, 163}, {103, 117, 20, 164}, {103, 117, 20, 201}}},
		{Region: "Asia", Country: "Malaysia", City: "Kuala Lumpur", IPs: []net.IP{{104, 250, 160, 4}}},
		{Region: "Asia", Country: "Mongolia", City: "Ulaanbaatar", IPs: []net.IP{{206, 123, 153, 131}}},
		{Region: "Asia", Country: "Pakistan", City: "Islamabad", IPs: []net.IP{{104, 250, 187, 3}}},
		{Region: "Asia", Country: "Papua New Guinea", City: "Port Moresby", IPs: []net.IP{{206, 123, 155, 131}}},
		{Region: "Asia", Country: "Philippines", City: "Manila", IPs: []net.IP{{129, 227, 119, 84}}},
		{Region: "Asia", Country: "Sri Lanka", City: "Colombo", IPs: []net.IP{{206, 123, 154, 4}}},
		{Region: "Asia", Country: "Taiwan", City: "Taipei", IPs: []net.IP{{128, 1, 155, 178}}},
		{Region: "Asia", Country: "Tajikistan", City: "Dushanbe", IPs: []net.IP{{206, 123, 151, 4}}},
		{Region: "Asia", Country: "Thailand", City: "Bangkok", IPs: []net.IP{{104, 37, 6, 4}}},
		{Region: "Asia", Country: "Turkey", City: "Istanbul", IPs: []net.IP{{185, 220, 58, 3}}},
		{Region: "Asia", Country: "Turkmenistan", City: "Ashgabat", IPs: []net.IP{{206, 123, 152, 131}}},
		{Region: "Asia", Country: "Uzbekistan", City: "Tashkent", IPs: []net.IP{{206, 123, 150, 131}}},
		{Region: "Asia", Country: "Vietnam", City: "Hanoi", IPs: []net.IP{{192, 253, 249, 132}}},
		{Region: "Europe", Country: "Albania", City: "Tirane", IPs: []net.IP{{46, 243, 224, 2}}},
		{Region: "Europe", Country: "Armenia", City: "Yerevan", IPs: []net.IP{{172, 94, 35, 4}}},
		{Region: "Europe", Country: "Austria", City: "Vienna", IPs: []net.IP{{172, 94, 109, 2}}},
		{Region: "Europe", Country: "Belgium", City: "Brussels", IPs: []net.IP{{185, 210, 217, 147}}},
		{Region: "Europe", Country: "Bosnia and Herzegovina", City: "Sarajevo", IPs: []net.IP{{104, 250, 169, 122}}},
		{Region: "Europe", Country: "Bulgaria", City: "Sofia", IPs: []net.IP{{217, 138, 221, 114}}},
		{Region: "Europe", Country: "Croatia", City: "Zagreb", IPs: []net.IP{{104, 250, 163, 2}}},
		{Region: "Europe", Country: "Cyprus", City: "Nicosia", IPs: []net.IP{{188, 72, 119, 4}}},
		{Region: "Europe", Country: "Denmark", City: "Copenhagen", IPs: []net.IP{{89, 45, 7, 5}}},
		{Region: "Europe", Country: "Estonia", City: "Tallinn", IPs: []net.IP{{185, 166, 87, 2}, {188, 72, 111, 4}}},
		{Region: "Europe", Country: "France", City: "Paris", IPs: []net.IP{{172, 94, 53, 2}, {172, 111, 219, 2}}},
		{Region: "Europe", Country: "Georgia", City: "Tbilisi", IPs: []net.IP{{141, 101, 156, 2}}},
		{Region: "Europe", Country: "Germany", City: "Frankfurt", IPs: []net.IP{{172, 94, 8, 4}}},
		{Region: "Europe", Country: "Germany", City: "Munich", IPs: []net.IP{{172, 94, 8, 4}}},
		{Region: "Europe", Country: "Germany", City: "Nuremberg", IPs: []net.IP{{172, 94, 125, 2}}},
		{Region: "Europe", Country: "Greece", City: "Thessaloniki", IPs: []net.IP{{172, 94, 109, 2}}},
		{Region: "Europe", Country: "Hungary", City: "Budapest", IPs: []net.IP{{172, 111, 129, 2}, {188, 72, 125, 126}}},
		{Region: "Europe", Country: "Iceland", City: "Reykjavik", IPs: []net.IP{{192, 253, 250, 1}}},
		{Region: "Europe", Country: "Ireland", City: "Dublin", IPs: []net.IP{{185, 210, 217, 147}}},
		{Region: "Europe", Country: "Isle of Man", City: "Onchan", IPs: []net.IP{{46, 243, 144, 2}}},
		{Region: "Europe", Country: "Italy", City: "Milano", IPs: []net.IP{{45, 9, 251, 2}}},
		{Region: "Europe", Country: "Latvia", City: "RIGA", IPs: []net.IP{{185, 118, 76, 5}}},
		{Region: "Europe", Country: "Liechtenstein", City: "Vaduz", IPs: []net.IP{{104, 250, 164, 4}}},
		{Region: "Europe", Country: "Lithuania", City: "Vilnius", IPs: []net.IP{{188, 72, 116, 3}}},
		{Region: "Europe", Country: "Luxembourg", City: "Luxembourg", IPs: []net.IP{{188, 72, 114, 2}}},
		{Region: "Europe", Country: "Malta", City: "Sliema", IPs: []net.IP{{46, 243, 241, 4}}},
		{Region: "Europe", Country: "Monaco", City: "Monaco", IPs: []net.IP{{104, 250, 168, 132}}},
		{Region: "Europe", Country: "Montenegro", City: "Podgorica", IPs: []net.IP{{104, 250, 165, 121}}},
		{Region: "Europe", Country: "Netherlands", City: "Amsterdam", IPs: []net.IP{{92, 119, 179, 195}}},
		{Region: "Europe", Country: "Norway", City: "Oslo", IPs: []net.IP{{82, 102, 22, 211}}},
		{Region: "Europe", Country: "Poland", City: "Warsaw", IPs: []net.IP{{5, 253, 206, 251}}},
		{Region: "Europe", Country: "Portugal", City: "Lisbon", IPs: []net.IP{{45, 74, 10, 1}}},
		{Region: "Europe", Country: "Romania", City: "Bucharest", IPs: []net.IP{{192, 253, 253, 2}}},
		{Region: "Europe", Country: "Serbia", City: "Niš", IPs: []net.IP{{104, 250, 166, 2}}},
		{Region: "Europe", Country: "Slovakia", City: "Bratislava", IPs: []net.IP{{188, 72, 112, 3}}},
		{Region: "Europe", Country: "Slovenia", City: "Ljubljana", IPs: []net.IP{{104, 243, 246, 129}}},
		{Region: "Europe", Country: "Spain", City: "Barcelona", IPs: []net.IP{{185, 230, 124, 147}}},
		{Region: "Europe", Country: "Sweden", City: "Stockholm", IPs: []net.IP{{45, 74, 46, 2}}},
		{Region: "Europe", Country: "Switzerland", City: "Zurich", IPs: []net.IP{{172, 111, 217, 2}}},
		{Region: "Europe", Country: "United Kingdom", City: "Gosport", IPs: []net.IP{{45, 74, 0, 2}, {45, 74, 62, 2}}},
		{Region: "Europe", Country: "United Kingdom", City: "London", IPs: []net.IP{{45, 74, 0, 2}, {45, 74, 62, 2}}},
		{Region: "Europe", Country: "United Kingdom", City: "Maidenhead", IPs: []net.IP{{172, 111, 183, 2}}},
		{Region: "Europe", Country: "United Kingdom", City: "Manchester", IPs: []net.IP{{172, 111, 183, 2}}},
		{Region: "Middle East", Country: "Bahrain", City: "Manama", IPs: []net.IP{{46, 243, 150, 4}}},
		{Region: "Middle East", Country: "Jordan", City: "Amman", IPs: []net.IP{{172, 111, 152, 3}}},
		{Region: "Middle East", Country: "Kuwait", City: "Kuwait", IPs: []net.IP{{206, 123, 146, 4}}},
		{Region: "Middle East", Country: "Oman", City: "Salalah", IPs: []net.IP{{46, 243, 148, 125}}},
		{Region: "Middle East", Country: "Qatar", City: "Doha", IPs: []net.IP{{46, 243, 147, 2}}},
		{Region: "Middle East", Country: "Saudi Arabia", City: "Jeddah", IPs: []net.IP{{45, 74, 1, 4}}},
		{Region: "Middle East", Country: "United Arab Emirates", City: "Dubai", IPs: []net.IP{{104, 37, 6, 4}}},
		{Region: "North America", Country: "Aruba", City: "Oranjestad", IPs: []net.IP{{104, 243, 246, 129}}},
		{Region: "North America", Country: "Barbados", City: "Bridgetown", IPs: []net.IP{{172, 94, 97, 2}}},
		{Region: "North America", Country: "Belize", City: "Belmopan", IPs: []net.IP{{104, 243, 241, 4}}},
		{Region: "North America", Country: "Bermuda", City: "Hamilton", IPs: []net.IP{{172, 94, 76, 2}}},
		{Region: "North America", Country: "Canada", City: "Montreal", IPs: []net.IP{{172, 94, 7, 2}}},
		{Region: "North America", Country: "Canada", City: "Toronto", IPs: []net.IP{{172, 94, 7, 2}}},
		{Region: "North America", Country: "Canada", City: "Vancouver", IPs: []net.IP{{107, 181, 177, 42}}},
		{Region: "North America", Country: "Cayman Islands", City: "George Town", IPs: []net.IP{{172, 94, 113, 2}}},
		{Region: "North America", Country: "Costa Rica", City: "San Jose", IPs: []net.IP{{104, 243, 245, 1}}},
		{Region: "North America", Country: "Dominica", City: "Roseau", IPs: []net.IP{{45, 74, 22, 2}}},
		{Region: "North America", Country: "Dominican Republic", City: "Santo Domingo", IPs: []net.IP{{45, 74, 23, 129}}},
		{Region: "North America", Country: "El Salvador", City: "San Salvador", IPs: []net.IP{{45, 74, 17, 129}}},
		{Region: "North America", Country: "Grenada", City: "St George's", IPs: []net.IP{{45, 74, 21, 129}}},
		{Region: "North America", Country: "Guatemala", City: "Guatemala", IPs: []net.IP{{45, 74, 17, 2}}},
		{Region: "North America", Country: "Haiti", City: "PORT-AU-PRINCE", IPs: []net.IP{{45, 74, 24, 2}}},
		{Region: "North America", Country: "Honduras", City: "TEGUCIGALPA", IPs: []net.IP{{45, 74, 18, 2}}},
		{Region: "North America", Country: "Jamaica", City: "Kingston", IPs: []net.IP{{104, 250, 182, 126}}},
		{Region: "North America", Country: "Mexico", City: "Mexico City", IPs: []net.IP{{104, 243, 243, 131}}},
		{Region: "North America", Country: "Montserrat", City: "plymouth", IPs: []net.IP{{45, 74, 26, 190}}},
		{Region: "North America", Country: "Puerto Rico", City: "San Juan", IPs: []net.IP{{104, 37, 2, 2}}},
		{Region: "North America", Country: "Saint Lucia", City: "Castries", IPs: []net.IP{{45, 74, 23, 2}}},
		{Region: "North America", Country: "The Bahamas", City: "Freeport", IPs: []net.IP{{104, 243, 242, 2}}},
		{Region: "North America", Country: "Trinidad and Tobago", City: "Port of Spain", IPs: []net.IP{{45, 74, 21, 2}}},
		{Region: "North America", Country: "Turks and Caicos Islands", City: "Balfour Town", IPs: []net.IP{{45, 74, 24, 129}}},
		{Region: "North America", Country: "United States", City: "Ashburn", IPs: []net.IP{{46, 243, 249, 2}}},
		{Region: "North America", Country: "United States", City: "Chicago", IPs: []net.IP{{46, 243, 249, 4}}},
		{Region: "North America", Country: "United States", City: "Columbus", IPs: []net.IP{{172, 94, 115, 2}}},
		{Region: "North America", Country: "United States", City: "Georgia", IPs: []net.IP{{141, 101, 168, 4}}},
		{Region: "North America", Country: "United States", City: "Houston", IPs: []net.IP{{172, 94, 1, 4}}},
		{Region: "North America", Country: "United States", City: "Los Angeles", IPs: []net.IP{{141, 101, 169, 4}}},
		{Region: "North America", Country: "United States", City: "Miami", IPs: []net.IP{{5, 254, 79, 114}}},
		{Region: "North America", Country: "United States", City: "New Jersey", IPs: []net.IP{{172, 94, 1, 4}}},
		{Region: "North America", Country: "United States", City: "New York", IPs: []net.IP{{172, 94, 1, 4}}},
		{Region: "North America", Country: "United States", City: "Phoenix", IPs: []net.IP{{172, 94, 26, 4}}},
		{Region: "North America", Country: "United States", City: "Salt Lake City", IPs: []net.IP{{141, 101, 168, 4}}},
		{Region: "North America", Country: "United States", City: "San Francisco", IPs: []net.IP{{172, 94, 1, 4}}},
		{Region: "North America", Country: "United States", City: "Seattle", IPs: []net.IP{{172, 94, 86, 2}}},
		{Region: "North America", Country: "United States", City: "Washington, D.C.", IPs: []net.IP{{141, 101, 169, 4}}},
		{Region: "Oceania", Country: "Australia", City: "Brisbane", IPs: []net.IP{{172, 111, 236, 2}}},
		{Region: "Oceania", Country: "Australia", City: "Melbourne", IPs: []net.IP{{118, 127, 62, 2}}},
		{Region: "Oceania", Country: "Australia", City: "Sydney", IPs: []net.IP{{192, 253, 241, 2}}},
		{Region: "Oceania", Country: "New Zealand", City: "Auckland", IPs: []net.IP{{43, 228, 156, 4}}},
		{Region: "South America", Country: "Argentina", City: "Buenos Aires", IPs: []net.IP{{104, 243, 244, 1}}},
		{Region: "South America", Country: "Bolivia", City: "Sucre", IPs: []net.IP{{172, 94, 77, 2}}},
		{Region: "South America", Country: "Brazil", City: "Sao Paulo", IPs: []net.IP{{104, 243, 244, 2}}},
		{Region: "South America", Country: "British Virgin Island", City: "Road Town", IPs: []net.IP{{104, 250, 184, 130}}},
		{Region: "South America", Country: "Chile", City: "Santiago", IPs: []net.IP{{191, 96, 183, 251}}},
		{Region: "South America", Country: "Colombia", City: "Bogota", IPs: []net.IP{{172, 111, 132, 1}}},
		{Region: "South America", Country: "Ecuador", City: "Quito", IPs: []net.IP{{104, 250, 180, 126}}},
		{Region: "South America", Country: "Guyana", City: "Georgetown", IPs: []net.IP{{45, 74, 20, 129}}},
		{Region: "South America", Country: "Panama", City: "Panama City", IPs: []net.IP{{104, 243, 243, 131}}},
		{Region: "South America", Country: "Paraguay", City: "Asuncion", IPs: []net.IP{{45, 74, 19, 129}}},
		{Region: "South America", Country: "Peru", City: "Lima", IPs: []net.IP{{172, 111, 131, 1}}},
		{Region: "South America", Country: "Suriname", City: "Paramaribo", IPs: []net.IP{{45, 74, 20, 4}}},
	}
}
