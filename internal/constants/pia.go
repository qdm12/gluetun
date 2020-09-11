package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

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

func PIAServers() []models.PIAServer {
	return []models.PIAServer{
		{Region: "AU Melbourne", IPs: []net.IP{{27, 50, 74, 184}}},
		{Region: "AU Perth", IPs: []net.IP{{43, 250, 205, 170}}},
		{Region: "AU Sydney", IPs: []net.IP{{103, 2, 196, 167}}},
		{Region: "Algeria", IPs: []net.IP{{45, 133, 91, 210}}},
		{Region: "Andorra", IPs: []net.IP{{45, 139, 49, 241}}},
		{Region: "Argentina", IPs: []net.IP{{190, 106, 134, 82}}},
		{Region: "Armenia", IPs: []net.IP{{45, 139, 50, 232}}},
		{Region: "Austria", IPs: []net.IP{{156, 146, 60, 14}}},
		{Region: "Bahamas", IPs: []net.IP{{45, 132, 143, 206}}},
		{Region: "Bangladesh", IPs: []net.IP{{45, 132, 142, 210}}},
		{Region: "Belgium", IPs: []net.IP{{5, 253, 205, 147}}},
		{Region: "Bulgaria", IPs: []net.IP{{217, 138, 221, 130}}},
		{Region: "CA Montreal", IPs: []net.IP{{172, 98, 71, 13}}},
		{Region: "CA Toronto", IPs: []net.IP{{66, 115, 142, 81}}},
		{Region: "Cambodia", IPs: []net.IP{{188, 215, 235, 103}}},
		{Region: "China", IPs: []net.IP{{45, 132, 193, 234}}},
		{Region: "Cyprus", IPs: []net.IP{{45, 132, 137, 235}}},
		{Region: "Czech Republic", IPs: []net.IP{{212, 102, 39, 194}}},
		{Region: "DE Berlin", IPs: []net.IP{{89, 36, 76, 69}}},
		{Region: "DE Frankfurt", IPs: []net.IP{{185, 216, 33, 164}}},
		{Region: "Denmark", IPs: []net.IP{{188, 126, 94, 124}}},
		{Region: "Egypt", IPs: []net.IP{{188, 214, 122, 119}}},
		{Region: "Finland", IPs: []net.IP{{188, 126, 89, 10}}},
		{Region: "France", IPs: []net.IP{{156, 146, 63, 210}}},
		{Region: "Georgia", IPs: []net.IP{{45, 132, 138, 236}}},
		{Region: "Greenland", IPs: []net.IP{{45, 131, 209, 233}}},
		{Region: "Hungary", IPs: []net.IP{{217, 138, 192, 222}}},
		{Region: "Iceland", IPs: []net.IP{{45, 133, 193, 85}}},
		{Region: "India", IPs: []net.IP{{103, 26, 205, 251}}},
		{Region: "Iran", IPs: []net.IP{{45, 131, 4, 208}}},
		{Region: "Ireland", IPs: []net.IP{{5, 157, 13, 41}}},
		{Region: "Isle of Man", IPs: []net.IP{{45, 132, 140, 213}}},
		{Region: "Israel", IPs: []net.IP{{185, 77, 248, 10}}},
		{Region: "Italy", IPs: []net.IP{{156, 146, 41, 77}}},
		{Region: "Japan", IPs: []net.IP{{156, 146, 34, 164}}},
		{Region: "Kazakhstan", IPs: []net.IP{{45, 133, 88, 231}}},
		{Region: "Liechtenstein", IPs: []net.IP{{45, 139, 48, 236}}},
		{Region: "Luxembourg", IPs: []net.IP{{92, 223, 89, 80}}},
		{Region: "Macao", IPs: []net.IP{{45, 137, 197, 207}}},
		{Region: "Malta", IPs: []net.IP{{45, 137, 198, 235}}},
		{Region: "Mexico", IPs: []net.IP{{77, 81, 142, 5}}},
		{Region: "Moldova", IPs: []net.IP{{178, 175, 129, 40}}},
		{Region: "Monaco", IPs: []net.IP{{45, 137, 199, 237}}},
		{Region: "Mongolia", IPs: []net.IP{{45, 139, 51, 211}}},
		{Region: "Montenegro", IPs: []net.IP{{45, 131, 208, 206}}},
		{Region: "Morocco", IPs: []net.IP{{45, 131, 211, 234}}},
		{Region: "Netherlands", IPs: []net.IP{{37, 235, 101, 73}}},
		{Region: "New Zealand", IPs: []net.IP{{43, 250, 207, 70}}},
		{Region: "Nigeria", IPs: []net.IP{{45, 137, 196, 208}}},
		{Region: "Norway", IPs: []net.IP{{46, 246, 122, 82}}},
		{Region: "Panama", IPs: []net.IP{{45, 131, 210, 206}}},
		{Region: "Philippines", IPs: []net.IP{{188, 214, 125, 138}}},
		{Region: "Poland", IPs: []net.IP{{217, 138, 209, 243}}},
		{Region: "Qatar", IPs: []net.IP{{45, 131, 7, 209}}},
		{Region: "Romania", IPs: []net.IP{{185, 45, 15, 22}}},
		{Region: "Saudi Arabia", IPs: []net.IP{{45, 131, 6, 208}}},
		{Region: "Serbia", IPs: []net.IP{{37, 120, 193, 248}}},
		{Region: "Singapore", IPs: []net.IP{{156, 146, 57, 123}}},
		{Region: "South Africa", IPs: []net.IP{{154, 16, 93, 35}}},
		{Region: "Spain", IPs: []net.IP{{195, 181, 167, 42}}},
		{Region: "Sri Lanka", IPs: []net.IP{{45, 132, 136, 232}}},
		{Region: "Sweden", IPs: []net.IP{{46, 246, 3, 150}}},
		{Region: "Switzerland", IPs: []net.IP{{212, 102, 37, 77}}},
		{Region: "Taiwan", IPs: []net.IP{{188, 214, 106, 70}}},
		{Region: "Turkey", IPs: []net.IP{{188, 213, 34, 87}}},
		{Region: "UK London", IPs: []net.IP{{37, 235, 96, 26}}},
		{Region: "UK Manchester", IPs: []net.IP{{193, 239, 84, 60}}},
		{Region: "US Atlanta", IPs: []net.IP{{195, 181, 171, 76}}},
		{Region: "US California", IPs: []net.IP{{37, 235, 108, 19}}},
		{Region: "US Chicago", IPs: []net.IP{{154, 21, 28, 111}}},
		{Region: "US Denver", IPs: []net.IP{{70, 39, 126, 143}}},
		{Region: "US Florida", IPs: []net.IP{{37, 235, 98, 18}}},
		{Region: "US Houston", IPs: []net.IP{{74, 81, 92, 147}}},
		{Region: "US New Jersey", IPs: []net.IP{{37, 235, 103, 75}}},
		{Region: "US New York", IPs: []net.IP{{156, 146, 55, 213}}},
		{Region: "US Seattle", IPs: []net.IP{{156, 146, 48, 14}}},
		{Region: "US Silicon Valley", IPs: []net.IP{{154, 21, 212, 228}}},
		{Region: "US Texas", IPs: []net.IP{{154, 29, 131, 17}}},
		{Region: "US Washington DC", IPs: []net.IP{{70, 32, 5, 172}}},
		{Region: "US West", IPs: []net.IP{{193, 37, 254, 239}}},
		{Region: "Ukraine", IPs: []net.IP{{62, 149, 20, 51}}},
		{Region: "United Arab Emirates", IPs: []net.IP{{45, 131, 5, 233}}},
		{Region: "Venezuela", IPs: []net.IP{{45, 133, 89, 212}}},
		{Region: "Vietnam", IPs: []net.IP{{188, 214, 152, 67}}},
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

func PIAOldServers() []models.PIAServer {
	return []models.PIAServer{
		{Region: "AU Melbourne", IPs: []net.IP{{27, 50, 82, 131}, {27, 50, 82, 133}, {43, 250, 204, 85}, {43, 250, 204, 87}, {43, 250, 204, 89}, {43, 250, 204, 91}, {43, 250, 204, 93}, {43, 250, 204, 95}, {43, 250, 204, 97}, {43, 250, 204, 99}, {43, 250, 204, 105}, {43, 250, 204, 107}, {43, 250, 204, 109}, {43, 250, 204, 111}, {43, 250, 204, 113}, {43, 250, 204, 115}, {43, 250, 204, 117}, {43, 250, 204, 119}, {43, 250, 204, 123}, {43, 250, 204, 125}}},
		{Region: "AU Perth", IPs: []net.IP{{43, 250, 205, 59}, {43, 250, 205, 93}, {43, 250, 205, 95}}},
		{Region: "AU Sydney", IPs: []net.IP{{27, 50, 68, 23}, {27, 50, 77, 251}, {27, 50, 81, 117}, {103, 13, 102, 117}, {103, 13, 102, 121}, {103, 13, 102, 123}, {103, 13, 102, 127}, {118, 127, 60, 51}, {221, 121, 145, 131}, {221, 121, 145, 135}, {221, 121, 145, 137}, {221, 121, 145, 143}, {221, 121, 145, 145}, {221, 121, 145, 147}, {221, 121, 145, 151}, {221, 121, 145, 159}, {221, 121, 146, 203}, {221, 121, 146, 217}, {221, 121, 148, 221}, {221, 121, 152, 215}}},
		{Region: "Albania", IPs: []net.IP{{31, 171, 154, 114}}},
		{Region: "Argentina", IPs: []net.IP{{190, 106, 134, 100}}},
		{Region: "Austria", IPs: []net.IP{{89, 187, 168, 6}, {156, 146, 60, 129}}},
		{Region: "Belgium", IPs: []net.IP{{77, 243, 191, 18}, {77, 243, 191, 19}, {77, 243, 191, 20}, {77, 243, 191, 21}, {77, 243, 191, 22}, {77, 243, 191, 23}, {185, 104, 186, 26}, {185, 232, 21, 26}, {185, 232, 21, 27}, {185, 232, 21, 28}, {185, 232, 21, 29}}},
		{Region: "Bosnia and Herzegovina", IPs: []net.IP{{185, 164, 35, 54}}},
		{Region: "Bulgaria", IPs: []net.IP{{217, 138, 221, 66}}},
		{Region: "CA Montreal", IPs: []net.IP{{172, 98, 71, 194}, {199, 36, 223, 130}, {199, 36, 223, 194}}},
		{Region: "CA Ontario", IPs: []net.IP{{162, 219, 176, 194}, {162, 253, 128, 98}, {184, 75, 208, 18}, {184, 75, 208, 34}, {184, 75, 208, 66}, {184, 75, 208, 90}, {184, 75, 208, 114}, {184, 75, 208, 122}, {184, 75, 208, 130}, {184, 75, 208, 170}, {184, 75, 208, 202}, {184, 75, 210, 18}, {184, 75, 210, 66}, {184, 75, 210, 106}, {184, 75, 210, 194}, {184, 75, 214, 18}, {184, 75, 215, 18}, {184, 75, 215, 26}, {184, 75, 215, 66}, {184, 75, 215, 74}}},
		{Region: "CA Toronto", IPs: []net.IP{{66, 115, 142, 130}, {66, 115, 145, 199}, {172, 98, 92, 66}, {172, 98, 92, 130}, {172, 98, 92, 194}}},
		{Region: "CA Vancouver", IPs: []net.IP{{162, 216, 47, 66}, {162, 216, 47, 194}, {172, 98, 89, 130}, {172, 98, 89, 194}}},
		{Region: "Czech Republic", IPs: []net.IP{{212, 102, 39, 1}}},
		{Region: "DE Berlin", IPs: []net.IP{{185, 230, 127, 226}, {185, 230, 127, 227}, {185, 230, 127, 228}, {185, 230, 127, 229}, {185, 230, 127, 230}, {185, 230, 127, 231}, {185, 230, 127, 232}, {185, 230, 127, 236}, {185, 230, 127, 237}, {185, 230, 127, 239}, {185, 230, 127, 240}, {185, 230, 127, 241}, {185, 230, 127, 243}, {193, 176, 86, 125}, {193, 176, 86, 138}, {193, 176, 86, 146}, {193, 176, 86, 154}, {193, 176, 86, 158}, {193, 176, 86, 162}, {193, 176, 86, 174}}},
		{Region: "DE Frankfurt", IPs: []net.IP{{195, 181, 170, 225}, {195, 181, 170, 239}, {195, 181, 170, 240}, {195, 181, 170, 241}, {195, 181, 170, 242}, {195, 181, 170, 243}, {195, 181, 170, 244}, {212, 102, 57, 138}}},
		{Region: "Denmark", IPs: []net.IP{{188, 126, 94, 34}}},
		{Region: "Estonia", IPs: []net.IP{{77, 247, 111, 82}, {77, 247, 111, 98}, {77, 247, 111, 114}}},
		{Region: "Finland", IPs: []net.IP{{188, 126, 89, 4}}},
		{Region: "France", IPs: []net.IP{{156, 146, 63, 1}}},
		{Region: "Greece", IPs: []net.IP{{154, 57, 3, 91}, {154, 57, 3, 106}, {154, 57, 3, 145}}},
		{Region: "Hungary", IPs: []net.IP{{185, 128, 26, 18}, {185, 128, 26, 19}, {185, 128, 26, 20}, {185, 128, 26, 21}, {185, 128, 26, 22}, {185, 128, 26, 23}, {185, 128, 26, 24}}},
		{Region: "Iceland", IPs: []net.IP{{45, 133, 193, 50}, {45, 133, 193, 66}}},
		{Region: "India", IPs: []net.IP{{150, 242, 12, 155}, {150, 242, 12, 171}, {150, 242, 12, 187}}},
		{Region: "Ireland", IPs: []net.IP{{23, 92, 127, 2}, {23, 92, 127, 10}, {23, 92, 127, 18}, {23, 92, 127, 34}, {23, 92, 127, 42}, {23, 92, 127, 50}}},
		{Region: "Israel", IPs: []net.IP{{31, 168, 172, 142}, {31, 168, 172, 143}, {31, 168, 172, 145}, {31, 168, 172, 146}}},
		{Region: "Italy", IPs: []net.IP{{156, 146, 41, 129}, {156, 146, 41, 193}}},
		{Region: "Japan", IPs: []net.IP{{156, 146, 34, 1}, {156, 146, 34, 65}}},
		{Region: "Latvia", IPs: []net.IP{{46, 183, 217, 34}, {46, 183, 218, 130}, {46, 183, 218, 146}}},
		{Region: "Lithuania", IPs: []net.IP{{85, 206, 165, 96}, {85, 206, 165, 112}, {85, 206, 165, 128}}},
		{Region: "Luxembourg", IPs: []net.IP{{92, 223, 89, 134}, {92, 223, 89, 135}, {92, 223, 89, 136}, {92, 223, 89, 137}, {92, 223, 89, 138}, {92, 223, 89, 140}}},
		{Region: "Moldova", IPs: []net.IP{{178, 17, 172, 242}, {178, 17, 173, 194}, {178, 175, 128, 34}}},
		{Region: "Netherlands", IPs: []net.IP{{212, 102, 35, 103}}},
		{Region: "New Zealand", IPs: []net.IP{{43, 250, 207, 1}, {43, 250, 207, 3}}},
		{Region: "North Macedonia", IPs: []net.IP{{185, 225, 28, 130}}},
		{Region: "Norway", IPs: []net.IP{{46, 246, 122, 34}, {46, 246, 122, 162}}},
		{Region: "Poland", IPs: []net.IP{{185, 244, 214, 195}, {185, 244, 214, 197}, {185, 244, 214, 198}, {185, 244, 214, 199}}},
		{Region: "Portugal", IPs: []net.IP{{89, 26, 241, 86}, {89, 26, 241, 102}, {89, 26, 241, 130}}},
		{Region: "Romania", IPs: []net.IP{{86, 105, 25, 70}, {86, 105, 25, 75}, {86, 105, 25, 76}, {86, 105, 25, 77}, {94, 176, 148, 35}, {143, 244, 54, 1}, {185, 210, 218, 99}, {185, 210, 218, 101}, {185, 210, 218, 103}, {185, 210, 218, 104}}},
		{Region: "Serbia", IPs: []net.IP{{37, 120, 193, 226}}},
		{Region: "Singapore", IPs: []net.IP{{156, 146, 56, 193}, {156, 146, 57, 38}, {156, 146, 57, 235}, {156, 146, 57, 244}}},
		{Region: "Slovakia", IPs: []net.IP{{37, 120, 221, 98}}},
		{Region: "South Africa", IPs: []net.IP{{102, 165, 20, 133}}},
		{Region: "Spain", IPs: []net.IP{{212, 102, 49, 185}, {212, 102, 49, 251}}},
		{Region: "Sweden", IPs: []net.IP{{46, 246, 3, 253}, {46, 246, 3, 254}}},
		{Region: "Switzerland", IPs: []net.IP{{156, 146, 62, 129}, {156, 146, 62, 193}, {212, 102, 36, 1}, {212, 102, 36, 166}}},
		{Region: "Turkey", IPs: []net.IP{{185, 195, 79, 34}, {185, 195, 79, 82}}},
		{Region: "UAE", IPs: []net.IP{{45, 9, 250, 46}}},
		{Region: "UK London", IPs: []net.IP{{212, 102, 52, 1}, {212, 102, 52, 134}, {212, 102, 53, 129}}},
		{Region: "UK Manchester", IPs: []net.IP{{89, 238, 137, 36}, {89, 238, 137, 37}, {89, 238, 137, 38}, {89, 238, 137, 39}, {89, 238, 139, 52}, {89, 238, 139, 53}, {89, 238, 139, 54}, {89, 238, 139, 55}, {89, 238, 139, 56}, {89, 238, 139, 57}, {89, 238, 139, 58}, {89, 249, 67, 220}}},
		{Region: "UK Southampton", IPs: []net.IP{{31, 24, 226, 141}, {31, 24, 226, 147}, {31, 24, 226, 188}, {31, 24, 226, 189}, {31, 24, 226, 203}, {31, 24, 226, 205}, {31, 24, 226, 206}, {31, 24, 226, 220}, {31, 24, 226, 222}, {31, 24, 226, 223}, {31, 24, 226, 225}, {31, 24, 226, 226}, {31, 24, 226, 228}, {31, 24, 226, 232}, {31, 24, 226, 235}, {31, 24, 226, 244}, {31, 24, 226, 245}, {31, 24, 226, 246}, {31, 24, 226, 252}, {31, 24, 226, 254}}},
		{Region: "US Atlanta", IPs: []net.IP{{66, 115, 169, 195}, {66, 115, 169, 197}, {66, 115, 169, 199}, {66, 115, 169, 203}, {66, 115, 169, 206}, {66, 115, 169, 207}, {66, 115, 169, 208}, {66, 115, 169, 211}, {66, 115, 169, 214}, {156, 146, 46, 1}, {156, 146, 46, 134}, {156, 146, 46, 198}, {156, 146, 47, 11}}},
		{Region: "US California", IPs: []net.IP{{37, 235, 108, 144}, {89, 187, 187, 129}, {89, 187, 187, 162}, {91, 207, 175, 194}, {91, 207, 175, 197}, {91, 207, 175, 198}, {91, 207, 175, 199}, {91, 207, 175, 200}, {91, 207, 175, 203}, {91, 207, 175, 206}, {91, 207, 175, 209}, {91, 207, 175, 210}, {91, 207, 175, 211}}},
		{Region: "US Chicago", IPs: []net.IP{{156, 146, 50, 1}, {156, 146, 50, 65}, {156, 146, 50, 134}, {156, 146, 50, 198}, {156, 146, 51, 11}, {212, 102, 58, 113}, {212, 102, 59, 54}, {212, 102, 59, 129}}},
		{Region: "US Dallas", IPs: []net.IP{{156, 146, 38, 65}, {156, 146, 38, 161}, {156, 146, 39, 1}, {156, 146, 39, 6}, {156, 146, 52, 6}, {156, 146, 52, 70}, {156, 146, 52, 139}, {156, 146, 52, 203}, {174, 127, 114, 53}, {174, 127, 114, 54}, {174, 127, 114, 56}, {174, 127, 114, 65}, {174, 127, 114, 66}, {174, 127, 114, 67}, {174, 127, 114, 71}, {174, 127, 114, 74}, {174, 127, 114, 75}, {174, 127, 114, 76}, {174, 127, 114, 77}, {174, 127, 114, 80}}},
		{Region: "US Denver", IPs: []net.IP{{174, 128, 225, 2}, {174, 128, 225, 98}, {174, 128, 225, 106}, {174, 128, 225, 186}, {174, 128, 236, 98}, {174, 128, 242, 234}, {174, 128, 242, 242}, {174, 128, 242, 250}, {174, 128, 243, 98}, {174, 128, 244, 66}, {174, 128, 244, 74}, {174, 128, 245, 122}, {174, 128, 250, 18}, {174, 128, 250, 26}, {198, 148, 82, 82}, {199, 115, 97, 202}, {199, 115, 98, 234}, {199, 115, 101, 178}, {199, 115, 101, 186}, {199, 115, 102, 146}}},
		{Region: "US East", IPs: []net.IP{{156, 146, 58, 201}, {156, 146, 58, 202}, {156, 146, 58, 203}, {156, 146, 58, 204}, {156, 146, 58, 205}, {156, 146, 58, 206}, {156, 146, 58, 207}, {156, 146, 58, 208}, {156, 146, 58, 209}, {193, 37, 253, 109}, {193, 37, 253, 114}, {193, 37, 253, 117}, {193, 37, 253, 133}, {194, 59, 251, 12}, {194, 59, 251, 24}, {194, 59, 251, 49}, {194, 59, 251, 53}, {194, 59, 251, 80}, {194, 59, 251, 93}, {194, 59, 251, 104}}},
		{Region: "US Florida", IPs: []net.IP{{156, 146, 42, 1}, {156, 146, 42, 65}, {156, 146, 42, 134}, {156, 146, 42, 198}, {156, 146, 43, 11}, {156, 146, 43, 75}, {193, 37, 252, 14}, {193, 37, 252, 15}, {193, 37, 252, 18}, {193, 37, 252, 19}, {193, 37, 252, 20}, {193, 37, 252, 21}, {193, 37, 252, 22}, {193, 37, 252, 24}, {193, 37, 252, 25}, {193, 37, 252, 26}, {193, 37, 252, 27}, {212, 102, 61, 19}, {212, 102, 61, 83}}},
		{Region: "US Houston", IPs: []net.IP{{74, 81, 88, 26}, {74, 81, 88, 42}, {74, 81, 88, 66}, {74, 81, 88, 74}, {205, 251, 148, 66}}},
		{Region: "US Las Vegas", IPs: []net.IP{{162, 251, 236, 2}, {162, 251, 236, 3}, {162, 251, 236, 4}, {162, 251, 236, 5}, {162, 251, 236, 6}, {162, 251, 236, 8}, {162, 251, 236, 9}, {199, 127, 56, 82}, {199, 127, 56, 83}, {199, 127, 56, 84}, {199, 127, 56, 86}, {199, 127, 56, 87}, {199, 127, 56, 88}, {199, 127, 56, 89}, {199, 127, 56, 90}, {199, 127, 56, 91}}},
		{Region: "US New York City", IPs: []net.IP{{156, 146, 36, 225}, {156, 146, 36, 240}, {156, 146, 37, 129}, {156, 146, 55, 198}, {156, 146, 58, 1}, {156, 146, 58, 134}, {173, 244, 217, 37}, {209, 95, 50, 50}, {209, 95, 50, 58}, {209, 95, 50, 60}, {209, 95, 50, 62}, {209, 95, 50, 64}, {209, 95, 50, 65}, {209, 95, 50, 66}, {209, 95, 50, 67}, {209, 95, 50, 68}, {209, 95, 50, 69}, {209, 95, 50, 84}, {209, 95, 50, 85}, {209, 95, 50, 87}}},
		{Region: "US Seattle", IPs: []net.IP{{84, 17, 41, 7}, {84, 17, 41, 10}, {84, 17, 41, 20}, {84, 17, 41, 22}, {84, 17, 41, 25}, {84, 17, 41, 27}, {84, 17, 41, 30}, {84, 17, 41, 38}, {84, 17, 41, 40}, {84, 17, 41, 41}, {84, 17, 41, 50}, {84, 17, 41, 53}, {84, 17, 41, 56}, {84, 17, 41, 58}, {84, 17, 41, 63}, {84, 17, 41, 92}, {84, 17, 41, 93}, {84, 17, 41, 95}, {212, 102, 46, 193}, {212, 102, 47, 134}}},
		{Region: "US Silicon Valley", IPs: []net.IP{{199, 116, 118, 130}, {199, 116, 118, 133}, {199, 116, 118, 136}, {199, 116, 118, 140}, {199, 116, 118, 158}, {199, 116, 118, 170}, {199, 116, 118, 172}, {199, 116, 118, 174}, {199, 116, 118, 178}, {199, 116, 118, 180}, {199, 116, 118, 184}, {199, 116, 118, 202}, {199, 116, 118, 204}, {199, 116, 118, 212}, {199, 116, 118, 219}, {199, 116, 118, 233}, {199, 116, 118, 239}, {199, 116, 118, 240}, {199, 116, 118, 244}, {199, 116, 118, 246}}},
		{Region: "US Washington DC", IPs: []net.IP{{70, 32, 0, 46}, {70, 32, 0, 47}, {70, 32, 0, 52}, {70, 32, 0, 53}, {70, 32, 0, 65}, {70, 32, 0, 68}, {70, 32, 0, 69}, {70, 32, 0, 70}, {70, 32, 0, 71}, {70, 32, 0, 72}, {70, 32, 0, 73}, {70, 32, 0, 75}, {70, 32, 0, 76}, {70, 32, 0, 113}, {70, 32, 0, 114}, {70, 32, 0, 115}, {70, 32, 0, 118}, {70, 32, 0, 119}, {70, 32, 0, 139}, {70, 32, 0, 173}}},
		{Region: "US West", IPs: []net.IP{{104, 200, 151, 6}, {104, 200, 151, 7}, {104, 200, 151, 9}, {104, 200, 151, 10}, {104, 200, 151, 12}, {104, 200, 151, 16}, {104, 200, 151, 17}, {104, 200, 151, 21}, {104, 200, 151, 49}, {104, 200, 151, 51}, {104, 200, 151, 56}, {104, 200, 151, 74}, {104, 200, 151, 78}, {104, 200, 151, 79}, {104, 200, 151, 81}, {104, 200, 151, 82}, {104, 200, 151, 84}, {104, 200, 151, 85}, {104, 200, 151, 87}, {104, 200, 151, 89}}},
		{Region: "Ukraine", IPs: []net.IP{{62, 149, 20, 10}, {62, 149, 20, 40}}},
	}
}

const (
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
