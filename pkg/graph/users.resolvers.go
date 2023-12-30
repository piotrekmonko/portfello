package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"fmt"

	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/graph/model"
)

// UserCreate is the resolver for the userCreate field.
func (r *mutationResolver) UserCreate(ctx context.Context, newUser model.NewUser) (*auth.User, error) {
	return r.AuthService.CreateUser(ctx, newUser.Email, newUser.DisplayName, auth.Roles{auth.RoleUser})
}

// AdminCreate is the resolver for the adminCreate field.
func (r *mutationResolver) AdminCreate(ctx context.Context, newAdmin model.NewUser) (*auth.User, error) {
	return r.AuthService.CreateUser(ctx, newAdmin.Email, newAdmin.DisplayName, auth.Roles{auth.RoleSuperAdmin})
}

// UserAssignRoles is the resolver for the userAssignRoles field.
func (r *mutationResolver) UserAssignRoles(ctx context.Context, email string, newRoles []auth.RoleID) ([]auth.RoleID, error) {
	user := auth.GetCtxUser(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	userRoles := auth.Roles(user.Roles)
	requestedRoles := auth.Roles(newRoles)

	switch {
	case requestedRoles.Has(auth.RoleAdmin) && !userRoles.Has(auth.RoleSuperAdmin):
		return nil, fmt.Errorf("only super administrators may grant admin role")
	case requestedRoles.Has(auth.RoleSuperAdmin) && !userRoles.Has(auth.RoleSuperAdmin):
		return nil, fmt.Errorf("only super administrators may grant super admin role")
	}

	return r.AuthService.AssignRoles(ctx, email, newRoles)
}

// GetUserRoles is the resolver for the getUserRoles field.
func (r *queryResolver) GetUserRoles(ctx context.Context, userID string) ([]auth.RoleID, error) {
	return r.AuthService.GetUserRoles(ctx, userID)
}

// ListUsers is the resolver for the listUsers field.
func (r *queryResolver) ListUsers(ctx context.Context) ([]*auth.User, error) {
	return r.AuthService.GetUsers(ctx)
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, email string) (*auth.User, error) {
	return r.AuthService.GetUser(ctx, email)
}
