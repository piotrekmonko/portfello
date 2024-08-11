package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LocalProvider struct {
	db   dao.DBInterface
	log  logz.Logger
	conf *conf.Auth0
}

var _ Provider = (*LocalProvider)(nil)

func NewLocalProvider(log logz.Logger, dao dao.DBInterface, conf *conf.Auth0) *LocalProvider {
	return &LocalProvider{
		db:   dao,
		log:  log.Named("prov.local"),
		conf: conf,
	}
}

func userFromLocal(u *dao.LocalUser) *User {
	return &User{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Roles:       RolesFromString(u.Roles),
		CreatedAt:   u.CreatedAt,
		pwdHash:     u.Pwdhash,
	}
}

func (p *LocalProvider) ProviderName() string {
	return conf.AuthProviderLocal
}

func (p *LocalProvider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	usr, err := p.db.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot find user by email", "email", email)
	}

	return userFromLocal(usr), nil
}

func (p *LocalProvider) GetUserByID(ctx context.Context, userID string) (*User, error) {
	usr, err := p.db.LocalUserGetByID(ctx, userID)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot find user", "userID", userID)
	}

	return userFromLocal(usr), nil
}

func (p *LocalProvider) ListUsers(ctx context.Context) ([]*User, int, error) {
	usrList, err := p.db.LocalUserList(ctx)
	if err != nil {
		return nil, -1, p.log.Errorw(ctx, err, "cannot list users")
	}

	out := make([]*User, len(usrList))
	for i := range usrList {
		out[i] = userFromLocal(usrList[i])
	}

	return out, len(out), nil
}

func (p *LocalProvider) CreateUser(ctx context.Context, email string, name string, roles Roles) (*User, error) {
	tx, rollbacker, err := p.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbacker()

	err = tx.LocalUserInsert(ctx, &dao.LocalUserInsertParams{
		ID:          uuid.NewString(),
		Email:       email,
		DisplayName: name,
		Roles:       roles.ToString(),
		CreatedAt:   time.Now().UTC(),
		Pwdhash:     "", // set initial pass to empty prevents login
	})
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot insert user with email", "email", email)
	}

	usr, err := tx.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot retrieve user with email", "email", email)
	}

	return userFromLocal(usr), tx.Commit(ctx)
}

func (p *LocalProvider) AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error) {
	tx, rollbacker, err := p.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbacker()

	err = tx.LocalUserUpdate(ctx, Roles(roles).ToString(), email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot insert user with email", "email", email)
	}

	usr, err := tx.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot retrieve user with email", "email", email)
	}

	return RolesFromString(usr.Roles), tx.Commit(ctx)
}

// CheckPassword compares pass to pwdhash stored in db. Used only in LocalProvider.
func (p *LocalProvider) CheckPassword(_ context.Context, usr *User, pass string) error {
	if usr.pwdHash == "" {
		return fmt.Errorf("user has not set their password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(usr.pwdHash), []byte(pass))
	if err != nil {
		return ErrInvalidPassword
	}

	return nil
}

// SetPassword sets a users stored in db. Used only in LocalProvider. Use empty pass to prevent login.
func (p *LocalProvider) SetPassword(ctx context.Context, usr *User, pass string) error {
	if pass == "" {
		err := p.db.LocalUserSetPass(ctx, "", usr.GetEmail())
		if err != nil {
			return p.log.Errorw(ctx, err, "cannot set empty password")
		}
	}

	var newPass []byte
	newPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return p.log.Errorw(ctx, err, "cannot use this password")
	}

	err = p.db.LocalUserSetPass(ctx, string(newPass), usr.GetEmail())
	if err != nil {
		return p.log.Errorw(ctx, err, "cannot update password")
	}

	return nil
}

func (p *LocalProvider) ValidateToken(ctx context.Context, tokenString string) (userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.conf.ClientSecret), nil
	})
	if err != nil {
		return "", p.log.Errorw(ctx, err, "cannot parse token")
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")
}

func (p *LocalProvider) IssueToken(ctx context.Context, email string, scope Roles) (string, error) {
	// Create the claims
	now := time.Now().UTC()
	claims := JwtClaims{
		Scope: scope.ToString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now.UTC()),
			NotBefore: jwt.NewNumericDate(now.UTC()),
			Issuer:    p.conf.Provider,
			Subject:   email,
			Audience:  []string{email},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(p.conf.ClientSecret))
	return signedToken, p.log.Errorw(ctx, err, "cannot sign token")
}
