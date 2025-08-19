package pmtud

import (
	"bytes"
	"errors"
	"fmt"

	"golang.org/x/net/icmp"
)

var (
	ErrICMPNextHopMTUTooLow  = errors.New("ICMP Next Hop MTU is too low")
	ErrICMPNextHopMTUTooHigh = errors.New("ICMP Next Hop MTU is too high")
)

func checkMTU(mtu, minMTU, physicalLinkMTU int) (err error) {
	switch {
	case mtu < minMTU:
		return fmt.Errorf("%w: %d", ErrICMPNextHopMTUTooLow, mtu)
	case mtu > physicalLinkMTU:
		return fmt.Errorf("%w: %d is larger than physical link MTU %d",
			ErrICMPNextHopMTUTooHigh, mtu, physicalLinkMTU)
	default:
		return nil
	}
}

func checkInvokingReplyIDMatch(icmpProtocol int, received []byte,
	outboundMessage *icmp.Message,
) (match bool, err error) {
	inboundMessage, err := icmp.ParseMessage(icmpProtocol, received)
	if err != nil {
		return false, fmt.Errorf("parsing invoking packet: %w", err)
	}
	inboundBody, ok := inboundMessage.Body.(*icmp.Echo)
	if !ok {
		return false, fmt.Errorf("%w: %T", ErrICMPBodyUnsupported, inboundMessage.Body)
	}
	outboundBody := outboundMessage.Body.(*icmp.Echo) //nolint:forcetypeassert
	return inboundBody.ID == outboundBody.ID, nil
}

var ErrICMPIDMismatch = errors.New("ICMP id mismatch")

func checkEchoReply(icmpProtocol int, received []byte,
	outboundMessage *icmp.Message, truncatedBody bool,
) (err error) {
	inboundMessage, err := icmp.ParseMessage(icmpProtocol, received)
	if err != nil {
		return fmt.Errorf("parsing invoking packet: %w", err)
	}
	inboundBody, ok := inboundMessage.Body.(*icmp.Echo)
	if !ok {
		return fmt.Errorf("%w: %T", ErrICMPBodyUnsupported, inboundMessage.Body)
	}
	outboundBody := outboundMessage.Body.(*icmp.Echo) //nolint:forcetypeassert
	if inboundBody.ID != outboundBody.ID {
		return fmt.Errorf("%w: sent id %d and received id %d",
			ErrICMPIDMismatch, outboundBody.ID, inboundBody.ID)
	}
	err = checkEchoBodies(outboundBody.Data, inboundBody.Data, truncatedBody)
	if err != nil {
		return fmt.Errorf("checking sent and received bodies: %w", err)
	}
	return nil
}

var ErrICMPEchoDataMismatch = errors.New("ICMP data mismatch")

func checkEchoBodies(sent, received []byte, receivedTruncated bool) (err error) {
	if len(received) > len(sent) {
		return fmt.Errorf("%w: sent %d bytes and received %d bytes",
			ErrICMPEchoDataMismatch, len(sent), len(received))
	}
	if receivedTruncated {
		sent = sent[:len(received)]
	}
	if !bytes.Equal(received, sent) {
		return fmt.Errorf("%w: sent %x and received %x",
			ErrICMPEchoDataMismatch, sent, received)
	}
	return nil
}
