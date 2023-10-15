package auth

import (
	"context"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/piotrekmonko/portfello/pkg/config"
	"net/http"
	"net/url"
	"time"
)

type Auth0Provider struct {
	client       *management.Management
	jwtValidator *validator.Validator
	jwtProvider  *jwks.CachingProvider
}

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
		jwtProvider:  jwtProvider,
		jwtValidator: jwtValidator,
		client:       client,
	}, nil
}

func (a *Auth0Provider) GetUserByEmail(ctx context.Context, email string) (*management.User, error) {
	matchingUsers, err := a.client.User.ListByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("cannot list users by email: %w", err)
	}

	return matchingUsers[0], nil
}
