package service

import (
	"fmt"
	"os"
)

func (s *Service) writePortForwardedFile(port uint16) (err error) {
	filepath := *s.settings.UserSettings.Filepath
	s.logger.Info("writing port file " + filepath)
	const perms = os.FileMode(0644)
	err = os.WriteFile(filepath, []byte(fmt.Sprint(port)), perms)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	err = os.Chown(filepath, s.puid, s.pgid)
	if err != nil {
		return fmt.Errorf("chowning file: %w", err)
	}

	return nil
}
