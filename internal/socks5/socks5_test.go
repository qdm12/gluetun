package socks5

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/qdm12/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/proxy"
)

func Test(t *testing.T) {
	server := New(Settings{
		Username: "test",
		Password: "test",
		Address:  ":8000",
		Logger:   log.New(),
	})

	runErr, startErr := server.Start(context.Background())
	require.NoError(t, startErr)

	select {
	case err := <-runErr:
		require.NoError(t, err)
	default:
	}

	t.Log("SlEEPING")
	time.Sleep(15 * time.Second)
	t.Log("Done sleeping")

	err := server.Stop()
	require.NoError(t, err)
}

func backendServer(listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("Test"))
	conn.Close()
	listener.Close()
}

func TestRead(t *testing.T) {
	// backend server which we'll use SOCKS5 to connect to
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	backendServerPort := listener.Addr().(*net.TCPAddr).Port
	go backendServer(listener)

	// SOCKS5 server
	server := New(Settings{
		Address: ":0",
	})
	_, err = server.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = server.Stop()
		assert.NoError(t, err)
	})
	socks5Port := server.listeningAddress().(*net.TCPAddr).Port

	addr := fmt.Sprintf("localhost:%d", socks5Port)
	socksDialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		t.Fatal(err)
	}

	addr = fmt.Sprintf("localhost:%d", backendServerPort)
	conn, err := socksDialer.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 4)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "Test" {
		t.Fatalf("got: %q want: Test", buf)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadPassword(t *testing.T) {
	// backend server which we'll use SOCKS5 to connect to
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	backendServerPort := ln.Addr().(*net.TCPAddr).Port
	go backendServer(ln)

	auth := &proxy.Auth{User: "foo", Password: "bar"}

	server := Server{
		logger:   log.New(),
		username: auth.User,
		password: auth.Password,
		address:  ":0",
	}
	_, err = server.Start(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		err = server.Stop()
		assert.NoError(t, err)
	})

	addr := fmt.Sprintf("localhost:%d", server.listeningAddress().(*net.TCPAddr).Port)

	if d, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct); err != nil {
		t.Fatal(err)
	} else {
		if _, err := d.Dial("tcp", addr); err == nil {
			t.Fatal("expected no-auth dial error")
		}
	}

	badPwd := &proxy.Auth{User: "foo", Password: "not right"}
	if d, err := proxy.SOCKS5("tcp", addr, badPwd, proxy.Direct); err != nil {
		t.Fatal(err)
	} else {
		if _, err := d.Dial("tcp", addr); err == nil {
			t.Fatal("expected bad password dial error")
		}
	}

	badUsr := &proxy.Auth{User: "not right", Password: "bar"}
	if d, err := proxy.SOCKS5("tcp", addr, badUsr, proxy.Direct); err != nil {
		t.Fatal(err)
	} else {
		if _, err := d.Dial("tcp", addr); err == nil {
			t.Fatal("expected bad username dial error")
		}
	}

	socksDialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		t.Fatal(err)
	}

	addr = fmt.Sprintf("localhost:%d", backendServerPort)
	conn, err := socksDialer.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "Test" {
		t.Fatalf("got: %q want: Test", buf)
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}
