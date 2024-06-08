package service_test

import (
	"testing"

	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertUser(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		userService  = injector.AdminUserService
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		password     = "User@123"
		organization = mock.RandomString(10)
	)

	userIDHex, err := userService.InsertUser(ctx, &username, &email, &password, &organization)
	assert.NoError(t, err)
	assert.NotEmpty(t, userIDHex)

	t.Logf("User ID: %s", userIDHex)

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	assert.NoError(t, err)
	user, err := userService.GetUser(ctx, &userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	t.Logf("User: %+v", user)
}

func TestGetUserList(t *testing.T) {
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		userService = injector.AdminUserService
		page        = int64(1)
		pageSize    = int64(10)
		desc        = false
		role        = "USER"
		query       = "a"
	)

	userList, err := userService.GetUserList(ctx, &page, &pageSize, &desc, &role, nil, nil, nil, nil, &query)
	assert.NoError(t, err)
	assert.NotNil(t, userList)
	t.Logf("User List: %+v", userList)

	userList, err = userService.GetUserList(ctx, &page, &pageSize, &desc, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, userList)
	assert.Equal(t, pageSize, int64(len(userList.UserList)))
	t.Logf("User List: %+v", userList)
}

func TestUpdateUser(t *testing.T) {
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		userService  = injector.AdminUserService
		userID       = injector.UserDaoMock.RandomUserID()
		username     = mock.RandomString(10)
		email        = mock.RandomString(10) + "@user.com"
		organization = mock.RandomString(10)
	)
	err := userService.UpdateUser(ctx, &userID, &username, &email, &organization)
	assert.NoError(t, err)

	user, err := injector.UserDao.GetUserByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, organization, user.Organization)

	t.Logf("User: %+v", user)
}

func TestDeleteUser(t *testing.T) {
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		userService = injector.AdminUserService
		userID      = injector.UserDaoMock.RandomUserID()
	)

	err := userService.DeleteUser(ctx, &userID)
	assert.NoError(t, err)

	user, err := injector.UserDao.GetUserByID(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, user)
}
