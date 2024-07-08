package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/piotrekmonko/portfello/mocks/github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	knownPass = "correct password"
	knownHash = "$2a$10$1n6Vu30YHSSQ.N5LyIbwY.jLwRLNHgjJG8YN5ZXjKbqNnkQ7cfOXO"
)

var sentinelError = fmt.Errorf("error from database")

func newLocalProvider(t *testing.T) (*LocalProvider, *logz.TestLogger, *conf.Config) {
	testLogger := logz.NewTestLogger(t)
	testConf := conf.NewTestConfig()
	prov := NewLocalProvider(testLogger, nil, &testConf.Auth)
	return prov, testLogger, testConf
}

func TestLocalProvider_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	prov, _, _ := newLocalProvider(t)
	var mockUser dao.LocalUser
	require.Nil(t, gofakeit.Struct(&mockUser))

	testDao := mock_dao.NewMockDBInterface(t)
	testDao.EXPECT().LocalUserGetByEmail(ctx, mockUser.DisplayName).Return(nil, sentinelError).Once()
	prov.db = testDao
	got, err := prov.GetUserByEmail(ctx, mockUser.DisplayName)
	require.True(t, errors.Is(err, sentinelError))

	testDao2 := mock_dao.NewMockDBInterface(t)
	testDao2.EXPECT().LocalUserGetByEmail(ctx, mockUser.Email).Return(&mockUser, nil).Once()
	prov.db = testDao2
	got, err = prov.GetUserByEmail(ctx, mockUser.Email)
	require.Nil(t, err)
	assert.Equal(t, mockUser.Email, got.Email)
	assert.Equal(t, mockUser.DisplayName, got.DisplayName)
	assert.Equal(t, []RoleID{RoleUser}, got.Roles)
}

func TestLocalProvider_ListUsers(t *testing.T) {
	ctx := context.Background()
	prov, _, _ := newLocalProvider(t)
	var mockUsers []*dao.LocalUser
	gofakeit.Slice(&mockUsers)
	require.True(t, len(mockUsers) > 0)

	testDao := mock_dao.NewMockDBInterface(t)
	testDao.EXPECT().LocalUserList(ctx).Return(nil, sentinelError).Once()
	prov.db = testDao
	got, count, err := prov.ListUsers(ctx)
	require.Equal(t, -1, count)
	require.True(t, errors.Is(err, sentinelError))

	testDao2 := mock_dao.NewMockDBInterface(t)
	testDao2.EXPECT().LocalUserList(ctx).Return(nil, nil).Once()
	prov.db = testDao2
	got, count, err = prov.ListUsers(ctx)
	require.Equal(t, 0, count)
	require.Nil(t, err)

	testDao3 := mock_dao.NewMockDBInterface(t)
	testDao3.EXPECT().LocalUserList(ctx).Return(mockUsers, nil).Once()
	prov.db = testDao3
	got, count, err = prov.ListUsers(ctx)
	require.Nil(t, err)
	require.Equal(t, len(mockUsers), count)
	for i, mockUser := range mockUsers {
		assert.Equal(t, mockUser.Email, got[i].Email)
		assert.Equal(t, mockUser.DisplayName, got[i].DisplayName)
		assert.Equal(t, []RoleID{RoleUser}, got[i].Roles)
	}
}

func TestLocalProvider_CreateUser(t *testing.T) {
	ctx := context.Background()
	prov, _, _ := newLocalProvider(t)
	var mockUser dao.LocalUser
	require.Nil(t, gofakeit.Struct(&mockUser))

	testDao := mock_dao.NewMockDBInterface(t)
	testDao.EXPECT().BeginTx(ctx).Return(testDao, func() {}, nil).Once()
	testDao.EXPECT().LocalUserInsert(ctx, mock.Anything).Return(nil).Once()
	testDao.EXPECT().LocalUserGetByEmail(ctx, mockUser.Email).Return(nil, sentinelError).Once()
	prov.db = testDao
	got, err := prov.CreateUser(ctx, mockUser.Email, mockUser.DisplayName, RolesFromString(mockUser.Roles))
	require.True(t, errors.Is(err, sentinelError))
	require.Nil(t, got)

	testDao2 := mock_dao.NewMockDBInterface(t)
	testDao2.EXPECT().BeginTx(ctx).Return(testDao2, func() {}, nil).Once()
	testDao2.EXPECT().LocalUserInsert(ctx, mock.Anything).Return(nil).Once()
	testDao2.EXPECT().LocalUserGetByEmail(ctx, mockUser.Email).Return(&mockUser, nil).Once()
	testDao2.EXPECT().Commit(ctx).Return(nil).Once()
	prov.db = testDao2
	got, err = prov.CreateUser(ctx, mockUser.Email, mockUser.DisplayName, RolesFromString(mockUser.Roles))
	require.Nil(t, err)
	assert.Equal(t, mockUser.Email, got.Email)
	assert.Equal(t, mockUser.DisplayName, got.DisplayName)
	assert.Equal(t, []RoleID{RoleUser}, got.Roles)
}

