package netlink

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . NetLinker

var _ NetLinker = (*NetLink)(nil)

type NetLinker interface {
	Addresser
	Linker
	Router
	Ruler
	WireguardChecker
}
