package openvpn

import (
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/logging"
)

func (l *looper) collectLines(wg *sync.WaitGroup, stdout, stderr <-chan string) {
	defer wg.Done()
	var line string
	var ok, errLine bool

	for {
		errLine = false
		select {
		case line, ok = <-stdout:
		case line, ok = <-stderr:
			errLine = true
		}
		if !ok {
			return
		}
		line, level := processLogLine(line)
		if len(line) == 0 {
			continue // filtered out
		}
		if errLine {
			level = logging.LevelError
		}
		switch level {
		case logging.LevelDebug:
			l.logger.Debug(line)
		case logging.LevelInfo:
			l.logger.Info(line)
		case logging.LevelWarn:
			l.logger.Warn(line)
		case logging.LevelError:
			l.logger.Error(line)
		}
		if strings.Contains(line, "Initialization Sequence Completed") {
			l.tunnelReady <- struct{}{}
		}
	}
}

func processLogLine(s string) (filtered string, level logging.Level) {
	for _, ignored := range []string{
		"WARNING: you are using user/group/chroot/setcon without persist-tun -- this may cause restarts to fail",
		"NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
	} {
		if s == ignored {
			return "", logging.LevelDebug
		}
	}
	switch {
	case strings.HasPrefix(s, "NOTE: "):
		filtered = strings.TrimPrefix(s, "NOTE: ")
		level = logging.LevelInfo
	case strings.HasPrefix(s, "WARNING: "):
		filtered = strings.TrimPrefix(s, "WARNING: ")
		level = logging.LevelWarn
	case strings.HasPrefix(s, "Options error: "):
		filtered = strings.TrimPrefix(s, "Options error: ")
		level = logging.LevelError
	case s == "Initialization Sequence Completed":
		return color.HiGreenString(s), logging.LevelInfo
	case s == "AUTH: Received control message: AUTH_FAILED":
		filtered = s + `

Your credentials might be wrong 🤨

💡 If you use Private Internet Access, check https://github.com/qdm12/gluetun/issues/265

`
		level = logging.LevelError
	case strings.Contains(s, "TLS Error: TLS key negotiation failed to occur within 60 seconds (check your network connectivity)"): //nolint:lll
		filtered = s + `
🚒🚒🚒🚒🚒🚨🚨🚨🚨🚨🚨🚒🚒🚒🚒🚒
That error usually happens because either:

1. The VPN server IP address you are trying to connect to is no longer valid 🔌
   Update your server information using https://github.com/qdm12/gluetun/wiki/Updating-Servers

2. The VPN server crashed 💥, try changing your VPN servers filtering options such as REGION

3. Your Internet connection is not working 🤯, ensure it works

4. Something else ➡️ https://github.com/qdm12/gluetun/issues/new/choose
`
	default:
		filtered = s
		level = logging.LevelInfo
	}
	filtered = constants.ColorOpenvpn().Sprintf(filtered)
	return filtered, level
}
