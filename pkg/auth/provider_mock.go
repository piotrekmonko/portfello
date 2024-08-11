package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"strings"
	"time"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
)

// MockProvider is only for test use.
type MockProvider struct {
	Users []*User
}

var _ Provider = (*MockProvider)(nil)

func NewMockProvider() (*MockProvider, error) {
	return &MockProvider{Users: []*User{
		{
			ID:           "u1",
			DisplayName:  "User One",
			Email:        "user.one@example.com",
			Roles:        []RoleID{RoleSuperAdmin},
			CreatedAt:    time.Now(),
			Registration: nil,
			key:          nil,
			pwdHash:      "123",
		},
		{
			ID:           "u2",
			DisplayName:  "User Two",
			Email:        "user.two@example.com",
			Roles:        []RoleID{RoleAdmin},
			CreatedAt:    time.Now(),
			Registration: nil,
			key:          nil,
			pwdHash:      "123",
		},
		{
			ID:           "u3",
			DisplayName:  "User Three",
			Email:        "user.three@example.com",
			Roles:        []RoleID{RoleUser},
			CreatedAt:    time.Now(),
			Registration: nil,
			key:          nil,
			pwdHash:      "123",
		},
	}}, nil
}

func (m *MockProvider) ProviderName() string {
	return conf.AuthProviderMock
}

func (m *MockProvider) GetUserByID(_ context.Context, userID string) (*User, error) {
	for _, usr := range m.Users {
		if usr.ID == userID {
			return usr, nil
		}
	}

	return nil, ErrUserNotFound
}

func (m *MockProvider) GetUserByEmail(_ context.Context, email string) (*User, error) {
	for i := range m.Users {
		if strings.EqualFold(m.Users[i].GetEmail(), email) {
			return m.Users[i], nil
		}
	}

	return nil, ErrUserNotFound
}

func (m *MockProvider) ListUsers(_ context.Context) ([]*User, int, error) {
	return m.Users, len(m.Users), nil
}

func (m *MockProvider) CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error) {
	existingUser, err := m.GetUserByEmail(ctx, email)
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	} else if existingUser != nil {
		return existingUser, nil
	}

	nu := &User{
		ID:           email,
		DisplayName:  name,
		Email:        email,
		Roles:        roles,
		CreatedAt:    time.Now(),
		Registration: nil,
		key:          nil,
	}
	m.Users = append(m.Users, nu)
	return nu, nil
}

func (m *MockProvider) AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error) {
	usr, err := m.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	usr.Roles = roles
	return usr.Roles, nil
}

// ValidateToken expects token to be in format "mocktoken:<user@email.com>".
func (m *MockProvider) ValidateToken(ctx context.Context, token string) (userID string, err error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 || parts[0] != "mocktoken" {
		return "", fmt.Errorf("invalid token")
	}

	usr, err := m.GetUserByEmail(ctx, parts[1])
	if err != nil {
		return "", err
	}

	return usr.ID, nil
}

func (m *MockProvider) IssueToken(_ context.Context, email string, _ Roles) (token string, err error) {
	return fmt.Sprintf("mocktoken:%s", email), nil
}

func (m *MockProvider) CheckPassword(_ context.Context, usr *User, pass string) error {
	if usr.pwdHash != pass {
		return ErrInvalidPassword
	}

	return nil
}

func (m *MockProvider) SetPassword(_ context.Context, usr *User, pass string) error {
	usr.pwdHash = pass
	return nil
}
