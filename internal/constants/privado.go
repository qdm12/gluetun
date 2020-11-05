package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

//nolint:lll
const (
	PrivadoCertificate = "MIIFKDCCAxCgAwIBAgIJAMtrmqZxIV/OMA0GCSqGSIb3DQEBDQUAMBIxEDAOBgNVBAMMB1ByaXZhZG8wHhcNMjAwMTA4MjEyODQ1WhcNMzUwMTA5MjEyODQ1WjASMRAwDgYDVQQDDAdQcml2YWRvMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxPwOgiwNJzZTnKIXwAB0TSu/Lu2qt2U2I8obtQjwhi/7OrfmbmYykSdro70al2XPhnwAGGdCxW6LDnp0UN/IOhD11mgBPo14f5CLkBQjSJ6VN5miPbvK746LsNZl9H8rQGvDuPo4CG9BfPZMiDRGlsMxij/jztzgT1gmuxQ7WHfFRcNzBas1dHa9hV/d3TU6/t47x4SE/ljdcCtJiu7Zn6ODKQoys3mB7Luz2ngqUJWvkqsg+E4+3eJ0M8Hlbn5TPaRJBID7DAdYo6Vs6xGCYr981ThFcmoIQ10js10yANrrfGAzd03b3TnLAgko0uQMHjliMZL6L8sWOPHxyxJI0us88SFh4UgcFyRHKHPKux7w24SxAlZUYoUcTHp9VjG5XvDKYxzgV2RdM4ulBGbQRQ3y3/CyddsyQYMvA55Ets0LfPaBvDIcct70iXijGsdvlX1du3ArGpG7Vaje/RU4nbbGT6HYRdt5YyZfof288ukMOSj20nVcmS+c/4tqsxSerRb1aq5LOi1IemSkTMeC5gCbexk+L1vl7NT/58sxjGmu5bXwnvev/lIItfi2AlITrfUSEv19iDMKkeshwn/+sFJBMWYyluP+yJ56yR+MWoXvLlSWphLDTqq19yx3BZn0P1tgbXoR0g8PTdJFcz8z3RIb7myVLYulV1oGG/3rka0CAwEAAaOBgDB+MB0GA1UdDgQWBBTFtJkZCVDuDAD6k5bJzefjJdO3DTBCBgNVHSMEOzA5gBTFtJkZCVDuDAD6k5bJzefjJdO3DaEWpBQwEjEQMA4GA1UEAwwHUHJpdmFkb4IJAMtrmqZxIV/OMAwGA1UdEwQFMAMBAf8wCwYDVR0PBAQDAgEGMA0GCSqGSIb3DQEBDQUAA4ICAQB7MUSXMeBb9wlSv4sUaT1JHEwE26nlBw+TKmezfuPU5pBlY0LYr6qQZY95DHqsRJ7ByUzGUrGo17dNGXlcuNc6TAaQQEDRPo6y+LVh2TWMk15TUMI+MkqryJtCret7xGvDigKYMJgBy58HN3RAVr1B7cL9youwzLgc2Y/NcFKvnQJKeiIYAJ7g0CcnJiQvgZTS7xdwkEBXfsngmUCIG320DLPEL+Ze0HiUrxwWljMRya6i40AeH3Zu2i532xX1wV5+cjA4RJWIKg6ri/Q54iFGtZrA9/nc6y9uoQHkmz8cGyVUmJxFzMrrIICVqUtVRxLhkTMe4UzwRWTBeGgtW4tS0yq1QonAKfOyjgRw/CeY55D2UGvnAFZdTadtYXS4Alu2P9zdwoEk3fzHiVmDjqfJVr5wz9383aABUFrPI3nz6ed/Z6LZflKh1k+DUDEp8NxU4klUULWsSOKoa5zGX51G8cdHxwQLImXvtGuN5eSR8jCTgxFZhdps/xes4KkyfIz9FMYG748M+uOTgKITf4zdJ9BAyiQaOufVQZ8WjhWzWk9YHec9VqPkzpWNGkVjiRI5ewuXwZzZ164tMv2hikBXSuUCnFz37/ZNwGlDi0oBdDszCk2GxccdFHHaCSmpjU5MrdJ+5IhtTKGeTx+US2hTIVHQFIO99DmacxSYvLNcSQ=="
)

func PrivadoCityChoices() (choices []string) {
	servers := PrivadoServers()
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return choices
}

