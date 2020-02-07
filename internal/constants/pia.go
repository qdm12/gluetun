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
	PIAOpenVPNURL     models.URL = "https://www.privateinternetaccess.com/openvpn"
	PIAPortForwardURL models.URL = "http://209.222.18.222:2000"
)
