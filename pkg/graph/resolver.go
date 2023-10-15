package graph

import (
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/dao"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DbQueries   *dao.Queries
	AuthService *auth.Service
}
