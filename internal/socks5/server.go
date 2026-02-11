package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	sockslib "github.com/armon/go-socks5"
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func newServer(settings settings.Socks5, logger Logger) (server *sockslib.Server, listener net.Listener, err error) {
	conf := &sockslib.Config{
		Logger: newSocks5Logger(logger, settings),
	}

	if settings.User != nil && *settings.User != "" {
		creds := sockslib.StaticCredentials{
			*settings.User: *settings.Password,
		}
		conf.AuthMethods = []sockslib.Authenticator{
			&sockslib.UserPassAuthenticator{Credentials: creds},
		}
		conf.Credentials = creds
	}

	server, err = sockslib.New(conf)
	if err != nil {
		return nil, nil, fmt.Errorf("creating SOCKS5 server: %w", err)
	}

	listener, err = net.Listen("tcp", settings.ListeningAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("listening on %s: %w", settings.ListeningAddress, err)
	}

	return server, listener, nil
}

func newSocks5Logger(logger Logger, settings settings.Socks5) *log.Logger {
	if settings.Log == nil || !*settings.Log {
		return log.New(io.Discard, "", 0)
	}
	return log.New(&socks5LogWriter{logger: logger}, "", 0)
}

type socks5LogWriter struct {
	logger Logger
}

func (w *socks5LogWriter) Write(p []byte) (n int, err error) {
	message := strings.TrimSpace(string(p))
	if message != "" {
		w.logger.Info(message)
	}
	return len(p), nil
}
