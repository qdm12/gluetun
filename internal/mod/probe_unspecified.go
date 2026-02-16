//go:build !linux

package mod

func Probe(moduleName string) error {
	panic("not implemented")
}
