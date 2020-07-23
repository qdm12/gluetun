package constants

import (
	"net"
	"sort"

	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	MullvadCertificate = "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"
)

func MullvadCountryChoices() (choices []string) {
	uniqueChoices := map[string]struct{}{}
	for _, server := range MullvadServers() {
		uniqueChoices[server.Country] = struct{}{}
	}
	for choice := range uniqueChoices {
		choices = append(choices, choice)
	}
	sort.Slice(choices, func(i, j int) bool {
		return choices[i] < choices[j]
	})
	return choices
}

func MullvadCityChoices() (choices []string) {
	uniqueChoices := map[string]struct{}{}
	for _, server := range MullvadServers() {
		uniqueChoices[server.City] = struct{}{}
	}
	for choice := range uniqueChoices {
		choices = append(choices, choice)
	}
	sort.Slice(choices, func(i, j int) bool {
		return choices[i] < choices[j]
	})
	return choices
}

func MullvadISPChoices() (choices []string) {
	uniqueChoices := map[string]struct{}{}
	for _, server := range MullvadServers() {
		uniqueChoices[server.ISP] = struct{}{}
	}
	for choice := range uniqueChoices {
		choices = append(choices, choice)
	}
	sort.Slice(choices, func(i, j int) bool {
		return choices[i] < choices[j]
	})
	return choices
}

