package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *Service) runUpCommand(ctx context.Context, ports []uint16) (err error) {
	// run command replacing {{PORTS}} with the ports (space separated)
	portStrings := make([]string, len(ports))
	for i, port := range ports {
		portStrings[i] = fmt.Sprint(int(port))
	}
	portsString := strings.Join(portStrings, ",")

	rawCommand := strings.ReplaceAll(s.settings.Command, "{{PORTS}}", portsString)
	s.logger.Info("running port forward command " + rawCommand)
	command := strings.Split(rawCommand, " ")
	// it is a user input and we trust it so we can ignore the gosec warning
	cmd := exec.CommandContext(ctx, command[0], command[1:]...) // #nosec G204
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
