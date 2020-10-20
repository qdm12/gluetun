package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	PIAEncryptionPresetNormal = "normal"
	PIAEncryptionPresetStrong = "strong"
	PiaX509CRLNormal          = "MIICWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAQEAQZo9X97ci8EcPYu/uK2HB152OZbeZCINmYyluLDOdcSvg6B5jI+ffKN3laDvczsG6CxmY3jNyc79XVpEYUnq4rT3FfveW1+Ralf+Vf38HdpwB8EWB4hZlQ205+21CALLvZvR8HcPxC9KEnev1mU46wkTiov0EKc+EdRxkj5yMgv0V2Reze7AP+NQ9ykvDScH4eYCsmufNpIjBLhpLE2cuZZXBLcPhuRzVoU3l7A9lvzG9mjA5YijHJGHNjlWFqyrn1CfYS6koa4TGEPngBoAziWRbDGdhEgJABHrpoaFYaL61zqyMR6jC0K2ps9qyZAN74LEBedEfK7tBOzWMwr58A=="
	PiaX509CRLStrong          = "MIIDWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAgEAppFfEpGsasjB1QgJcosGpzbf2kfRhM84o2TlqY1ua+Gi5TMdKydA3LJcNTjlI9a0TYAJfeRX5IkpoglSUuHuJgXhP3nEvX10mjXDpcu/YvM8TdE5JV2+EGqZ80kFtBeOq94WcpiVKFTR4fO+VkOK9zwspFfb1cNs9rHvgJ1QMkRUF8PpLN6AkntHY0+6DnigtSaKqldqjKTDTv2OeH3nPoh80SGrt0oCOmYKfWTJGpggMGKvIdvU3vH9+EuILZKKIskt+1dwdfA5Bkz1GLmiQG7+9ZZBQUjBG9Dos4hfX/rwJ3eU8oUIm4WoTz9rb71SOEuUUjP5NPy9HNx2vx+cVvLsTF4ZDZaUztW9o9JmIURDtbeyqxuHN3prlPWB6aj73IIm2dsDQvs3XXwRIxs8NwLbJ6CyEuvEOVCskdM8rdADWx1J0lRNlOJ0Z8ieLLEmYAA834VN1SboB6wJIAPxQU3rcBhXqO9y8aa2oRMg8NxZ5gr+PnKVMqag1x0IxbIgLxtkXQvxXxQHEMSODzvcOfK/nBRBsqTj30P+R87sU8titOoxNeRnBDRNhdEy/QGAqGh62ShPpQUCJdnKRiRTjnil9hMQHevoSuFKeEMO30FQL7BZyo37GFU+q1WPCplVZgCP9hC8Rn5K2+f6KLFo5bhtowSmu+GY1yZtg+RTtsA="
	PIACertificateNormal      = "MIIFqzCCBJOgAwIBAgIJAKZ7D5Yv87qDMA0GCSqGSIb3DQEBDQUAMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTAeFw0xNDA0MTcxNzM1MThaFw0zNDA0MTIxNzM1MThaMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPXDL1L9tX6DGf36liA7UBTy5I869z0UVo3lImfOs/GSiFKPtInlesP65577nd7UNzzXlH/P/CnFPdBWlLp5ze3HRBCc/Avgr5CdMRkEsySL5GHBZsx6w2cayQ2EcRhVTwWpcdldeNO+pPr9rIgPrtXqT4SWViTQRBeGM8CDxAyTopTsobjSiYZCF9Ta1gunl0G/8Vfp+SXfYCC+ZzWvP+L1pFhPRqzQQ8k+wMZIovObK1s+nlwPaLyayzw9a8sUnvWB/5rGPdIYnQWPgoNlLN9HpSmsAcw2z8DXI9pIxbr74cb3/HSfuYGOLkRqrOk6h4RCOfuWoTrZup1uEOn+fw8CAwEAAaOCAVQwggFQMB0GA1UdDgQWBBQv63nQ/pJAt5tLy8VJcbHe22ZOsjCCAR8GA1UdIwSCARYwggESgBQv63nQ/pJAt5tLy8VJcbHe22ZOsqGB7qSB6zCB6DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRMwEQYDVQQHEwpMb3NBbmdlbGVzMSAwHgYDVQQKExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UECxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAMTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQpExdQcml2YXRlIEludGVybmV0IEFjY2VzczEvMC0GCSqGSIb3DQEJARYgc2VjdXJlQHByaXZhdGVpbnRlcm5ldGFjY2Vzcy5jb22CCQCmew+WL/O6gzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBDQUAA4IBAQAna5PgrtxfwTumD4+3/SYvwoD66cB8IcK//h1mCzAduU8KgUXocLx7QgJWo9lnZ8xUryXvWab2usg4fqk7FPi00bED4f4qVQFVfGfPZIH9QQ7/48bPM9RyfzImZWUCenK37pdw4Bvgoys2rHLHbGen7f28knT2j/cbMxd78tQc20TIObGjo8+ISTRclSTRBtyCGohseKYpTS9himFERpUgNtefvYHbn70mIOzfOJFTVqfrptf9jXa9N8Mpy3ayfodz1wiqdteqFXkTYoSDctgKMiZ6GdocK9nMroQipIQtpnwd4yBDWIyC6Bvlkrq5TQUtYDQ8z9v+DMO6iwyIDRiU"
	PIACertificateStrong      = "MIIHqzCCBZOgAwIBAgIJAJ0u+vODZJntMA0GCSqGSIb3DQEBDQUAMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTAeFw0xNDA0MTcxNzQwMzNaFw0zNDA0MTIxNzQwMzNaMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBALVkhjumaqBbL8aSgj6xbX1QPTfTd1qHsAZd2B97m8Vw31c/2yQgZNf5qZY0+jOIHULNDe4R9TIvyBEbvnAg/OkPw8n/+ScgYOeH876VUXzjLDBnDb8DLr/+w9oVsuDeFJ9KV2UFM1OYX0SnkHnrYAN2QLF98ESK4NCSU01h5zkcgmQ+qKSfA9Ny0/UpsKPBFqsQ25NvjDWFhCpeqCHKUJ4Be27CDbSl7lAkBuHMPHJs8f8xPgAbHRXZOxVCpayZ2SNDfCwsnGWpWFoMGvdMbygngCn6jA/W1VSFOlRlfLuuGe7QFfDwA0jaLCxuWt/BgZylp7tAzYKR8lnWmtUCPm4+BtjyVDYtDCiGBD9Z4P13RFWvJHw5aapx/5W/CuvVyI7pKwvc2IT+KPxCUhH1XI8ca5RN3C9NoPJJf6qpg4g0rJH3aaWkoMRrYvQ+5PXXYUzjtRHImghRGd/ydERYoAZXuGSbPkm9Y/p2X8unLcW+F0xpJD98+ZI+tzSsI99Zs5wijSUGYr9/j18KHFTMQ8n+1jauc5bCCegN27dPeKXNSZ5riXFL2XX6BkY68y58UaNzmeGMiUL9BOV1iV+PMb7B7PYs7oFLjAhh0EdyvfHkrh/ZV9BEhtFa7yXp8XR0J6vz1YV9R6DYJmLjOEbhU8N0gc3tZm4Qz39lIIG6w3FDAgMBAAGjggFUMIIBUDAdBgNVHQ4EFgQUrsRtyWJftjpdRM0+925Y6Cl08SUwggEfBgNVHSMEggEWMIIBEoAUrsRtyWJftjpdRM0+925Y6Cl08SWhge6kgeswgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tggkAnS7684Nkme0wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQ0FAAOCAgEAJsfhsPk3r8kLXLxY+v+vHzbr4ufNtqnL9/1Uuf8NrsCtpXAoyZ0YqfbkWx3NHTZ7OE9ZRhdMP/RqHQE1p4N4Sa1nZKhTKasV6KhHDqSCt/dvEm89xWm2MVA7nyzQxVlHa9AkcBaemcXEiyT19XdpiXOP4Vhs+J1R5m8zQOxZlV1GtF9vsXmJqWZpOVPmZ8f35BCsYPvv4yMewnrtAC8PFEK/bOPeYcKN50bol22QYaZuLfpkHfNiFTnfMh8sl/ablPyNY7DUNiP5DRcMdIwmfGQxR5WEQoHL3yPJ42LkB5zs6jIm26DGNXfwura/mi105+ENH1CaROtRYwkiHb08U6qLXXJz80mWJkT90nr8Asj35xN2cUppg74nG3YVav/38P48T56hG1NHbYF5uOCske19F6wi9maUoto/3vEr0rnXJUp2KODmKdvBI7co245lHBABWikk8VfejQSlCtDBXn644ZMtAdoxKNfR2WTFVEwJiyd1Fzx0yujuiXDROLhISLQDRjVVAvawrAtLZWYK31bY7KlezPlQnl/D9Asxe85l8jO5+0LdJ6VyOs/Hd4w52alDW/MFySDZSfQHMTIc30hLBJ8OnCEIvluVQQ2UQvoW+no177N9L2Y+M9TcTA62ZyMXShHQGeh20rb4kK8f+iFX8NxtdHVSkxMEFSfDDyQ="
)

