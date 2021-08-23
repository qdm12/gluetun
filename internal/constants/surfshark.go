package constants

import (
	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	SurfsharkCertificate        = "MIIFTTCCAzWgAwIBAgIJAMs9S3fqwv+mMA0GCSqGSIb3DQEBCwUAMD0xCzAJBgNVBAYTAlZHMRIwEAYDVQQKDAlTdXJmc2hhcmsxGjAYBgNVBAMMEVN1cmZzaGFyayBSb290IENBMB4XDTE4MDMxNDA4NTkyM1oXDTI4MDMxMTA4NTkyM1owPTELMAkGA1UEBhMCVkcxEjAQBgNVBAoMCVN1cmZzaGFyazEaMBgGA1UEAwwRU3VyZnNoYXJrIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDEGMNj0aisM63oSkmVJyZPaYX7aPsZtzsxo6m6p5Wta3MGASoryRsBuRaH6VVa0fwbI1nw5ubyxkuaNa4v3zHVwuSq6F1p8S811+1YP1av+jqDcMyojH0ujZSHIcb/i5LtaHNXBQ3qN48Cc7sqBnTIIFpmb5HthQ/4pW+a82b1guM5dZHsh7q+LKQDIGmvtMtO1+NEnmj81BApFayiaD1ggvwDI4x7o/Y3ksfWSCHnqXGyqzSFLh8QuQrTmWUm84YHGFxoI1/8AKdIyVoB6BjcaMKtKs/pbctk6vkzmYf0XmGovDKPQF6MwUekchLjB5gSBNnptSQ9kNgnTLqi0OpSwI6ixX52Ksva6UM8P01ZIhWZ6ua/T/tArgODy5JZMW+pQ1A6L0b7egIeghpwKnPRG+5CzgO0J5UE6gv000mqbmC3CbiS8xi2xuNgruAyY2hUOoV9/BuBev8ttE5ZCsJH3YlG6NtbZ9hPc61GiBSx8NJnX5QHyCnfic/X87eST/amZsZCAOJ5v4EPSaKrItt+HrEFWZQIq4fJmHJNNbYvWzCE08AL+5/6Z+lxb/Bm3dapx2zdit3x2e+miGHekuiE8lQWD0rXD4+T+nDRi3X+kyt8Ex/8qRiUfrisrSHFzVMRungIMGdO9O/zCINFrb7wahm4PqU2f12Z9TRCOTXciQIDAQABo1AwTjAdBgNVHQ4EFgQUYRpbQwyDahLMN3F2ony3+UqOYOgwHwYDVR0jBBgwFoAUYRpbQwyDahLMN3F2ony3+UqOYOgwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAn9zV7F/XVnFNZhHFrt0ZS1Yqz+qM9CojLmiyblMFh0p7t+Hh+VKVgMwrz0LwDH4UsOosXA28eJPmech6/bjfymkoXISy/NUSTFpUChGO9RabGGxJsT4dugOw9MPaIVZffny4qYOc/rXDXDSfF2b+303lLPI43y9qoe0oyZ1vtk/UKG75FkWfFUogGNbpOkuz+et5Y0aIEiyg0yh6/l5Q5h8+yom0HZnREHhqieGbkaGKLkyu7zQ4D4tRK/mBhd8nv+09GtPEG+D5LPbabFVxKjBMP4Vp24WuSUOqcGSsURHevawPVBfgmsxf1UCjelaIwngdh6WfNCRXa5QQPQTKubQvkvXONCDdhmdXQccnRX1nJWhPYi0onffvjsWUfztRypsKzX4dvM9k7xnIcGSGEnCC4RCgt1UiZIj7frcCMssbA6vJ9naM0s7JF7N3VKeHJtqe1OCRHMYnWUZt9vrqX6IoIHlZCoLlv39wFW9QNxelcAOCVbD+19MZ0ZXt7LitjIqe7yF5WxDQN4xru087FzQ4Hfj7eH1SNLLyKZkA1eecjmRoi/OoqAt7afSnwtQLtMUc2bQDg6rHt5C0e4dCLqP/9PGZTSJiwmtRHJ/N5qYWIh9ju83APvLm/AGBTR2pXmj9G3KdVOkpIC7L35dI623cSEC3Q3UZutsEm/UplsM="
	SurfsharkOpenvpnStaticKeyV1 = "b02cb1d7c6fee5d4f89b8de72b51a8d0c7b282631d6fc19be1df6ebae9e2779e6d9f097058a31c97f57f0c35526a44ae09a01d1284b50b954d9246725a1ead1ff224a102ed9ab3da0152a15525643b2eee226c37041dc55539d475183b889a10e18bb94f079a4a49888da566b99783460ece01daaf93548beea6c827d9674897e7279ff1a19cb092659e8c1860fbad0db4ad0ad5732f1af4655dbd66214e552f04ed8fd0104e1d4bf99c249ac229ce169d9ba22068c6c0ab742424760911d4636aafb4b85f0c952a9ce4275bc821391aa65fcd0d2394f006e3fba0fd34c4bc4ab260f4b45dec3285875589c97d3087c9134d3a3aa2f904512e85aa2dc2202498"
)

