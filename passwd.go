package google

import (
	"recast.sh/v0/core"
	google "recast.sh/v0/provider/google/impl"
)

var passwdScope *googleVM

func (vm *googleVM) Passwd(fn func()) core.BaseVM {
	if passwdScope != nil {
		panic(core.ErrUnitsNested)
	}
	passwdScope = vm
	defer func() {
		passwdScope = nil
	}()
	fn()
	return vm
}

func User(name string) core.BaseUser {
	return &googleUser{
		User: passwdScope.AddUser(google.User{
			Name: name,
		}),
	}
}

type googleUser struct {
	*google.User
}

func (u *googleUser) GetName() string {
	return u.User.Name
}

func (u *googleUser) PasswordHash(hash string) core.BaseUser {
	u.User.PasswordHash = hash
	return u
}
func (u *googleUser) SSHAuthorizedKeys(keys ...string) core.BaseUser {
	u.User.SSHAuthorizedKeys = append(u.User.SSHAuthorizedKeys, keys...)
	return u
}

func (u *googleUser) Uid(id uint) core.BaseUser {
	u.User.Uid = id
	return u
}
func (u *googleUser) Homedir(home string) core.BaseUser {
	u.User.Homedir = home
	return u
}
func (u *googleUser) NoCreateHome(b bool) core.BaseUser {
	u.User.NoCreateHome = b
	return u
}
func (u *googleUser) PrimaryGroup(g string) core.BaseUser {
	u.User.PrimaryGroup = g
	return u
}
func (u *googleUser) Groups(g ...string) core.BaseUser {
	u.User.Groups = append(u.User.Groups, g...)
	return u
}
func (u *googleUser) NoUserGroup(b bool) core.BaseUser {
	u.User.NoUserGroup = b
	return u
}
func (u *googleUser) System(b bool) core.BaseUser {
	u.User.System = b
	return u
}
func (u *googleUser) NoLogInit(b bool) core.BaseUser {
	u.User.NoLogInit = b
	return u
}
func (u *googleUser) Shell(shell string) core.BaseUser {
	u.User.Shell = shell
	return u
}
