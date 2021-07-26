package openvpn

import (
	"os"
	"strings"
)

func (l *Loop) writeOpenvpnConf(lines []string) error {
	file, err := os.OpenFile(l.targetConfPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	return file.Close()
}
