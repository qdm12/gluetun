package constants

import (
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// PIAEncryptionNormal is the normal level of encryption for communication with PIA servers
	PIAEncryptionNormal models.PIAEncryption = "normal"
	// PIAEncryptionStrong is the strong level of encryption for communication with PIA servers
	PIAEncryptionStrong models.PIAEncryption = "strong"
)

const (
	PIAX509CRL_NORMAL     = "MIICWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAQEAQZo9X97ci8EcPYu/uK2HB152OZbeZCINmYyluLDOdcSvg6B5jI+ffKN3laDvczsG6CxmY3jNyc79XVpEYUnq4rT3FfveW1+Ralf+Vf38HdpwB8EWB4hZlQ205+21CALLvZvR8HcPxC9KEnev1mU46wkTiov0EKc+EdRxkj5yMgv0V2Reze7AP+NQ9ykvDScH4eYCsmufNpIjBLhpLE2cuZZXBLcPhuRzVoU3l7A9lvzG9mjA5YijHJGHNjlWFqyrn1CfYS6koa4TGEPngBoAziWRbDGdhEgJABHrpoaFYaL61zqyMR6jC0K2ps9qyZAN74LEBedEfK7tBOzWMwr58A=="
	PIAX509CRL_STRONG     = "MIIDWDCCAUAwDQYJKoZIhvcNAQENBQAwgegxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTETMBEGA1UEBxMKTG9zQW5nZWxlczEgMB4GA1UEChMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAsTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQDExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEKRMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxLzAtBgkqhkiG9w0BCQEWIHNlY3VyZUBwcml2YXRlaW50ZXJuZXRhY2Nlc3MuY29tFw0xNjA3MDgxOTAwNDZaFw0zNjA3MDMxOTAwNDZaMCYwEQIBARcMMTYwNzA4MTkwMDQ2MBECAQYXDDE2MDcwODE5MDA0NjANBgkqhkiG9w0BAQ0FAAOCAgEAppFfEpGsasjB1QgJcosGpzbf2kfRhM84o2TlqY1ua+Gi5TMdKydA3LJcNTjlI9a0TYAJfeRX5IkpoglSUuHuJgXhP3nEvX10mjXDpcu/YvM8TdE5JV2+EGqZ80kFtBeOq94WcpiVKFTR4fO+VkOK9zwspFfb1cNs9rHvgJ1QMkRUF8PpLN6AkntHY0+6DnigtSaKqldqjKTDTv2OeH3nPoh80SGrt0oCOmYKfWTJGpggMGKvIdvU3vH9+EuILZKKIskt+1dwdfA5Bkz1GLmiQG7+9ZZBQUjBG9Dos4hfX/rwJ3eU8oUIm4WoTz9rb71SOEuUUjP5NPy9HNx2vx+cVvLsTF4ZDZaUztW9o9JmIURDtbeyqxuHN3prlPWB6aj73IIm2dsDQvs3XXwRIxs8NwLbJ6CyEuvEOVCskdM8rdADWx1J0lRNlOJ0Z8ieLLEmYAA834VN1SboB6wJIAPxQU3rcBhXqO9y8aa2oRMg8NxZ5gr+PnKVMqag1x0IxbIgLxtkXQvxXxQHEMSODzvcOfK/nBRBsqTj30P+R87sU8titOoxNeRnBDRNhdEy/QGAqGh62ShPpQUCJdnKRiRTjnil9hMQHevoSuFKeEMO30FQL7BZyo37GFU+q1WPCplVZgCP9hC8Rn5K2+f6KLFo5bhtowSmu+GY1yZtg+RTtsA="
	PIACertificate_NORMAL = "MIIFqzCCBJOgAwIBAgIJAKZ7D5Yv87qDMA0GCSqGSIb3DQEBDQUAMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTAeFw0xNDA0MTcxNzM1MThaFw0zNDA0MTIxNzM1MThaMIHoMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExEzARBgNVBAcTCkxvc0FuZ2VsZXMxIDAeBgNVBAoTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQLExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UEAxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBCkTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMS8wLQYJKoZIhvcNAQkBFiBzZWN1cmVAcHJpdmF0ZWludGVybmV0YWNjZXNzLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPXDL1L9tX6DGf36liA7UBTy5I869z0UVo3lImfOs/GSiFKPtInlesP65577nd7UNzzXlH/P/CnFPdBWlLp5ze3HRBCc/Avgr5CdMRkEsySL5GHBZsx6w2cayQ2EcRhVTwWpcdldeNO+pPr9rIgPrtXqT4SWViTQRBeGM8CDxAyTopTsobjSiYZCF9Ta1gunl0G/8Vfp+SXfYCC+ZzWvP+L1pFhPRqzQQ8k+wMZIovObK1s+nlwPaLyayzw9a8sUnvWB/5rGPdIYnQWPgoNlLN9HpSmsAcw2z8DXI9pIxbr74cb3/HSfuYGOLkRqrOk6h4RCOfuWoTrZup1uEOn+fw8CAwEAAaOCAVQwggFQMB0GA1UdDgQWBBQv63nQ/pJAt5tLy8VJcbHe22ZOsjCCAR8GA1UdIwSCARYwggESgBQv63nQ/pJAt5tLy8VJcbHe22ZOsqGB7qSB6zCB6DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRMwEQYDVQQHEwpMb3NBbmdlbGVzMSAwHgYDVQQKExdQcml2YXRlIEludGVybmV0IEFjY2VzczEgMB4GA1UECxMXUHJpdmF0ZSBJbnRlcm5ldCBBY2Nlc3MxIDAeBgNVBAMTF1ByaXZhdGUgSW50ZXJuZXQgQWNjZXNzMSAwHgYDVQQpExdQcml2YXRlIEludGVybmV0IEFjY2VzczEvMC0GCSqGSIb3DQEJARYgc2VjdXJlQHByaXZhdGVpbnRlcm5ldGFjY2Vzcy5jb22CCQCmew+WL/O6gzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBDQUAA4IBAQAna5PgrtxfwTumD4+3/SYvwoD66cB8IcK//h1mCzAduU8KgUXocLx7QgJWo9lnZ8xUryXvWab2usg4fqk7FPi00bED4f4qVQFVfGfPZIH9QQ7/48bPM9RyfzImZWUCenK37pdw4Bvgoys2rHLHbGen7f28knT2j/cbMxd78tQc20TIObGjo8+ISTRclSTRBtyCGohseKYpTS9himFERpUgNtefvYHbn70mIOzfOJFTVqfrptf9jXa9N8Mpy3ayfodz1wiqdteqFXkTYoSDctgKMiZ6GdocK9nMroQipIQtpnwd4yBDWIyC6Bvlkrq5TQUtYDQ8z9v+DMO6iwyIDRiU"
	PIACertificate_STRONG = "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"
)

