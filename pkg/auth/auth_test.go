package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRolesFromString(t *testing.T) {
	tests := []struct {
		in  string
		out Roles
	}{
		{
			in:  "",
			out: Roles{},
		},
		{
			in:  RoleUser.String(),
			out: Roles{RoleUser},
		},
		{
			in:  RoleUser.String() + ";" + RoleAdmin.String(),
			out: Roles{RoleUser, RoleAdmin},
		},
		{
			in:  RoleSuperAdmin.String() + ";invalid-role;" + RoleAdmin.String(),
			out: Roles{RoleSuperAdmin, RoleAdmin},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := RolesFromString(tt.in)
			assert.Equal(t, tt.out, out)
		})
	}
}
