package dao_test

import (
	"testing"
	"time"

	"fiber-admin/internal/pkg/domain/entity"
	"fiber-admin/pkg/utils/crypt"
	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userID   primitive.ObjectID
	username string
	email    string
)

func TestInsertUser(t *testing.T) {
	// t.Skip("Skip TestInsertUser")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		userDao     = injector.UserDao
		password, _ = crypt.Hash("Admin@123")
		role        = "ADMIN"
		org         = "Fiber Admin"
		err         error
	)
	username = mock.RandomString(10)
	email = mock.RandomString(10) + "@gmail.com"
	assert.NoError(t, err)
	userID, err = userDao.InsertUser(ctx, username, email, password, role, org)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	user, err := userDao.GetUserByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, org, user.Organization)
	assert.True(t, crypt.Compare("Admin@123", user.Password))
	_, err = userDao.InsertUser(ctx, username, email, password, role, org)
	assert.Error(t, err)
	t.Logf("Error: %v", err)
}

func TestGetUser(t *testing.T) {
	// t.Skip("Skip TestGetUser")
	var (
		injector = wire.GetInjector()
		ctx      = injector.Ctx
		userDao  = injector.UserDao
		err      error
	)
	user, err := userDao.GetUserByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.UserID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Role)
	assert.NotEmpty(t, user.Organization)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	assert.False(t, user.Deleted)

	_, _ = userDao.GetUserByID(ctx, userID)

	userCache, err := wire.GetInjector().Cache.Get(ctx, "dao:user:userID:"+user.UserID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, userCache)
	assert.NotEmpty(t, userCache)
	t.Logf("User wire.GetInjector().Cache: %v", *userCache)

	user, err = userDao.GetUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.UserID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Role)
	assert.NotEmpty(t, user.Organization)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	assert.False(t, user.Deleted)
	_, _ = userDao.GetUserByEmail(ctx, email)

	user, err = userDao.GetUserByUsername(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.UserID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Role)
	assert.NotEmpty(t, user.Organization)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	assert.False(t, user.Deleted)
	_, _ = userDao.GetUserByUsername(ctx, username)

	userNil, err := userDao.GetUserByID(ctx, primitive.NewObjectID())
	assert.Error(t, err)
	assert.Nil(t, userNil)

	userNil, err = userDao.GetUserByEmail(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, userNil)

	userNil, err = userDao.GetUserByUsername(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, userNil)
}

func TestGetUserList(t *testing.T) {
	// t.Skip("Skip TestGetUserList")
	var (
		injector           = wire.GetInjector()
		ctx                = injector.Ctx
		userDao            = injector.UserDao
		organization       = "FOO"
		role               = "USER"
		createStartTime    = time.Now().Add(-time.Hour)
		createEndTime      = time.Now()
		updateStartTime    = time.Now().Add(-time.Hour)
		updateEndTime      = time.Now()
		lastLoginStartTime = time.Now().Add(-time.Hour)
		lastLoginEndTime   = time.Now()
		query              = "Fo"
	)
	userList, count, err := userDao.GetUserList(
		ctx, 0, 10, false, nil, nil, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotEmpty(t, *count)
	assert.NotEmpty(t, userList)
	assert.NotNil(t, userList)
	assert.Equal(t, 10, len(userList))
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, &organization, nil, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Organization: %s", organization)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, &organization, &role, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Organization: %s", organization)
	t.Logf("Role: %s", role)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, nil, &createStartTime,
		&createEndTime, nil, nil, nil,
		nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Create time start: %v", createStartTime)
	t.Logf("Create time end: %v", createEndTime)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, nil, nil,
		nil, &updateStartTime, &updateEndTime, nil,
		nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Update time start: %v", updateStartTime)
	t.Logf("Update time end: %v", updateEndTime)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, nil, nil,
		nil, nil, nil, &lastLoginStartTime,
		&lastLoginEndTime, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Last login time start: %v", lastLoginStartTime)
	t.Logf("Last login time end: %v", lastLoginEndTime)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, nil, nil,
		nil, nil, nil, nil,
		nil, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Query: %s", query)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, &organization, &role, &createStartTime,
		&createEndTime, &updateStartTime, &updateEndTime, &lastLoginStartTime,
		&lastLoginEndTime, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotNil(t, userList)
	t.Logf("Organization: %s", organization)
	t.Logf("Role: %s", role)
	t.Logf("Create time start: %v", createStartTime)
	t.Logf("Create time end: %v", createEndTime)
	t.Logf("Update time start: %v", updateStartTime)
	t.Logf("Update time end: %v", updateEndTime)
	t.Logf("Last login time start: %v", lastLoginStartTime)
	t.Logf("Last login time end: %v", lastLoginEndTime)
	t.Logf("Query: %s", query)
	t.Logf("User count: %d", *count)
	t.Logf("User list: %v", userList)
	t.Logf("=====================================")
}

