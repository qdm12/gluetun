package constants

import (
	"net"

	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	WindscribeCertificate        = "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"
	WindscribeOpenvpnStaticKeyV1 = "5801926a57ac2ce27e3dfd1dd6ef82042d82bd4f3f0021296f57734f6f1ea714a6623845541c4b0c3dea0a050fe6746cb66dfab14cda27e5ae09d7c155aa554f399fa4a863f0e8c1af787e5c602a801d3a2ec41e395a978d56729457fe6102d7d9e9119aa83643210b33c678f9d4109e3154ac9c759e490cb309b319cf708cae83ddadc3060a7a26564d1a24411cd552fe6620ea16b755697a4fc5e6e9d0cfc0c5c4a1874685429046a424c026db672e4c2c492898052ba59128d46200b40f880027a8b6610a4d559bdc9346d33a0a6b08e75c7fd43192b162bfd0aef0c716b31584827693f676f9a5047123466f0654eade34972586b31c6ce7e395f4b478cb"
)

func WindscribeRegionChoices() (choices []string) {
	for _, server := range WindscribeServers() {
		choices = append(choices, string(server.Region))
	}
	return choices
}

func WindscribeServers() []models.WindscribeServer {
	return []models.WindscribeServer{
		{Region: models.WindscribeRegion("Albania"), IPs: []net.IP{{31, 171, 152, 179}}},
		{Region: models.WindscribeRegion("Argentina"), IPs: []net.IP{{167, 250, 6, 121}, {190, 105, 236, 19}, {190, 105, 236, 32}, {190, 105, 236, 50}}},
		{Region: models.WindscribeRegion("Australia"), IPs: []net.IP{{43, 245, 160, 35}, {45, 121, 208, 160}, {45, 121, 209, 160}, {45, 121, 210, 208}, {103, 62, 50, 208}, {103, 77, 233, 67}, {103, 77, 234, 211}, {103, 108, 92, 83}, {116, 90, 72, 243}, {116, 206, 228, 67}, {116, 206, 229, 131}}},
		{Region: models.WindscribeRegion("Austria"), IPs: []net.IP{{89, 187, 168, 66}, {217, 64, 127, 11}}},
		{Region: models.WindscribeRegion("Azerbaijan"), IPs: []net.IP{{85, 132, 61, 123}}},
		{Region: models.WindscribeRegion("Belgium"), IPs: []net.IP{{185, 232, 21, 131}, {194, 187, 251, 147}}},
		{Region: models.WindscribeRegion("Bosnia"), IPs: []net.IP{{185, 99, 3, 24}}},
		{Region: models.WindscribeRegion("Brazil"), IPs: []net.IP{{177, 54, 144, 68}, {177, 67, 80, 59}, {189, 1, 172, 12}}},
		{Region: models.WindscribeRegion("Bulgaria"), IPs: []net.IP{{185, 94, 192, 35}}},
		{Region: models.WindscribeRegion("Canada East"), IPs: []net.IP{{23, 154, 160, 177}, {66, 70, 148, 80}, {104, 227, 235, 129}, {104, 254, 92, 11}, {104, 254, 92, 91}, {144, 168, 163, 160}, {144, 168, 163, 193}, {184, 75, 212, 91}, {192, 190, 19, 65}, {192, 190, 19, 97}, {198, 8, 85, 195}, {198, 8, 85, 210}, {199, 204, 208, 158}}},
		{Region: models.WindscribeRegion("Canada West"), IPs: []net.IP{{104, 218, 61, 1}, {104, 218, 61, 33}, {162, 221, 207, 95}, {208, 78, 41, 1}, {208, 78, 41, 131}, {208, 78, 41, 163}}},
		{Region: models.WindscribeRegion("Colombia"), IPs: []net.IP{{138, 121, 203, 203}, {138, 186, 141, 155}}},
		{Region: models.WindscribeRegion("Croatia"), IPs: []net.IP{{85, 10, 56, 252}}},
		{Region: models.WindscribeRegion("Cyprus"), IPs: []net.IP{{157, 97, 132, 43}}},
		{Region: models.WindscribeRegion("Czech Republic"), IPs: []net.IP{{185, 156, 174, 11}, {185, 246, 210, 2}}},
		{Region: models.WindscribeRegion("Denmark"), IPs: []net.IP{{134, 90, 149, 147}, {185, 206, 224, 195}}},
		{Region: models.WindscribeRegion("Estona"), IPs: []net.IP{{46, 22, 211, 251}, {196, 196, 216, 131}}},
		{Region: models.WindscribeRegion("Fake Antarctica"), IPs: []net.IP{{23, 154, 160, 212}, {23, 154, 160, 222}}},
		{Region: models.WindscribeRegion("Finland"), IPs: []net.IP{{185, 112, 82, 227}, {194, 34, 133, 82}}},
		{Region: models.WindscribeRegion("France"), IPs: []net.IP{{45, 89, 174, 35}, {82, 102, 18, 35}, {84, 17, 42, 2}, {84, 17, 42, 34}, {185, 156, 173, 187}}},
		{Region: models.WindscribeRegion("Georgia"), IPs: []net.IP{{188, 93, 90, 83}}},
		{Region: models.WindscribeRegion("Germany"), IPs: []net.IP{{45, 87, 212, 51}, {89, 249, 65, 19}, {185, 130, 184, 195}, {195, 181, 170, 66}, {195, 181, 175, 98}, {217, 138, 194, 115}}},
		{Region: models.WindscribeRegion("Greece"), IPs: []net.IP{{78, 108, 38, 155}, {185, 226, 64, 111}, {188, 123, 126, 146}}},
		{Region: models.WindscribeRegion("Hong Kong"), IPs: []net.IP{{84, 17, 57, 114}, {103, 10, 197, 99}}},
		{Region: models.WindscribeRegion("Hungary"), IPs: []net.IP{{185, 104, 187, 43}}},
		{Region: models.WindscribeRegion("Iceland"), IPs: []net.IP{{82, 221, 139, 38}, {185, 165, 170, 2}}},
		{Region: models.WindscribeRegion("India"), IPs: []net.IP{{103, 205, 140, 227}, {169, 38, 68, 188}, {169, 38, 72, 12}, {169, 38, 72, 14}}},
		{Region: models.WindscribeRegion("Indonesia"), IPs: []net.IP{{45, 127, 134, 91}}},
		{Region: models.WindscribeRegion("Ireland"), IPs: []net.IP{{185, 24, 232, 146}, {185, 104, 219, 2}}},
		{Region: models.WindscribeRegion("Israel"), IPs: []net.IP{{160, 116, 0, 27}, {185, 191, 205, 139}}},
		{Region: models.WindscribeRegion("Italy"), IPs: []net.IP{{37, 120, 135, 83}, {37, 120, 207, 19}, {84, 17, 59, 66}, {87, 101, 94, 195}, {89, 40, 182, 3}}},
		{Region: models.WindscribeRegion("Japan"), IPs: []net.IP{{89, 187, 161, 114}, {193, 148, 16, 243}}},
		{Region: models.WindscribeRegion("Latvia"), IPs: []net.IP{{85, 254, 72, 23}}},
		{Region: models.WindscribeRegion("Lithuania"), IPs: []net.IP{{85, 206, 163, 225}}},
		{Region: models.WindscribeRegion("Macedonia"), IPs: []net.IP{{185, 225, 28, 51}}},
		{Region: models.WindscribeRegion("Malaysia"), IPs: []net.IP{{103, 106, 250, 31}, {103, 212, 69, 232}}},
		{Region: models.WindscribeRegion("Mexico"), IPs: []net.IP{{143, 255, 57, 67}, {190, 103, 179, 211}, {190, 103, 179, 217}, {201, 131, 125, 107}}},
		{Region: models.WindscribeRegion("Moldova"), IPs: []net.IP{{178, 175, 144, 123}}},
		{Region: models.WindscribeRegion("Netherlands"), IPs: []net.IP{{37, 120, 192, 19}, {46, 166, 143, 98}, {72, 11, 157, 35}, {72, 11, 157, 67}, {84, 17, 46, 2}, {185, 212, 171, 131}, {185, 253, 96, 3}}},
		{Region: models.WindscribeRegion("New Zealand"), IPs: []net.IP{{103, 62, 49, 113}, {103, 108, 94, 163}}},
		{Region: models.WindscribeRegion("Norway"), IPs: []net.IP{{37, 120, 203, 67}, {185, 206, 225, 131}}},
		{Region: models.WindscribeRegion("Philippines"), IPs: []net.IP{{103, 103, 0, 118}}},
		{Region: models.WindscribeRegion("Poland"), IPs: []net.IP{{5, 133, 11, 116}, {84, 17, 55, 98}, {185, 244, 214, 35}}},
		{Region: models.WindscribeRegion("Portugal"), IPs: []net.IP{{94, 46, 13, 215}, {185, 15, 21, 66}}},
		{Region: models.WindscribeRegion("Romania"), IPs: []net.IP{{89, 46, 103, 147}, {91, 207, 102, 147}}},
		{Region: models.WindscribeRegion("Russia"), IPs: []net.IP{{94, 242, 62, 19}, {94, 242, 62, 67}, {95, 213, 193, 195}, {95, 213, 193, 227}, {185, 22, 175, 132}, {188, 124, 42, 99}, {188, 124, 42, 115}}},
		{Region: models.WindscribeRegion("Serbia"), IPs: []net.IP{{141, 98, 103, 19}}},
		{Region: models.WindscribeRegion("Singapore"), IPs: []net.IP{{82, 102, 25, 131}, {89, 187, 162, 130}, {103, 62, 48, 224}, {185, 200, 117, 163}}},
		{Region: models.WindscribeRegion("Slovakia"), IPs: []net.IP{{185, 245, 85, 3}}},
		{Region: models.WindscribeRegion("Slovenia"), IPs: []net.IP{{146, 247, 24, 207}}},
		{Region: models.WindscribeRegion("South Africa"), IPs: []net.IP{{129, 232, 167, 211}, {165, 73, 248, 91}, {197, 242, 156, 53}, {197, 242, 157, 235}}},
		{Region: models.WindscribeRegion("South Korea"), IPs: []net.IP{{103, 212, 223, 3}, {218, 232, 76, 179}}},
		{Region: models.WindscribeRegion("Spain"), IPs: []net.IP{{37, 120, 142, 227}, {89, 238, 178, 43}, {185, 253, 99, 131}, {217, 138, 218, 99}}},
		{Region: models.WindscribeRegion("Sweden"), IPs: []net.IP{{31, 13, 191, 67}, {195, 181, 166, 129}}},
		{Region: models.WindscribeRegion("Switzerland"), IPs: []net.IP{{31, 7, 57, 242}, {37, 120, 213, 163}, {84, 17, 53, 2}, {89, 187, 165, 98}, {185, 156, 175, 179}}},
		{Region: models.WindscribeRegion("Thailand"), IPs: []net.IP{{27, 254, 130, 221}}},
		{Region: models.WindscribeRegion("Tunisia"), IPs: []net.IP{{41, 231, 5, 23}}},
		{Region: models.WindscribeRegion("Turkey"), IPs: []net.IP{{45, 123, 118, 156}, {45, 123, 119, 11}, {79, 98, 131, 43}, {176, 53, 113, 163}, {185, 125, 33, 227}}},
		{Region: models.WindscribeRegion("Ukraine"), IPs: []net.IP{{45, 141, 156, 11}, {45, 141, 156, 50}}},
		{Region: models.WindscribeRegion("United Arab Emirates"), IPs: []net.IP{{45, 9, 249, 43}}},
		{Region: models.WindscribeRegion("United kingdom"), IPs: []net.IP{{2, 58, 29, 17}, {81, 92, 207, 69}, {84, 17, 50, 130}, {89, 238, 131, 131}, {89, 238, 135, 133}, {89, 238, 150, 229}, {185, 212, 168, 133}, {212, 102, 63, 32}, {212, 102, 63, 62}, {217, 138, 254, 51}}},
		{Region: models.WindscribeRegion("US Central"), IPs: []net.IP{{67, 212, 238, 196}, {69, 12, 94, 67}, {104, 129, 18, 3}, {104, 223, 92, 163}, {107, 150, 31, 67}, {107, 150, 31, 131}, {107, 161, 86, 131}, {107, 182, 234, 240}, {161, 129, 70, 195}, {162, 222, 198, 67}, {172, 241, 26, 78}, {172, 241, 131, 129}, {198, 12, 76, 211}, {198, 54, 128, 116}, {198, 55, 125, 195}, {199, 115, 96, 83}, {204, 44, 112, 131}, {206, 217, 139, 19}, {206, 217, 143, 131}}},
		{Region: models.WindscribeRegion("Us East"), IPs: []net.IP{{23, 82, 136, 93}, {23, 83, 91, 170}, {23, 105, 170, 139}, {23, 226, 141, 195}, {38, 132, 118, 227}, {67, 21, 32, 145}, {68, 235, 35, 12}, {68, 235, 50, 227}, {86, 106, 87, 83}, {104, 168, 34, 147}, {104, 223, 127, 195}, {107, 150, 29, 131}, {142, 234, 200, 176}, {156, 96, 59, 102}, {162, 222, 195, 67}, {167, 160, 167, 195}, {167, 160, 172, 3}, {173, 44, 36, 67}, {173, 208, 45, 33}, {185, 232, 22, 195}, {198, 12, 64, 35}, {198, 147, 22, 225}, {206, 217, 128, 3}, {206, 217, 129, 227}, {217, 138, 255, 179}}},
		{Region: models.WindscribeRegion("US West"), IPs: []net.IP{{23, 83, 130, 166}, {23, 83, 131, 187}, {23, 94, 74, 99}, {37, 120, 147, 163}, {64, 120, 2, 174}, {66, 115, 176, 3}, {82, 102, 30, 67}, {89, 187, 185, 34}, {89, 187, 187, 98}, {104, 129, 3, 67}, {104, 129, 3, 163}, {104, 129, 56, 67}, {104, 129, 56, 131}, {104, 152, 222, 33}, {167, 88, 60, 227}, {167, 88, 60, 243}, {172, 241, 214, 202}, {172, 241, 250, 131}, {172, 255, 125, 141}, {185, 236, 200, 35}, {192, 3, 20, 51}, {198, 12, 116, 195}, {198, 23, 242, 147}, {209, 58, 129, 121}, {216, 45, 53, 131}, {217, 138, 217, 51}, {217, 138, 217, 211}}},
		{Region: models.WindscribeRegion("Vietnam"), IPs: []net.IP{{103, 9, 76, 197}, {103, 9, 79, 186}, {103, 9, 79, 219}}},
		{Region: models.WindscribeRegion("Windflix CA"), IPs: []net.IP{{104, 218, 60, 111}, {104, 254, 92, 99}}},
		{Region: models.WindscribeRegion("Windflix JP"), IPs: []net.IP{{5, 181, 235, 67}}},
		{Region: models.WindscribeRegion("Windflix UK"), IPs: []net.IP{{45, 9, 248, 3}, {81, 92, 200, 85}, {89, 47, 62, 83}}},
		{Region: models.WindscribeRegion("Windflix US"), IPs: []net.IP{{23, 105, 170, 130}, {23, 105, 170, 151}, {185, 232, 22, 131}, {204, 44, 112, 67}, {217, 138, 206, 211}}},
	}
}
