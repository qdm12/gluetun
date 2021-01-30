package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
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

//nolint:lll
func PurevpnServers() []models.PurevpnServer {
	return []models.PurevpnServer{
		{Country: "Australia", Region: "New South Wales", City: "Sydney", IPs: []net.IP{{192, 253, 241, 4}, {43, 245, 161, 84}}},
		{Country: "Australia", Region: "Western Australia", City: "Perth", IPs: []net.IP{{172, 94, 123, 4}}},
		{Country: "Austria", Region: "Lower Austria", City: "Langenzersdorf", IPs: []net.IP{{172, 94, 109, 4}}},
		{Country: "Austria", Region: "Vienna", City: "Vienna", IPs: []net.IP{{217, 64, 127, 251}}},
		{Country: "Belgium", Region: "Flanders", City: "Zaventem", IPs: []net.IP{{172, 111, 223, 4}}},
		{Country: "Bulgaria", Region: "Sofia-Capital", City: "Sofia", IPs: []net.IP{{217, 138, 221, 121}}},
		{Country: "Canada", Region: "Alberta", City: "Calgary", IPs: []net.IP{{172, 94, 34, 4}}},
		{Country: "Canada", Region: "Ontario", City: "Toronto", IPs: []net.IP{{104, 200, 138, 196}}},
		{Country: "France", Region: "Île-de-France", City: "Paris", IPs: []net.IP{{89, 40, 183, 178}}},
		{Country: "Germany", Region: "Hesse", City: "Frankfurt am Main", IPs: []net.IP{{82, 102, 16, 110}}},
		{Country: "Greece", Region: "Central Macedonia", City: "Thessaloníki", IPs: []net.IP{{178, 21, 169, 244}}},
		{Country: "Hong Kong", Region: "Central and Western", City: "Hong Kong", IPs: []net.IP{{103, 109, 103, 60}, {43, 226, 231, 4}}},
		{Country: "Hong Kong", Region: "Kowloon City", City: "Kowloon", IPs: []net.IP{{36, 255, 97, 3}}},
		{Country: "Italy", Region: "Trentino-Alto Adige", City: "Trento", IPs: []net.IP{{172, 111, 173, 3}}},
		{Country: "Japan", Region: "Ōsaka", City: "Osaka", IPs: []net.IP{{172, 94, 56, 4}}},
		{Country: "Malaysia", Region: "Kuala Lumpur", City: "Kuala Lumpur", IPs: []net.IP{{103, 55, 10, 133}}},
		{Country: "Netherlands", Region: "North Holland", City: "Amsterdam", IPs: []net.IP{{79, 142, 64, 51}}},
		{Country: "Norway", Region: "Oslo", City: "Oslo", IPs: []net.IP{{82, 102, 22, 212}}},
		{Country: "Poland", Region: "Mazovia", City: "Warsaw", IPs: []net.IP{{5, 253, 206, 251}}},
		{Country: "Portugal", Region: "Lisbon", City: "Lisbon", IPs: []net.IP{{5, 154, 174, 3}}},
		{Country: "Russian Federation", Region: "Moscow", City: "Moscow", IPs: []net.IP{{206, 123, 139, 4}}},
		{Country: "Singapore", Region: "Singapore", City: "Singapore", IPs: []net.IP{{37, 120, 208, 147}, {129, 227, 107, 242}}},
		{Country: "South Africa", Region: "Gauteng", City: "Johannesburg", IPs: []net.IP{{102, 165, 3, 34}}},
		{Country: "Spain", Region: "Madrid", City: "Madrid", IPs: []net.IP{{217, 138, 218, 210}}},
		{Country: "Sweden", Region: "Stockholm", City: "Kista", IPs: []net.IP{{172, 111, 246, 4}}},
		{Country: "Switzerland", Region: "Zurich", City: "Zürich", IPs: []net.IP{{45, 12, 222, 106}}},
		{Country: "Taiwan", Region: "Taiwan", City: "Taipei", IPs: []net.IP{{128, 1, 155, 178}}},
		{Country: "United Kingdom", Region: "England", City: "Birmingham", IPs: []net.IP{{188, 72, 89, 4}}},
		{Country: "United Kingdom", Region: "England", City: "London", IPs: []net.IP{{193, 9, 113, 70}, {104, 37, 6, 4}, {45, 141, 154, 189}}},
		{Country: "United States", Region: "Arkansas", City: "Hot Springs", IPs: []net.IP{{172, 111, 147, 4}}},
		{Country: "United States", Region: "Florida", City: "Miami", IPs: []net.IP{{86, 106, 87, 178}}},
		{Country: "United States", Region: "Illinois", City: "Lincolnshire", IPs: []net.IP{{141, 101, 149, 4}, {141, 101, 149, 4}, {141, 101, 149, 4}, {141, 101, 149, 4}}},
		{Country: "United States", Region: "Massachusetts", City: "Newton", IPs: []net.IP{{104, 243, 244, 2}}},
		{Country: "United States", Region: "New Mexico", City: "Rio Rancho", IPs: []net.IP{{104, 243, 243, 131}}},
		{Country: "United States", Region: "New York", City: "New York City", IPs: []net.IP{{172, 111, 149, 4}}},
		{Country: "United States", Region: "Texas", City: "Dallas", IPs: []net.IP{{172, 94, 1, 4}, {172, 94, 1, 4}, {172, 94, 1, 4}, {172, 94, 1, 4}, {172, 94, 1, 4}, {172, 94, 1, 4}, {208, 84, 155, 104}}},
		{Country: "United States", Region: "Virginia", City: "Reston", IPs: []net.IP{{5, 254, 77, 27}}},
		{Country: "Vietnam", Region: "Ho Chi Minh", City: "Ho Chi Minh City", IPs: []net.IP{{192, 253, 249, 132}}},
	}
}
