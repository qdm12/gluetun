package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *Service) runPortForwardedScript(ports []uint16) (err error) {
	// run bash script with ports as arguments
	portStrings := make([]string, len(ports))
	for i, port := range ports {
		portStrings[i] = fmt.Sprint(int(port))
	}
	portsString := strings.Join(portStrings, " ")

	scriptPath := s.settings.Scriptpath
	s.logger.Info("running port forward script " + scriptPath)
	cmd := exec.Command(scriptPath, portsString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running script: %w", err)
	}

	return nil
}
