package pmtud

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/netip"
	"time"

	"golang.org/x/net/icmp"
)

// If the physicalLinkMTU is zero, it defaults to 1500 which is the ethernet standard MTU.
// If the pingTimeout is zero, it defaults to 1 second.
// If the logger is nil, a no-op logger is used.
func PathMTUDiscover(ctx context.Context, ip netip.Addr,
	physicalLinkMTU int, pingTimeout time.Duration, logger Logger) (
	mtu int, err error,
) {
	if physicalLinkMTU == 0 {
		const ethernetStandardMTU = 1500
		physicalLinkMTU = ethernetStandardMTU
	}
	if pingTimeout == 0 {
		pingTimeout = time.Second
	}
	if logger == nil {
		logger = &noopLogger{}
	}

	if ip.Is4() {
		mtu, err = findIPv4NextHopMTU(ctx, ip, physicalLinkMTU, pingTimeout, logger)
		switch {
		case err == nil:
			return mtu, nil
		case errors.Is(err, net.ErrClosed): // blackhole
		default:
			return 0, fmt.Errorf("finding IPv4 next hop MTU: %w", err)
		}
	} else {
		mtu, err = getIPv6PacketTooBig(ctx, ip, physicalLinkMTU, pingTimeout, logger)
		switch {
		case err == nil:
			return mtu, nil
		case errors.Is(err, net.ErrClosed): // blackhole
		default:
			return 0, fmt.Errorf("getting IPv6 packet-too-big message: %w", err)
		}
	}

	// Fall back method: send echo requests with different packet
	// sizes and check which ones succeed to find the maximum MTU.
	minMTU := minIPv4MTU
	if ip.Is6() {
		minMTU = minIPv6MTU
	}
	return pmtudMultiSizes(ctx, ip, minMTU, physicalLinkMTU, pingTimeout, logger)
}

type pmtudTestUnit struct {
	mtu       int
	echoID    uint16
	sentBytes int
	ok        bool
}

func pmtudMultiSizes(ctx context.Context, ip netip.Addr,
	minMTU, maxPossibleMTU int, pingTimeout time.Duration,
	logger Logger,
) (maxMTU int, err error) {
	var ipVersion string
	var conn net.PacketConn
	if ip.Is4() {
		ipVersion = "v4"
		conn, err = listenICMPv4(ctx)
	} else {
		ipVersion = "v6"
		conn, err = listenICMPv6(ctx)
	}
	if err != nil {
		return 0, fmt.Errorf("listening for ICMP packets: %w", err)
	}

	mtusToTest := makeMTUsToTest(minMTU, maxPossibleMTU)
	if len(mtusToTest) == 0 {
		return minMTU, nil
	}
	logger.Debug("testing the following MTUs: %v", mtusToTest)

	tests := make([]pmtudTestUnit, len(mtusToTest))
	for i := range mtusToTest {
		tests[i] = pmtudTestUnit{mtu: mtusToTest[i]}
	}

	timedCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	go func() {
		<-timedCtx.Done()
		conn.Close()
	}()

	for i := range tests {
		id, message := buildMessageToSend(ipVersion, tests[i].mtu)
		tests[i].echoID = id

		encodedMessage, err := message.Marshal(nil)
		if err != nil {
			return 0, fmt.Errorf("encoding ICMP message: %w", err)
		}
		tests[i].sentBytes = len(encodedMessage)

		_, err = conn.WriteTo(encodedMessage, &net.IPAddr{IP: ip.AsSlice()})
		if err != nil {
			return 0, fmt.Errorf("writing ICMP message: %w", err)
		}
	}

	err = collectReplies(conn, ipVersion, tests)
	switch {
	case err == nil: // max possible MTU is working
		return tests[len(tests)-1].mtu, nil
	case err != nil && errors.Is(err, net.ErrClosed):
		// we have timeouts (IPv4 testing or IPv6 PMTUD blackholes)
		// so find the highest MTU which worked.
		// Note we start from index len(tests) - 2 since the max MTU
		// cannot be working if we had a timeout.
		for i := len(tests) - 2; i >= 0; i-- { //nolint:mnd
			if tests[i].ok {
				return pmtudMultiSizes(ctx, ip, tests[i].mtu, tests[i+1].mtu,
					pingTimeout, logger)
			}
		}

		// All MTUs failed.
		if tests[0].mtu == minMTU+1 {
			return minMTU, nil
		}
		// Re-test with MTUs between the minimum MTU
		// and the smallest next MTU we tested.
		return pmtudMultiSizes(ctx, ip, minMTU, tests[0].mtu,
			pingTimeout, logger)
	case err != nil:
		return 0, fmt.Errorf("collecting ICMP echo replies: %w", err)
	default:
		panic("unreachable")
	}
}

