package auth

import (
	"context"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"github.com/piotrekmonko/portfello/pkg/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Auth0Provider struct {
	client       *management.Management
	jwtValidator *validator.Validator
	jwtProvider  *jwks.CachingProvider
	log          *log.Logger
	config       config.Auth0
}

var _ Provider = (*Auth0Provider)(nil)

func NewAuth0Provider(ctx context.Context, c *config.Config) (*Auth0Provider, error) {
	// initialize auth0 management API client
	client, err := management.New(
		c.Auth.Domain,
		management.WithClientCredentials(ctx, c.Auth.ClientID, c.Auth.ClientSecret),
		management.WithClient(&http.Client{Timeout: time.Second * 5}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot configure management api client: %w", err)
	}

	// initialize JWT key validator
	issuerURL, err := url.Parse(c.Auth.Domain + "/")
	if err != nil {
		return nil, fmt.Errorf("cannot parse issuer domain: %w", err)
	}

	jwtProvider := jwks.NewCachingProvider(issuerURL, time.Hour)
	jwtValidator, err := validator.New(
		jwtProvider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{c.Auth.Audience},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create JWT validator: %w", err)
	}

	return &Auth0Provider{
		log:          log.New(os.Stderr, "auth0", log.LstdFlags),
		jwtProvider:  jwtProvider,
		jwtValidator: jwtValidator,
		client:       client,
		config:       c.Auth,
	}, nil
}

func (a *Auth0Provider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	matchingUsers, err := a.client.User.ListByEmail(ctx, email)
	if err != nil || len(matchingUsers) == 0 {
		return nil, fmt.Errorf("cannot list users(%s) by email: %w", email, err)
	}

	auth0User := matchingUsers[0]
	roles, err := a.client.User.Roles(ctx, auth0User.GetID())
	if err != nil {
		return nil, fmt.Errorf("cannot read user(%s) roles: %w", email, err)
	}

	roleSlice := make(Roles, len(roles.Roles))
	for _, role := range roles.Roles {
		roleSlice = append(roleSlice, RoleID(role.GetName()))
	}

	return &User{
		ID:          auth0User.GetID(),
		Email:       auth0User.GetEmail(),
		DisplayName: auth0User.GetName(),
		CreatedAt:   auth0User.GetCreatedAt(),
		Roles:       roleSlice,
	}, nil
}

func (a *Auth0Provider) ListUsers(ctx context.Context) ([]*User, int, error) {
	var wg sync.WaitGroup
	uList, err := a.client.User.List(ctx, management.IncludeTotals(true))
	if err != nil {
		return nil, 0, fmt.Errorf("cannot list users: %w", err)
	}

	wg.Add(len(uList.Users))
	out := make([]*User, len(uList.Users))
	for i, auth0User := range uList.Users {
		out[i] = &User{
			ID:          auth0User.GetID(),
			Email:       auth0User.GetEmail(),
			DisplayName: auth0User.GetName(),
			CreatedAt:   auth0User.GetCreatedAt(),
			Roles:       Roles{},
		}

		go func(userID string, index int) {
			defer wg.Done()

			roles, err := a.client.User.Roles(ctx, userID)
			if err != nil {
				a.log.Printf("cannot read user(%s) roles: %v\n", userID, err)
			}

			roleSlice := make(Roles, len(roles.Roles))
			for _, role := range roles.Roles {
				roleSlice = append(roleSlice, RoleID(role.GetName()))
			}

			out[index].Roles = roleSlice
		}(auth0User.GetID(), i)
	}

	wg.Wait()

	return out, uList.Total, nil
}

func (a *Auth0Provider) CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error) {
	conn, err := a.client.Connection.Read(ctx, a.config.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch Auth0 db connection: %w", err)
	}

	initialPassword := uuid.NewString() + strings.ToUpper(uuid.NewString())
	userReq := &management.User{
		Connection: conn.Name, // note: database connection is referenced through NAME, not by ID.
		Email:      &email,
		Name:       &name,
		Password:   &initialPassword,
	}
	err = a.client.User.Create(ctx, userReq)
	if err != nil {
		if mngmtErr, isMngmtErr := err.(management.Error); isMngmtErr && mngmtErr.Status() == http.StatusConflict {
			a.log.Printf("user(%s) already exists in Auth0", email)
			existingUser, err := a.GetUserByEmail(ctx, email)
			if err != nil {
				return nil, fmt.Errorf("failed at reusing existing Auth0 user(%s): %w", email, err)
			}
			userReq.ID = &existingUser.ID
		} else {
			return nil, fmt.Errorf("cannot create Auth0 user(%s): %w", email, err)
		}
	}

	existingRoles, err := a.client.Role.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot list roles: %w", err)
	}

	auth0Roles := make([]*management.Role, len(roles))
	for i, role := range roles {
		for _, existingRole := range existingRoles.Roles {
			if existingRole.GetName() == string(role) {
				auth0Roles[i] = existingRole
				break
			}
		}

		if auth0Roles[i] == nil {
			return nil, fmt.Errorf("cannot find role %s, please run 'provision systems auth0' first", role)
		}
	}

	if err = a.client.User.AssignRoles(ctx, userReq.GetID(), auth0Roles); err != nil {
		return nil, fmt.Errorf("cannot assign roles to user(%s): %w", email, err)
	}

	return a.GetUserByEmail(ctx, email)
}

func (a *Auth0Provider) AssignRoles(ctx context.Context, auth0UserID string, roles []RoleID) ([]RoleID, error) {
	auth0Roles := make([]*management.Role, len(roles))
	for i, role := range roles {
		auth0Roles[i] = &management.Role{Name: role.StrPtr()}
	}
	if err := a.client.User.AssignRoles(ctx, auth0UserID, auth0Roles); err != nil {
		return nil, fmt.Errorf("cannot assign roles to user(%s): %w", auth0UserID, err)
	}

	return roles, nil
}
