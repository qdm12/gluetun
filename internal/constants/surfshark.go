package constants

import (
	"net"

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
	return choices
}

func SurfsharkServers() []models.SurfsharkServer {
	return []models.SurfsharkServer{
		{Region: "Albania", IPs: []net.IP{{31, 171, 153, 99}, {31, 171, 153, 131}, {31, 171, 154, 101}, {31, 171, 154, 149}, {31, 171, 154, 163}, {31, 171, 154, 165}}},
		{Region: "Argentina Buenos Aires", IPs: []net.IP{{91, 206, 168, 9}, {91, 206, 168, 13}, {91, 206, 168, 21}, {91, 206, 168, 31}, {91, 206, 168, 41}, {91, 206, 168, 45}, {91, 206, 168, 54}, {91, 206, 168, 62}}},
		{Region: "Australia Melbourne", IPs: []net.IP{{103, 192, 80, 11}, {103, 192, 80, 147}, {103, 192, 80, 227}, {144, 48, 38, 139}, {144, 48, 38, 149}, {144, 48, 38, 181}}},
		{Region: "Australia Perth", IPs: []net.IP{{45, 248, 78, 43}, {124, 150, 139, 27}, {124, 150, 139, 35}, {124, 150, 139, 123}, {124, 150, 139, 125}, {124, 150, 139, 179}}},
		{Region: "Australia Sydney", IPs: []net.IP{{45, 125, 247, 45}, {180, 149, 228, 115}, {180, 149, 228, 149}, {180, 149, 228, 165}, {180, 149, 228, 171}, {180, 149, 228, 173}}},
		{Region: "Austria", IPs: []net.IP{{5, 253, 207, 85}, {37, 120, 212, 133}, {37, 120, 212, 139}, {89, 187, 168, 44}, {89, 187, 168, 46}, {89, 187, 168, 54}, {89, 187, 168, 56}}},
		{Region: "Belgium", IPs: []net.IP{{5, 253, 205, 99}, {5, 253, 205, 101}, {37, 120, 218, 29}, {89, 249, 73, 197}, {91, 90, 123, 149}, {91, 90, 123, 173}, {91, 90, 123, 213}, {185, 104, 186, 77}}},
		{Region: "Brazil", IPs: []net.IP{{45, 231, 207, 72}, {191, 96, 13, 39}, {191, 96, 13, 41}, {191, 96, 13, 43}, {191, 96, 73, 210}, {191, 96, 73, 214}, {191, 96, 73, 216}, {191, 96, 73, 228}, {191, 96, 73, 232}}},
		{Region: "Bulgaria", IPs: []net.IP{{37, 120, 152, 37}, {37, 120, 152, 195}, {217, 138, 202, 19}, {217, 138, 202, 21}}},
		{Region: "Canada Montreal", IPs: []net.IP{{172, 98, 82, 243}, {198, 8, 85, 5}, {198, 8, 85, 69}, {198, 8, 85, 72}, {198, 8, 85, 89}, {198, 8, 85, 163}}},
		{Region: "Canada Toronto", IPs: []net.IP{{68, 71, 244, 131}, {68, 71, 244, 197}, {68, 71, 244, 212}, {68, 71, 244, 220}, {104, 200, 138, 147}, {104, 200, 138, 154}}},
		{Region: "Canada Toronto mp001", IPs: []net.IP{{138, 197, 151, 26}}},
		{Region: "Chile", IPs: []net.IP{{31, 169, 121, 3}, {31, 169, 121, 5}}},
		{Region: "Colombia", IPs: []net.IP{{45, 129, 32, 3}, {45, 129, 32, 8}, {45, 129, 32, 20}, {45, 129, 32, 22}, {45, 129, 32, 27}, {45, 129, 32, 29}, {45, 129, 32, 32}, {45, 129, 32, 38}}},
		{Region: "Costa Rica", IPs: []net.IP{{176, 227, 241, 24}, {176, 227, 241, 26}, {176, 227, 241, 29}, {176, 227, 241, 31}, {176, 227, 241, 33}, {176, 227, 241, 35}}},
		{Region: "Croatia", IPs: []net.IP{{85, 10, 51, 89}, {85, 10, 51, 91}, {89, 164, 99, 111}, {89, 164, 99, 134}, {89, 164, 99, 136}}},
		{Region: "Cyprus", IPs: []net.IP{{195, 47, 194, 36}, {195, 47, 194, 54}, {195, 47, 194, 56}, {195, 47, 194, 59}, {195, 47, 194, 61}, {195, 47, 194, 64}, {195, 47, 194, 66}, {195, 47, 194, 70}}},
		{Region: "Czech Republic", IPs: []net.IP{{185, 180, 14, 149}, {193, 9, 112, 181}, {217, 138, 199, 179}, {217, 138, 220, 133}, {217, 138, 220, 149}, {217, 138, 220, 163}}},
		{Region: "Estonia", IPs: []net.IP{{165, 231, 163, 21}, {185, 174, 159, 51}, {185, 174, 159, 53}, {185, 174, 159, 67}}},
		{Region: "Finland", IPs: []net.IP{{196, 244, 191, 43}, {196, 244, 191, 91}, {196, 244, 191, 101}, {196, 244, 191, 165}}},
		{Region: "France Bordeaux", IPs: []net.IP{{185, 108, 106, 26}, {185, 108, 106, 69}, {185, 108, 106, 102}, {185, 108, 106, 106}, {185, 108, 106, 180}, {185, 108, 106, 186}}},
		{Region: "France Marseilles", IPs: []net.IP{{185, 166, 84, 19}, {185, 166, 84, 29}, {185, 166, 84, 51}, {185, 166, 84, 55}, {185, 166, 84, 57}, {185, 166, 84, 61}, {185, 166, 84, 63}, {185, 166, 84, 83}}},
		{Region: "France Paris", IPs: []net.IP{{45, 89, 174, 61}, {84, 17, 60, 250}, {84, 247, 51, 253}, {143, 244, 56, 232}, {143, 244, 57, 101}, {143, 244, 57, 103}}},
		{Region: "Germany Berlin", IPs: []net.IP{{37, 120, 217, 179}, {193, 29, 106, 21}, {193, 29, 106, 83}, {193, 29, 106, 149}, {193, 29, 106, 187}, {217, 138, 216, 219}}},
		{Region: "Germany Frankfurt am Main", IPs: []net.IP{{37, 120, 197, 13}, {89, 187, 169, 104}, {138, 199, 19, 132}, {138, 199, 19, 147}, {138, 199, 19, 157}, {138, 199, 19, 190}, {156, 146, 33, 73}, {156, 146, 33, 87}}},
		{Region: "Germany Frankfurt am Main st001", IPs: []net.IP{{45, 87, 212, 179}}},
		{Region: "Germany Frankfurt am Main st002", IPs: []net.IP{{45, 87, 212, 181}}},
		{Region: "Germany Frankfurt am Main st003", IPs: []net.IP{{45, 87, 212, 183}}},
		{Region: "Germany Frankfurt am Main st004", IPs: []net.IP{{195, 181, 174, 226}}},
		{Region: "Germany Frankfurt am Main st005", IPs: []net.IP{{195, 181, 174, 228}}},
		{Region: "Germany Frankfurt mp001", IPs: []net.IP{{46, 101, 189, 14}}},
		{Region: "Germany Munich", IPs: []net.IP{{79, 143, 191, 141}, {79, 143, 191, 231}, {178, 238, 231, 53}, {178, 238, 231, 55}}},
		{Region: "Germany Nuremberg", IPs: []net.IP{{62, 171, 149, 162}, {62, 171, 151, 154}, {62, 171, 151, 156}, {62, 171, 151, 158}, {62, 171, 151, 160}, {95, 111, 253, 65}, {144, 91, 123, 50}, {144, 91, 123, 52}}},
		{Region: "Germany Singapour", IPs: []net.IP{{159, 89, 14, 157}}},
		{Region: "Greece", IPs: []net.IP{{194, 150, 167, 28}, {194, 150, 167, 32}, {194, 150, 167, 38}, {194, 150, 167, 40}, {194, 150, 167, 44}, {194, 150, 167, 48}}},
		{Region: "Hong Kong", IPs: []net.IP{{84, 17, 37, 156}, {84, 17, 57, 66}, {84, 17, 57, 71}, {212, 102, 42, 194}, {212, 102, 42, 196}, {212, 102, 42, 211}}},
		{Region: "Iceland", IPs: []net.IP{{82, 221, 128, 156}, {82, 221, 128, 169}, {82, 221, 143, 62}, {82, 221, 143, 64}, {82, 221, 143, 69}, {82, 221, 143, 71}}},
		{Region: "India Chennai", IPs: []net.IP{{103, 94, 27, 99}, {103, 94, 27, 101}, {103, 94, 27, 115}, {103, 94, 27, 227}, {103, 108, 117, 147}}},
		{Region: "India Indore", IPs: []net.IP{{103, 39, 132, 187}, {103, 39, 132, 189}, {103, 39, 134, 59}, {103, 39, 134, 61}, {137, 59, 52, 107}, {137, 59, 52, 109}}},
		{Region: "Indonesia", IPs: []net.IP{{103, 120, 66, 214}, {103, 120, 66, 216}, {103, 120, 66, 221}, {103, 120, 66, 227}}},
		{Region: "Ireland", IPs: []net.IP{{5, 157, 13, 67}, {5, 157, 13, 91}, {5, 157, 13, 93}, {5, 157, 13, 107}, {5, 157, 13, 117}, {5, 157, 13, 133}, {185, 108, 128, 161}, {217, 138, 222, 43}}},
		{Region: "Israel", IPs: []net.IP{{87, 239, 255, 111}, {87, 239, 255, 114}, {87, 239, 255, 116}, {87, 239, 255, 121}}},
		{Region: "Italy Milan", IPs: []net.IP{{37, 120, 201, 71}, {45, 9, 251, 167}, {84, 17, 58, 134}, {84, 17, 58, 146}, {212, 102, 54, 152}, {212, 102, 54, 167}, {212, 102, 54, 170}, {212, 102, 54, 177}}},
		{Region: "Italy Rome", IPs: []net.IP{{82, 102, 26, 115}, {87, 101, 94, 213}, {185, 217, 71, 21}, {185, 217, 71, 51}, {185, 217, 71, 213}, {185, 217, 71, 243}, {217, 138, 219, 237}, {217, 138, 219, 253}}},
		{Region: "Japan Tokyo", IPs: []net.IP{{45, 87, 213, 87}, {45, 87, 213, 103}, {84, 17, 34, 26}, {89, 187, 161, 4}, {89, 187, 161, 241}, {138, 199, 22, 130}}},
		{Region: "Japan Tokyo st001", IPs: []net.IP{{45, 87, 213, 19}}},
		{Region: "Japan Tokyo st002", IPs: []net.IP{{45, 87, 213, 21}}},
		{Region: "Japan Tokyo st003", IPs: []net.IP{{45, 87, 213, 23}}},
		{Region: "Japan Tokyo st004", IPs: []net.IP{{217, 138, 212, 19}}},
		{Region: "Japan Tokyo st005", IPs: []net.IP{{217, 138, 212, 21}}},
		{Region: "Japan Tokyo st006", IPs: []net.IP{{82, 102, 28, 123}}},
		{Region: "Japan Tokyo st007", IPs: []net.IP{{82, 102, 28, 125}}},
		{Region: "Japan Tokyo st008", IPs: []net.IP{{89, 187, 161, 12}}},
		{Region: "Japan Tokyo st009", IPs: []net.IP{{89, 187, 161, 14}}},
		{Region: "Japan Tokyo st010", IPs: []net.IP{{89, 187, 161, 17}}},
		{Region: "Japan Tokyo st011", IPs: []net.IP{{89, 187, 161, 19}}},
		{Region: "Japan Tokyo st012", IPs: []net.IP{{89, 187, 161, 7}}},
		{Region: "Japan Tokyo st013", IPs: []net.IP{{89, 187, 161, 9}}},
		{Region: "Kazakhstan", IPs: []net.IP{{95, 57, 207, 200}}},
		{Region: "Korea", IPs: []net.IP{{45, 130, 137, 3}, {45, 130, 137, 10}, {45, 130, 137, 16}, {45, 130, 137, 26}, {45, 130, 137, 32}, {45, 130, 137, 46}, {45, 130, 137, 48}, {45, 130, 137, 50}}},
		{Region: "Latvia", IPs: []net.IP{{188, 92, 78, 140}, {188, 92, 78, 142}, {188, 92, 78, 145}, {188, 92, 78, 150}}},
		{Region: "Luxembourg", IPs: []net.IP{{185, 153, 151, 73}, {185, 153, 151, 80}, {185, 153, 151, 98}, {185, 153, 151, 100}, {185, 153, 151, 116}, {185, 153, 151, 118}, {185, 153, 151, 126}, {185, 153, 151, 160}}},
		{Region: "Malaysia", IPs: []net.IP{{42, 0, 30, 158}, {42, 0, 30, 164}, {42, 0, 30, 179}, {42, 0, 30, 181}, {42, 0, 30, 183}, {42, 0, 30, 209}}},
		{Region: "Mexico City Mexico", IPs: []net.IP{{194, 41, 112, 14}, {194, 41, 112, 30}, {194, 41, 112, 33}, {194, 41, 112, 35}, {194, 41, 112, 37}, {194, 41, 112, 39}}},
		{Region: "Moldova", IPs: []net.IP{{178, 175, 128, 235}, {178, 175, 128, 237}}},
		{Region: "Netherlands Amsterdam", IPs: []net.IP{{81, 19, 208, 56}, {81, 19, 209, 59}, {89, 46, 223, 54}, {89, 46, 223, 60}, {89, 46, 223, 84}, {143, 244, 42, 74}, {178, 239, 173, 51}, {212, 102, 35, 216}}},
		{Region: "Netherlands Amsterdam mp001", IPs: []net.IP{{188, 166, 43, 117}}},
		{Region: "Nigeria", IPs: []net.IP{{102, 165, 23, 4}, {102, 165, 23, 6}, {102, 165, 23, 42}, {102, 165, 23, 44}}},
		{Region: "North Macedonia", IPs: []net.IP{{185, 225, 28, 67}, {185, 225, 28, 83}, {185, 225, 28, 91}, {185, 225, 28, 99}, {185, 225, 28, 101}, {185, 225, 28, 107}, {185, 225, 28, 109}, {185, 225, 28, 245}}},
		{Region: "Norway", IPs: []net.IP{{45, 12, 223, 197}, {45, 12, 223, 213}, {84, 247, 50, 27}, {84, 247, 50, 29}, {91, 219, 215, 53}, {91, 219, 215, 69}, {95, 174, 66, 37}, {95, 174, 66, 41}}},
		{Region: "Paraguay", IPs: []net.IP{{181, 40, 18, 47}, {181, 40, 18, 59}, {186, 16, 32, 163}, {186, 16, 32, 168}, {186, 16, 32, 173}}},
		{Region: "Philippines", IPs: []net.IP{{45, 134, 224, 3}, {45, 134, 224, 8}, {45, 134, 224, 18}, {45, 134, 224, 20}}},
		{Region: "Poland Gdansk", IPs: []net.IP{{5, 133, 14, 198}, {5, 187, 49, 147}, {5, 187, 53, 53}, {5, 187, 53, 55}, {178, 255, 44, 69}, {178, 255, 45, 187}}},
		{Region: "Poland Warsaw", IPs: []net.IP{{5, 253, 206, 67}, {5, 253, 206, 71}, {5, 253, 206, 227}, {5, 253, 206, 229}, {84, 17, 55, 132}, {84, 17, 55, 134}, {185, 246, 208, 77}, {185, 246, 208, 105}}},
		{Region: "Portugal Loule", IPs: []net.IP{{176, 61, 146, 97}, {176, 61, 146, 108}, {176, 61, 146, 113}, {176, 61, 146, 118}}},
		{Region: "Portugal Porto", IPs: []net.IP{{194, 39, 127, 171}, {194, 39, 127, 191}, {194, 39, 127, 193}, {194, 39, 127, 231}, {194, 39, 127, 233}, {194, 39, 127, 240}, {194, 39, 127, 244}}},
		{Region: "Romania", IPs: []net.IP{{45, 89, 175, 55}, {86, 106, 137, 147}, {185, 102, 217, 157}, {185, 102, 217, 159}, {185, 102, 217, 167}, {185, 102, 217, 169}, {185, 102, 217, 194}, {185, 102, 217, 196}}},
		{Region: "Russia St. Petersburg", IPs: []net.IP{{185, 246, 88, 66}, {185, 246, 88, 118}}},
		{Region: "Serbia", IPs: []net.IP{{37, 120, 193, 51}, {152, 89, 160, 213}, {152, 89, 160, 215}}},
		{Region: "Singapore", IPs: []net.IP{{89, 187, 162, 184}, {89, 187, 162, 186}, {89, 187, 163, 130}, {89, 187, 163, 134}, {89, 187, 163, 136}, {89, 187, 163, 195}, {89, 187, 163, 197}, {89, 187, 163, 207}}},
		{Region: "Singapore in", IPs: []net.IP{{128, 199, 193, 35}}},
		{Region: "Singapore mp001", IPs: []net.IP{{206, 189, 94, 229}}},
		{Region: "Singapore st001", IPs: []net.IP{{217, 138, 201, 91}}},
		{Region: "Singapore st002", IPs: []net.IP{{217, 138, 201, 93}}},
		{Region: "Singapore st003", IPs: []net.IP{{84, 247, 49, 19}}},
		{Region: "Singapore st004", IPs: []net.IP{{84, 247, 49, 21}}},
		{Region: "Slovekia", IPs: []net.IP{{37, 120, 221, 3}, {185, 76, 8, 210}, {185, 76, 8, 212}, {185, 76, 8, 215}, {185, 76, 8, 217}, {193, 37, 255, 35}, {193, 37, 255, 37}, {193, 37, 255, 39}}},
		{Region: "Slovenia", IPs: []net.IP{{195, 158, 249, 36}, {195, 158, 249, 38}, {195, 158, 249, 40}, {195, 158, 249, 46}}},
		{Region: "South Africa", IPs: []net.IP{{102, 165, 47, 132}, {154, 16, 93, 51}, {154, 16, 93, 53}, {154, 127, 49, 230}, {154, 127, 49, 232}, {154, 127, 50, 138}}},
		{Region: "Spain Barcelona", IPs: []net.IP{{37, 120, 142, 179}, {37, 120, 142, 181}, {185, 188, 61, 7}, {185, 188, 61, 23}, {185, 188, 61, 37}, {185, 188, 61, 41}}},
		{Region: "Spain Madrid", IPs: []net.IP{{37, 120, 148, 229}, {89, 37, 95, 11}, {89, 37, 95, 27}, {188, 208, 141, 18}, {188, 208, 141, 100}, {212, 102, 48, 4}, {212, 102, 48, 18}, {212, 102, 48, 20}}},
		{Region: "Spain Valencia", IPs: []net.IP{{196, 196, 150, 67}, {196, 196, 150, 71}, {196, 196, 150, 83}, {196, 196, 150, 85}}},
		{Region: "Sweden", IPs: []net.IP{{185, 76, 9, 34}, {185, 76, 9, 39}, {185, 76, 9, 41}, {185, 76, 9, 51}, {185, 76, 9, 55}, {185, 76, 9, 57}}},
		{Region: "Switzerland", IPs: []net.IP{{45, 12, 222, 243}, {84, 17, 53, 86}, {84, 17, 53, 166}, {84, 17, 53, 210}, {84, 17, 53, 219}, {84, 17, 53, 223}, {156, 146, 62, 41}, {156, 146, 62, 56}}},
		{Region: "Taiwan", IPs: []net.IP{{2, 58, 242, 43}, {2, 58, 242, 157}, {103, 152, 151, 5}, {103, 152, 151, 19}, {103, 152, 151, 69}, {103, 152, 151, 83}}},
		{Region: "Turkey Istanbul", IPs: []net.IP{{107, 150, 95, 149}, {107, 150, 95, 157}, {107, 150, 95, 163}, {107, 150, 95, 165}}},
		{Region: "UK Glasgow", IPs: []net.IP{{185, 108, 105, 5}, {185, 108, 105, 7}, {185, 108, 105, 38}, {185, 108, 105, 151}, {185, 108, 105, 153}, {185, 108, 105, 170}, {185, 108, 105, 174}, {185, 108, 105, 182}}},
		{Region: "UK London", IPs: []net.IP{{37, 10, 114, 70}, {89, 35, 29, 71}, {185, 16, 206, 116}, {185, 44, 76, 55}, {185, 44, 78, 90}, {185, 114, 224, 119}, {185, 141, 206, 182}, {188, 240, 71, 179}}},
		{Region: "UK London mp001", IPs: []net.IP{{206, 189, 119, 92}}},
		{Region: "UK London st001", IPs: []net.IP{{217, 146, 82, 83}}},
		{Region: "UK London st002", IPs: []net.IP{{185, 134, 22, 80}}},
		{Region: "UK London st003", IPs: []net.IP{{185, 134, 22, 92}}},
		{Region: "UK London st004", IPs: []net.IP{{185, 44, 76, 186}}},
		{Region: "UK London st005", IPs: []net.IP{{185, 44, 76, 188}}},
		{Region: "UK Manchester", IPs: []net.IP{{37, 120, 200, 5}, {37, 120, 200, 117}, {89, 238, 130, 235}, {91, 90, 121, 131}, {91, 90, 121, 149}, {194, 37, 98, 37}, {194, 37, 98, 219}, {217, 138, 196, 3}}},
		{Region: "US Bend", IPs: []net.IP{{45, 43, 14, 73}, {45, 43, 14, 75}, {45, 43, 14, 85}, {45, 43, 14, 93}, {45, 43, 14, 95}, {45, 43, 14, 105}, {154, 16, 168, 186}}},
		{Region: "US Boston", IPs: []net.IP{{173, 237, 207, 32}, {173, 237, 207, 42}, {173, 237, 207, 60}, {192, 34, 83, 230}, {192, 34, 83, 236}, {199, 217, 107, 20}}},
		{Region: "US Charlotte", IPs: []net.IP{{154, 16, 171, 195}, {154, 16, 171, 197}, {154, 16, 171, 206}, {155, 254, 29, 165}, {155, 254, 31, 182}, {192, 154, 253, 67}, {192, 154, 254, 135}}},
		{Region: "US Chicago", IPs: []net.IP{{74, 119, 146, 181}, {107, 152, 100, 26}, {143, 244, 60, 167}, {143, 244, 60, 169}, {184, 170, 250, 72}, {184, 170, 250, 154}}},
		{Region: "US Dallas", IPs: []net.IP{{66, 115, 177, 133}, {66, 115, 177, 138}, {66, 115, 177, 146}, {66, 115, 177, 151}, {66, 115, 177, 153}, {66, 115, 177, 156}, {89, 187, 175, 165}, {212, 102, 40, 76}}},
		{Region: "US Denver", IPs: []net.IP{{174, 128, 245, 149}, {212, 102, 44, 68}, {212, 102, 44, 71}, {212, 102, 44, 83}, {212, 102, 44, 91}, {212, 102, 44, 98}}},
		{Region: "US Gahanna", IPs: []net.IP{{104, 244, 208, 37}, {104, 244, 208, 107}, {104, 244, 209, 53}, {104, 244, 209, 101}, {104, 244, 210, 115}, {104, 244, 211, 141}}},
		{Region: "US Houston", IPs: []net.IP{{104, 148, 30, 37}, {104, 148, 30, 83}, {199, 10, 64, 67}, {199, 10, 64, 69}, {199, 10, 64, 99}, {199, 10, 64, 179}}},
		{Region: "US Kansas City", IPs: []net.IP{{63, 141, 236, 243}, {63, 141, 236, 245}, {69, 30, 249, 123}, {173, 208, 149, 197}, {173, 208, 202, 59}, {173, 208, 202, 61}, {198, 204, 231, 147}, {198, 204, 231, 149}}},
		{Region: "US Las Vegas", IPs: []net.IP{{45, 89, 173, 203}, {79, 110, 54, 125}, {79, 110, 54, 131}, {89, 187, 187, 147}, {89, 187, 187, 149}, {185, 242, 5, 155}, {185, 242, 5, 211}, {185, 242, 5, 213}}},
		{Region: "US Latham", IPs: []net.IP{{45, 43, 19, 74}, {45, 43, 19, 84}, {45, 43, 19, 90}, {154, 16, 169, 3}, {154, 16, 169, 7}}},
		{Region: "US Los Angeles", IPs: []net.IP{{84, 17, 45, 249}, {138, 199, 9, 193}, {138, 199, 9, 199}, {138, 199, 9, 209}, {172, 83, 44, 83}, {184, 170, 243, 215}, {192, 111, 134, 69}, {192, 111, 134, 202}}},
		{Region: "US Miami", IPs: []net.IP{{89, 187, 173, 201}, {107, 181, 164, 211}, {172, 83, 42, 3}, {172, 83, 42, 5}, {172, 83, 42, 55}, {172, 83, 42, 141}}},
		{Region: "US New York City", IPs: []net.IP{{84, 17, 35, 71}, {84, 17, 35, 86}, {138, 199, 40, 169}, {138, 199, 40, 179}, {172, 98, 75, 35}, {192, 40, 59, 227}, {192, 40, 59, 240}, {199, 36, 221, 85}}},
		{Region: "US New York City mp001", IPs: []net.IP{{45, 55, 60, 159}}},
		{Region: "US New York City st001", IPs: []net.IP{{92, 119, 177, 19}}},
		{Region: "US New York City st002", IPs: []net.IP{{92, 119, 177, 21}}},
		{Region: "US New York City st003", IPs: []net.IP{{92, 119, 177, 23}}},
		{Region: "US New York City st004", IPs: []net.IP{{193, 148, 18, 51}}},
		{Region: "US New York City st005", IPs: []net.IP{{193, 148, 18, 53}}},
		{Region: "US Orlando", IPs: []net.IP{{66, 115, 182, 74}, {198, 147, 22, 83}, {198, 147, 22, 85}, {198, 147, 22, 87}, {198, 147, 22, 131}, {198, 147, 22, 147}, {198, 147, 22, 163}, {198, 147, 22, 211}}},
		{Region: "US Phoenix", IPs: []net.IP{{107, 181, 184, 117}, {199, 58, 187, 3}, {199, 58, 187, 5}, {199, 58, 187, 8}, {199, 58, 187, 15}, {199, 58, 187, 69}}},
		{Region: "US Saint Louis", IPs: []net.IP{{148, 72, 169, 209}, {148, 72, 169, 211}, {148, 72, 169, 213}, {148, 72, 174, 36}, {148, 72, 174, 38}, {148, 72, 174, 48}}},
		{Region: "US Salt Lake City", IPs: []net.IP{{104, 200, 131, 165}, {104, 200, 131, 167}, {104, 200, 131, 172}, {104, 200, 131, 229}, {104, 200, 131, 233}, {104, 200, 131, 245}}},
		{Region: "US San Francisco", IPs: []net.IP{{107, 181, 166, 55}, {185, 124, 240, 143}, {185, 124, 240, 151}, {185, 124, 240, 161}, {185, 124, 240, 173}, {198, 8, 81, 37}}},
		{Region: "US San Francisco mp001", IPs: []net.IP{{165, 232, 53, 25}}},
		{Region: "US Tampa", IPs: []net.IP{{209, 216, 92, 200}, {209, 216, 92, 205}, {209, 216, 92, 210}, {209, 216, 92, 215}, {209, 216, 92, 220}, {209, 216, 92, 227}}},
		{Region: "Ukraine", IPs: []net.IP{{45, 9, 238, 23}, {45, 9, 238, 38}, {176, 107, 185, 71}, {176, 107, 185, 73}}},
		{Region: "United Arab Emirates", IPs: []net.IP{{45, 9, 249, 243}, {45, 9, 249, 247}, {45, 9, 250, 101}, {176, 125, 231, 5}, {176, 125, 231, 13}, {176, 125, 231, 27}}},
		{Region: "Vietnam", IPs: []net.IP{{202, 143, 110, 29}, {202, 143, 110, 36}}},
	}
}
