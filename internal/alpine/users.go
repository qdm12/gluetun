package alpine

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

// CreateUser creates a user in Alpine with the given UID.
func (a *Alpine) CreateUser(username string, uid int) (createdUsername string, err error) {
	UIDStr := strconv.Itoa(uid)
	u, err := a.lookupID(UIDStr)
	_, unknownUID := err.(user.UnknownUserIdError)
	if err != nil && !unknownUID {
		return "", err
	}

	if u != nil {
		if u.Username == username {
			return "", nil
		}
		return u.Username, nil
	}

	u, err = a.lookup(username)
	_, unknownUsername := err.(user.UnknownUserError)
	if err != nil && !unknownUsername {
		return "", err
	}

	if u != nil {
		return "", fmt.Errorf("%w: with name %s for ID %s instead of %d",
			ErrUserAlreadyExists, username, u.Uid, uid)
	}

	file, err := os.OpenFile(a.passwdPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s:x:%d:::/dev/null:/sbin/nologin\n", username, uid)
	_, err = file.WriteString(s)
	if err != nil {
		_ = file.Close()
		return "", err
	}

	return username, file.Close()
}
