package openvpn

import (
	"context"
	"strings"
)

func streamLines(ctx context.Context, done chan<- struct{},
	logger Logger, stdout, stderr chan string,
	tunnelReady chan<- struct{}) {
	defer close(done)

	var line string

	for {
		errLine := false
		select {
		case <-ctx.Done():
			// Context should only be canceled after stdout and stderr are done
			// being written to.
			close(stdout)
			close(stderr)
			return
		case line = <-stdout:
		case line = <-stderr:
			errLine = true
		}
		line, level := processLogLine(line)
		if line == "" {
			continue // filtered out
		}
		if errLine {
			level = levelError
		}
		switch level {
		case levelInfo:
			logger.Info(line)
		case levelWarn:
			logger.Warn(line)
		case levelError:
			logger.Error(line)
		}
		if strings.Contains(line, "Initialization Sequence Completed") {
			// do not close tunnelReady in case the initialization
			// happens multiple times without Openvpn restarting
			tunnelReady <- struct{}{}
		}
	}
}
