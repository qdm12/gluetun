package pmtud

import (
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var _ net.PacketConn = &ipv4Wrapper{}

// ipv4Wrapper is a wrapper around ipv4.PacketConn to implement
// the net.PacketConn interface. It's only used for Darwin or iOS.
type ipv4Wrapper struct {
	ipv4Conn *ipv4.PacketConn
}

func ipv4ToNetPacketConn(ipv4 *ipv4.PacketConn) *ipv4Wrapper {
	return &ipv4Wrapper{ipv4Conn: ipv4}
}

func (i *ipv4Wrapper) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, _, addr, err = i.ipv4Conn.ReadFrom(p)
	return n, addr, err
}

func (i *ipv4Wrapper) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return i.ipv4Conn.WriteTo(p, nil, addr)
}

func (i *ipv4Wrapper) Close() error {
	return i.ipv4Conn.Close()
}

func (i *ipv4Wrapper) LocalAddr() net.Addr {
	return i.ipv4Conn.LocalAddr()
}

func (i *ipv4Wrapper) SetDeadline(t time.Time) error {
	return i.ipv4Conn.SetDeadline(t)
}

func (i *ipv4Wrapper) SetReadDeadline(t time.Time) error {
	return i.ipv4Conn.SetReadDeadline(t)
}

func (i *ipv4Wrapper) SetWriteDeadline(t time.Time) error {
	return i.ipv4Conn.SetWriteDeadline(t)
}
