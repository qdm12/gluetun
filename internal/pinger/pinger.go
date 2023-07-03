package pinger

import (
	"fmt"
	"net"
	"time"
)

const (
	icmpv6EchoRequestHeader = "\x80\x00\x00\x00\x00\x00\x00\x00"
)

func Ping() (bool, error) {
	ipAddr, err := net.ResolveIPAddr("ip6", "ipv6.test-ipv6.com")
	if err != nil {
		return false, fmt.Errorf("failed to resolve IP address: %w", err)
	}
	conn, err := net.Dial("udp6", ipAddr.IP.String())
	if err != nil {
		return false, fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			fmt.Printf("Error closing UDP connection: %v\n", cerr)
		}
	}()

	// Set the timeout for receiving ICMPv6 echo reply packets
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return false, fmt.Errorf("failed to set read deadline: %w", err)
	}

	// Prepare the ICMPv6 echo request packet
	icmpv6EchoRequestPacket := icmpv6EchoRequestHeader

	// Send the ICMPv6 echo request packet
	_, err = conn.Write([]byte(icmpv6EchoRequestPacket))
	if err != nil {
		return false, fmt.Errorf("failed to send ICMPv6 echo request packet: %w", err)
	}

	// Receive the ICMPv6 echo reply packet
	replyPacket := make([]byte, 1500) // Maximum packet size
	rSize, err := conn.Read(replyPacket)
	if err != nil {
		return false, fmt.Errorf("failed to receive ICMPv6 echo reply packet: %w", err)
	}

	if rSize < 8 {
		return false, fmt.Errorf("invalid ICMPv6 echo reply packet (too short)")
	}

	if replyPacket[0] != 129 {
		return false, fmt.Errorf("invalid ICMPv6 echo reply packet (type = %d)", replyPacket[0])
	}

	if replyPacket[1] != 0 {
		return false, fmt.Errorf("invalid ICMPv6 echo reply packet (code = %d)", replyPacket[1])
	}
	if rSize < 8+len(icmpv6EchoRequestHeader) {
		return false, fmt.Errorf("invalid ICMPv6 echo reply packet (too short payload)")
	}

	if string(replyPacket[8:8+len(icmpv6EchoRequestHeader)]) != icmpv6EchoRequestHeader {
		return false, fmt.Errorf("invalid ICMPv6 echo reply packet (payload does not match)")
	}

	fmt.Println("IPv6 ping successful!")
	return true, nil
}