//nolint:gomnd
func PrivadoServers() []models.PrivadoServer {
	return []models.PrivadoServer{
		{City: "akl", Number: 1, IPs: []net.IP{{23, 254, 104, 114}}},
		{City: "akl", Number: 2, IPs: []net.IP{{23, 254, 104, 120}}},
		{City: "akl", Number: 3, IPs: []net.IP{{23, 254, 104, 51}}},
		{City: "ams", Number: 1, IPs: []net.IP{{91, 148, 224, 10}}},
		{City: "ams", Number: 10, IPs: []net.IP{{91, 148, 228, 20}}},
		{City: "ams", Number: 11, IPs: []net.IP{{91, 148, 228, 30}}},
		{City: "ams", Number: 12, IPs: []net.IP{{91, 148, 228, 40}}},
		{City: "ams", Number: 13, IPs: []net.IP{{91, 148, 228, 50}}},
		{City: "ams", Number: 14, IPs: []net.IP{{91, 148, 228, 60}}},
		{City: "ams", Number: 15, IPs: []net.IP{{91, 148, 228, 70}}},
		{City: "ams", Number: 16, IPs: []net.IP{{91, 148, 228, 80}}},
		{City: "ams", Number: 2, IPs: []net.IP{{91, 148, 224, 20}}},
		{City: "ams", Number: 3, IPs: []net.IP{{91, 148, 224, 30}}},
		{City: "ams", Number: 4, IPs: []net.IP{{91, 148, 224, 40}}},
		{City: "ams", Number: 5, IPs: []net.IP{{91, 148, 224, 50}}},
		{City: "ams", Number: 6, IPs: []net.IP{{91, 148, 224, 60}}},
		{City: "ams", Number: 7, IPs: []net.IP{{91, 148, 224, 70}}},
		{City: "ams", Number: 8, IPs: []net.IP{{91, 148, 224, 80}}},
		{City: "ams", Number: 9, IPs: []net.IP{{91, 148, 228, 10}}},
		{City: "arn", Number: 1, IPs: []net.IP{{86, 106, 103, 67}}},
		{City: "arn", Number: 2, IPs: []net.IP{{86, 106, 103, 74}}},
		{City: "arn", Number: 3, IPs: []net.IP{{86, 106, 103, 81}}},
		{City: "ath", Number: 1, IPs: []net.IP{{188, 123, 126, 61}}},
		{City: "ath", Number: 2, IPs: []net.IP{{188, 123, 126, 64}}},
		{City: "ath", Number: 3, IPs: []net.IP{{188, 123, 126, 68}}},
		{City: "ath", Number: 4, IPs: []net.IP{{188, 123, 126, 72}}},
		{City: "beg", Number: 1, IPs: []net.IP{{89, 38, 224, 19}}},
		{City: "beg", Number: 2, IPs: []net.IP{{89, 38, 224, 25}}},
		{City: "bkk", Number: 1, IPs: []net.IP{{119, 59, 111, 3}}},
		{City: "bkk", Number: 2, IPs: []net.IP{{119, 59, 111, 11}}},
		{City: "bom", Number: 1, IPs: []net.IP{{103, 26, 204, 61}}},
		{City: "bom", Number: 2, IPs: []net.IP{{103, 26, 204, 70}}},
		{City: "bru", Number: 1, IPs: []net.IP{{217, 138, 211, 163}}},
		{City: "bru", Number: 2, IPs: []net.IP{{217, 138, 211, 170}}},
		{City: "bru", Number: 3, IPs: []net.IP{{217, 138, 211, 177}}},
		{City: "bru", Number: 4, IPs: []net.IP{{217, 138, 211, 184}}},
		{City: "bts", Number: 1, IPs: []net.IP{{37, 120, 221, 227}}},
		{City: "bts", Number: 2, IPs: []net.IP{{37, 120, 221, 233}}},
		{City: "bud", Number: 1, IPs: []net.IP{{185, 128, 26, 194}}},
		{City: "bud", Number: 2, IPs: []net.IP{{185, 128, 26, 200}}},
		{City: "cdg", Number: 1, IPs: []net.IP{{89, 40, 183, 99}}},
		{City: "cdg", Number: 2, IPs: []net.IP{{89, 40, 183, 106}}},
		{City: "cdg", Number: 3, IPs: []net.IP{{89, 40, 183, 113}}},
		{City: "cdg", Number: 4, IPs: []net.IP{{89, 40, 183, 120}}},
		{City: "cph", Number: 1, IPs: []net.IP{{2, 58, 46, 35}}},
		{City: "cph", Number: 2, IPs: []net.IP{{2, 58, 46, 42}}},
		{City: "cph", Number: 3, IPs: []net.IP{{2, 58, 46, 49}}},
		{City: "cph", Number: 4, IPs: []net.IP{{2, 58, 46, 56}}},
		{City: "dca", Number: 1, IPs: []net.IP{{85, 12, 61, 10}}},
		{City: "dca", Number: 13, IPs: []net.IP{{185, 247, 68, 3}}},
		{City: "dca", Number: 14, IPs: []net.IP{{185, 247, 68, 10}}},
		{City: "dca", Number: 15, IPs: []net.IP{{185, 247, 68, 17}}},
		{City: "dca", Number: 16, IPs: []net.IP{{185, 247, 68, 24}}},
		{City: "dca", Number: 2, IPs: []net.IP{{85, 12, 61, 20}}},
		{City: "dca", Number: 3, IPs: []net.IP{{85, 12, 61, 30}}},
		{City: "dca", Number: 4, IPs: []net.IP{{85, 12, 61, 40}}},
		{City: "dca", Number: 5, IPs: []net.IP{{85, 12, 61, 50}}},
		{City: "dca", Number: 6, IPs: []net.IP{{85, 12, 61, 60}}},
		{City: "dca", Number: 7, IPs: []net.IP{{85, 12, 61, 70}}},
		{City: "dca", Number: 8, IPs: []net.IP{{85, 12, 61, 80}}},
		{City: "dfw", Number: 1, IPs: []net.IP{{23, 105, 32, 243}}},
		{City: "dfw", Number: 2, IPs: []net.IP{{23, 105, 32, 244}}},
		{City: "dub", Number: 1, IPs: []net.IP{{84, 247, 48, 227}}},
		{City: "dub", Number: 2, IPs: []net.IP{{84, 247, 48, 234}}},
		{City: "dub", Number: 3, IPs: []net.IP{{84, 247, 48, 241}}},
		{City: "dub", Number: 4, IPs: []net.IP{{84, 247, 48, 248}}},
		{City: "eze", Number: 1, IPs: []net.IP{{168, 205, 93, 211}}},
		{City: "eze", Number: 2, IPs: []net.IP{{168, 205, 93, 217}}},
		{City: "fra", Number: 1, IPs: []net.IP{{91, 148, 232, 10}}},
		{City: "fra", Number: 2, IPs: []net.IP{{91, 148, 232, 20}}},
		{City: "fra", Number: 3, IPs: []net.IP{{91, 148, 232, 30}}},
		{City: "fra", Number: 4, IPs: []net.IP{{91, 148, 232, 40}}},
		{City: "fra", Number: 5, IPs: []net.IP{{91, 148, 233, 7}}},
		{City: "fra", Number: 6, IPs: []net.IP{{91, 148, 233, 8}}},
		{City: "fra", Number: 7, IPs: []net.IP{{91, 148, 233, 9}}},
		{City: "fra", Number: 8, IPs: []net.IP{{91, 148, 233, 10}}},
		{City: "gru", Number: 1, IPs: []net.IP{{177, 54, 145, 193}}},
		{City: "gru", Number: 2, IPs: []net.IP{{177, 54, 145, 197}}},
		{City: "hel", Number: 1, IPs: []net.IP{{194, 34, 134, 219}}},
		{City: "hel", Number: 2, IPs: []net.IP{{194, 34, 134, 227}}},
		{City: "hkg", Number: 1, IPs: []net.IP{{209, 58, 185, 88}}},
		{City: "hkg", Number: 2, IPs: []net.IP{{209, 58, 185, 97}}},
		{City: "hkg", Number: 3, IPs: []net.IP{{209, 58, 185, 108}}},
		{City: "hkg", Number: 4, IPs: []net.IP{{209, 58, 185, 120}}},
		{City: "icn", Number: 1, IPs: []net.IP{{169, 56, 73, 146}}},
		{City: "icn", Number: 2, IPs: []net.IP{{169, 56, 73, 153}}},
		{City: "iev", Number: 1, IPs: []net.IP{{176, 103, 52, 40}}},
		{City: "iev", Number: 2, IPs: []net.IP{{176, 103, 53, 40}}},
		{City: "ist", Number: 1, IPs: []net.IP{{185, 84, 183, 3}}},
		{City: "ist", Number: 2, IPs: []net.IP{{185, 84, 183, 4}}},
		{City: "jfk", Number: 1, IPs: []net.IP{{217, 138, 208, 99}}},
		{City: "jfk", Number: 2, IPs: []net.IP{{217, 138, 208, 106}}},
		{City: "jfk", Number: 3, IPs: []net.IP{{217, 138, 208, 113}}},
		{City: "jfk", Number: 4, IPs: []net.IP{{217, 138, 208, 120}}},
		{City: "jnb", Number: 1, IPs: []net.IP{{172, 107, 93, 131}}},
		{City: "jnb", Number: 2, IPs: []net.IP{{172, 107, 93, 137}}},
		{City: "lax", Number: 10, IPs: []net.IP{{45, 152, 182, 234}}},
		{City: "lax", Number: 11, IPs: []net.IP{{45, 152, 182, 241}}},
		{City: "lax", Number: 12, IPs: []net.IP{{45, 152, 182, 248}}},
		{City: "lax", Number: 9, IPs: []net.IP{{45, 152, 182, 227}}},
		{City: "lis", Number: 1, IPs: []net.IP{{89, 26, 243, 153}}},
		{City: "lis", Number: 2, IPs: []net.IP{{89, 26, 243, 154}}},
		{City: "lon", Number: 1, IPs: []net.IP{{217, 138, 195, 163}}},
		{City: "lon", Number: 2, IPs: []net.IP{{217, 138, 195, 170}}},
		{City: "lon", Number: 3, IPs: []net.IP{{217, 138, 195, 177}}},
		{City: "lon", Number: 4, IPs: []net.IP{{217, 138, 195, 184}}},
		{City: "mad", Number: 1, IPs: []net.IP{{217, 138, 218, 131}}},
		{City: "man", Number: 1, IPs: []net.IP{{217, 138, 196, 131}}},
		{City: "man", Number: 2, IPs: []net.IP{{217, 138, 196, 138}}},
		{City: "man", Number: 3, IPs: []net.IP{{217, 138, 196, 145}}},
		{City: "man", Number: 4, IPs: []net.IP{{217, 138, 196, 152}}},
		{City: "mex", Number: 1, IPs: []net.IP{{169, 57, 96, 52}}},
		{City: "mex", Number: 2, IPs: []net.IP{{169, 57, 96, 57}}},
		{City: "mia", Number: 1, IPs: []net.IP{{86, 106, 87, 131}}},
		{City: "mia", Number: 2, IPs: []net.IP{{86, 106, 87, 138}}},
		{City: "mia", Number: 3, IPs: []net.IP{{86, 106, 87, 145}}},
		{City: "mia", Number: 4, IPs: []net.IP{{86, 106, 87, 152}}},
		{City: "mxp", Number: 1, IPs: []net.IP{{89, 40, 182, 195}}},
		{City: "mxp", Number: 2, IPs: []net.IP{{89, 40, 182, 201}}},
		{City: "nrt", Number: 1, IPs: []net.IP{{217, 138, 252, 3}}},
		{City: "nrt", Number: 2, IPs: []net.IP{{217, 138, 252, 10}}},
		{City: "nrt", Number: 3, IPs: []net.IP{{217, 138, 252, 17}}},
		{City: "nrt", Number: 4, IPs: []net.IP{{217, 138, 252, 24}}},
		{City: "ord", Number: 1, IPs: []net.IP{{23, 108, 95, 129}}},
		{City: "ord", Number: 2, IPs: []net.IP{{23, 108, 95, 167}}},
		{City: "osl", Number: 1, IPs: []net.IP{{84, 247, 50, 115}}},
		{City: "osl", Number: 2, IPs: []net.IP{{84, 247, 50, 119}}},
		{City: "osl", Number: 3, IPs: []net.IP{{84, 247, 50, 123}}},
		{City: "otp", Number: 1, IPs: []net.IP{{89, 46, 102, 179}}},
		{City: "otp", Number: 2, IPs: []net.IP{{89, 46, 102, 185}}},
		{City: "phx", Number: 1, IPs: []net.IP{{91, 148, 236, 10}}},
		{City: "phx", Number: 2, IPs: []net.IP{{91, 148, 236, 20}}},
		{City: "phx", Number: 3, IPs: []net.IP{{91, 148, 236, 30}}},
		{City: "phx", Number: 4, IPs: []net.IP{{91, 148, 236, 40}}},
		{City: "phx", Number: 5, IPs: []net.IP{{91, 148, 236, 50}}},
		{City: "phx", Number: 6, IPs: []net.IP{{91, 148, 236, 60}}},
		{City: "phx", Number: 7, IPs: []net.IP{{91, 148, 236, 70}}},
		{City: "phx", Number: 8, IPs: []net.IP{{91, 148, 236, 80}}},
		{City: "prg", Number: 1, IPs: []net.IP{{185, 216, 35, 99}}},
		{City: "prg", Number: 2, IPs: []net.IP{{185, 216, 35, 105}}},
		{City: "rix", Number: 1, IPs: []net.IP{{109, 248, 149, 35}}},
		{City: "rix", Number: 2, IPs: []net.IP{{109, 248, 149, 40}}},
		{City: "rkv", Number: 1, IPs: []net.IP{{82, 221, 131, 78}}},
		{City: "rkv", Number: 2, IPs: []net.IP{{82, 221, 131, 127}}},
		{City: "sea", Number: 1, IPs: []net.IP{{23, 81, 208, 96}}},
		{City: "sea", Number: 2, IPs: []net.IP{{23, 81, 208, 104}}},
		{City: "sin", Number: 1, IPs: []net.IP{{92, 119, 178, 131}}},
		{City: "sin", Number: 2, IPs: []net.IP{{92, 119, 178, 138}}},
		{City: "sin", Number: 3, IPs: []net.IP{{92, 119, 178, 145}}},
		{City: "sin", Number: 4, IPs: []net.IP{{92, 119, 178, 152}}},
		{City: "sof", Number: 1, IPs: []net.IP{{217, 138, 221, 163}}},
		{City: "sof", Number: 2, IPs: []net.IP{{217, 138, 221, 169}}},
		{City: "stl", Number: 1, IPs: []net.IP{{148, 72, 170, 145}}},
		{City: "stl", Number: 2, IPs: []net.IP{{148, 72, 172, 82}}},
		{City: "syd", Number: 1, IPs: []net.IP{{93, 115, 35, 35}}},
		{City: "syd", Number: 2, IPs: []net.IP{{93, 115, 35, 42}}},
		{City: "syd", Number: 3, IPs: []net.IP{{93, 115, 35, 49}}},
		{City: "syd", Number: 4, IPs: []net.IP{{93, 115, 35, 56}}},
		{City: "vie", Number: 1, IPs: []net.IP{{5, 253, 207, 227}}},
		{City: "vie", Number: 2, IPs: []net.IP{{5, 253, 207, 234}}},
		{City: "vie", Number: 3, IPs: []net.IP{{5, 253, 207, 241}}},
		{City: "vie", Number: 4, IPs: []net.IP{{5, 253, 207, 248}}},
		{City: "vno", Number: 1, IPs: []net.IP{{185, 64, 104, 176}}},
		{City: "vno", Number: 2, IPs: []net.IP{{185, 64, 104, 180}}},
		{City: "waw", Number: 1, IPs: []net.IP{{217, 138, 209, 163}}},
		{City: "waw", Number: 2, IPs: []net.IP{{217, 138, 209, 164}}},
		{City: "waw", Number: 3, IPs: []net.IP{{217, 138, 209, 165}}},
		{City: "waw", Number: 4, IPs: []net.IP{{217, 138, 209, 166}}},
		{City: "yul", Number: 1, IPs: []net.IP{{217, 138, 213, 67}}},
		{City: "yul", Number: 2, IPs: []net.IP{{217, 138, 213, 74}}},
		{City: "yul", Number: 3, IPs: []net.IP{{217, 138, 213, 81}}},
		{City: "yul", Number: 4, IPs: []net.IP{{217, 138, 213, 88}}},
		{City: "yvr", Number: 1, IPs: []net.IP{{71, 19, 248, 57}}},
		{City: "yvr", Number: 2, IPs: []net.IP{{71, 19, 248, 113}}},
		{City: "yyz", Number: 3, IPs: []net.IP{{199, 189, 27, 19}}},
		{City: "zrh", Number: 1, IPs: []net.IP{{185, 156, 175, 195}}},
		{City: "zrh", Number: 2, IPs: []net.IP{{185, 156, 175, 202}}},
		{City: "zrh", Number: 3, IPs: []net.IP{{185, 156, 175, 209}}},
		{City: "zrh", Number: 4, IPs: []net.IP{{185, 156, 175, 216}}},
	}
}