const (
	AUMelbourne     models.PIARegion = "AU Melbourne"
	AUPerth         models.PIARegion = "AU Perth"
	AUSydney        models.PIARegion = "AU Sydney"
	Austria         models.PIARegion = "Austria"
	Belgium         models.PIARegion = "Belgium"
	CAMontreal      models.PIARegion = "CA Montreal"
	CAToronto       models.PIARegion = "CA Toronto"
	CAVancouver     models.PIARegion = "CA Vancouver"
	CzechRepublic   models.PIARegion = "Czech Republic"
	DEBerlin        models.PIARegion = "DE Berlin"
	DEFrankfurt     models.PIARegion = "DE Frankfurt"
	Denmark         models.PIARegion = "Denmark"
	Finland         models.PIARegion = "Finland"
	France          models.PIARegion = "France"
	HongKong        models.PIARegion = "Hong Kong"
	Hungary         models.PIARegion = "Hungary"
	India           models.PIARegion = "India"
	Ireland         models.PIARegion = "Ireland"
	Israel          models.PIARegion = "Israel"
	Italy           models.PIARegion = "Italy"
	Japan           models.PIARegion = "Japan"
	Luxembourg      models.PIARegion = "Luxembourg"
	Mexico          models.PIARegion = "Mexico"
	Netherlands     models.PIARegion = "Netherlands"
	NewZealand      models.PIARegion = "New Zealand"
	Norway          models.PIARegion = "Norway"
	Poland          models.PIARegion = "Poland"
	Romania         models.PIARegion = "Romania"
	Singapore       models.PIARegion = "Singapore"
	Spain           models.PIARegion = "Spain"
	Sweden          models.PIARegion = "Sweden"
	Switzerland     models.PIARegion = "Switzerland"
	UAE             models.PIARegion = "UAE"
	UKLondon        models.PIARegion = "UK London"
	UKManchester    models.PIARegion = "UK Manchester"
	UKSouthampton   models.PIARegion = "UK Southampton"
	USAtlanta       models.PIARegion = "US Atlanta"
	USCalifornia    models.PIARegion = "US California"
	USChicago       models.PIARegion = "US Chicago"
	USDenver        models.PIARegion = "US Denver"
	USEast          models.PIARegion = "US East"
	USFlorida       models.PIARegion = "US Florida"
	USHouston       models.PIARegion = "US Houston"
	USLasVegas      models.PIARegion = "US Las Vegas"
	USNewYorkCity   models.PIARegion = "US New York City"
	USSeattle       models.PIARegion = "US Seattle"
	USSiliconValley models.PIARegion = "US Silicon Valley"
	USTexas         models.PIARegion = "US Texas"
	USWashingtonDC  models.PIARegion = "US Washington DC"
	USWest          models.PIARegion = "US West"
)

