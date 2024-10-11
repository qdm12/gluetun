package publicip

import (
	"io/fs"
	"os"
)

func persistPublicIP(path string, content string, puid, pgid int) error {
	const permission = fs.FileMode(0o644)
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, permission)
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
