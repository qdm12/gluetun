package netlink

var _ NetLinker = (*NetLink)(nil)

type NetLinker interface {
	Addresser
	Linker
	Router
	Ruler
}
