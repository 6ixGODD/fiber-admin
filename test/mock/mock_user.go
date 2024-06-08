package mock

import (
	"context"
	"math/rand"

	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserDaoMock is a mock for UserDao
type UserDaoMock struct {
	UserMap map[primitive.ObjectID]*entity.UserModel
	UserIDs []primitive.ObjectID
	UserDao mods.UserDao
}

// NewUserDaoMock returns a new UserDaoMock
func NewUserDaoMock(userDao mods.UserDao) *UserDaoMock {
	return &UserDaoMock{
		UserMap: make(map[primitive.ObjectID]*entity.UserModel),
		UserDao: userDao,
	}
}

// NewUserDaoMockWithRandomData returns a new UserDaoMock with n random users
func NewUserDaoMockWithRandomData(n int, userDao mods.UserDao) *UserDaoMock {
	userDaoMock := NewUserDaoMock(userDao)
	for i := 0; i < n; i++ {
		user := userDaoMock.GenerateUserModel()
		userDaoMock.UserMap[user.UserID] = user
		userDaoMock.UserIDs = append(userDaoMock.UserIDs, user.UserID)
	}
	return userDaoMock
}

// Create mocks the Create method
func (m *UserDaoMock) Create(user *entity.UserModel) {
	m.UserMap[user.UserID] = user
}

// Get mocks the Get method
func (m *UserDaoMock) Get(userID primitive.ObjectID) (*entity.UserModel, error) {
	user, ok := m.UserMap[userID]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *UserDaoMock) RandomUserID() primitive.ObjectID {
	return m.UserIDs[rand.Intn(len(m.UserIDs))]
}

// GenerateUserModel generates a new UserModel
func (m *UserDaoMock) GenerateUserModel() *entity.UserModel {
	username, email, password, role, organization := GenerateUser()
	userID, err := m.UserDao.InsertUser(context.Background(), username, email, password, role, organization)
	if err != nil {
		panic(err)
	}
	if err = m.UserDao.UpdateUserLastLogin(context.Background(), userID); err != nil {
		panic(err)
	}
	user, err := m.UserDao.GetUserByID(context.Background(), userID)
	if err != nil {
		panic(err)
	}
	return user
}

func (m *UserDaoMock) Delete() {
	for _, userID := range m.UserIDs {
		_ = m.UserDao.DeleteUser(context.Background(), userID)
	}
}

// GenerateUser generates a new user, returns username, email, password, role, organization
func GenerateUser() (string, string, string, string, string) {
	return RandomString(10),
		RandomString(10) + "@test.com",
		RandomString(10),
		RandomEnum([]string{"ADMIN", "USER"}),
		RandomEnum([]string{"FOO", "BAR", "BAZ", "QUX"})
}