const (
	PIASubdomainAUMelbourne     string = "au-melbourne"
	PIASubdomainAUPerth         string = "au-perth"
	PIASubdomainAUSydney        string = "au-sydney"
	PIASubdomainAustria         string = "austria"
	PIASubdomainBelgium         string = "belgium"
	PIASubdomainCAMontreal      string = "ca-montreal"
	PIASubdomainCAToronto       string = "ca-toronto"
	PIASubdomainCAVancouver     string = "ca-vancouver"
	PIASubdomainCzechRepublic   string = "czech"
	PIASubdomainDEBerlin        string = "de-berlin"
	PIASubdomainDEFrankfurt     string = "de-frankfurt"
	PIASubdomainDenmark         string = "denmark"
	PIASubdomainFinland         string = "fi"
	PIASubdomainFrance          string = "france"
	PIASubdomainHongKong        string = "hk"
	PIASubdomainHungary         string = "hungary"
	PIASubdomainIndia           string = "in"
	PIASubdomainIreland         string = "ireland"
	PIASubdomainIsrael          string = "israel"
	PIASubdomainItaly           string = "italy"
	PIASubdomainJapan           string = "japan"
	PIASubdomainLuxembourg      string = "lu"
	PIASubdomainMexico          string = "mexico"
	PIASubdomainNetherlands     string = "nl"
	PIASubdomainNewZealand      string = "nz"
	PIASubdomainNorway          string = "no"
	PIASubdomainPoland          string = "poland"
	PIASubdomainRomania         string = "ro"
	PIASubdomainSingapore       string = "sg"
	PIASubdomainSpain           string = "spain"
	PIASubdomainSweden          string = "sweden"
	PIASubdomainSwitzerland     string = "swiss"
	PIASubdomainUAE             string = "ae"
	PIASubdomainUKLondon        string = "uk-london"
	PIASubdomainUKManchester    string = "uk-manchester"
	PIASubdomainUKSouthampton   string = "uk-southampton"
	PIASubdomainUSAtlanta       string = "us-atlanta"
	PIASubdomainUSCalifornia    string = "us-california"
	PIASubdomainUSChicago       string = "us-chicago"
	PIASubdomainUSDenver        string = "us-denver"
	PIASubdomainUSEast          string = "us-east"
	PIASubdomainUSFlorida       string = "us-florida"
	PIASubdomainUSHouston       string = "us-houston"
	PIASubdomainUSLasVegas      string = "us-lasvegas"
	PIASubdomainUSNewYorkCity   string = "us-newyorkcity"
	PIASubdomainUSSeattle       string = "us-seattle"
	PIASubdomainUSSiliconValley string = "us-siliconvalley"
	PIASubdomainUSTexas         string = "us-texas"
	PIASubdomainUSWashingtonDC  string = "us-washingtondc"
	PIASubdomainUSWest          string = "us-west"
)

var PIARegionToSubdomainMapping = map[models.PIARegion]string{
	AUMelbourne:     PIASubdomainAUMelbourne,
	AUPerth:         PIASubdomainAUPerth,
	AUSydney:        PIASubdomainAUSydney,
	Austria:         PIASubdomainAustria,
	Belgium:         PIASubdomainBelgium,
	CAMontreal:      PIASubdomainCAMontreal,
	CAToronto:       PIASubdomainCAToronto,
	CAVancouver:     PIASubdomainCAVancouver,
	CzechRepublic:   PIASubdomainCzechRepublic,
	DEBerlin:        PIASubdomainDEBerlin,
	DEFrankfurt:     PIASubdomainDEFrankfurt,
	Denmark:         PIASubdomainDenmark,
	Finland:         PIASubdomainFinland,
	France:          PIASubdomainFrance,
	HongKong:        PIASubdomainHongKong,
	Hungary:         PIASubdomainHungary,
	India:           PIASubdomainIndia,
	Ireland:         PIASubdomainIreland,
	Israel:          PIASubdomainIsrael,
	Italy:           PIASubdomainItaly,
	Japan:           PIASubdomainJapan,
	Luxembourg:      PIASubdomainLuxembourg,
	Mexico:          PIASubdomainMexico,
	Netherlands:     PIASubdomainNetherlands,
	NewZealand:      PIASubdomainNewZealand,
	Norway:          PIASubdomainNorway,
	Poland:          PIASubdomainPoland,
	Romania:         PIASubdomainRomania,
	Singapore:       PIASubdomainSingapore,
	Spain:           PIASubdomainSpain,
	Sweden:          PIASubdomainSweden,
	Switzerland:     PIASubdomainSwitzerland,
	UAE:             PIASubdomainUAE,
	UKLondon:        PIASubdomainUKLondon,
	UKManchester:    PIASubdomainUKManchester,
	UKSouthampton:   PIASubdomainUKSouthampton,
	USAtlanta:       PIASubdomainUSAtlanta,
	USCalifornia:    PIASubdomainUSCalifornia,
	USChicago:       PIASubdomainUSChicago,
	USDenver:        PIASubdomainUSDenver,
	USEast:          PIASubdomainUSEast,
	USFlorida:       PIASubdomainUSFlorida,
	USHouston:       PIASubdomainUSHouston,
	USLasVegas:      PIASubdomainUSLasVegas,
	USNewYorkCity:   PIASubdomainUSNewYorkCity,
	USSeattle:       PIASubdomainUSSeattle,
	USSiliconValley: PIASubdomainUSSiliconValley,
	USTexas:         PIASubdomainUSTexas,
	USWashingtonDC:  PIASubdomainUSWashingtonDC,
	USWest:          PIASubdomainUSWest,
}

const (
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