func SurfsharkRegionChoices() (choices []string) {
	servers := SurfsharkServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func SurfsharkCountryChoices() (choices []string) {
	servers := SurfsharkServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func SurfsharkCityChoices() (choices []string) {
	servers := SurfsharkServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

// TODO remove in v4.
func SurfsharkRetroLocChoices() (choices []string) {
	locationData := SurfsharkLocationData()
	choices = make([]string, len(locationData))
	for i := range locationData {
		choices[i] = locationData[i].RetroLoc
	}
	return makeUnique(choices)
}

func SurfsharkHostToLocation() (hostToLocation map[string]models.SurfsharkLocationData) {
	locationData := SurfsharkLocationData()
	hostToLocation = make(map[string]models.SurfsharkLocationData, len(locationData))
	for _, data := range locationData {
		hostToLocation[data.Hostname] = data
	}
	return hostToLocation
}

// TODO remove retroRegion and servers from API in v4.
func SurfsharkLocationData() (data []models.SurfsharkLocationData) {
	//nolint:lll
	return []models.SurfsharkLocationData{
		{Region: "Asia Pacific", Country: "Australia", City: "Adelaide", RetroLoc: "Australia Adelaide", Hostname: "au-adl.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Australia", City: "Brisbane", RetroLoc: "Australia Brisbane", Hostname: "au-bne.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Australia", City: "Melbourne", RetroLoc: "Australia Melbourne", Hostname: "au-mel.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Australia", City: "Perth", RetroLoc: "Australia Perth", Hostname: "au-per.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Australia", City: "Sydney", RetroLoc: "Australia Sydney", Hostname: "au-syd.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Azerbaijan", City: "Baku", RetroLoc: "Azerbaijan", Hostname: "az-bak.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Hong Kong", City: "Hong Kong", RetroLoc: "Hong Kong", Hostname: "hk-hkg.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "India", City: "Chennai", RetroLoc: "India Chennai", Hostname: "in-chn.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "India", City: "Indore", RetroLoc: "India Indore", Hostname: "in-idr.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "India", City: "Mumbai", RetroLoc: "India Mumbai", Hostname: "in-mum.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Indonesia", City: "Jakarta", RetroLoc: "Indonesia", Hostname: "id-jak.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st001", Hostname: "jp-tok-st001.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st002", Hostname: "jp-tok-st002.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st003", Hostname: "jp-tok-st003.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st004", Hostname: "jp-tok-st004.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st005", Hostname: "jp-tok-st005.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st006", Hostname: "jp-tok-st006.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st007", Hostname: "jp-tok-st007.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st008", Hostname: "jp-tok-st008.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st009", Hostname: "jp-tok-st009.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st010", Hostname: "jp-tok-st010.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st011", Hostname: "jp-tok-st011.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st012", Hostname: "jp-tok-st012.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo st013", Hostname: "jp-tok-st013.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Japan", City: "Tokyo", RetroLoc: "Japan Tokyo", Hostname: "jp-tok.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Kazakhstan", City: "Oral", RetroLoc: "Kazakhstan", Hostname: "kz-ura.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Malaysia", City: "Kuala Lumpur", RetroLoc: "Malaysia", Hostname: "my-kul.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "New Zealand", City: "Auckland", RetroLoc: "New Zealand", Hostname: "nz-akl.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Philippines", City: "Manila", RetroLoc: "Philippines", Hostname: "ph-mnl.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore Hong Kong", City: "Hong Kong", RetroLoc: "Singapore Hong Kong", Hostname: "sg-hk.prod.surfshark.com", MultiHop: true},
		{Region: "Asia Pacific", Country: "Singapore in", City: "", RetroLoc: "Singapore in", Hostname: "sg-in.prod.surfshark.com", MultiHop: true},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore mp001", Hostname: "sg-sng-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore st001", Hostname: "sg-sng-st001.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore st002", Hostname: "sg-sng-st002.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore st003", Hostname: "sg-sng-st003.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore st004", Hostname: "sg-sng-st004.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Singapore", City: "Singapore", RetroLoc: "Singapore", Hostname: "sg-sng.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "South Korea", City: "Seoul", RetroLoc: "Korea", Hostname: "kr-seo.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Taiwan", City: "Taichung City", RetroLoc: "Taiwan", Hostname: "tw-tai.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Thailand", City: "Bangkok", RetroLoc: "Thailand", Hostname: "th-bkk.prod.surfshark.com", MultiHop: false},
		{Region: "Asia Pacific", Country: "Vietnam", City: "Ho Chi Minh City", RetroLoc: "Vietnam", Hostname: "vn-hcm.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Albania", City: "Tirana", RetroLoc: "Albania", Hostname: "al-tia.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Austria", City: "Vienna", RetroLoc: "Austria", Hostname: "at-vie.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Belgium", City: "Brussels", RetroLoc: "Belgium", Hostname: "be-bru.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Bosnia and Herzegovina", City: "Sarajevo", RetroLoc: "Bosnia and Herzegovina", Hostname: "ba-sjj.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Bulgaria", City: "Sofia", RetroLoc: "Bulgaria", Hostname: "bg-sof.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Croatia", City: "Zagreb", RetroLoc: "Croatia", Hostname: "hr-zag.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Cyprus", City: "Nicosia", RetroLoc: "Cyprus", Hostname: "cy-nic.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Czech Republic", City: "Prague", RetroLoc: "Czech Republic", Hostname: "cz-prg.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Denmark", City: "Copenhagen", RetroLoc: "Denmark", Hostname: "dk-cph.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Estonia", City: "Tallinn", RetroLoc: "Estonia", Hostname: "ee-tll.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Finland", City: "Helsinki", RetroLoc: "Finland", Hostname: "fi-hel.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "France Sweden", City: "", RetroLoc: "France Sweden", Hostname: "fr-se.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "France", City: "Bordeaux", RetroLoc: "France Bordeaux", Hostname: "fr-bod.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "France", City: "Marseille", RetroLoc: "France Marseilles", Hostname: "fr-mrs.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "France", City: "Paris", RetroLoc: "France Paris", Hostname: "fr-par.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany Singapour", City: "", RetroLoc: "Germany Singapour", Hostname: "de-sg.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "Germany UK", City: "", RetroLoc: "Germany UK", Hostname: "de-uk.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "Germany", City: "Berlin", RetroLoc: "Germany Berlin", Hostname: "de-ber.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main st001", Hostname: "de-fra-st001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main st002", Hostname: "de-fra-st002.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main st003", Hostname: "de-fra-st003.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main st004", Hostname: "de-fra-st004.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main st005", Hostname: "de-fra-st005.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt am Main", Hostname: "de-fra.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Germany", City: "Frankfurt am Main", RetroLoc: "Germany Frankfurt mp001", Hostname: "de-fra-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Greece", City: "Athens", RetroLoc: "Greece", Hostname: "gr-ath.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Hungary", City: "Budapest", RetroLoc: "Hungary", Hostname: "hu-bud.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Iceland", City: "Reykjavik", RetroLoc: "Iceland", Hostname: "is-rkv.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "India UK", City: "", RetroLoc: "India UK", Hostname: "in-uk.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "Ireland", City: "Dublin", RetroLoc: "Ireland", Hostname: "ie-dub.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Italy", City: "Milan", RetroLoc: "Italy Milan", Hostname: "it-mil.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Italy", City: "Rome", RetroLoc: "Italy Rome", Hostname: "it-rom.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Latvia", City: "Riga", RetroLoc: "Latvia", Hostname: "lv-rig.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Luxembourg", City: "Luxembourg", RetroLoc: "Luxembourg", Hostname: "lu-ste.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Moldova", City: "Chisinau", RetroLoc: "Moldova", Hostname: "md-chi.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Netherlands", City: "Amsterdam", RetroLoc: "Netherlands Amsterdam mp001", Hostname: "nl-ams-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Netherlands", City: "Amsterdam", RetroLoc: "Netherlands Amsterdam st001", Hostname: "nl-ams-st001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Netherlands", City: "Amsterdam", RetroLoc: "Netherlands Amsterdam", Hostname: "nl-ams.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "North Macedonia", City: "Skopje", RetroLoc: "North Macedonia", Hostname: "mk-skp.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Norway", City: "Oslo", RetroLoc: "Norway", Hostname: "no-osl.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Poland", City: "Gdansk", RetroLoc: "Poland Gdansk", Hostname: "pl-gdn.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Poland", City: "Warsaw", RetroLoc: "Poland Warsaw", Hostname: "pl-waw.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Portugal", City: "Lisbon", RetroLoc: "Portugal Lisbon", Hostname: "pt-lis.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Portugal", City: "Porto", RetroLoc: "Portugal Porto", Hostname: "pt-opo.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Romania", City: "Bucharest", RetroLoc: "Romania", Hostname: "ro-buc.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Russia", City: "Moscow", RetroLoc: "Russia Moscow", Hostname: "ru-mos.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Serbia", City: "Belgrade", RetroLoc: "Serbia", Hostname: "rs-beg.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Singapore Netherlands", City: "", RetroLoc: "Singapore Netherlands", Hostname: "sg-nl.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "Slovakia", City: "Bratislava", RetroLoc: "Slovekia", Hostname: "sk-bts.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Slovenia", City: "Ljubljana", RetroLoc: "Slovenia", Hostname: "si-lju.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Spain", City: "Barcelona", RetroLoc: "Spain Barcelona", Hostname: "es-bcn.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Spain", City: "Madrid", RetroLoc: "Spain Madrid", Hostname: "es-mad.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Spain", City: "Valencia", RetroLoc: "Spain Valencia", Hostname: "es-vlc.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Sweden", City: "Stockholm", RetroLoc: "Sweden", Hostname: "se-sto.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Switzerland", City: "Zurich", RetroLoc: "Switzerland", Hostname: "ch-zur.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "Turkey", City: "Istanbul", RetroLoc: "Turkey Istanbul", Hostname: "tr-ist.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "UK France", City: "", RetroLoc: "UK France", Hostname: "uk-fr.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "UK Germany", City: "", RetroLoc: "UK Germany", Hostname: "uk-de.prod.surfshark.com", MultiHop: true},
		{Region: "Europe", Country: "Ukraine", City: "Kyiv", RetroLoc: "Ukraine", Hostname: "ua-iev.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "Glasgow", RetroLoc: "UK Glasgow", Hostname: "uk-gla.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London mp001", Hostname: "uk-lon-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London st001", Hostname: "uk-lon-st001.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London st002", Hostname: "uk-lon-st002.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London st003", Hostname: "uk-lon-st003.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London st004", Hostname: "uk-lon-st004.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London st005", Hostname: "uk-lon-st005.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "London", RetroLoc: "UK London", Hostname: "uk-lon.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "United Kingdom", City: "Manchester", RetroLoc: "UK Manchester", Hostname: "uk-man.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "US Netherlands", City: "", RetroLoc: "US Netherlands", Hostname: "us-nl.prod.surfshark.com", MultiHop: false},
		{Region: "Europe", Country: "US Portugal", City: "", RetroLoc: "US Portugal", Hostname: "us-pt.prod.surfshark.com", MultiHop: true},
		{Region: "Middle East and Africa", Country: "Israel", City: "Tel Aviv", RetroLoc: "Israel", Hostname: "il-tlv.prod.surfshark.com", MultiHop: false},
		{Region: "Middle East and Africa", Country: "Nigeria", City: "Lagos", RetroLoc: "Nigeria", Hostname: "ng-lag.prod.surfshark.com", MultiHop: false},
		{Region: "Middle East and Africa", Country: "South Africa", City: "Johannesburg", RetroLoc: "South Africa", Hostname: "za-jnb.prod.surfshark.com", MultiHop: false},
		{Region: "Middle East and Africa", Country: "United Arab Emirates", City: "Dubai", RetroLoc: "United Arab Emirates", Hostname: "ae-dub.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Argentina", City: "Buenos Aires", RetroLoc: "Argentina Buenos Aires", Hostname: "ar-bua.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Australia US", City: "", RetroLoc: "Australia US", Hostname: "au-us.prod.surfshark.com", MultiHop: true},
		{Region: "The Americas", Country: "Brazil", City: "Sao Paulo", RetroLoc: "Brazil", Hostname: "br-sao.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Canada US", City: "", RetroLoc: "Canada US", Hostname: "ca-us.prod.surfshark.com", MultiHop: true},
		{Region: "The Americas", Country: "Canada", City: "Montreal", RetroLoc: "Canada Montreal", Hostname: "ca-mon.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Canada", City: "Toronto", RetroLoc: "Canada Toronto mp001", Hostname: "ca-tor-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Canada", City: "Toronto", RetroLoc: "Canada Toronto", Hostname: "ca-tor.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Canada", City: "Vancouver", RetroLoc: "Canada Vancouver", Hostname: "ca-van.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Chile", City: "Santiago", RetroLoc: "Chile", Hostname: "cl-san.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Colombia", City: "Bogota", RetroLoc: "Colombia", Hostname: "co-bog.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Costa Rica", City: "San Jose", RetroLoc: "Costa Rica", Hostname: "cr-sjn.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Mexico", City: "Mexico City", RetroLoc: "Mexico City Mexico", Hostname: "mx-mex.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "Netherlands US", City: "", RetroLoc: "Netherlands US", Hostname: "nl-us.prod.surfshark.com", MultiHop: true},
		{Region: "The Americas", Country: "United States", City: "Atlanta", RetroLoc: "US Atlanta", Hostname: "us-atl.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Bend", RetroLoc: "US Bend", Hostname: "us-bdn.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Boston", RetroLoc: "US Boston", Hostname: "us-bos.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Buffalo", RetroLoc: "US Buffalo", Hostname: "us-buf.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Charlotte", RetroLoc: "US Charlotte", Hostname: "us-clt.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Chicago", RetroLoc: "US Chicago", Hostname: "us-chi.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Dallas", RetroLoc: "US Dallas", Hostname: "us-dal.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Denver", RetroLoc: "US Denver", Hostname: "us-den.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Detroit", RetroLoc: "US Gahanna", Hostname: "us-dtw.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Houston", RetroLoc: "US Houston", Hostname: "us-hou.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Kansas City", RetroLoc: "US Kansas City", Hostname: "us-kan.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Las Vegas", RetroLoc: "US Las Vegas", Hostname: "us-las.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Latham", RetroLoc: "US Latham", Hostname: "us-ltm.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Los Angeles", RetroLoc: "US Los Angeles", Hostname: "us-lax.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Manassas", RetroLoc: "US Maryland", Hostname: "us-mnz.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Miami", RetroLoc: "US Miami", Hostname: "us-mia.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City mp001", Hostname: "us-nyc-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City st001", Hostname: "us-nyc-st001.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City st002", Hostname: "us-nyc-st002.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City st003", Hostname: "us-nyc-st003.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City st004", Hostname: "us-nyc-st004.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City st005", Hostname: "us-nyc-st005.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "New York", RetroLoc: "US New York City", Hostname: "us-nyc.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Orlando", RetroLoc: "US Orlando", Hostname: "us-orl.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Phoenix", RetroLoc: "US Phoenix", Hostname: "us-phx.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Salt Lake City", RetroLoc: "US Salt Lake City", Hostname: "us-slc.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "San Francisco", RetroLoc: "US San Francisco mp001", Hostname: "us-sfo-mp001.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "San Francisco", RetroLoc: "US San Francisco", Hostname: "us-sfo.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Seattle", RetroLoc: "US Seatle", Hostname: "us-sea.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "St. Louis", RetroLoc: "US Saint Louis", Hostname: "us-stl.prod.surfshark.com", MultiHop: false},
		{Region: "The Americas", Country: "United States", City: "Tampa", RetroLoc: "US Tampa", Hostname: "us-tpa.prod.surfshark.com", MultiHop: false},
	}
}

func SurfsharkHostnameChoices() (choices []string) {
	servers := SurfsharkServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

// SurfsharkServers returns a slice of all the server information for Surfshark.
func SurfsharkServers() (servers []models.SurfsharkServer) {
	servers = make([]models.SurfsharkServer, len(allServers.Surfshark.Servers))
	copy(servers, allServers.Surfshark.Servers)
	return servers
}
