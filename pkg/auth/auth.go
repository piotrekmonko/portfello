package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type Service struct {
	provider Provider
	cUsers   *cache.Cache[*management.User]
}

type Provider interface {
	GetUserByEmail(ctx context.Context, email string) (*management.User, error)
}

func New(p Provider) *Service {
	gocacheClient := gocache.New(50*time.Minute, 100*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	return &Service{provider: p,
		cUsers: cache.New[*management.User](gocacheStore),
	}
}

func (s *Service) GetUsers(ctx context.Context) {

}

func (s *Service) GetUser(ctx context.Context, userEmail string) (*User, error) {
	auth0User, err := s.cUsers.Get(ctx, userEmail)
	if err != nil && !errors.Is(err, store.NotFound{}) {
		return nil, fmt.Errorf("cannot reach user cache: %w", err)
	}

	if auth0User == nil {
		auth0User, err = s.provider.GetUserByEmail(ctx, userEmail)
		if err != nil {
			return nil, fmt.Errorf("cannot find user in auth0: %w", err)
		}

		err = s.cUsers.Set(ctx, userEmail, auth0User)
		if err != nil {
			return nil, fmt.Errorf("cannot save auth0 user in cache: %w", err)
		}
	}

	return &User{
		ID:          auth0User.GetID(),
		Email:       auth0User.GetEmail(),
		DisplayName: auth0User.GetName(),
		CreatedAt:   auth0User.GetCreatedAt(),
	}, nil
}
