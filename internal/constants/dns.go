package constants

type DNSProvider uint8

const (
	Cloudflare DNSProvider = iota
	Google
	Quad9
	CleanBrowsing
)
