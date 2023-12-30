package auth

import (
	"context"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
)

type Provider interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, int, error)
	CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error)
	AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error)
	ValidateToken(ctx context.Context, token string) (userID string, err error)
}

// NewProvider builds correct provider based on config.
func NewProvider(ctx context.Context, log logz.Logger, conf *config.Config, dao *dao.DAO) (Provider, error) {
	switch conf.Auth.Provider {
	case config.AuthProviderLocal:
		return NewLocalProvider(log.Named("auth"), dao), nil
	case config.AuthProviderAuth0:
		return NewAuth0Provider(ctx, log.Named("auth"), &conf.Auth)
	default:
		return nil, fmt.Errorf("unsupported auth provider configuration valu: %s", conf.Auth.Provider)
	}
}
