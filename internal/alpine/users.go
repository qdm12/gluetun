package alpine

import (
	"fmt"
	"os"
	"os/user"
)

// CreateUser creates a user in Alpine with the given UID.
func (c *configurator) CreateUser(username string, uid int) (createdUsername string, err error) {
	UIDStr := fmt.Sprintf("%d", uid)
	u, err := c.lookupUID(UIDStr)
	_, unknownUID := err.(user.UnknownUserIdError)
	if err != nil && !unknownUID {
		return "", fmt.Errorf("cannot create user: %w", err)
	} else if u != nil {
		if u.Username == username {
			return "", nil
		}
		return u.Username, nil
	}
	u, err = c.lookupUser(username)
	_, unknownUsername := err.(user.UnknownUserError)
	if err != nil && !unknownUsername {
		return "", fmt.Errorf("cannot create user: %w", err)
	} else if u != nil {
		return "", fmt.Errorf("cannot create user: user with name %s already exists for ID %s instead of %d",
			username, u.Uid, uid)
	}
	file, err := c.openFile("/etc/passwd", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}
	s := fmt.Sprintf("%s:x:%d:::/dev/null:/sbin/nologin\n", username, uid)
	_, err = file.WriteString(s)
	if err != nil {
		_ = file.Close()
		return "", err
	}
	return username, file.Close()
}
