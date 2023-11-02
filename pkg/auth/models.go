package auth

import (
	"crypto"
	"github.com/go-acme/lego/v4/registration"
	"time"
)

type RoleID string

func (i RoleID) StrPtr() *string {
	return (*string)(&i)
}

const (
	RoleUser       RoleID = "user"
	RoleAdmin      RoleID = "admin"
	RoleSuperAdmin RoleID = "super"
)

type Roles []RoleID

func (r Roles) Has(role RoleID) bool {
	if r == nil {
		return false
	}

	for i := range r {
		if r[i] == RoleSuperAdmin || r[i] == role {
			return true
		}
	}

	return false
}

type User struct {
	// Basic data kept and fetched from auth provider
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Roles       []RoleID  `json:"roles"`
	CreatedAt   time.Time `json:"created_at"`

	// Extra data to support LE cert, stored in local db
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
