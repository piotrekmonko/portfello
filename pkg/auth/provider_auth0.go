package auth

import (
	"context"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Auth0Claims struct {
	Scope string `json:"scope"`
}

func (c Auth0Claims) Validate(_ context.Context) error {
	return nil
}

type Auth0Provider struct {
	log          logz.Logger
	manager      *management.Management
	jwtValidator *validator.Validator
	jwtProvider  *jwks.CachingProvider
	config       *conf.Auth0
}

var _ Provider = (*Auth0Provider)(nil)

func NewAuth0Provider(ctx context.Context, log logz.Logger, conf *conf.Auth0) (*Auth0Provider, error) {
	// Initialize auth0 management API client
	client, err := management.New(
		conf.Domain,
		management.WithClientCredentials(ctx, conf.ClientID, conf.ClientSecret),
		management.WithClient(&http.Client{Timeout: time.Second * 5}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot configure management api client: %w", err)
	}

	// Initialize JWT key validator
	issuerURL, err := url.Parse(conf.Domain + "/")
	if err != nil {
		return nil, fmt.Errorf("cannot parse issuer domain: %w", err)
	}

	jwtProvider := jwks.NewCachingProvider(issuerURL, time.Hour)
	jwtValidator, err := validator.New(
		jwtProvider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{conf.Audience},
		validator.WithAllowedClockSkew(time.Minute),
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &Auth0Claims{}
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create JWT validator: %w", err)
	}

	return &Auth0Provider{
		log:          log.Named("prov.auth0"),
		jwtProvider:  jwtProvider,
		jwtValidator: jwtValidator,
		manager:      client,
		config:       conf,
	}, nil
}

func (a *Auth0Provider) ValidateToken(ctx context.Context, token string) (string, error) {
	validatedToken, err := a.jwtValidator.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}

	userID := validatedToken.(validator.ValidatedClaims).RegisteredClaims.Subject
	// TODO: Validate custom claims, needs configuration on Auth0 end.
	// claims := validatedToken.(validator.ValidatedClaims).CustomClaims

	return userID, nil
}

func (a *Auth0Provider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	matchingUsers, err := a.manager.User.ListByEmail(ctx, email)
	if err != nil || len(matchingUsers) == 0 {
		return nil, a.log.Errorw(ctx, err, "cannot find user by email='%s'", email)
	}

	auth0User := matchingUsers[0]
	roles, err := a.manager.User.Roles(ctx, auth0User.GetID())
	if err != nil {
		return nil, a.log.Errorw(ctx, err, "cannot read user roles for email='%s'", email)
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
	uList, err := a.manager.User.List(ctx, management.IncludeTotals(true))
	if err != nil {
		return nil, -1, a.log.Errorw(ctx, err, "cannot list users")
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

			roles, err := a.manager.User.Roles(ctx, userID)
			if err != nil {
				_ = a.log.Errorw(ctx, err, "cannot read user roles for userID='%s'", userID)
				return
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
	conn, err := a.manager.Connection.Read(ctx, a.config.ConnectionID)
	if err != nil {
		return nil, a.log.Errorw(ctx, err, "cannot fetch Auth0 db connection")
	}

	initialPassword := uuid.NewString() + strings.ToUpper(uuid.NewString())
	userReq := &management.User{
		Connection: conn.Name, // note: database connection is referenced through NAME, not by ID.
		Email:      &email,
		Name:       &name,
		Password:   &initialPassword,
	}
	err = a.manager.User.Create(ctx, userReq)
	if err != nil {
		if mngmtErr, isMngmtErr := err.(management.Error); isMngmtErr && mngmtErr.Status() == http.StatusConflict {
			a.log.Warnw(ctx, "user(%s) already exists in Auth0", email)
			existingUser, err := a.GetUserByEmail(ctx, email)
			if err != nil {
				return nil, a.log.Errorw(ctx, err, "failed at reusing existing Auth0 user(%s)", email)
			}
			userReq.ID = &existingUser.ID
		} else {
			return nil, a.log.Errorw(ctx, err, "cannot create Auth0 user(%s)", email)
		}
	}

	existingRoles, err := a.manager.Role.List(ctx)
	if err != nil {
		return nil, a.log.Errorw(ctx, err, "cannot list roles")
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
			return nil, a.log.Errorw(ctx, fmt.Errorf("please run 'provision systems auth0' first"), "cannot find role %s", role)
		}
	}

	if err = a.manager.User.AssignRoles(ctx, userReq.GetID(), auth0Roles); err != nil {
		return nil, a.log.Errorw(ctx, err, "cannot assign roles to user(%s)", email)
	}

	return a.GetUserByEmail(ctx, email)
}

func (a *Auth0Provider) AssignRoles(ctx context.Context, auth0UserID string, roles []RoleID) ([]RoleID, error) {
	auth0Roles := make([]*management.Role, len(roles))
	for i, role := range roles {
		auth0Roles[i] = &management.Role{Name: role.StrPtr()}
	}
	if err := a.manager.User.AssignRoles(ctx, auth0UserID, auth0Roles); err != nil {
		return nil, a.log.Errorw(ctx, err, "cannot assign roles to user(%s)", auth0UserID)
	}

	return roles, nil
}
