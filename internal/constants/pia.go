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
	PIACertificateStrong      = "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"
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
		{Region: "AU Melbourne", IPs: []net.IP{{43, 250, 204, 97}}},
		{Region: "AU Perth", IPs: []net.IP{{43, 250, 205, 59}}},
		{Region: "AU Sydney", IPs: []net.IP{{221, 121, 146, 203}}},
		{Region: "Albania", IPs: []net.IP{{31, 171, 154, 130}}},
		{Region: "Argentina", IPs: []net.IP{{190, 106, 134, 80}}},
		{Region: "Austria", IPs: []net.IP{{185, 216, 34, 229}}},
		{Region: "Belgium", IPs: []net.IP{{185, 232, 21, 29}}},
		{Region: "Bosnia and Herzegovina", IPs: []net.IP{{185, 164, 35, 55}}},
		{Region: "Bulgaria", IPs: []net.IP{{217, 138, 221, 82}}},
		{Region: "CA Montreal", IPs: []net.IP{{199, 229, 249, 159}}},
		{Region: "CA Ontario", IPs: []net.IP{{184, 75, 213, 218}}},
		{Region: "CA Toronto", IPs: []net.IP{{172, 98, 67, 85}}},
		{Region: "CA Vancouver", IPs: []net.IP{{172, 83, 40, 25}}},
		{Region: "Czech Republic", IPs: []net.IP{{185, 242, 6, 27}}},
		{Region: "DE Berlin", IPs: []net.IP{{193, 176, 86, 123}}},
		{Region: "DE Frankfurt", IPs: []net.IP{{185, 220, 70, 147}}},
		{Region: "Denmark", IPs: []net.IP{{82, 102, 20, 181}}},
		{Region: "Estonia", IPs: []net.IP{{77, 247, 111, 98}}},
		{Region: "Finland", IPs: []net.IP{{196, 244, 191, 146}}},
		{Region: "France", IPs: []net.IP{{194, 187, 249, 47}}},
		{Region: "Greece", IPs: []net.IP{{154, 57, 3, 91}}},
		{Region: "Hungary", IPs: []net.IP{{185, 128, 26, 19}}},
		{Region: "Iceland", IPs: []net.IP{{213, 167, 139, 66}}},
		{Region: "India", IPs: []net.IP{{150, 242, 12, 155}}},
		{Region: "Ireland", IPs: []net.IP{{23, 92, 127, 34}}},
		{Region: "Israel", IPs: []net.IP{{31, 168, 172, 145}}},
		{Region: "Italy", IPs: []net.IP{{82, 102, 21, 217}}},
		{Region: "Japan", IPs: []net.IP{{156, 146, 34, 65}}},
		{Region: "Latvia", IPs: []net.IP{{109, 248, 149, 2}}},
		{Region: "Lithuania", IPs: []net.IP{{85, 206, 165, 160}}},
		{Region: "Luxembourg", IPs: []net.IP{{92, 223, 89, 137}}},
		{Region: "Moldova", IPs: []net.IP{{178, 17, 172, 242}}},
		{Region: "Netherlands", IPs: []net.IP{{46, 166, 190, 227}}},
		{Region: "New Zealand", IPs: []net.IP{{43, 250, 207, 3}}},
		{Region: "North Macedonia", IPs: []net.IP{{185, 225, 28, 130}}},
		{Region: "Norway", IPs: []net.IP{{82, 102, 27, 52}}},
		{Region: "Poland", IPs: []net.IP{{185, 244, 214, 198}}},
		{Region: "Portugal", IPs: []net.IP{{89, 26, 241, 102}}},
		{Region: "Romania", IPs: []net.IP{{185, 210, 218, 98}}},
		{Region: "Serbia", IPs: []net.IP{{37, 120, 193, 242}}},
		{Region: "Singapore", IPs: []net.IP{{37, 120, 208, 82}}},
		{Region: "Slovakia", IPs: []net.IP{{37, 120, 221, 82}}},
		{Region: "South Africa", IPs: []net.IP{{102, 165, 20, 133}}},
		{Region: "Spain", IPs: []net.IP{{185, 230, 124, 52}}},
		{Region: "Sweden", IPs: []net.IP{{45, 12, 220, 170}}},
		{Region: "Switzerland", IPs: []net.IP{{91, 132, 136, 45}}},
		{Region: "Turkey", IPs: []net.IP{{185, 195, 79, 82}}},
		{Region: "UAE", IPs: []net.IP{{45, 9, 250, 46}}},
		{Region: "UK London", IPs: []net.IP{{89, 238, 154, 229}}},
		{Region: "UK Manchester", IPs: []net.IP{{89, 238, 139, 9}}},
		{Region: "UK Southampton", IPs: []net.IP{{31, 24, 226, 234}}},
		{Region: "US Atlanta", IPs: []net.IP{{156, 146, 47, 11}}},
		{Region: "US California", IPs: []net.IP{{91, 207, 175, 169}}},
		{Region: "US Chicago", IPs: []net.IP{{212, 102, 58, 113}}},
		{Region: "US Dallas", IPs: []net.IP{{156, 146, 52, 70}}},
		{Region: "US Denver", IPs: []net.IP{{174, 128, 242, 250}}},
		{Region: "US East", IPs: []net.IP{{193, 37, 253, 120}}},
		{Region: "US Florida", IPs: []net.IP{{193, 37, 252, 126}}},
		{Region: "US Houston", IPs: []net.IP{{205, 251, 150, 194}}},
		{Region: "US Las Vegas", IPs: []net.IP{{199, 127, 56, 119}}},
		{Region: "US New York City", IPs: []net.IP{{156, 146, 54, 53}}},
		{Region: "US Seattle", IPs: []net.IP{{156, 146, 48, 65}}},
		{Region: "US Silicon Valley", IPs: []net.IP{{199, 116, 118, 181}}},
		{Region: "US Washington DC", IPs: []net.IP{{70, 32, 0, 75}}},
		{Region: "US West", IPs: []net.IP{{104, 200, 151, 44}}},
		{Region: "Ukraine", IPs: []net.IP{{62, 149, 20, 50}}},
	}
}

const (
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
