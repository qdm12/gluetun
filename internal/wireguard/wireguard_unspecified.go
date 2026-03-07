//go:build !linux

package wireguard

import (
	"net"
	"os"
)

func uapiOpen(name string) (*os.File, error) {
	panic("not implemented")
}

func uapiListen(interfaceName string, uapiFile *os.File) (net.Listener, error) {
	panic("not implemented")
}