func PIAGeoChoices() (choices []string) {
	servers := PIAServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return choices
}

//nolint:lll
func PIAServers() []models.PIAServer {
	return []models.PIAServer{
		{Region: "AU Melbourne", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "melbourne405", IPs: []net.IP{{103, 2, 198, 108}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "melbourne405", IPs: []net.IP{{103, 2, 198, 103}}}},
		{Region: "AU Perth", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "perth404", IPs: []net.IP{{43, 250, 205, 186}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "perth404", IPs: []net.IP{{43, 250, 205, 188}}}},
		{Region: "AU Sydney", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "sydney405", IPs: []net.IP{{27, 50, 76, 132}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "sydney405", IPs: []net.IP{{27, 50, 76, 132}}}},
		{Region: "Albania", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "tirana401", IPs: []net.IP{{31, 171, 154, 131}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "tirana401", IPs: []net.IP{{31, 171, 154, 137}}}},
		{Region: "Algeria", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "algiers402", IPs: []net.IP{{45, 133, 91, 209}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "algiers402", IPs: []net.IP{{45, 133, 91, 227}}}},
		{Region: "Andorra", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "andorra401", IPs: []net.IP{{45, 139, 49, 232}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "andorra401", IPs: []net.IP{{45, 139, 49, 238}}}},
		{Region: "Argentina", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "buenosaires401", IPs: []net.IP{{190, 106, 134, 92}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "buenosaires401", IPs: []net.IP{{190, 106, 134, 89}}}},
		{Region: "Armenia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "armenia402", IPs: []net.IP{{45, 139, 50, 229}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "armenia402", IPs: []net.IP{{45, 139, 50, 213}}}},
		{Region: "Austria", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "vienna403", IPs: []net.IP{{156, 146, 60, 104}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "vienna403", IPs: []net.IP{{156, 146, 60, 100}}}},
		{Region: "Bahamas", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "bahamas402", IPs: []net.IP{{45, 132, 143, 206}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "bahamas402", IPs: []net.IP{{45, 132, 143, 229}}}},
		{Region: "Belgium", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "brussels403", IPs: []net.IP{{5, 253, 205, 147}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "brussels403", IPs: []net.IP{{5, 253, 205, 153}}}},
		{Region: "Bosnia and Herzegovina", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "sarajevo401", IPs: []net.IP{{185, 212, 111, 76}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "sarajevo401", IPs: []net.IP{{185, 212, 111, 77}}}},
		{Region: "Brazil", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "saopaolo402", IPs: []net.IP{{188, 241, 177, 56}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "saopaolo402", IPs: []net.IP{{188, 241, 177, 51}}}},
		{Region: "Bulgaria", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "sofia401", IPs: []net.IP{{217, 138, 221, 131}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "sofia401", IPs: []net.IP{{217, 138, 221, 133}}}},
		{Region: "CA Montreal", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "montreal403", IPs: []net.IP{{172, 98, 71, 62}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "montreal403", IPs: []net.IP{{172, 98, 71, 59}}}},
		{Region: "CA Ontario", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "ontario402", IPs: []net.IP{{172, 83, 47, 138}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "ontario402", IPs: []net.IP{{172, 83, 47, 196}}}},
		{Region: "CA Toronto", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "toronto405", IPs: []net.IP{{172, 83, 47, 250}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "toronto405", IPs: []net.IP{{172, 83, 47, 251}}}},
		{Region: "CA Vancouver", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "vancouver407", IPs: []net.IP{{172, 98, 89, 70}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "vancouver407", IPs: []net.IP{{172, 98, 89, 18}}}},
		{Region: "Cambodia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "cambodia401", IPs: []net.IP{{188, 215, 235, 105}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "cambodia401", IPs: []net.IP{{188, 215, 235, 102}}}},
		{Region: "China", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "china403", IPs: []net.IP{{86, 107, 104, 212}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "china403", IPs: []net.IP{{86, 107, 104, 216}}}},
		{Region: "Cyprus", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "cyprus402", IPs: []net.IP{{45, 132, 137, 220}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "cyprus402", IPs: []net.IP{{45, 132, 137, 225}}}},
		{Region: "Czech Republic", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "prague402", IPs: []net.IP{{212, 102, 39, 148}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "prague402", IPs: []net.IP{{212, 102, 39, 149}}}},
		{Region: "DE Berlin", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "berlin410", IPs: []net.IP{{89, 36, 76, 153}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "berlin410", IPs: []net.IP{{89, 36, 76, 149}}}},
		{Region: "DE Frankfurt", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "frankfurt406", IPs: []net.IP{{212, 102, 57, 96}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "frankfurt406", IPs: []net.IP{{212, 102, 57, 106}}}},
		{Region: "Denmark", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "copenhagen402", IPs: []net.IP{{188, 126, 94, 93}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "copenhagen402", IPs: []net.IP{{188, 126, 94, 93}}}},
		{Region: "Egypt", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "cairo401", IPs: []net.IP{{188, 214, 122, 106}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "cairo401", IPs: []net.IP{{188, 214, 122, 104}}}},
		{Region: "Estonia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "talinn402", IPs: []net.IP{{95, 153, 31, 73}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "talinn402", IPs: []net.IP{{95, 153, 31, 73}}}},
		{Region: "Finland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "helsinki402", IPs: []net.IP{{188, 126, 89, 45}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "helsinki402", IPs: []net.IP{{188, 126, 89, 45}}}},
		{Region: "France", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "paris402", IPs: []net.IP{{156, 146, 63, 159}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "paris402", IPs: []net.IP{{156, 146, 63, 159}}}},
		{Region: "Georgia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "georgia401", IPs: []net.IP{{45, 132, 138, 245}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "georgia401", IPs: []net.IP{{45, 132, 138, 236}}}},
		{Region: "Greece", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "athens401", IPs: []net.IP{{154, 57, 3, 80}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "athens401", IPs: []net.IP{{154, 57, 3, 84}}}},
		{Region: "Greenland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "greenland402", IPs: []net.IP{{45, 131, 209, 222}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "greenland402", IPs: []net.IP{{45, 131, 209, 208}}}},
		{Region: "Hong Kong", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "hongkong402", IPs: []net.IP{{86, 107, 104, 234}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "hongkong402", IPs: []net.IP{{86, 107, 104, 240}}}},
		{Region: "Hungary", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "budapest402", IPs: []net.IP{{86, 106, 74, 121}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "budapest402", IPs: []net.IP{{86, 106, 74, 125}}}},
		{Region: "Iceland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "reykjavik402", IPs: []net.IP{{45, 133, 193, 86}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "reykjavik402", IPs: []net.IP{{45, 133, 193, 86}}}},
		{Region: "India", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "mumbai405", IPs: []net.IP{{45, 120, 139, 97}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "mumbai405", IPs: []net.IP{{45, 120, 139, 97}}}},
		{Region: "Iran", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "iran402", IPs: []net.IP{{45, 131, 4, 219}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "iran402", IPs: []net.IP{{45, 131, 4, 218}}}},
		{Region: "Ireland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "dublin404", IPs: []net.IP{{193, 56, 252, 28}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "dublin404", IPs: []net.IP{{193, 56, 252, 24}}}},
		{Region: "Isle of Man", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "douglas401", IPs: []net.IP{{45, 132, 140, 236}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "douglas401", IPs: []net.IP{{45, 132, 140, 244}}}},
		{Region: "Israel", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "jerusalem401", IPs: []net.IP{{185, 77, 248, 19}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "jerusalem401", IPs: []net.IP{{185, 77, 248, 17}}}},
		{Region: "Italy", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "milano402", IPs: []net.IP{{156, 146, 41, 20}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "milano402", IPs: []net.IP{{156, 146, 41, 42}}}},
		{Region: "Japan", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "tokyo401", IPs: []net.IP{{156, 146, 34, 135}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "tokyo401", IPs: []net.IP{{156, 146, 34, 157}}}},
		{Region: "Kazakhstan", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "kazakhstan402", IPs: []net.IP{{45, 133, 88, 209}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "kazakhstan402", IPs: []net.IP{{45, 133, 88, 229}}}},
		{Region: "Latvia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "riga401", IPs: []net.IP{{109, 248, 149, 12}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "riga401", IPs: []net.IP{{109, 248, 149, 12}}}},
		{Region: "Liechtenstein", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "liechtenstein401", IPs: []net.IP{{45, 139, 48, 236}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "liechtenstein401", IPs: []net.IP{{45, 139, 48, 242}}}},
		{Region: "Lithuania", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "vilnius401", IPs: []net.IP{{85, 206, 165, 163}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "vilnius401", IPs: []net.IP{{85, 206, 165, 163}}}},
		{Region: "Luxembourg", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "luxembourg401", IPs: []net.IP{{92, 223, 89, 74}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "luxembourg401", IPs: []net.IP{{92, 223, 89, 78}}}},
		{Region: "Macedonia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "macedonia401", IPs: []net.IP{{185, 225, 28, 115}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "macedonia401", IPs: []net.IP{{185, 225, 28, 115}}}},
		{Region: "Malta", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "malta401", IPs: []net.IP{{45, 137, 198, 238}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "malta401", IPs: []net.IP{{45, 137, 198, 244}}}},
		{Region: "Mexico", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "mexico403", IPs: []net.IP{{77, 81, 142, 8}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "mexico403", IPs: []net.IP{{77, 81, 142, 7}}}},
		{Region: "Moldova", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "chisinau401", IPs: []net.IP{{178, 175, 129, 43}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "chisinau401", IPs: []net.IP{{178, 175, 129, 44}}}},
		{Region: "Monaco", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "monaco402", IPs: []net.IP{{45, 137, 199, 226}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "monaco402", IPs: []net.IP{{45, 137, 199, 218}}}},
		{Region: "Montenegro", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "montenegro402", IPs: []net.IP{{45, 131, 208, 212}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "montenegro402", IPs: []net.IP{{45, 131, 208, 212}}}},
		{Region: "Morocco", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "morocco401", IPs: []net.IP{{45, 131, 211, 233}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "morocco401", IPs: []net.IP{{45, 131, 211, 248}}}},
		{Region: "Netherlands", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "amsterdam416", IPs: []net.IP{{212, 102, 35, 136}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "amsterdam416", IPs: []net.IP{{212, 102, 35, 136}}}},
		{Region: "New Zealand", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "newzealand403", IPs: []net.IP{{43, 250, 207, 89}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "newzealand403", IPs: []net.IP{{43, 250, 207, 94}}}},
		{Region: "Norway", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "oslo403", IPs: []net.IP{{46, 246, 122, 124}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "oslo403", IPs: []net.IP{{46, 246, 122, 99}}}},
		{Region: "Panama", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "panama401", IPs: []net.IP{{45, 131, 210, 248}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "panama401", IPs: []net.IP{{45, 131, 210, 231}}}},
		{Region: "Philippines", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "philippines401", IPs: []net.IP{{188, 214, 125, 142}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "philippines401", IPs: []net.IP{{188, 214, 125, 142}}}},
		{Region: "Poland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "warsaw402", IPs: []net.IP{{194, 110, 114, 13}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "warsaw402", IPs: []net.IP{{194, 110, 114, 13}}}},
		{Region: "Portugal", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "lisbon401", IPs: []net.IP{{89, 26, 241, 72}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "lisbon401", IPs: []net.IP{{89, 26, 241, 76}}}},
		{Region: "Qatar", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "qatar401", IPs: []net.IP{{45, 131, 7, 234}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "qatar401", IPs: []net.IP{{45, 131, 7, 232}}}},
		{Region: "Romania", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "romania408", IPs: []net.IP{{143, 244, 54, 93}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "romania408", IPs: []net.IP{{143, 244, 54, 92}}}},
		{Region: "Saudi Arabia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "saudiarabia401", IPs: []net.IP{{45, 131, 6, 238}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "saudiarabia401", IPs: []net.IP{{45, 131, 6, 231}}}},
		{Region: "Serbia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "belgrade401", IPs: []net.IP{{37, 120, 193, 254}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "belgrade401", IPs: []net.IP{{37, 120, 193, 254}}}},
		{Region: "Singapore", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "singapore401", IPs: []net.IP{{156, 146, 57, 210}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "singapore401", IPs: []net.IP{{156, 146, 57, 190}}}},
		{Region: "Slovakia", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "bratislava401", IPs: []net.IP{{37, 120, 221, 93}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "bratislava401", IPs: []net.IP{{37, 120, 221, 83}}}},
		{Region: "South Africa", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "johannesburg401", IPs: []net.IP{{154, 16, 93, 46}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "johannesburg401", IPs: []net.IP{{154, 16, 93, 44}}}},
		{Region: "Spain", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "madrid402", IPs: []net.IP{{212, 102, 49, 33}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "madrid402", IPs: []net.IP{{212, 102, 49, 29}}}},
		{Region: "Sri Lanka", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "srilanka402", IPs: []net.IP{{45, 132, 136, 224}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "srilanka402", IPs: []net.IP{{45, 132, 136, 216}}}},
		{Region: "Sweden", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "stockholm404", IPs: []net.IP{{195, 246, 120, 140}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "stockholm404", IPs: []net.IP{{195, 246, 120, 116}}}},
		{Region: "Switzerland", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "zurich404", IPs: []net.IP{{212, 102, 37, 104}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "zurich404", IPs: []net.IP{{212, 102, 37, 84}}}},
		{Region: "Taiwan", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "taiwan401", IPs: []net.IP{{188, 214, 106, 76}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "taiwan401", IPs: []net.IP{{188, 214, 106, 71}}}},
		{Region: "Turkey", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "istanbul401", IPs: []net.IP{{188, 213, 34, 71}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "istanbul401", IPs: []net.IP{{188, 213, 34, 76}}}},
		{Region: "UK London", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "london412", IPs: []net.IP{{37, 235, 96, 109}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "london412", IPs: []net.IP{{37, 235, 96, 109}}}},
		{Region: "UK Manchester", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "manchester460", IPs: []net.IP{{37, 120, 159, 136}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "manchester460", IPs: []net.IP{{37, 120, 159, 122}}}},
		{Region: "UK Southampton", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "southampton401", IPs: []net.IP{{143, 244, 37, 223}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "southampton401", IPs: []net.IP{{143, 244, 37, 189}}}},
		{Region: "US Atlanta", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "atlanta421", IPs: []net.IP{{154, 21, 21, 77}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "atlanta421", IPs: []net.IP{{154, 21, 21, 70}}}},
		{Region: "US California", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "losangeles401", IPs: []net.IP{{37, 235, 107, 62}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "losangeles401", IPs: []net.IP{{37, 235, 107, 17}}}},
		{Region: "US Chicago", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "chicago416", IPs: []net.IP{{154, 21, 114, 12}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "chicago416", IPs: []net.IP{{154, 21, 114, 12}}}},
		{Region: "US Denver", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "denver402", IPs: []net.IP{{70, 39, 126, 157}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "denver402", IPs: []net.IP{{70, 39, 126, 175}}}},
		{Region: "US East", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "newjersey402", IPs: []net.IP{{37, 235, 103, 74}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "newjersey402", IPs: []net.IP{{37, 235, 103, 131}}}},
		{Region: "US Florida", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "miami405", IPs: []net.IP{{37, 235, 98, 169}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "miami405", IPs: []net.IP{{37, 235, 98, 188}}}},
		{Region: "US Houston", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "houston418", IPs: []net.IP{{205, 251, 154, 205}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "houston418", IPs: []net.IP{{205, 251, 154, 208}}}},
		{Region: "US Las Vegas", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "lasvegas402", IPs: []net.IP{{45, 89, 173, 178}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "lasvegas402", IPs: []net.IP{{45, 89, 173, 181}}}},
		{Region: "US New York", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "newyork403", IPs: []net.IP{{156, 146, 54, 108}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "newyork403", IPs: []net.IP{{156, 146, 54, 63}}}},
		{Region: "US Seattle", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "seattle417", IPs: []net.IP{{154, 21, 20, 187}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "seattle417", IPs: []net.IP{{154, 21, 20, 169}}}},
		{Region: "US Silicon Valley", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "siliconvalley401", IPs: []net.IP{{154, 21, 212, 40}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "siliconvalley401", IPs: []net.IP{{154, 21, 212, 14}}}},
		{Region: "US Texas", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "dallas401", IPs: []net.IP{{156, 146, 53, 180}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "dallas401", IPs: []net.IP{{156, 146, 53, 186}}}},
		{Region: "US Washington DC", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "washington412", IPs: []net.IP{{23, 105, 168, 143}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "washington412", IPs: []net.IP{{23, 105, 168, 150}}}},
		{Region: "US West", PortForward: false, OpenvpnUDP: models.PIAServerOpenvpn{CN: "phoenix407", IPs: []net.IP{{184, 170, 241, 67}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "phoenix407", IPs: []net.IP{{184, 170, 241, 121}}}},
		{Region: "Ukraine", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "kiev402", IPs: []net.IP{{62, 149, 20, 23}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "kiev402", IPs: []net.IP{{62, 149, 20, 22}}}},
		{Region: "United Arab Emirates", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "dubai403", IPs: []net.IP{{217, 138, 193, 146}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "dubai403", IPs: []net.IP{{217, 138, 193, 148}}}},
		{Region: "Venezuela", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "venezuela402", IPs: []net.IP{{45, 133, 89, 217}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "venezuela402", IPs: []net.IP{{45, 133, 89, 217}}}},
		{Region: "Vietnam", PortForward: true, OpenvpnUDP: models.PIAServerOpenvpn{CN: "vietnam401", IPs: []net.IP{{188, 214, 152, 76}}}, OpenvpnTCP: models.PIAServerOpenvpn{CN: "vietnam401", IPs: []net.IP{{188, 214, 152, 70}}}},
	}
}

func PIAOldGeoChoices() (choices []string) {
	servers := PIAOldServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return choices
}

//nolint:lll
func PIAOldServers() []models.PIAOldServer {
	return []models.PIAOldServer{
		{Region: "AU Melbourne", IPs: []net.IP{{27, 50, 82, 131}, {43, 250, 204, 105}, {43, 250, 204, 107}, {43, 250, 204, 109}, {43, 250, 204, 111}, {43, 250, 204, 113}, {43, 250, 204, 115}, {43, 250, 204, 117}, {43, 250, 204, 119}, {43, 250, 204, 123}, {43, 250, 204, 125}}},
		{Region: "AU Perth", IPs: []net.IP{{43, 250, 205, 59}, {43, 250, 205, 91}, {43, 250, 205, 93}, {43, 250, 205, 95}}},
		{Region: "AU Sydney", IPs: []net.IP{{27, 50, 68, 23}, {27, 50, 70, 87}, {27, 50, 77, 251}, {27, 50, 81, 117}, {103, 13, 102, 123}, {103, 13, 102, 127}, {118, 127, 60, 51}, {221, 121, 145, 135}, {221, 121, 145, 137}, {221, 121, 145, 145}, {221, 121, 145, 147}, {221, 121, 145, 159}, {221, 121, 146, 203}, {221, 121, 148, 221}, {221, 121, 152, 215}}},
		{Region: "Albania", IPs: []net.IP{{31, 171, 154, 114}}},
		{Region: "Argentina", IPs: []net.IP{{190, 106, 134, 100}}},
		{Region: "Austria", IPs: []net.IP{{89, 187, 168, 6}, {156, 146, 60, 129}}},
		{Region: "Belgium", IPs: []net.IP{{77, 243, 191, 18}, {77, 243, 191, 19}, {77, 243, 191, 20}, {185, 232, 21, 26}}},
		{Region: "Bosnia and Herzegovina", IPs: []net.IP{{185, 164, 35, 54}}},
		{Region: "Bulgaria", IPs: []net.IP{{217, 138, 221, 66}}},
		{Region: "CA Montreal", IPs: []net.IP{{172, 98, 71, 194}, {199, 36, 223, 130}, {199, 36, 223, 194}}},
		{Region: "CA Ontario", IPs: []net.IP{{162, 219, 176, 26}, {162, 219, 176, 42}, {184, 75, 208, 2}, {184, 75, 208, 90}, {184, 75, 208, 114}, {184, 75, 208, 122}, {184, 75, 208, 130}, {184, 75, 208, 146}, {184, 75, 208, 170}, {184, 75, 208, 202}, {184, 75, 210, 18}, {184, 75, 210, 98}, {184, 75, 210, 106}, {184, 75, 213, 186}, {184, 75, 213, 218}, {184, 75, 214, 18}, {184, 75, 215, 18}, {184, 75, 215, 26}, {184, 75, 215, 66}, {184, 75, 215, 74}}},
		{Region: "CA Toronto", IPs: []net.IP{{66, 115, 142, 130}, {66, 115, 145, 199}, {172, 98, 92, 66}, {172, 98, 92, 130}, {172, 98, 92, 194}}},
		{Region: "CA Vancouver", IPs: []net.IP{{162, 216, 47, 66}, {162, 216, 47, 194}, {172, 98, 89, 130}, {172, 98, 89, 194}}},
		{Region: "Czech Republic", IPs: []net.IP{{212, 102, 39, 1}}},
		{Region: "DE Berlin", IPs: []net.IP{{185, 230, 127, 238}, {193, 176, 86, 122}, {193, 176, 86, 123}, {193, 176, 86, 134}, {193, 176, 86, 178}, {194, 36, 108, 6}}},
		{Region: "DE Frankfurt", IPs: []net.IP{{195, 181, 170, 239}, {195, 181, 170, 240}, {195, 181, 170, 241}, {195, 181, 170, 242}, {195, 181, 170, 243}, {195, 181, 170, 244}, {212, 102, 57, 138}}},
		{Region: "Denmark", IPs: []net.IP{{188, 126, 94, 34}}},
		{Region: "Estonia", IPs: []net.IP{{77, 247, 111, 82}, {77, 247, 111, 98}, {77, 247, 111, 114}, {77, 247, 111, 130}}},
		{Region: "Finland", IPs: []net.IP{{188, 126, 89, 4}, {188, 126, 89, 194}}},
		{Region: "France", IPs: []net.IP{{156, 146, 63, 1}, {156, 146, 63, 65}}},
		{Region: "Greece", IPs: []net.IP{{154, 57, 3, 91}, {154, 57, 3, 106}, {154, 57, 3, 145}}},
		{Region: "Hungary", IPs: []net.IP{{185, 128, 26, 18}, {185, 128, 26, 19}, {185, 128, 26, 20}, {185, 128, 26, 21}, {185, 128, 26, 22}, {185, 128, 26, 23}, {185, 128, 26, 24}, {185, 189, 114, 98}}},
		{Region: "Iceland", IPs: []net.IP{{45, 133, 193, 50}}},
		{Region: "India", IPs: []net.IP{{45, 120, 139, 108}, {45, 120, 139, 109}, {150, 242, 12, 155}, {150, 242, 12, 171}, {150, 242, 12, 187}}},
		{Region: "Ireland", IPs: []net.IP{{193, 56, 252, 210}, {193, 56, 252, 226}, {193, 56, 252, 242}, {193, 56, 252, 250}, {193, 56, 252, 251}, {193, 56, 252, 252}}},
		{Region: "Israel", IPs: []net.IP{{31, 168, 172, 142}, {31, 168, 172, 143}, {31, 168, 172, 145}, {31, 168, 172, 146}}},
		{Region: "Italy", IPs: []net.IP{{156, 146, 41, 129}, {156, 146, 41, 193}}},
		{Region: "Japan", IPs: []net.IP{{156, 146, 34, 1}, {156, 146, 34, 65}}},
		{Region: "Latvia", IPs: []net.IP{{46, 183, 217, 34}, {46, 183, 218, 130}, {46, 183, 218, 146}}},
		{Region: "Lithuania", IPs: []net.IP{{85, 206, 165, 96}, {85, 206, 165, 112}, {85, 206, 165, 128}}},
		{Region: "Luxembourg", IPs: []net.IP{{92, 223, 89, 133}, {92, 223, 89, 134}, {92, 223, 89, 135}, {92, 223, 89, 136}, {92, 223, 89, 137}, {92, 223, 89, 138}, {92, 223, 89, 140}, {92, 223, 89, 142}}},
		{Region: "Moldova", IPs: []net.IP{{178, 17, 172, 242}, {178, 17, 173, 194}, {178, 175, 128, 34}}},
		{Region: "Netherlands", IPs: []net.IP{{89, 187, 174, 198}, {212, 102, 35, 101}, {212, 102, 35, 102}, {212, 102, 35, 103}, {212, 102, 35, 104}}},
		{Region: "New Zealand", IPs: []net.IP{{43, 250, 207, 1}, {43, 250, 207, 3}}},
		{Region: "North Macedonia", IPs: []net.IP{{185, 225, 28, 130}}},
		{Region: "Norway", IPs: []net.IP{{46, 246, 122, 34}, {46, 246, 122, 162}}},
		{Region: "Poland", IPs: []net.IP{{185, 244, 214, 195}, {185, 244, 214, 196}, {185, 244, 214, 197}, {185, 244, 214, 198}, {185, 244, 214, 199}, {185, 244, 214, 200}}},
		{Region: "Portugal", IPs: []net.IP{{89, 26, 241, 86}, {89, 26, 241, 102}, {89, 26, 241, 130}}},
		{Region: "Romania", IPs: []net.IP{{86, 105, 25, 69}, {86, 105, 25, 70}, {86, 105, 25, 74}, {86, 105, 25, 75}, {86, 105, 25, 76}, {86, 105, 25, 77}, {86, 105, 25, 78}, {89, 33, 8, 38}, {89, 33, 8, 42}, {93, 115, 7, 70}, {94, 176, 148, 35}, {143, 244, 54, 1}, {185, 45, 12, 126}, {185, 210, 218, 98}, {185, 210, 218, 99}, {185, 210, 218, 100}, {185, 210, 218, 101}, {185, 210, 218, 102}, {185, 210, 218, 105}, {188, 240, 220, 26}}},
		{Region: "Serbia", IPs: []net.IP{{37, 120, 193, 226}}},
		{Region: "Singapore", IPs: []net.IP{{156, 146, 56, 193}, {156, 146, 57, 38}, {156, 146, 57, 235}, {156, 146, 57, 244}}},
		{Region: "Slovakia", IPs: []net.IP{{37, 120, 221, 82}, {37, 120, 221, 98}}},
		{Region: "South Africa", IPs: []net.IP{{102, 165, 20, 133}}},
		{Region: "Spain", IPs: []net.IP{{212, 102, 49, 185}, {212, 102, 49, 251}}},
		{Region: "Sweden", IPs: []net.IP{{46, 246, 3, 254}}},
		{Region: "Switzerland", IPs: []net.IP{{156, 146, 62, 193}, {212, 102, 36, 1}, {212, 102, 36, 166}, {212, 102, 37, 240}, {212, 102, 37, 241}, {212, 102, 37, 242}, {212, 102, 37, 243}}},
		{Region: "Turkey", IPs: []net.IP{{185, 195, 79, 34}, {185, 195, 79, 82}}},
		{Region: "UAE", IPs: []net.IP{{45, 9, 250, 46}}},
		{Region: "UK London", IPs: []net.IP{{212, 102, 52, 1}}},
		{Region: "UK Manchester", IPs: []net.IP{{89, 238, 137, 36}, {89, 238, 137, 37}, {89, 238, 137, 38}, {89, 238, 137, 39}, {89, 238, 139, 52}, {89, 238, 139, 53}, {89, 238, 139, 54}, {89, 238, 139, 55}, {89, 238, 139, 56}, {89, 238, 139, 57}, {89, 238, 139, 58}, {89, 249, 67, 220}}},
		{Region: "UK Southampton", IPs: []net.IP{{143, 244, 36, 58}, {143, 244, 37, 1}, {143, 244, 38, 1}, {143, 244, 38, 60}, {143, 244, 38, 119}}},
		{Region: "US Atlanta", IPs: []net.IP{{156, 146, 46, 1}, {156, 146, 46, 134}, {156, 146, 46, 198}, {156, 146, 47, 11}}},
		{Region: "US California", IPs: []net.IP{{37, 235, 108, 208}, {89, 187, 187, 129}, {89, 187, 187, 162}, {91, 207, 175, 194}, {91, 207, 175, 195}, {91, 207, 175, 197}, {91, 207, 175, 198}, {91, 207, 175, 199}, {91, 207, 175, 200}, {91, 207, 175, 205}, {91, 207, 175, 206}, {91, 207, 175, 207}, {91, 207, 175, 209}, {91, 207, 175, 210}, {91, 207, 175, 212}}},
		{Region: "US Chicago", IPs: []net.IP{{156, 146, 50, 1}, {156, 146, 50, 65}, {156, 146, 50, 134}, {156, 146, 50, 198}, {156, 146, 51, 11}, {212, 102, 58, 113}, {212, 102, 59, 54}, {212, 102, 59, 129}}},
		{Region: "US Dallas", IPs: []net.IP{{156, 146, 38, 65}, {156, 146, 38, 161}, {156, 146, 39, 1}, {156, 146, 39, 6}, {156, 146, 52, 6}, {156, 146, 52, 70}, {156, 146, 52, 139}, {156, 146, 52, 203}}},
		{Region: "US Denver", IPs: []net.IP{{70, 39, 77, 130}, {70, 39, 92, 2}, {70, 39, 113, 194}, {174, 128, 225, 2}, {174, 128, 226, 10}, {174, 128, 226, 18}, {174, 128, 227, 2}, {174, 128, 227, 226}, {174, 128, 236, 98}, {174, 128, 242, 234}, {174, 128, 242, 250}, {174, 128, 243, 98}, {174, 128, 244, 74}, {174, 128, 245, 122}, {174, 128, 246, 10}, {199, 115, 98, 146}, {199, 115, 98, 234}, {199, 115, 101, 178}, {199, 115, 101, 186}, {199, 115, 102, 146}}},
		{Region: "US East", IPs: []net.IP{{156, 146, 58, 202}, {156, 146, 58, 203}, {156, 146, 58, 204}, {156, 146, 58, 205}, {156, 146, 58, 207}, {156, 146, 58, 208}, {156, 146, 58, 209}, {193, 37, 253, 115}, {193, 37, 253, 134}, {194, 59, 251, 8}, {194, 59, 251, 11}, {194, 59, 251, 22}, {194, 59, 251, 28}, {194, 59, 251, 56}, {194, 59, 251, 62}, {194, 59, 251, 69}, {194, 59, 251, 82}, {194, 59, 251, 84}, {194, 59, 251, 91}, {194, 59, 251, 112}}},
		{Region: "US Florida", IPs: []net.IP{{193, 37, 252, 6}, {193, 37, 252, 7}, {193, 37, 252, 8}, {193, 37, 252, 9}, {193, 37, 252, 10}, {193, 37, 252, 11}, {193, 37, 252, 12}, {193, 37, 252, 14}, {193, 37, 252, 15}, {193, 37, 252, 16}, {193, 37, 252, 17}, {193, 37, 252, 18}, {193, 37, 252, 19}, {193, 37, 252, 20}, {193, 37, 252, 21}, {193, 37, 252, 23}, {193, 37, 252, 24}, {193, 37, 252, 25}, {193, 37, 252, 26}, {193, 37, 252, 27}}},
		{Region: "US Houston", IPs: []net.IP{{74, 81, 88, 26}, {74, 81, 88, 42}, {74, 81, 88, 66}, {74, 81, 88, 74}, {205, 251, 148, 66}, {205, 251, 148, 90}, {205, 251, 148, 98}, {205, 251, 148, 122}, {205, 251, 148, 130}, {205, 251, 148, 138}, {205, 251, 148, 186}, {205, 251, 150, 146}, {205, 251, 150, 170}}},
		{Region: "US Las Vegas", IPs: []net.IP{{79, 110, 53, 50}, {79, 110, 53, 66}, {79, 110, 53, 98}, {79, 110, 53, 114}, {79, 110, 53, 130}, {79, 110, 53, 146}, {79, 110, 53, 162}, {79, 110, 53, 178}, {79, 110, 53, 194}, {79, 110, 53, 210}, {162, 251, 236, 7}, {199, 127, 56, 83}, {199, 127, 56, 84}, {199, 127, 56, 87}, {199, 127, 56, 89}, {199, 127, 56, 90}}},
		{Region: "US New York City", IPs: []net.IP{{156, 146, 36, 225}, {156, 146, 37, 129}, {156, 146, 58, 1}, {156, 146, 58, 134}}},
		{Region: "US Seattle", IPs: []net.IP{{156, 146, 48, 65}, {156, 146, 48, 135}, {156, 146, 48, 200}, {156, 146, 49, 13}, {212, 102, 46, 129}, {212, 102, 46, 193}, {212, 102, 47, 134}}},
		{Region: "US Silicon Valley", IPs: []net.IP{{199, 116, 118, 130}, {199, 116, 118, 132}, {199, 116, 118, 134}, {199, 116, 118, 136}, {199, 116, 118, 145}, {199, 116, 118, 148}, {199, 116, 118, 149}, {199, 116, 118, 157}, {199, 116, 118, 166}, {199, 116, 118, 169}, {199, 116, 118, 172}}},
		{Region: "US Washington DC", IPs: []net.IP{{70, 32, 0, 46}, {70, 32, 0, 51}, {70, 32, 0, 53}, {70, 32, 0, 62}, {70, 32, 0, 64}, {70, 32, 0, 68}, {70, 32, 0, 69}, {70, 32, 0, 72}, {70, 32, 0, 76}, {70, 32, 0, 77}, {70, 32, 0, 106}, {70, 32, 0, 107}, {70, 32, 0, 114}, {70, 32, 0, 116}, {70, 32, 0, 120}, {70, 32, 0, 167}, {70, 32, 0, 168}, {70, 32, 0, 170}, {70, 32, 0, 172}, {70, 32, 0, 173}}},
		{Region: "US West", IPs: []net.IP{{184, 170, 241, 130}, {184, 170, 241, 194}, {184, 170, 242, 135}, {184, 170, 242, 199}}},
		{Region: "Ukraine", IPs: []net.IP{{62, 149, 20, 10}, {62, 149, 20, 40}}},
	}
}

const (
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
