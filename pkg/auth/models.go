package auth

import "time"

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
	for i := range r {
		if r[i] == role {
			return true
		}
	}

	return false
}

type User struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Roles       []RoleID  `json:"roles"`
	CreatedAt   time.Time `json:"created_at"`
}
