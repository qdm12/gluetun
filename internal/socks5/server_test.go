package socks5

import (
	"net"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/stretchr/testify/require"
)

type testLogger struct {
	infos []string
}

func (l *testLogger) Debug(string) {}
func (l *testLogger) Error(string) {}
func (l *testLogger) Info(s string) { l.infos = append(l.infos, s) }

func ptrTo[T any](value T) *T { return &value }

func Test_newServer_NoAuth(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	settings := settings.Socks5{
		Enabled:          ptrTo(true),
		ListeningAddress: "127.0.0.1:0",
		User:             ptrTo(""),
		Password:         ptrTo(""),
		Log:              ptrTo(false),
	}

	server, listener, err := newServer(settings, logger)
	require.NoError(t, err)
	defer listener.Close()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()

	conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	require.NoError(t, err)
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	// Client greeting: version 5, 1 method: no-auth (0x00)
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	require.NoError(t, err)

	resp := make([]byte, 2)
	_, err = conn.Read(resp)
	require.NoError(t, err)
	require.Equal(t, byte(0x05), resp[0])
	require.Equal(t, byte(0x00), resp[1])

	_ = listener.Close()
	<-serveErr
}

func Test_newServer_UserPassAuth(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	settings := settings.Socks5{
		Enabled:          ptrTo(true),
		ListeningAddress: "127.0.0.1:0",
		User:             ptrTo("user"),
		Password:         ptrTo("pass"),
		Log:              ptrTo(false),
	}

	server, listener, err := newServer(settings, logger)
	require.NoError(t, err)
	defer listener.Close()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()

	conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	require.NoError(t, err)
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	// Client greeting: version 5, 2 methods: no-auth (0x00), user/pass (0x02)
	_, err = conn.Write([]byte{0x05, 0x02, 0x00, 0x02})
	require.NoError(t, err)

	resp := make([]byte, 2)
	_, err = conn.Read(resp)
	require.NoError(t, err)
	require.Equal(t, byte(0x05), resp[0])
	require.Equal(t, byte(0x02), resp[1])

	// Username/password auth (RFC 1929)
	_, err = conn.Write([]byte{
		0x01, // auth version
		0x04, 'u', 's', 'e', 'r',
		0x04, 'p', 'a', 's', 's',
	})
	require.NoError(t, err)

	authResp := make([]byte, 2)
	_, err = conn.Read(authResp)
	require.NoError(t, err)
	require.Equal(t, byte(0x01), authResp[0])
	require.Equal(t, byte(0x00), authResp[1])

	_ = listener.Close()
	<-serveErr
}
