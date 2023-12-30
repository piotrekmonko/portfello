package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type Service struct {
	provider Provider
	cUsers   *cache.Cache[*User]
}

func New(p Provider) *Service {
	gocacheClient := gocache.New(50*time.Minute, 100*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	return &Service{provider: p,
		cUsers: cache.New[*User](gocacheStore),
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]*User, error) {
	users, _, err := s.provider.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot list users: %w", err)
	}

	return users, nil
}

func (s *Service) GetUser(ctx context.Context, userEmail string) (*User, error) {
	user, err := s.cUsers.Get(ctx, userEmail)
	if err != nil && !errors.Is(err, store.NotFound{}) {
		return nil, fmt.Errorf("cannot reach user cache: %w", err)
	}

	if user == nil {
		user, err = s.provider.GetUserByEmail(ctx, userEmail)
		if err != nil {
			return nil, fmt.Errorf("cannot find user in auth0: %w", err)
		}

		err = s.cUsers.Set(ctx, userEmail, user)
		if err != nil {
			return nil, fmt.Errorf("cannot save auth0 user in cache: %w", err)
		}
	}

	return user, nil
}

func (s *Service) GetUserRoles(ctx context.Context, id string) ([]RoleID, error) {
	u, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.Roles, nil
}

func (s *Service) CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error) {
	return s.provider.CreateUser(ctx, email, name, roles)
}

func (s *Service) AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error) {
	user, err := s.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}

	return s.provider.AssignRoles(ctx, user.ID, roles)
}

func (s *Service) HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role RoleID) (res interface{}, err error) {
	user := GetCtxUser(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}

	if Roles(user.Roles).Has(role) {
		return next(ctx)
	}

	return nil, nil
}