func TestUpdateUser(t *testing.T) {
	// t.Skip("Skip TestUpdateUser")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		userDao     = injector.UserDao
		username    = mock.RandomString(10)
		email       = mock.RandomString(10) + "@gmail.com"
		role        = "USER"
		org         = "Fiber Admin"
		password, _ = crypt.Hash("User@123")
		err         error
	)
	err = userDao.UpdateUser(ctx, userID, &username, &email, &password, &role, &org)
	assert.NoError(t, err)

	user, err := userDao.GetUserByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, org, user.Organization)
	assert.True(t, crypt.Compare("User@123", user.Password))
}

func TestDeleteUser(t *testing.T) {
	// t.Skip("Skip TestDeleteUser")
	var (
		injector = wire.GetInjector()
		ctx      = injector.Ctx
		userDao  = injector.UserDao
		err      error
	)
	err = userDao.SoftDeleteUser(ctx, userID)
	assert.NoError(t, err)

	user, err := userDao.GetUserByID(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, user)

	err = userDao.DeleteUser(ctx, userID)
	assert.NoError(t, err)

	user, err = userDao.GetUserByID(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestDeleteUserList(t *testing.T) {
	// t.Skip("Skip TestDeleteUserList")
	var (
		injector     = wire.GetInjector()
		ctx          = injector.Ctx
		userDao      = injector.UserDao
		organization = "FOO"
		role         = "USER"
	)

	count, err := userDao.SoftDeleteUserList(ctx, &organization, nil, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Organization: %s", organization)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	userList, count, err := userDao.GetUserList(
		ctx, 0, 10, false, &organization, nil, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.Empty(t, userList)

	count, err = userDao.DeleteUserList(ctx, nil, &role, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Role: %s", role)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, &role, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.Empty(t, userList)

	count, err = userDao.SoftDeleteUserList(ctx, nil, &role, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Role: %s", role)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, &role, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.Empty(t, userList)

	count, err = userDao.DeleteUserList(ctx, nil, &role, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Role: %s", role)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	userList, count, err = userDao.GetUserList(
		ctx, 0, 10, false, nil, &role, nil,
		nil, nil, nil, nil,
		nil, nil,
	)
	assert.Empty(t, userList)
}

func BenchmarkInsertUser(b *testing.B) {
	var (
		injector = wire.GetInjector()
		ctx      = injector.Ctx
		userDao  = injector.UserDao
	)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		usernameMock, emailMock, password, role, org := mock.GenerateUser()
		b.StartTimer()
		InsertUserID, err := userDao.InsertUser(ctx, usernameMock, emailMock, password, role, org)
		b.StopTimer()
		assert.NoError(b, err)
		assert.NotEmpty(b, InsertUserID)
		injector.UserDaoMock.Create(
			&entity.UserModel{
				UserID:       InsertUserID,
				Username:     usernameMock,
				Email:        emailMock,
				Password:     password,
				Role:         role,
				Organization: org,
			},
		)
		b.StartTimer()
	}
}

func BenchmarkGetUser(b *testing.B) {
	var (
		injector = wire.GetInjector()
		ctx      = injector.Ctx
		userDao  = injector.UserDao
	)
	for i := 0; i < b.N; i++ {
		userID := injector.UserDaoMock.RandomUserID()
		user, err := userDao.GetUserByID(ctx, userID)
		assert.NoError(b, err)
		assert.NotNil(b, user)
	}

	for i := 0; i < b.N; i++ {
		email := injector.UserDaoMock.UserMap[injector.UserDaoMock.RandomUserID()].Email
		user, err := userDao.GetUserByEmail(ctx, email)
		assert.NoError(b, err)
		assert.NotNil(b, user)
	}

	for i := 0; i < b.N; i++ {
		username := injector.UserDaoMock.UserMap[injector.UserDaoMock.RandomUserID()].Username
		user, err := userDao.GetUserByUsername(ctx, username)
		assert.NoError(b, err)
		assert.NotNil(b, user)
	}
}

func BenchmarkUpdateUser(b *testing.B) {

}

func BenchmarkDeleteUser(b *testing.B) {

}
