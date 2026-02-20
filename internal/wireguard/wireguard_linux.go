package wireguard

import (
	"net"
	"os"

	"golang.zx2c4.com/wireguard/ipc"
)

func uapiOpen(name string) (*os.File, error) {
	return ipc.UAPIOpen(name)
}

func uapiListen(interfaceName string, uapiFile *os.File) (net.Listener, error) {
	return ipc.UAPIListen(interfaceName, uapiFile)
}