func TestLocalProvider_AssignRoles(t *testing.T) {
	ctx := context.Background()
	prov, _, _ := newLocalProvider(t)
	var mockUser dao.LocalUser
	require.Nil(t, gofakeit.Struct(&mockUser))

	testDao := mock_dao.NewMockDBInterface(t)
	testDao.EXPECT().BeginTx(ctx).Return(testDao, func() {}, nil).Once()
	testDao.EXPECT().LocalUserUpdate(ctx, mockUser.Roles, mockUser.Email).Return(nil).Once()
	testDao.EXPECT().LocalUserGetByEmail(ctx, mockUser.Email).Return(nil, sentinelError).Once()
	prov.db = testDao
	got, err := prov.AssignRoles(ctx, mockUser.Email, RolesFromString(mockUser.Roles))
	require.True(t, errors.Is(err, sentinelError))
	require.Nil(t, got)

	testDao2 := mock_dao.NewMockDBInterface(t)
	testDao2.EXPECT().BeginTx(ctx).Return(testDao2, func() {}, nil).Once()
	testDao2.EXPECT().LocalUserUpdate(ctx, mockUser.Roles, mockUser.Email).Return(nil).Once()
	testDao2.EXPECT().LocalUserGetByEmail(ctx, mockUser.Email).Return(&mockUser, nil).Once()
	testDao2.EXPECT().Commit(ctx).Return(nil).Once()
	prov.db = testDao2
	got, err = prov.AssignRoles(ctx, mockUser.Email, RolesFromString(mockUser.Roles))
	require.Nil(t, err)
	assert.Equal(t, mockUser.Roles, Roles(got).ToString())
}

func TestLocalProvider_CheckPassword(t *testing.T) {
	ctx := context.Background()
	prov, _, _ := newLocalProvider(t)

	tests := []struct {
		hash string
		pass string
		err  string
	}{
		{
			hash: "",
			err:  "user has not set their password",
		},
		{
			hash: knownHash,
			pass: "abc",
			err:  "invalid password",
		},
		{
			hash: knownHash,
			pass: "",
			err:  "invalid password",
		},
		{
			hash: knownHash,
			pass: knownPass,
			err:  "",
		},
	}

	for _, tt := range tests {
		mockUser := &User{pwdHash: tt.hash}
		err := prov.CheckPassword(ctx, mockUser, tt.pass)
		if tt.err != "" {
			assert.EqualError(t, err, tt.err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestLocalProvider_SetPassword(t *testing.T) {
	ctx := context.Background()
	mockUser := &User{Email: gofakeit.Email()}

	tests := []struct {
		pass string
		hash string
		err  error
	}{
		{
			pass: "",
			err:  sentinelError,
		},
		{
			pass: "some pass",
			err:  sentinelError,
		},
		{
			pass: knownPass,
			err:  nil,
		},
	}

	for _, tt := range tests {
		testDao := mock_dao.NewMockDBInterface(t)
		testDao.EXPECT().LocalUserSetPass(ctx, mock.Anything, mockUser.Email).Return(tt.err).Once()
		prov, _, _ := newLocalProvider(t)
		prov.db = testDao
		err := prov.SetPassword(ctx, mockUser, tt.pass)
		if tt.err != nil {
			require.True(t, errors.Is(err, tt.err))
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestLocalProvider_IssueToken_ValidateToken(t *testing.T) {
	ctx := context.Background()
	email := gofakeit.Email()
	prov, _, _ := newLocalProvider(t)

	token, err := prov.IssueToken(ctx, email, Roles{RoleAdmin})
	require.Nil(t, err)

	userEmail, err := prov.ValidateToken(ctx, token)
	require.Nil(t, err)
	assert.Equal(t, email, userEmail)

	_, err = prov.ValidateToken(ctx, "invalid token")
	assert.EqualError(t, err, "cannot parse token: token contains an invalid number of segments")
}
