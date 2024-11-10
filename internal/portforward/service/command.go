package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/qdm12/gluetun/internal/command"
)

func runCommand(ctx context.Context, cmder Cmder, logger Logger,
	commandTemplate string, ports []uint16,
) (err error) {
	portStrings := make([]string, len(ports))
	for i, port := range ports {
		portStrings[i] = fmt.Sprint(int(port))
	}
	portsString := strings.Join(portStrings, ",")
	commandString := strings.ReplaceAll(commandTemplate, "{{PORTS}}", portsString)
	args, err := command.Split(commandString)
	if err != nil {
		return fmt.Errorf("parsing command: %w", err)
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec G204
	stdout, stderr, waitError, err := cmder.Start(cmd)
	if err != nil {
		return err
	}

	streamCtx, streamCancel := context.WithCancel(context.Background())
	streamDone := make(chan struct{})
	go streamLines(streamCtx, streamDone, logger, stdout, stderr)

	err = <-waitError
	streamCancel()
	<-streamDone
	return err
}

func streamLines(ctx context.Context, done chan<- struct{},
	logger Logger, stdout, stderr <-chan string,
) {
	defer close(done)

	var line string

	for {
		select {
		case <-ctx.Done():
			return
		case line = <-stdout:
			logger.Info(line)
		case line = <-stderr:
			logger.Error(line)
		}
	}
}
