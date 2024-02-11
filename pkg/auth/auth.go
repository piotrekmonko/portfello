package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocacheStore "github.com/eko/gocache/store/go_cache/v4"
	"github.com/golang-jwt/jwt/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"time"
)

type Service struct {
	provider Provider
	cUsers   *cache.Cache[*User]
}

func New(p Provider) *Service {
	gocacheClient := gocache.New(50*time.Minute, 100*time.Minute)
	cacheStore := gocacheStore.NewGoCache(gocacheClient)
	return &Service{provider: p,
		cUsers: cache.New[*User](cacheStore),
	}
}

func NewFromConfig(ctx context.Context, log *logz.Log, c *conf.Config, dbQuerier *dao.DAO) (*Service, error) {
	authProvider, err := NewProvider(ctx, log, c, dbQuerier)
	if err != nil {
		return nil, err
	}

	return New(authProvider), nil
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

func (s *Service) HasRole(ctx context.Context, _ interface{}, next graphql.Resolver, role RoleID) (res interface{}, err error) {
	user := GetCtxUser(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}

	if Roles(user.Roles).Has(role) {
		return next(ctx)
	}

	return nil, nil
}

type passChecker interface {
	CheckPassword(ctx context.Context, usr *User, pass string) error
	SetPassword(ctx context.Context, usr *User, pass string) error
}

// CheckPassword compares pass to pwdhash stored in db. Used only with LocalProvider.
func (s *Service) CheckPassword(ctx context.Context, usr *User, pass string) error {
	passCheckerService, isPassChecker := s.provider.(passChecker)
	if !isPassChecker {
		return fmt.Errorf("password login not available with this backend")
	}

	return passCheckerService.CheckPassword(ctx, usr, pass)
}

func (s *Service) SetPassword(ctx context.Context, usr *User, pass string) error {
	passCheckerService, isPassChecker := s.provider.(passChecker)
	if !isPassChecker {
		return fmt.Errorf("password change not available with this backend")
	}

	return passCheckerService.SetPassword(ctx, usr, pass)
}

func (s *Service) IssueToken(ctx context.Context, usr *User) (string, error) {
	return s.provider.IssueToken(ctx, usr.GetEmail(), usr.Roles)
}

type JwtClaims struct {
	jwt.RegisteredClaims
	// Scope holds the issuers roles. Should not be empty.
	Scope string `json:"scope"`
}

func (c JwtClaims) Validate(_ context.Context) error {
	return nil
}
