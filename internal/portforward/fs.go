package portforward

import (
	"fmt"
	"os"
)

func (l *Loop) removePortForwardedFile() {
	filepath := *l.state.GetSettings().Filepath
	l.logger.Info("removing port file " + filepath)
	if err := os.Remove(filepath); err != nil {
		l.logger.Error(err.Error())
	}
}

func (l *Loop) writePortForwardedFile(port uint16) {
	filepath := *l.state.GetSettings().Filepath
	l.logger.Info("writing port file " + filepath)
	if err := writePortForwardedToFile(filepath, port, l.puid, l.pgid); err != nil {
		l.logger.Error("writing port forwarded to file: " + err.Error())
	}
}

func writePortForwardedToFile(filepath string, port uint16, uid, gid int) (err error) {
	const perms = os.FileMode(0644)
	err = os.WriteFile(filepath, []byte(fmt.Sprint(port)), perms)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	err = os.Chown(filepath, uid, gid)
	if err != nil {
		return fmt.Errorf("chowning file: %w", err)
	}

	return nil
}
