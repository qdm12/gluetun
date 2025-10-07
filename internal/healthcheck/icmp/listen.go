package icmp

import (
	"context"
	"fmt"
	"net"
	"runtime"

	"golang.org/x/net/ipv4"
)

func listenICMPv4(ctx context.Context) (conn net.PacketConn, err error) {
	var listenConfig net.ListenConfig
	const listenAddress = ""
	packetConn, err := listenConfig.ListenPacket(ctx, "ip4:icmp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("listening for ICMP packets: %w", err)
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
		packetConn = ipv4ToNetPacketConn(ipv4.NewPacketConn(packetConn))
	}

	return packetConn, nil
}

func listenICMPv6(ctx context.Context) (conn net.PacketConn, err error) {
	var listenConfig net.ListenConfig
	const listenAddress = ""
	packetConn, err := listenConfig.ListenPacket(ctx, "ip6:ipv6-icmp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("listening for ICMPv6 packets: %w", err)
	}
	return packetConn, nil
}
