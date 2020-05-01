package constants

import (
	"net"

	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// PIAEncryptionNormal is the normal level of encryption for communication with PIA servers
	PIAEncryptionNormal models.PIAEncryption = "normal"
	// PIAEncryptionStrong is the strong level of encryption for communication with PIA servers
	PIAEncryptionStrong models.PIAEncryption = "strong"
)

const (
	PiaX509CRLNormal     = "MIICWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAQEAQZo9X97ci8EcPYu/uK2HB152OZbeZCINmYyluLDOdcSvg6B5jI+ffKN3laDvczsG6CxmY3jNyc79XVpEYUnq4rT3FfveW1+Ralf+Vf38HdpwB8EWB4hZlQ205+21CALLvZvR8HcPxC9KEnev1mU46wkTiov0EKc+EdRxkj5yMgv0V2Reze7AP+NQ9ykvDScH4eYCsmufNpIjBLhpLE2cuZZXBLcPhuRzVoU3l7A9lvzG9mjA5YijHJGHNjlWFqyrn1CfYS6koa4TGEPngBoAziWRbDGdhEgJABHrpoaFYaL61zqyMR6jC0K2ps9qyZAN74LEBedEfK7tBOzWMwr58A=="
	PiaX509CRLStrong     = "MIIDWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAgEAppFfEpGsasjB1QgJcosGpzbf2kfRhM84o2TlqY1ua+Gi5TMdKydA3LJcNTjlI9a0TYAJfeRX5IkpoglSUuHuJgXhP3nEvX10mjXDpcu/YvM8TdE5JV2+EGqZ80kFtBeOq94WcpiVKFTR4fO+VkOK9zwspFfb1cNs9rHvgJ1QMkRUF8PpLN6AkntHY0+6DnigtSaKqldqjKTDTv2OeH3nPoh80SGrt0oCOmYKfWTJGpggMGKvIdvU3vH9+EuILZKKIskt+1dwdfA5Bkz1GLmiQG7+9ZZBQUjBG9Dos4hfX/rwJ3eU8oUIm4WoTz9rb71SOEuUUjP5NPy9HNx2vx+cVvLsTF4ZDZaUztW9o9JmIURDtbeyqxuHN3prlPWB6aj73IIm2dsDQvs3XXwRIxs8NwLbJ6CyEuvEOVCskdM8rdADWx1J0lRNlOJ0Z8ieLLEmYAA834VN1SboB6wJIAPxQU3rcBhXqO9y8aa2oRMg8NxZ5gr+PnKVMqag1x0IxbIgLxtkXQvxXxQHEMSODzvcOfK/nBRBsqTj30P+R87sU8titOoxNeRnBDRNhdEy/QGAqGh62ShPpQUCJdnKRiRTjnil9hMQHevoSuFKeEMO30FQL7BZyo37GFU+q1WPCplVZgCP9hC8Rn5K2+f6KLFo5bhtowSmu+GY1yZtg+RTtsA="
	PIACertificateNormal = "MIIFqzCCBJOgAwIBAgIJAKZ7D5Yv87qDMA0GCSqGSIb3DQEBDQUAMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTAeFw0xNDA0MTcxNzM1MThaFw0zNDA0MTIxNzM1MThaMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPXDL1L9tX6DGf36liA7UBTy5I869z0UVo3lImfOs/GSiFKPtInlesP65577nd7UNzzXlH/P/CnFPdBWlLp5ze3HRBCc/Avgr5CdMRkEsySL5GHBZsx6w2cayQ2EcRhVTwWpcdldeNO+pPr9rIgPrtXqT4SWViTQRBeGM8CDxAyTopTsobjSiYZCF9Ta1gunl0G/8Vfp+SXfYCC+ZzWvP+L1pFhPRqzQQ8k+wMZIovObK1s+nlwPaLyayzw9a8sUnvWB/5rGPdIYnQWPgoNlLN9HpSmsAcw2z8DXI9pIxbr74cb3/HSfuYGOLkRqrOk6h4RCOfuWoTrZup1uEOn+fw8CAwEAAaOCAVQwggFQMB0GA1UdDgQWBBQv63nQ/pJAt5tLy8VJcbHe22ZOsjCCAR8GA1UdIwSCARYwggESgBQv63nQ/pJAt5tLy8VJcbHe22ZOsqGB7qSB6zCB6DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRMwEQYDVQQHEwpMb3NBbmdlbGVzMSAwHgYDVQQKExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UECxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAMTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQpExdQcml2YXRlIEludGVybmV0IEFjY2VzczEvMC0GCSqGSIb3DQEJARYgc2VjdXJlQHByaXZhdGVpbnRlcm5ldGFjY2Vzcy5jb22CCQCmew+WL/O6gzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBDQUAA4IBAQAna5PgrtxfwTumD4+3/SYvwoD66cB8IcK//h1mCzAduU8KgUXocLx7QgJWo9lnZ8xUryXvWab2usg4fqk7FPi00bED4f4qVQFVfGfPZIH9QQ7/48bPM9RyfzImZWUCenK37pdw4Bvgoys2rHLHbGen7f28knT2j/cbMxd78tQc20TIObGjo8+ISTRclSTRBtyCGohseKYpTS9himFERpUgNtefvYHbn70mIOzfOJFTVqfrptf9jXa9N8Mpy3ayfodz1wiqdteqFXkTYoSDctgKMiZ6GdocK9nMroQipIQtpnwd4yBDWIyC6Bvlkrq5TQUtYDQ8z9v+DMO6iwyIDRiU"
	PIACertificateStrong = "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"
)

