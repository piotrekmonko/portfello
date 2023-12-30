package auth

import (
	"context"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"time"
)

type LocalProvider struct {
	DbDAO *dao.DAO
	log   logz.Logger
}

var _ Provider = (*LocalProvider)(nil)

func NewLocalProvider(log logz.Logger, dao *dao.DAO) *LocalProvider {
	return &LocalProvider{
		DbDAO: dao,
		log:   log.Named("prov.local"),
	}
}

func userFromLocal(u *dao.LocalUser) *User {
	return &User{
		ID:          u.Email,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Roles:       RolesFromString(u.Roles),
		CreatedAt:   u.CreatedAt,
	}
}

func (p *LocalProvider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	usr, err := p.DbDAO.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot find user by email='%s'", email)
	}

	return userFromLocal(usr), nil
}

func (p *LocalProvider) ListUsers(ctx context.Context) ([]*User, int, error) {
	usrList, err := p.DbDAO.LocalUserList(ctx)
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
	tx, rollbacker, err := p.DbDAO.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbacker()

	err = tx.LocalUserInsert(ctx, email, name, roles.ToString(), time.Now().UTC())
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot insert user with email='%s'", email)
	}

	usr, err := tx.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot retrieve user with email='%s'", email)
	}

	return userFromLocal(usr), tx.Commit()
}

func (p *LocalProvider) AssignRoles(ctx context.Context, email string, roles []RoleID) ([]RoleID, error) {
	tx, rollbacker, err := p.DbDAO.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbacker()

	err = tx.LocalUserUpdate(ctx, Roles(roles).ToString(), email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot insert user with email='%s'", email)
	}

	usr, err := tx.LocalUserGetByEmail(ctx, email)
	if err != nil {
		return nil, p.log.Errorw(ctx, err, "cannot retrieve user with email='%s'", email)
	}

	return RolesFromString(usr.Roles), tx.Commit()
}

func (p *LocalProvider) ValidateToken(ctx context.Context, token string) (userID string, err error) {
	//TODO implement me
	panic("implement me")
}
