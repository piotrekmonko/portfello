package auth

import (
	"context"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"net/http"
)

type CtxKey int

const CtxUserKey CtxKey = 1

var ErrNotAuthorized = fmt.Errorf("not authorized")

// GetCtxUser returns the auth.User instance from context. Available only in routes wrapped in Service.Middleware.
func GetCtxUser(ctx context.Context) *User {
	ctxUser := ctx.Value(CtxUserKey)
	if ctxUser == nil {
		return nil
	}
	return ctxUser.(*User)
}

func setCtxUser(ctx context.Context, usr *User) context.Context {
	return context.WithValue(ctx, CtxUserKey, usr)
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := jwtmiddleware.AuthHeaderTokenExtractor(r)

		// Allow unauthenticated users in. HasRole handler will return error if access is made to protected resource.
		if err != nil || token == "" {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := s.provider.ValidateToken(ctx, token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusForbidden)
			return
		}

		// Get the user from the auth provider
		user, err := s.GetUserByID(ctx, userID)
		if err != nil {
			http.Error(w, `{"error":"invalid user token"}`, http.StatusForbidden)
			return
		}

		// And put them on context
		ctx = setCtxUser(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