func PIAGeoChoices() (regions []string) {
	for _, server := range PIAServers() {
		regions = append(regions, string(server.Region))
	}
	return regions
}

func PIAServers() []models.PIAServer {
	return []models.PIAServer{
		{
			Region: models.PIARegion("AU Melbourne"),
			IPs:    []net.IP{{43, 250, 204, 83}, {43, 250, 204, 85}, {43, 250, 204, 89}, {43, 250, 204, 105}, {43, 250, 204, 107}, {43, 250, 204, 111}, {43, 250, 204, 117}, {43, 250, 204, 119}, {43, 250, 204, 123}, {43, 250, 204, 125}, {43, 250, 204, 125}, {43, 250, 204, 85}, {43, 250, 204, 89}, {43, 250, 204, 105}, {43, 250, 204, 107}, {43, 250, 204, 111}, {43, 250, 204, 117}, {43, 250, 204, 119}, {43, 250, 204, 123}, {43, 250, 204, 83}},
		},
		{
			Region: models.PIARegion("AU Perth"),
			IPs:    []net.IP{{103, 231, 89, 2}, {103, 231, 89, 3}, {103, 231, 89, 11}, {103, 231, 89, 12}, {103, 231, 89, 13}, {103, 231, 89, 11}, {103, 231, 89, 13}, {103, 231, 89, 3}, {103, 231, 89, 12}, {103, 231, 89, 2}},
		},
		{
			Region: models.PIARegion("AU Sydney"),
			IPs:    []net.IP{{217, 138, 205, 99}, {217, 138, 205, 108}, {217, 138, 205, 109}, {217, 138, 205, 114}, {217, 138, 205, 118}, {217, 138, 205, 120}, {217, 138, 205, 195}, {217, 138, 205, 213}, {217, 138, 205, 206}, {217, 138, 205, 202}, {217, 138, 205, 99}, {217, 138, 205, 108}, {217, 138, 205, 109}, {217, 138, 205, 114}, {217, 138, 205, 118}, {217, 138, 205, 206}, {217, 138, 205, 195}, {217, 138, 205, 213}, {217, 138, 205, 120}, {217, 138, 205, 202}},
		},
		{
			Region: models.PIARegion("Austria"),
			IPs:    []net.IP{{185, 216, 34, 226}, {185, 216, 34, 227}, {185, 216, 34, 228}, {185, 216, 34, 229}, {185, 216, 34, 230}, {185, 216, 34, 231}, {185, 216, 34, 232}, {185, 216, 34, 236}, {185, 216, 34, 237}, {185, 216, 34, 238}},
		},
		{
			Region: models.PIARegion("Belgium"),
			IPs:    []net.IP{{77, 243, 191, 19}, {77, 243, 191, 20}, {77, 243, 191, 21}, {77, 243, 191, 22}, {77, 243, 191, 26}, {77, 243, 191, 27}, {185, 104, 186, 26}, {185, 232, 21, 26}, {185, 232, 21, 28}, {185, 232, 21, 29}},
		},
		{
			Region: models.PIARegion("CA Montreal"),
			IPs:    []net.IP{{199, 229, 249, 132}, {199, 229, 249, 141}, {199, 229, 249, 143}, {199, 229, 249, 144}, {199, 229, 249, 152}, {199, 229, 249, 158}, {199, 229, 249, 166}, {199, 229, 249, 189}, {199, 229, 249, 190}, {199, 229, 249, 196}},
		},
		{
			Region: models.PIARegion("CA Toronto"),
			IPs:    []net.IP{{162, 219, 176, 130}, {172, 98, 67, 61}, {172, 98, 67, 74}, {172, 98, 67, 89}, {172, 98, 67, 108}, {184, 75, 208, 114}, {184, 75, 210, 18}, {184, 75, 210, 98}, {184, 75, 210, 106}, {184, 75, 215, 66}, {162, 219, 176, 130}, {184, 75, 215, 66}, {172, 98, 67, 74}, {172, 98, 67, 89}, {184, 75, 210, 98}, {184, 75, 208, 114}, {184, 75, 210, 18}, {172, 98, 67, 61}, {184, 75, 210, 106}, {172, 98, 67, 108}},
		},
		{
			Region: models.PIARegion("CA Vancouver"),
			IPs:    []net.IP{{107, 181, 189, 82}, {172, 83, 40, 113}, {172, 83, 40, 114}, {107, 181, 189, 82}, {172, 83, 40, 114}, {172, 83, 40, 113}},
		},
		{
			Region: models.PIARegion("Czech Republic"),
			IPs:    []net.IP{{89, 238, 186, 226}, {89, 238, 186, 227}, {89, 238, 186, 228}, {89, 238, 186, 230}, {185, 216, 35, 66}, {185, 216, 35, 67}, {185, 242, 6, 30}, {185, 216, 35, 70}, {185, 242, 6, 29}, {185, 216, 35, 68}},
		},
		{
			Region: models.PIARegion("DE Berlin"),
			IPs:    []net.IP{{185, 230, 127, 228}, {185, 230, 127, 229}, {185, 230, 127, 234}, {185, 230, 127, 235}, {193, 176, 86, 123}, {193, 176, 86, 125}, {193, 176, 86, 130}, {193, 176, 86, 142}, {193, 176, 86, 150}, {193, 176, 86, 182}},
		},
		{
			Region: models.PIARegion("DE Frankfurt"),
			IPs:    []net.IP{{185, 220, 70, 131}, {185, 220, 70, 134}, {185, 220, 70, 142}, {185, 220, 70, 144}, {185, 220, 70, 149}, {185, 220, 70, 151}, {185, 220, 70, 152}, {185, 220, 70, 155}, {185, 220, 70, 162}, {185, 220, 70, 167}},
		},
		{
			Region: models.PIARegion("Denmark"),
			IPs:    []net.IP{{82, 102, 20, 164}, {82, 102, 20, 167}, {82, 102, 20, 168}, {82, 102, 20, 173}, {82, 102, 20, 175}, {82, 102, 20, 176}, {82, 102, 20, 177}, {82, 102, 20, 178}, {82, 102, 20, 181}, {82, 102, 20, 182}},
		},
		{
			Region: models.PIARegion("Finland"),
			IPs:    []net.IP{{196, 244, 191, 2}, {196, 244, 191, 10}, {196, 244, 191, 26}, {196, 244, 191, 50}, {196, 244, 191, 66}, {196, 244, 191, 98}, {196, 244, 191, 106}, {196, 244, 191, 114}, {196, 244, 191, 122}, {196, 244, 191, 146}},
		},
		{
			Region: models.PIARegion("France"),
			IPs:    []net.IP{{194, 99, 106, 150}, {194, 187, 249, 34}, {194, 187, 249, 38}, {194, 187, 249, 44}, {194, 187, 249, 54}, {194, 187, 249, 55}, {194, 187, 249, 180}, {194, 187, 249, 183}, {194, 187, 249, 190}, {194, 187, 249, 184}},
		},
		{
			Region: models.PIARegion("Hong Kong"),
			IPs:    []net.IP{{84, 17, 37, 1}, {84, 17, 37, 23}, {84, 17, 37, 45}, {119, 81, 135, 2}, {119, 81, 135, 29}, {119, 81, 135, 47}, {119, 81, 253, 229}, {119, 81, 253, 230}, {119, 81, 253, 241}, {119, 81, 253, 242}},
		},
		{
			Region: models.PIARegion("Hungary"),
			IPs:    []net.IP{{185, 128, 26, 18}, {185, 128, 26, 19}, {185, 128, 26, 20}, {185, 128, 26, 21}, {185, 128, 26, 22}, {185, 128, 26, 23}, {185, 128, 26, 24}, {185, 189, 114, 98}},
		},
		{
			Region: models.PIARegion("India"),
			IPs:    []net.IP{{150, 242, 12, 155}, {150, 242, 12, 171}, {150, 242, 12, 187}, {150, 242, 12, 155}, {150, 242, 12, 171}, {150, 242, 12, 187}},
		},
		{
			Region: models.PIARegion("Ireland"),
			IPs:    []net.IP{{23, 92, 127, 2}, {23, 92, 127, 18}, {23, 92, 127, 34}, {23, 92, 127, 42}, {23, 92, 127, 58}, {23, 92, 127, 66}, {23, 92, 127, 2}, {23, 92, 127, 18}, {23, 92, 127, 34}, {23, 92, 127, 66}, {23, 92, 127, 58}, {23, 92, 127, 42}},
		},
		{
			Region: models.PIARegion("Israel"),
			IPs:    []net.IP{{31, 168, 172, 136}, {31, 168, 172, 142}, {31, 168, 172, 143}, {31, 168, 172, 145}, {31, 168, 172, 146}, {31, 168, 172, 147}},
		},
		{
			Region: models.PIARegion("Italy"),
			IPs:    []net.IP{{82, 102, 21, 98}, {82, 102, 21, 210}, {82, 102, 21, 211}, {82, 102, 21, 212}, {82, 102, 21, 213}, {82, 102, 21, 214}, {82, 102, 21, 215}, {82, 102, 21, 216}, {82, 102, 21, 217}, {82, 102, 21, 218}, {82, 102, 21, 98}, {82, 102, 21, 210}, {82, 102, 21, 211}, {82, 102, 21, 212}, {82, 102, 21, 213}, {82, 102, 21, 214}, {82, 102, 21, 215}, {82, 102, 21, 216}, {82, 102, 21, 217}, {82, 102, 21, 218}},
		},
		{
			Region: models.PIARegion("Japan"),
			IPs:    []net.IP{{103, 208, 220, 130}, {103, 208, 220, 143}, {103, 208, 220, 133}, {103, 208, 220, 134}, {103, 208, 220, 136}, {103, 208, 220, 137}, {103, 208, 220, 138}, {103, 208, 220, 140}, {103, 208, 220, 141}, {103, 208, 220, 131}, {103, 208, 220, 130}, {103, 208, 220, 143}, {103, 208, 220, 133}, {103, 208, 220, 134}, {103, 208, 220, 131}, {103, 208, 220, 137}, {103, 208, 220, 136}, {103, 208, 220, 138}, {103, 208, 220, 141}, {103, 208, 220, 140}},
		},
		{
			Region: models.PIARegion("Luxembourg"),
			IPs:    []net.IP{{92, 223, 89, 133}, {92, 223, 89, 134}, {92, 223, 89, 135}, {92, 223, 89, 136}, {92, 223, 89, 137}, {92, 223, 89, 138}, {92, 223, 89, 140}, {92, 223, 89, 142}},
		},
		{
			Region: models.PIARegion("Mexico"),
			IPs:    []net.IP{{169, 57, 0, 211}, {169, 57, 0, 216}, {169, 57, 0, 218}, {169, 57, 0, 219}, {169, 57, 0, 221}, {169, 57, 0, 225}, {169, 57, 0, 229}, {169, 57, 0, 231}, {169, 57, 0, 233}, {169, 57, 0, 249}},
		},
		{
			Region: models.PIARegion("Netherlands"),
			IPs:    []net.IP{{46, 166, 137, 218}, {46, 166, 138, 138}, {46, 166, 138, 161}, {46, 166, 138, 172}, {46, 166, 188, 211}, {46, 166, 188, 218}, {46, 166, 190, 178}, {46, 166, 190, 223}, {46, 166, 190, 227}, {46, 166, 190, 230}},
		},
		{
			Region: models.PIARegion("New Zealand"),
			IPs:    []net.IP{{103, 231, 90, 171}, {103, 231, 90, 172}, {103, 231, 90, 173}, {103, 231, 91, 34}, {103, 231, 91, 35}, {103, 231, 91, 66}, {103, 231, 91, 67}, {103, 231, 91, 68}, {103, 231, 91, 69}, {103, 231, 91, 74}},
		},
		{
			Region: models.PIARegion("Norway"),
			IPs:    []net.IP{{82, 102, 27, 50}, {82, 102, 27, 55}, {82, 102, 27, 74}, {82, 102, 27, 76}, {82, 102, 27, 77}, {82, 102, 27, 78}, {82, 102, 27, 124}, {82, 102, 27, 125}, {82, 102, 27, 126}, {185, 253, 97, 228}},
		},
		{
			Region: models.PIARegion("Poland"),
			IPs:    []net.IP{{185, 244, 214, 14}, {185, 244, 214, 194}, {185, 244, 214, 195}, {185, 244, 214, 197}, {185, 244, 214, 198}, {185, 244, 214, 199}, {185, 244, 214, 200}, {185, 244, 214, 200}, {185, 244, 214, 194}, {185, 244, 214, 198}, {185, 244, 214, 195}, {185, 244, 214, 14}, {185, 244, 214, 199}, {185, 244, 214, 197}},
		},
		{
			Region: models.PIARegion("Romania"),
			IPs:    []net.IP{{86, 105, 25, 68}, {86, 105, 25, 69}, {89, 33, 8, 42}, {185, 210, 218, 98}, {185, 210, 218, 99}, {185, 210, 218, 101}, {185, 210, 218, 102}, {185, 210, 218, 104}, {185, 210, 218, 105}, {185, 210, 218, 108}},
		},
		{
			Region: models.PIARegion("Singapore"),
			IPs:    []net.IP{{37, 120, 208, 67}, {37, 120, 208, 69}, {37, 120, 208, 71}, {37, 120, 208, 72}, {37, 120, 208, 73}, {37, 120, 208, 74}, {37, 120, 208, 75}, {37, 120, 208, 76}, {37, 120, 208, 79}, {37, 120, 208, 80}, {37, 120, 208, 69}, {37, 120, 208, 80}, {37, 120, 208, 71}, {37, 120, 208, 72}, {37, 120, 208, 73}, {37, 120, 208, 74}, {37, 120, 208, 75}, {37, 120, 208, 76}, {37, 120, 208, 79}, {37, 120, 208, 67}},
		},
		{
			Region: models.PIARegion("Spain"),
			IPs:    []net.IP{{37, 120, 148, 86}, {194, 99, 104, 30}, {185, 230, 124, 51}, {185, 230, 124, 52}, {185, 230, 124, 53}, {185, 230, 124, 50}},
		},
		{
			Region: models.PIARegion("Sweden"),
			IPs:    []net.IP{{45, 12, 220, 175}, {45, 12, 220, 194}, {45, 12, 220, 204}, {45, 12, 220, 206}, {45, 83, 91, 35}, {45, 12, 220, 210}, {45, 12, 220, 228}, {45, 12, 220, 238}, {45, 12, 220, 239}, {45, 12, 220, 207}, {45, 12, 220, 175}, {45, 12, 220, 194}, {45, 12, 220, 204}, {45, 12, 220, 206}, {45, 83, 91, 35}, {45, 12, 220, 210}, {45, 12, 220, 228}, {45, 12, 220, 238}, {45, 12, 220, 239}, {45, 12, 220, 207}},
		},
		{
			Region: models.PIARegion("Switzerland"),
			IPs:    []net.IP{{82, 102, 24, 167}, {185, 156, 175, 85}, {185, 156, 175, 87}, {185, 212, 170, 179}, {185, 212, 170, 180}, {185, 212, 170, 182}, {185, 230, 125, 34}, {185, 230, 125, 42}, {195, 206, 105, 211}, {212, 102, 36, 1}},
		},
		{
			Region: models.PIARegion("UAE"),
			IPs:    []net.IP{{45, 9, 250, 42}, {45, 9, 250, 46}, {45, 9, 250, 62}, {45, 9, 250, 42}, {45, 9, 250, 46}, {45, 9, 250, 62}},
		},
		{
			Region: models.PIARegion("UK London"),
			IPs:    []net.IP{{89, 238, 150, 8}, {89, 238, 150, 24}, {89, 238, 154, 18}, {89, 238, 154, 119}, {89, 238, 154, 121}, {89, 238, 154, 168}, {89, 238, 154, 178}, {89, 238, 154, 179}, {89, 238, 154, 228}, {89, 238, 154, 244}, {89, 238, 154, 244}, {89, 238, 150, 24}, {89, 238, 154, 18}, {89, 238, 154, 119}, {89, 238, 154, 121}, {89, 238, 154, 168}, {89, 238, 154, 178}, {89, 238, 154, 179}, {89, 238, 154, 228}, {89, 238, 150, 8}},
		},
		{
			Region: models.PIARegion("UK Manchester"),
			IPs:    []net.IP{{89, 238, 137, 38}, {89, 238, 139, 4}, {89, 238, 139, 5}, {89, 238, 139, 8}, {89, 238, 139, 10}, {89, 238, 139, 12}, {89, 238, 139, 13}, {89, 238, 139, 52}, {89, 238, 139, 56}, {89, 238, 139, 58}},
		},
		{
			Region: models.PIARegion("UK Southampton"),
			IPs:    []net.IP{{31, 24, 226, 132}, {31, 24, 226, 138}, {31, 24, 226, 207}, {31, 24, 226, 219}, {31, 24, 226, 220}, {31, 24, 226, 223}, {31, 24, 226, 227}, {31, 24, 226, 241}, {31, 24, 231, 208}, {88, 202, 231, 118}},
		},
		{
			Region: models.PIARegion("US Atlanta"),
			IPs:    []net.IP{{66, 115, 168, 6}, {66, 115, 168, 14}, {66, 115, 168, 17}, {66, 115, 168, 23}, {66, 115, 168, 25}, {66, 115, 169, 199}, {66, 115, 169, 208}, {66, 115, 169, 217}, {66, 115, 169, 219}, {66, 115, 169, 230}, {66, 115, 168, 6}, {66, 115, 168, 14}, {66, 115, 168, 17}, {66, 115, 169, 230}, {66, 115, 168, 25}, {66, 115, 169, 199}, {66, 115, 169, 208}, {66, 115, 169, 217}, {66, 115, 169, 219}, {66, 115, 168, 23}},
		},
		{
			Region: models.PIARegion("US California"),
			IPs:    []net.IP{{91, 207, 175, 37}, {91, 207, 175, 60}, {91, 207, 175, 86}, {91, 207, 175, 109}, {185, 245, 87, 198}, {91, 207, 175, 167}, {91, 207, 175, 176}, {185, 245, 87, 171}, {91, 207, 175, 117}, {185, 245, 87, 195}, {91, 207, 175, 37}, {91, 207, 175, 60}, {91, 207, 175, 86}, {91, 207, 175, 109}, {91, 207, 175, 117}, {91, 207, 175, 167}, {91, 207, 175, 176}, {185, 245, 87, 171}, {185, 245, 87, 195}, {185, 245, 87, 198}},
		},
		{
			Region: models.PIARegion("US Chicago"),
			IPs:    []net.IP{{199, 116, 115, 130}, {199, 116, 115, 131}, {199, 116, 115, 147}, {199, 116, 115, 134}, {199, 116, 115, 135}, {199, 116, 115, 136}, {199, 116, 115, 138}, {199, 116, 115, 141}, {199, 116, 115, 144}, {199, 116, 115, 133}},
		},
		{
			Region: models.PIARegion("US Denver"),
			IPs:    []net.IP{{174, 128, 226, 10}, {174, 128, 226, 18}, {174, 128, 243, 106}, {174, 128, 243, 114}, {199, 115, 101, 186}, {198, 148, 88, 250}, {199, 115, 97, 202}, {199, 115, 98, 146}, {199, 115, 99, 218}, {174, 128, 250, 26}, {199, 115, 97, 202}, {174, 128, 226, 18}, {174, 128, 243, 106}, {174, 128, 243, 114}, {198, 148, 88, 250}, {174, 128, 226, 10}, {199, 115, 101, 186}, {199, 115, 98, 146}, {199, 115, 99, 218}, {174, 128, 250, 26}},
		},
		{
			Region: models.PIARegion("US East"),
			IPs:    []net.IP{{193, 37, 253, 82}, {193, 37, 253, 113}, {194, 59, 251, 22}, {194, 59, 251, 90}, {194, 59, 251, 109}, {194, 59, 251, 140}, {194, 59, 251, 148}, {194, 59, 251, 155}, {194, 59, 251, 218}, {194, 59, 251, 249}, {194, 59, 251, 101}, {194, 59, 251, 152}, {194, 59, 251, 156}, {194, 59, 251, 187}, {193, 37, 253, 67}},
		},
		{
			Region: models.PIARegion("US Florida"),
			IPs:    []net.IP{{193, 37, 252, 2}, {193, 37, 252, 3}, {193, 37, 252, 4}, {193, 37, 252, 5}, {193, 37, 252, 6}, {193, 37, 252, 7}, {193, 37, 252, 8}, {193, 37, 252, 9}, {193, 37, 252, 10}, {193, 37, 252, 11}, {193, 37, 252, 12}, {193, 37, 252, 13}, {193, 37, 252, 14}, {193, 37, 252, 15}, {193, 37, 252, 16}, {193, 37, 252, 17}, {193, 37, 252, 18}, {193, 37, 252, 19}, {193, 37, 252, 20}, {193, 37, 252, 21}, {193, 37, 252, 22}, {193, 37, 252, 23}, {193, 37, 252, 24}, {193, 37, 252, 25}, {193, 37, 252, 34}, {193, 37, 252, 35}, {193, 37, 252, 36}, {193, 37, 252, 37}, {193, 37, 252, 38}, {193, 37, 252, 39}, {193, 37, 252, 40}, {193, 37, 252, 41}, {193, 37, 252, 42}, {193, 37, 252, 43}, {193, 37, 252, 44}, {193, 37, 252, 45}, {193, 37, 252, 46}, {193, 37, 252, 47}, {193, 37, 252, 48}, {193, 37, 252, 49}, {193, 37, 252, 50}, {193, 37, 252, 51}, {193, 37, 252, 52}, {193, 37, 252, 53}, {193, 37, 252, 54}, {193, 37, 252, 55}, {193, 37, 252, 56}, {193, 37, 252, 57}, {193, 37, 252, 58}, {193, 37, 252, 59}, {193, 37, 252, 60}, {193, 37, 252, 61}, {193, 37, 252, 62}, {193, 37, 252, 66}, {193, 37, 252, 67}, {193, 37, 252, 174}, {193, 37, 252, 69}, {193, 37, 252, 70}, {193, 37, 252, 74}, {193, 37, 252, 75}, {193, 37, 252, 76}, {193, 37, 252, 77}, {193, 37, 252, 78}, {193, 37, 252, 82}, {193, 37, 252, 86}, {193, 37, 252, 98}, {193, 37, 252, 99}, {193, 37, 252, 100}, {193, 37, 252, 101}, {193, 37, 252, 102}, {193, 37, 252, 103}, {193, 37, 252, 104}, {193, 37, 252, 105}, {193, 37, 252, 106}, {193, 37, 252, 107}, {193, 37, 252, 108}, {193, 37, 252, 109}, {193, 37, 252, 110}, {193, 37, 252, 111}, {193, 37, 252, 112}, {193, 37, 252, 113}, {193, 37, 252, 114}, {193, 37, 252, 115}, {193, 37, 252, 116}, {193, 37, 252, 117}, {193, 37, 252, 118}, {193, 37, 252, 119}, {193, 37, 252, 120}, {193, 37, 252, 121}, {193, 37, 252, 122}, {193, 37, 252, 123}, {193, 37, 252, 124}, {193, 37, 252, 125}, {193, 37, 252, 126}, {193, 37, 252, 170}, {193, 37, 252, 68}},
		},
		{
			Region: models.PIARegion("US Houston"),
			IPs:    []net.IP{{74, 81, 88, 18}, {74, 81, 88, 26}, {74, 81, 88, 34}, {74, 81, 88, 58}, {74, 81, 88, 74}, {74, 81, 88, 130}, {205, 251, 148, 82}, {205, 251, 148, 186}, {205, 251, 150, 194}, {205, 251, 150, 242}},
		},
		{
			Region: models.PIARegion("US Las Vegas"),
			IPs:    []net.IP{{162, 251, 236, 3}, {162, 251, 236, 4}, {162, 251, 236, 5}, {162, 251, 236, 7}, {199, 127, 56, 83}, {199, 127, 56, 84}, {199, 127, 56, 86}, {199, 127, 56, 89}, {199, 127, 56, 90}, {199, 127, 56, 116}},
		},
		{
			Region: models.PIARegion("US New York City"),
			IPs:    []net.IP{{107, 182, 231, 23}, {107, 182, 231, 37}, {173, 244, 223, 122}, {209, 95, 50, 15}, {209, 95, 50, 47}, {209, 95, 50, 48}, {209, 95, 50, 50}, {209, 95, 50, 53}, {209, 95, 50, 138}, {209, 95, 50, 158}, {209, 95, 50, 158}, {107, 182, 231, 23}, {173, 244, 223, 122}, {209, 95, 50, 47}, {209, 95, 50, 50}, {209, 95, 50, 48}, {209, 95, 50, 53}, {107, 182, 231, 37}, {209, 95, 50, 138}, {209, 95, 50, 15}},
		},
		{
			Region: models.PIARegion("US Seattle"),
			IPs:    []net.IP{{104, 200, 154, 4}, {104, 200, 154, 7}, {104, 200, 154, 11}, {104, 200, 154, 38}, {104, 200, 154, 55}, {104, 200, 154, 58}, {104, 200, 154, 67}, {104, 200, 154, 74}, {104, 200, 154, 77}, {104, 200, 154, 99}, {104, 200, 154, 74}, {104, 200, 154, 7}, {104, 200, 154, 11}, {104, 200, 154, 4}, {104, 200, 154, 55}, {104, 200, 154, 58}, {104, 200, 154, 38}, {104, 200, 154, 99}, {104, 200, 154, 77}, {104, 200, 154, 67}},
		},
		{
			Region: models.PIARegion("US Silicon Valley"),
			IPs:    []net.IP{{199, 116, 118, 131}, {199, 116, 118, 141}, {199, 116, 118, 170}, {199, 116, 118, 183}, {199, 116, 118, 221}, {199, 116, 118, 227}, {199, 116, 118, 228}, {199, 116, 118, 229}, {199, 116, 118, 238}, {199, 116, 118, 246}},
		},
		{
			Region: models.PIARegion("US Texas"),
			IPs:    []net.IP{{162, 216, 46, 15}, {162, 216, 46, 30}, {162, 216, 46, 48}, {162, 216, 46, 55}, {162, 216, 46, 68}, {162, 216, 46, 95}, {162, 216, 46, 116}, {162, 216, 46, 142}, {162, 216, 46, 166}, {162, 216, 46, 174}},
		},
		{
			Region: models.PIARegion("US Washington DC"),
			IPs:    []net.IP{{70, 32, 0, 31}, {70, 32, 0, 50}, {70, 32, 0, 53}, {70, 32, 0, 54}, {70, 32, 0, 66}, {70, 32, 0, 101}, {70, 32, 0, 141}, {70, 32, 0, 172}, {70, 32, 0, 173}, {70, 32, 0, 175}, {70, 32, 0, 31}, {70, 32, 0, 50}, {70, 32, 0, 53}, {70, 32, 0, 175}, {70, 32, 0, 66}, {70, 32, 0, 173}, {70, 32, 0, 141}, {70, 32, 0, 172}, {70, 32, 0, 54}, {70, 32, 0, 101}},
		},
		{
			Region: models.PIARegion("US West"),
			IPs:    []net.IP{{104, 200, 151, 6}, {104, 200, 151, 8}, {104, 200, 151, 12}, {104, 200, 151, 23}, {104, 200, 151, 40}, {104, 200, 151, 44}, {104, 200, 151, 45}, {104, 200, 151, 54}, {104, 200, 151, 74}, {104, 200, 151, 89}},
		},
	}
}

const (
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