func MullvadServers() []models.MullvadServer {
	return []models.MullvadServer{
		{
			Country:     "united arab emirates",
			City:        "dubai",
			ISP:         "m247",
			IPs:         []net.IP{{45, 9, 249, 34}},
			DefaultPort: 1194,
		},
		{
			Country:     "albania",
			City:        "tirana",
			ISP:         "iregister",
			IPs:         []net.IP{{31, 171, 154, 210}},
			DefaultPort: 1197,
		},
		{
			Country:     "austria",
			City:        "wien",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 155, 250}, {217, 64, 127, 138}, {217, 64, 127, 202}},
			DefaultPort: 1196,
		},
		{
			Country:     "australia",
			City:        "adelaide",
			ISP:         "intergrid",
			IPs:         []net.IP{{116, 206, 231, 58}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "brisbane",
			ISP:         "intergrid",
			IPs:         []net.IP{{43, 245, 160, 162}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "canberra",
			ISP:         "intergrid",
			IPs:         []net.IP{{116, 206, 229, 98}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "melbourne",
			ISP:         "intergrid",
			IPs:         []net.IP{{116, 206, 228, 202}, {116, 206, 228, 242}, {116, 206, 230, 98}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "perth",
			ISP:         "intergrid",
			IPs:         []net.IP{{103, 77, 235, 66}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "sydney",
			ISP:         "intergrid",
			IPs:         []net.IP{{43, 245, 162, 130}, {103, 77, 232, 130}, {103, 77, 232, 146}},
			DefaultPort: 1300,
		},
		{
			Country:     "australia",
			City:        "sydney",
			ISP:         "m247",
			IPs:         []net.IP{{217, 138, 204, 82}, {217, 138, 204, 98}, {217, 138, 204, 66}},
			DefaultPort: 1300,
		},
		{
			Country:     "belgium",
			City:        "brussels",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 218, 146}, {37, 120, 218, 138}, {91, 207, 57, 50}, {37, 120, 143, 138}, {185, 104, 186, 202}},
			DefaultPort: 1300,
		},
		{
			Country:     "bulgaria",
			City:        "sofia",
			ISP:         "m247",
			IPs:         []net.IP{{185, 94, 192, 42}, {185, 94, 192, 66}},
			DefaultPort: 1300,
		},
		{
			Country:     "brazil",
			City:        "sao paulo",
			ISP:         "qnax",
			IPs:         []net.IP{{191, 101, 62, 178}},
			DefaultPort: 1301,
		},
		{
			Country:     "brazil",
			City:        "sao paulo",
			ISP:         "heficed",
			IPs:         []net.IP{{177, 67, 80, 186}},
			DefaultPort: 1300,
		},
		{
			Country:     "canada",
			City:        "montreal",
			ISP:         "m247",
			IPs:         []net.IP{{139, 28, 218, 114}, {217, 138, 200, 194}, {217, 138, 200, 186}, {87, 101, 92, 146}, {176, 113, 74, 178}, {37, 120, 205, 114}, {87, 101, 92, 138}, {37, 120, 205, 122}, {217, 138, 200, 210}, {217, 138, 200, 202}},
			DefaultPort: 1300,
		},
		{
			Country:     "canada",
			City:        "toronto",
			ISP:         "amanah",
			IPs:         []net.IP{{184, 75, 214, 130}, {162, 219, 176, 250}},
			DefaultPort: 1300,
		},
		{
			Country:     "canada",
			City:        "vancouver",
			ISP:         "100tb",
			IPs:         []net.IP{{172, 83, 40, 34}, {172, 83, 40, 38}},
			DefaultPort: 1300,
		},
		{
			Country:     "canada",
			City:        "vancouver",
			ISP:         "esecuredata",
			IPs:         []net.IP{{71, 19, 249, 81}, {176, 113, 74, 186}, {71, 19, 248, 240}},
			DefaultPort: 1300,
		},
		{
			Country:     "switzerland",
			City:        "zurich",
			ISP:         "31173",
			IPs:         []net.IP{{193, 32, 127, 82}, {193, 32, 127, 81}, {193, 32, 127, 83}, {193, 32, 127, 84}},
			Owned:       true,
			DefaultPort: 1301,
		},
		{
			Country:     "zwitzerland",
			City:        "zurich",
			ISP:         "m247",
			IPs:         []net.IP{{185, 212, 170, 50}, {185, 183, 104, 82}, {185, 9, 18, 98}, {82, 102, 24, 130}, {82, 102, 24, 186}, {185, 212, 170, 162}, {185, 9, 18, 114}},
			DefaultPort: 1301,
		},
		{
			Country:     "switzerland",
			City:        "zurich",
			ISP:         "privateLayer",
			IPs:         []net.IP{{179, 43, 128, 170}, {81, 17, 20, 34}},
			DefaultPort: 1301,
		},
		{
			Country:     "czech republic",
			City:        "prague",
			ISP:         "m247",
			IPs:         []net.IP{{217, 138, 199, 82}, {217, 138, 199, 74}},
			DefaultPort: 1197,
		},
		{
			Country:     "germany",
			City:        "frankfurt",
			ISP:         "31173",
			IPs:         []net.IP{{185, 213, 155, 132}, {185, 213, 155, 140}, {185, 213, 155, 136}, {185, 213, 155, 133}, {185, 213, 155, 144}, {185, 213, 155, 143}, {185, 213, 155, 138}, {185, 213, 155, 142}, {185, 213, 155, 139}, {185, 213, 155, 135}, {185, 213, 155, 145}, {185, 213, 155, 137}, {185, 213, 155, 131}, {185, 213, 155, 134}, {185, 213, 155, 141}},
			Owned:       true,
			DefaultPort: 1197,
		},
		{
			Country:     "germany",
			City:        "frankfurt",
			ISP:         "m247",
			IPs:         []net.IP{{82, 102, 16, 90}, {185, 104, 184, 186}, {77, 243, 183, 202}},
			DefaultPort: 1197,
		},
		{
			Country:     "denmark",
			City:        "copenhagen",
			ISP:         "31173",
			IPs:         []net.IP{{141, 98, 254, 71}, {141, 98, 254, 72}},
			Owned:       true,
			DefaultPort: 1195,
		},
		{
			Country:     "denmark",
			City:        "copenhagen",
			ISP:         "m247",
			IPs:         []net.IP{{185, 206, 224, 114}, {185, 206, 224, 119}},
			DefaultPort: 1195,
		},
		{
			Country:     "denmark",
			City:        "copenhagen",
			ISP:         "blix",
			IPs:         []net.IP{{134, 90, 149, 138}},
			DefaultPort: 1195,
		},
		{
			Country:     "denmark",
			City:        "copenhagen",
			ISP:         "asergo",
			IPs:         []net.IP{{82, 103, 140, 213}},
			DefaultPort: 1195,
		},
		{
			Country:     "spain",
			City:        "madrid",
			ISP:         "m247",
			IPs:         []net.IP{{195, 206, 107, 146}, {45, 152, 183, 42}, {89, 238, 178, 74}, {45, 152, 183, 26}, {89, 238, 178, 34}},
			DefaultPort: 1195,
		},
		{
			Country:     "finland",
			City:        "helsinki",
			ISP:         "creanova",
			IPs:         []net.IP{{185, 204, 1, 174}, {185, 204, 1, 176}, {185, 212, 149, 201}, {185, 204, 1, 175}, {185, 204, 1, 173}, {185, 204, 1, 172}, {185, 204, 1, 171}},
			Owned:       true,
			DefaultPort: 1196,
		},
		{
			Country:     "france",
			City:        "paris",
			ISP:         "31173",
			IPs:         []net.IP{{193, 32, 126, 83}, {193, 32, 126, 82}, {193, 32, 126, 81}, {193, 32, 126, 84}},
			Owned:       true,
			DefaultPort: 1301,
		},
		{
			Country:     "france",
			City:        "paris",
			ISP:         "m247",
			IPs:         []net.IP{{185, 189, 113, 82}, {185, 156, 173, 218}, {185, 128, 25, 162}},
			DefaultPort: 1301,
		},
		{
			Country:     "uk",
			City:        "london",
			ISP:         "31173",
			IPs:         []net.IP{{141, 98, 252, 133}, {141, 98, 252, 139}, {141, 98, 252, 137}, {141, 98, 252, 143}, {141, 98, 252, 142}, {141, 98, 252, 132}, {141, 98, 252, 134}, {141, 98, 252, 140}, {141, 98, 252, 141}, {141, 98, 252, 136}, {141, 98, 252, 144}, {141, 98, 252, 131}, {141, 98, 252, 135}, {141, 98, 252, 138}},
			Owned:       true,
			DefaultPort: 1196,
		},
		{
			Country:     "uk",
			City:        "london",
			ISP:         "m247",
			IPs:         []net.IP{{185, 200, 118, 105}, {185, 212, 168, 244}},
			DefaultPort: 1196,
		},
		{
			Country:     "uk",
			City:        "manchester",
			ISP:         "m247",
			IPs:         []net.IP{{89, 238, 130, 66}, {81, 92, 205, 10}, {89, 238, 130, 74}, {81, 92, 205, 18}, {81, 92, 205, 26}, {89, 238, 183, 244}, {89, 238, 132, 36}, {217, 151, 98, 68}, {37, 120, 159, 164}, {89, 238, 183, 60}},
			DefaultPort: 1196,
		},
		{
			Country:     "greece",
			City:        "athens",
			ISP:         "aweb",
			IPs:         []net.IP{{185, 226, 67, 168}},
			DefaultPort: 1302,
		},
		{
			Country:     "hong kong",
			City:        "hong kong",
			ISP:         "leaseweb",
			IPs:         []net.IP{{209, 58, 185, 53}, {209, 58, 184, 146}},
			DefaultPort: 1194,
		},
		{
			Country:     "hungary",
			City:        "budapest",
			ISP:         "m247",
			IPs:         []net.IP{{185, 94, 190, 138}, {185, 189, 114, 10}},
			DefaultPort: 1300,
		},
		{
			Country:     "ireland",
			City:        "dublin",
			ISP:         "m247",
			IPs:         []net.IP{{217, 138, 222, 90}, {217, 138, 222, 82}},
			DefaultPort: 1197,
		},
		{
			Country:     "israel",
			City:        "tel aviv",
			ISP:         "hqserv",
			IPs:         []net.IP{{185, 191, 207, 210}},
			DefaultPort: 1301,
		},
		{
			Country:     "italy",
			City:        "milan",
			ISP:         "m247",
			IPs:         []net.IP{{217, 138, 197, 106}, {217, 64, 113, 180}, {217, 138, 197, 98}, {217, 138, 197, 114}, {217, 64, 113, 183}},
			DefaultPort: 1300,
		},
		{
			Country:     "japan",
			City:        "tokyo",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 210, 138}, {193, 148, 16, 218}, {37, 120, 210, 146}, {185, 242, 4, 50}, {37, 120, 210, 122}},
			DefaultPort: 1300,
		},
		{
			Country:     "luxembourg",
			City:        "luxembourg",
			ISP:         "evoluso",
			IPs:         []net.IP{{92, 223, 89, 160}, {92, 223, 89, 182}},
			DefaultPort: 1301,
		},
		{
			Country:     "latvia",
			City:        "riga",
			ISP:         "makonix",
			IPs:         []net.IP{{31, 170, 22, 2}},
			DefaultPort: 1300,
		},
		{
			Country:     "moldova",
			City:        "chisinau",
			ISP:         "trabia",
			IPs:         []net.IP{{178, 175, 142, 194}},
			DefaultPort: 1197,
		},
		{
			Country:     "netherlands",
			City:        "amsterdam",
			ISP:         "31173",
			IPs:         []net.IP{{185, 65, 134, 139}, {185, 65, 134, 133}, {185, 65, 134, 148}, {185, 65, 134, 147}, {185, 65, 134, 141}, {185, 65, 134, 140}, {185, 65, 134, 145}, {185, 65, 134, 132}, {185, 65, 134, 146}, {185, 65, 134, 143}, {185, 65, 134, 134}, {185, 65, 134, 136}, {185, 65, 134, 135}, {185, 65, 134, 142}, {185, 65, 134, 144}},
			Owned:       true,
			DefaultPort: 1194,
		},
		{
			Country:     "norway",
			City:        "oslo",
			ISP:         "blix",
			IPs:         []net.IP{{91, 90, 44, 13}, {91, 90, 44, 18}, {91, 90, 44, 12}, {91, 90, 44, 15}, {91, 90, 44, 16}, {91, 90, 44, 17}, {91, 90, 44, 14}, {91, 90, 44, 11}},
			Owned:       true,
			DefaultPort: 1302,
		},
		{
			Country:     "new zealand",
			City:        "auckland",
			ISP:         "intergrid",
			IPs:         []net.IP{{103, 231, 91, 114}},
			DefaultPort: 1195,
		},
		{
			Country:     "poland",
			City:        "warsaw",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 211, 202}, {37, 120, 156, 162}, {185, 244, 214, 210}, {37, 120, 211, 186}, {185, 244, 214, 215}, {37, 120, 211, 194}},
			DefaultPort: 1194,
		},
		{
			Country:     "portugal",
			City:        "lisbon",
			ISP:         "dotsi",
			IPs:         []net.IP{{5, 206, 231, 214}},
			DefaultPort: 1194,
		},
		{
			Country:     "romania",
			City:        "bucharest",
			ISP:         "m247",
			IPs:         []net.IP{{89, 40, 181, 146}, {185, 181, 100, 202}, {89, 40, 181, 82}, {185, 45, 13, 10}, {89, 40, 181, 210}},
			DefaultPort: 1301,
		},
		{
			Country:     "serbia",
			City:        "belgrade",
			ISP:         "m247",
			IPs:         []net.IP{{141, 98, 103, 50}},
			DefaultPort: 1301,
		},
		{
			Country:     "serbia",
			City:        "nis",
			ISP:         "ninet",
			IPs:         []net.IP{{176, 104, 107, 118}},
			DefaultPort: 1301,
		},
		{
			Country:     "sweden",
			City:        "gothenburg",
			ISP:         "31173",
			IPs:         []net.IP{{185, 213, 154, 139}, {185, 213, 154, 141}, {185, 213, 154, 140}, {185, 213, 154, 132}, {185, 213, 154, 135}, {185, 213, 154, 138}, {185, 213, 154, 133}, {185, 213, 154, 131}, {185, 213, 154, 134}, {185, 213, 154, 142}, {185, 213, 154, 137}},
			Owned:       true,
			DefaultPort: 1302,
		},
		{
			Country:     "sweden",
			City:        "helsingborg",
			ISP:         "31173",
			IPs:         []net.IP{{185, 213, 152, 133}, {185, 213, 152, 132}, {185, 213, 152, 138}, {185, 213, 152, 131}, {185, 213, 152, 137}},
			Owned:       true,
			DefaultPort: 1302,
		},
		{
			Country:     "sweden",
			City:        "malmo",
			ISP:         "31173",
			IPs:         []net.IP{{193, 138, 218, 138}, {45, 83, 220, 87}, {141, 98, 255, 94}, {141, 98, 255, 85}, {141, 98, 255, 87}, {141, 98, 255, 92}, {45, 83, 220, 84}, {141, 98, 255, 86}, {45, 83, 220, 81}, {193, 138, 218, 135}, {193, 138, 218, 131}, {193, 138, 218, 136}, {141, 98, 255, 88}, {141, 98, 255, 91}, {193, 138, 218, 133}, {45, 83, 220, 89}, {45, 83, 220, 88}, {141, 98, 255, 84}, {141, 98, 255, 89}, {193, 138, 218, 134}, {45, 83, 220, 86}, {141, 98, 255, 83}, {45, 83, 220, 85}, {141, 98, 255, 90}, {141, 98, 255, 93}, {193, 138, 218, 132}, {193, 138, 218, 137}, {45, 83, 220, 91}},
			Owned:       true,
			DefaultPort: 1302,
		},
		{
			Country:     "sweden",
			City:        "stockholm",
			ISP:         "31173",
			IPs:         []net.IP{{185, 65, 135, 150}, {185, 65, 135, 153}, {185, 65, 135, 151}, {185, 65, 135, 149}, {185, 65, 135, 141}, {185, 65, 135, 144}, {185, 65, 135, 145}, {185, 65, 135, 140}, {185, 65, 135, 134}, {185, 65, 135, 139}, {185, 65, 135, 131}, {185, 65, 135, 152}, {185, 65, 135, 146}, {185, 65, 135, 138}, {185, 65, 135, 143}, {185, 65, 135, 135}, {185, 65, 135, 154}, {185, 65, 135, 136}, {185, 65, 135, 133}, {185, 65, 135, 132}},
			Owned:       true,
			DefaultPort: 1302,
		},
		{
			Country:     "singapore",
			City:        "singapore",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 208, 218}, {37, 120, 208, 234}, {37, 120, 208, 226}, {185, 128, 24, 50}},
			DefaultPort: 1196,
		},
		{
			Country:     "singapore",
			City:        "singapore",
			ISP:         "leaseweb",
			IPs:         []net.IP{{103, 254, 153, 82}},
			DefaultPort: 1196,
		},
		{
			Country:     "usa",
			City:        "atlanta",
			ISP:         "100tb",
			IPs:         []net.IP{{208, 84, 153, 142}, {107, 152, 108, 62}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "atlanta",
			ISP:         "quadranet",
			IPs:         []net.IP{{104, 129, 24, 242}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "atlanta",
			ISP:         "micfo",
			IPs:         []net.IP{{155, 254, 96, 2}, {155, 254, 96, 18}, {155, 254, 96, 34}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "chicago",
			ISP:         "tzulo",
			IPs:         []net.IP{{68, 235, 43, 18}, {68, 235, 43, 26}, {68, 235, 43, 42}, {68, 235, 43, 50}, {68, 235, 43, 58}, {68, 235, 43, 66}, {68, 235, 43, 74}}, // 3 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "chicago",
			ISP:         "quadranet",
			IPs:         []net.IP{}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "dallas",
			ISP:         "quadranet",
			IPs:         []net.IP{{96, 44, 145, 18}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "dallas",
			ISP:         "100tb",
			IPs:         []net.IP{{104, 200, 142, 50}, {107, 152, 102, 106}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "denver",
			ISP:         "tzulo",
			IPs:         []net.IP{{198, 54, 128, 74}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "los angeles",
			ISP:         "m247",
			IPs:         []net.IP{{45, 152, 182, 66}, {45, 152, 182, 74}, {45, 83, 89, 162}, {185, 230, 126, 146}}, // 7 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "los angeles",
			ISP:         "tzulo",
			IPs:         []net.IP{{198, 54, 129, 74}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "los angeles",
			ISP:         "100tb",
			IPs:         []net.IP{{104, 200, 152, 66}, {107, 181, 168, 130}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "los angeles",
			ISP:         "choopa",
			IPs:         []net.IP{{104, 238, 143, 58}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "miami",
			ISP:         "m247",
			IPs:         []net.IP{{37, 120, 215, 130}, {193, 37, 252, 138}, {193, 37, 252, 154}, {37, 120, 215, 138}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "miami",
			ISP:         "100tb",
			IPs:         []net.IP{{172, 98, 76, 114}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "miami",
			ISP:         "micfo",
			IPs:         []net.IP{}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "new york",
			ISP:         "m247",
			IPs:         []net.IP{{185, 232, 22, 66}, {185, 232, 22, 98}, {193, 148, 18, 250}, {185, 232, 22, 10}, {217, 138, 206, 10}, {193, 148, 18, 218}, {193, 148, 18, 226}, {193, 148, 18, 194}, {87, 101, 95, 98}, {87, 101, 95, 114}, {87, 101, 95, 122}, {212, 103, 48, 226}, {176, 113, 72, 226}, {217, 138, 198, 250}, {217, 138, 206, 58}}, // 5 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "new york",
			ISP:         "100tb",
			IPs:         []net.IP{{107, 182, 226, 206}, {107, 182, 226, 218}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "phoenix",
			ISP:         "100tb",
			IPs:         []net.IP{}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "phoenix",
			ISP:         "micfo",
			IPs:         []net.IP{{192, 200, 24, 82}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "piscataway",
			ISP:         "choopa",
			IPs:         []net.IP{{108, 61, 78, 138}, {108, 61, 48, 115}, {66, 55, 147, 59}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "seattle",
			ISP:         "100tb",
			IPs:         []net.IP{{104, 200, 129, 202}, {104, 200, 129, 150}, {104, 200, 129, 110}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "seattle",
			ISP:         "micfo",
			IPs:         []net.IP{{104, 128, 136, 146}},
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "san francisco",
			ISP:         "micfo",
			IPs:         []net.IP{{209, 209, 238, 34}}, // 1 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "salt lake city",
			ISP:         "100tb",
			IPs:         []net.IP{{107, 182, 238, 229}, {107, 182, 235, 233}, {67, 212, 238, 236}, {67, 212, 238, 237}, {67, 212, 238, 239}, {107, 182, 239, 185}, {107, 182, 239, 170}}, // 2 missing
			DefaultPort: 1194,
		},
		{
			Country:     "usa",
			City:        "secaucus",
			ISP:         "quadranet",
			IPs:         []net.IP{{23, 226, 131, 154}}, // 1 missing
			DefaultPort: 1194,
		},
	}
}
