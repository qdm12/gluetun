package natpmp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"
)

var (
	ErrGatewayIPUnspecified = errors.New("gateway IP is unspecified")
	ErrConnectionTimeout    = errors.New("connection timeout")
)

func (c *Client) rpc(ctx context.Context, gateway netip.Addr,
	request []byte, responseSize uint) (
	response []byte, err error) {
	if gateway.IsUnspecified() || !gateway.IsValid() {
		return nil, fmt.Errorf("%w", ErrGatewayIPUnspecified)
	}

	err = checkRequest(request)
	if err != nil {
		return nil, fmt.Errorf("checking request: %w", err)
	}

	gatewayAddress := &net.UDPAddr{
		IP:   gateway.AsSlice(),
		Port: int(c.serverPort),
	}

	connection, err := net.DialUDP("udp", nil, gatewayAddress)
	if err != nil {
		return nil, fmt.Errorf("dialing udp: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	endGoroutineDone := make(chan struct{})
	defer func() {
		cancel()
		<-endGoroutineDone
	}()
	go func() {
		defer close(endGoroutineDone)
		// Context is canceled either by the parent context or
		// when this function returns.
		<-ctx.Done()
		closeErr := connection.Close()
		if closeErr == nil {
			return
		}
		if err == nil {
			err = fmt.Errorf("closing connection: %w", closeErr)
			return
		}
		err = fmt.Errorf("%w; closing connection: %w", err, closeErr)
	}()

	const maxResponseSize = 16
	response = make([]byte, maxResponseSize)

	// Retry duration doubles on every network error
	// Note it does not double if the source IP mismatches the gateway IP.
	retryDuration := c.initialRetry

	var totalRetryDuration time.Duration

	var retryCount uint
	for retryCount = 0; retryCount < c.maxRetries; retryCount++ {
		deadline := time.Now().Add(retryDuration)
		err = connection.SetDeadline(deadline)
		if err != nil {
			return nil, fmt.Errorf("setting connection deadline: %w", err)
		}

		_, err = connection.Write(request)
		if err != nil {
			return nil, fmt.Errorf("writing to connection: %w", err)
		}

		bytesRead, receivedRemoteAddress, err := connection.ReadFromUDP(response)
		if err != nil {
			if ctx.Err() != nil {
				return nil, fmt.Errorf("reading from udp connection: %w", ctx.Err())
			}
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				totalRetryDuration += retryDuration
				retryDuration *= 2
				continue
			}
			return nil, fmt.Errorf("reading from udp connection: %w", err)
		}

		if !receivedRemoteAddress.IP.Equal(gatewayAddress.IP) {
			// Upon receiving a response packet, the client MUST check the source IP
			// address, and silently discard the packet if the address is not the
			// address of the gateway to which the request was sent.
			continue
		}

		response = response[:bytesRead]
		break
	}

	if retryCount == c.maxRetries {
		return nil, fmt.Errorf("%w: after %s",
			ErrConnectionTimeout, totalRetryDuration)
	}

	// Opcodes between 0 and 127 are client requests.  Opcodes from 128 to
	// 255 are corresponding server responses.
	const operationCodeMask = 128
	expectedOperationCode := request[1] | operationCodeMask
	err = checkResponse(response, expectedOperationCode, responseSize)
	if err != nil {
		return nil, fmt.Errorf("checking response: %w", err)
	}

	return response, nil
}
