package service_test

import (
	"context"
	"testing"
	"time"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/pkg/utils/crypt"
	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		authService  = injector.CommonAuthService
		userDaoMock  = injector.UserDaoMock
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		password     = "User@123"
		role         = "USER"
		organization = "ORG"
	)
	passwordHash, err := crypt.Hash(password)
	assert.NoError(t, err)
	userID, err := userDaoMock.UserDao.InsertUser(
		ctx, username, email, passwordHash, role, organization,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	resp, err := authService.Login(ctx, &email, &password)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	t.Logf("Response Data: %+v", resp)
}

func TestLogout(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		authService  = injector.CommonAuthService
		userDaoMock  = injector.UserDaoMock
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		password     = "User@123"
		role         = "USER"
		organization = "ORG"
	)
	passwordHash, err := crypt.Hash(password)
	assert.NoError(t, err)
	userID, err := userDaoMock.UserDao.InsertUser(
		ctx, username, email, passwordHash, role, organization,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	resp, err := authService.Login(ctx, &email, &password)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	accessToken := resp.AccessToken
	err = authService.Logout(ctx, &accessToken)
	assert.NoError(t, err)

	cache, err := injector.Cache.Get(ctx, "token:blacklist:"+crypt.MD5(accessToken))
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	t.Logf("AccessToken: %s", accessToken)
}

func TestRefreshToken(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		authService  = injector.CommonAuthService
		userDaoMock  = injector.UserDaoMock
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		password     = "User@123"
		role         = "USER"
		organization = "ORG"
	)
	passwordHash, err := crypt.Hash(password)
	userID, err := userDaoMock.UserDao.InsertUser(
		ctx, username, email, passwordHash, role, organization,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	resp, err := authService.Login(ctx, &email, &password)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	refreshToken := resp.RefreshToken
	time.Sleep(4 * time.Second) // wait for token to expire
	respRefresh, err := authService.RefreshToken(ctx, &refreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, respRefresh)

	newAccessToken := respRefresh.AccessToken

	t.Logf("NewAccessToken: %s", newAccessToken)
}

func TestChangePassword(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		authService  = injector.CommonAuthService
		userDaoMock  = injector.UserDaoMock
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		password     = "123456"
		role         = "USER"
		organization = "ORG"
	)
	passwordHash, err := crypt.Hash(password)
	userID, err := userDaoMock.UserDao.InsertUser(
		ctx, username, email, passwordHash, role, organization,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	resp, err := authService.Login(ctx, &email, &password)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	ctx = context.WithValue(ctx, config.UserIDKey, userID.Hex())
	newPassword := "1234567"
	err = authService.ChangePassword(ctx, &password, &newPassword)
	assert.NoError(t, err)

	resp, err = authService.Login(ctx, &email, &newPassword)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	t.Logf("Response Data: %+v", resp)
}
