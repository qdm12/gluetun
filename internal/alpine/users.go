package alpine

import (
	"fmt"
	"os/user"
)

// CreateUser creates a user in Alpine with the given UID
func (c *configurator) CreateUser(username string, uid int) error {
	UIDStr := fmt.Sprintf("%d", uid)
	u, err := c.lookupUID(UIDStr)
	_, unknownUID := err.(user.UnknownUserIdError)
	if err != nil && !unknownUID {
		return fmt.Errorf("cannot create user: %w", err)
	} else if u != nil {
		if u.Username == username {
			return nil
		}
		return fmt.Errorf("user with ID %d exists with username %q instead of %q", uid, u.Username, username)
	}
	u, err = c.lookupUser(username)
	_, unknownUsername := err.(user.UnknownUserError)
	if err != nil && !unknownUsername {
		return fmt.Errorf("cannot create user: %w", err)
	} else if u != nil {
		return fmt.Errorf("cannot create user: user with name %s already exists for ID %s instead of %d", username, u.Uid, uid)
	}
	passwd, err := c.fileManager.ReadFile("/etc/passwd")
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	passwd = append(passwd, []byte(fmt.Sprintf("%s:x:%d:::/dev/null:/sbin/nologin\n", username, uid))...)

	if err := c.fileManager.WriteToFile("/etc/passwd", passwd); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	return nil
}
