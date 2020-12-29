package user

import osuser "os/user"

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . OSUser

type OSUser interface {
	LookupID(uid string) (*osuser.User, error)
	Lookup(username string) (*osuser.User, error)
}

func New() OSUser {
	return &osUser{}
}

type osUser struct{}

func (u *osUser) LookupID(uid string) (*osuser.User, error) {
	return osuser.LookupId(uid)
}

func (u *osUser) Lookup(username string) (*osuser.User, error) {
	return osuser.Lookup(username)
}
