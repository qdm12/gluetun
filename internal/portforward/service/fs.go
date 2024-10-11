package service

import (
	"fmt"
	"os"
	"strings"
)

func (s *Service) writePortForwardedFile(ports []uint16) (err error) {
	portStrings := make([]string, len(ports))
	for i, port := range ports {
		portStrings[i] = fmt.Sprint(int(port))
	}
	fileData := []byte(strings.Join(portStrings, "\n"))

	filepath := s.settings.Filepath
	s.logger.Info("writing port file " + filepath)
	const perms = os.FileMode(0o644)
	err = os.WriteFile(filepath, fileData, perms)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	err = os.Chown(filepath, s.puid, s.pgid)
	if err != nil {
		return fmt.Errorf("chowning file: %w", err)
	}

	return nil
}
