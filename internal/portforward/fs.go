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
	if err := writePortForwardedToFile(filepath, port); err != nil {
		l.logger.Error(err.Error())
	}
}

func writePortForwardedToFile(filepath string, port uint16) (err error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(fmt.Sprint(port)))
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
