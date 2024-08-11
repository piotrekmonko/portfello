package auth

import (
	"context"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
)

type Provider interface {
	ProviderName() string
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, int, error)
	CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error)
	AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error)
	ValidateToken(ctx context.Context, token string) (userID string, err error)
	IssueToken(ctx context.Context, email string, scope Roles) (token string, err error)
}

// NewProvider builds correct provider based on config.
func NewProvider(ctx context.Context, log logz.Logger, c *conf.Config, dao *dao.DAO) (Provider, error) {
	switch c.Auth.Provider {
	case conf.AuthProviderLocal:
		return NewLocalProvider(log.Named("auth"), dao, &c.Auth), nil
	case conf.AuthProviderAuth0:
		return NewAuth0Provider(ctx, log.Named("auth"), &c.Auth)
	case conf.AuthProviderMock:
		return NewMockProvider()
	default:
		return nil, fmt.Errorf("unsupported auth provider configuration valu: %s", c.Auth.Provider)
	}
}