// Create the MTU slice of length 8 such that:
// - the first element is the minMTU plus the step
// - the last element is the maxMTU
// - elements in-between are separated as close to each other
// - Don't make the minMTU part of the MTUs to test since
// it's assumed it's already working.
func makeMTUsToTest(minMTU, maxMTU int) (mtus []int) {
	const mtusLength = 8
	diff := maxMTU - minMTU
	switch {
	case minMTU > maxMTU:
		panic("minMTU > maxMTU")
	case diff <= mtusLength:
		mtus = make([]int, 0, diff)
		for mtu := minMTU + 1; mtu <= maxMTU; mtu++ {
			mtus = append(mtus, mtu)
		}
	case diff < mtusLength*2:
		step := float64(diff) / float64(mtusLength)
		mtus = make([]int, 0, mtusLength)
		for mtu := float64(minMTU) + step; len(mtus) < mtusLength-1; mtu += step {
			mtus = append(mtus, int(math.Round(mtu)))
		}
		mtus = append(mtus, maxMTU) // last element is the maxMTU
	default:
		step := diff / mtusLength
		mtus = make([]int, 0, mtusLength)
		for mtu := minMTU + step; len(mtus) < mtusLength-1; mtu += step {
			mtus = append(mtus, mtu)
		}
		mtus = append(mtus, maxMTU) // last element is the maxMTU
	}

	return mtus
}

func collectReplies(conn net.PacketConn, ipVersion string,
	tests []pmtudTestUnit,
) (err error) {
	echoIDToTestIndex := make(map[uint16]int, len(tests))
	for i, test := range tests {
		echoIDToTestIndex[test.echoID] = i
	}

	// The theoretical limit is 4GiB for IPv6 MTU path discovery jumbograms, but that would
	// create huge buffers which we don't really want to support anyway.
	// The standard frame maximum MTU is 1500 bytes, and there are Jumbo frames with
	// a conventional maximum of 9000 bytes. However, some manufacturers support up
	// 9216-20 = 9196 bytes for the maximum MTU. We thus use buffers of size 9196 to
	// match eventual Jumbo frames. More information at:
	// https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	const maxPossibleMTU = 9196
	buffer := make([]byte, maxPossibleMTU)

	idsFound := 0
	for idsFound < len(tests) {
		// Note we need to read the whole packet in one call to ReadFrom, so the buffer
		// must be large enough to read the entire reply packet. See:
		// https://groups.google.com/g/golang-nuts/c/5dy2Q4nPs08/m/KmuSQAGEtG4J
		bytesRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			return fmt.Errorf("reading from ICMP connection: %w", err)
		}
		packetBytes := buffer[:bytesRead]

		ipPacketLength := len(packetBytes)

		var icmpProtocol int
		switch ipVersion {
		case "v4":
			icmpProtocol = icmpv4Protocol
		case "v6":
			icmpProtocol = icmpv6Protocol
		default:
			panic(fmt.Sprintf("unknown IP version: %s", ipVersion))
		}

		// Parse the ICMP message
		// Note: this parsing works for a truncated 556 bytes ICMP reply packet.
		message, err := icmp.ParseMessage(icmpProtocol, packetBytes)
		if err != nil {
			return fmt.Errorf("parsing message: %w", err)
		}

		echoBody, ok := message.Body.(*icmp.Echo)
		if !ok {
			return fmt.Errorf("%w: %T", ErrICMPBodyUnsupported, message.Body)
		}

		id := uint16(echoBody.ID) //nolint:gosec
		testIndex, testing := echoIDToTestIndex[id]
		if !testing { // not an id we expected so ignore it
			// TODO log warning
			fmt.Println("ID not sent:", echoBody.ID, message.Code, message.Type, ipPacketLength)
			continue
		}
		fmt.Println("ID ok:", echoBody.ID, message.Code, message.Type, ipPacketLength)
		idsFound++
		sentBytes := tests[testIndex].sentBytes

		// echo reply should be at most the number of bytes sent,
		// and can be lower, more precisely 556 bytes, in case
		// the host we are reaching wants to stay out of trouble
		// and ensure its echo reply goes through without
		// fragmentation, see the following page:
		// https://datatracker.ietf.org/doc/html/rfc1122#page-59
		const conservativeReplyLength = 556
		truncated := ipPacketLength < sentBytes &&
			ipPacketLength == conservativeReplyLength
		// Check the packet size is the same if the reply is not truncated
		if !truncated && sentBytes != ipPacketLength {
			return fmt.Errorf("%w: sent %dB and received %dB",
				ErrICMPEchoDataMismatch, sentBytes, ipPacketLength)
		}
		// Truncated reply or matching reply size
		tests[testIndex].ok = true
	}
	return nil
}
