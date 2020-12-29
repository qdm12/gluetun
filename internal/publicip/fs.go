package publicip

import "github.com/qdm12/gluetun/internal/os"

func persistPublicIP(openFile os.OpenFileFunc,
	filepath string, content string, puid, pgid int) error {
	file, err := openFile(
		filepath,
		os.O_TRUNC|os.O_WRONLY|os.O_CREATE,
		0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Chown(puid, pgid); err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
