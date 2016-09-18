package impl

import (
	ignition "recast.sh/v0/ignition/impl"
)

type User struct {
	Name              string
	PasswordHash      string
	SSHAuthorizedKeys []string
	Uid               uint
	Homedir           string
	NoCreateHome      bool
	PrimaryGroup      string
	Groups            []string
	NoUserGroup       bool
	System            bool
	NoLogInit         bool
	Shell             string
}

func (u *User) ignitionUser() ignition.User {
	return ignition.User{
		Name:              u.Name,
		PasswordHash:      u.PasswordHash,
		SSHAuthorizedKeys: u.SSHAuthorizedKeys,
		Create: &ignition.UserCreate{
			Uid:          &u.Uid,
			Homedir:      u.Homedir,
			NoCreateHome: u.NoCreateHome,
			PrimaryGroup: u.PrimaryGroup,
			Groups:       u.Groups,
			NoUserGroup:  u.NoUserGroup,
			System:       u.System,
			NoLogInit:    u.NoLogInit,
			Shell:        u.Shell,
		},
	}
}
