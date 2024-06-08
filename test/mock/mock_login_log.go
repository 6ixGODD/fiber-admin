package mock

import (
	"context"
	"math/rand"

	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginLogDaoMock struct {
	LoginLogMap map[primitive.ObjectID]*entity.LoginLogModel
	LoginLogIDs []primitive.ObjectID
	LoginLogDao mods.LoginLogDao
	UserMock    UserDaoMock
}

func NewLoginLogDaoMock(loginLogDao mods.LoginLogDao, userMock *UserDaoMock) *LoginLogDaoMock {
	return &LoginLogDaoMock{
		LoginLogMap: make(map[primitive.ObjectID]*entity.LoginLogModel),
		LoginLogDao: loginLogDao,
		UserMock:    *userMock,
	}
}

func NewLoginLogDaoMockWithRandomData(n int, loginLogDao mods.LoginLogDao, userMock *UserDaoMock) *LoginLogDaoMock {
	loginLogDaoMock := NewLoginLogDaoMock(loginLogDao, userMock)
	for i := 0; i < n; i++ {
		loginLog := loginLogDaoMock.GenerateLoginLogModel()
		loginLogDaoMock.LoginLogMap[loginLog.LoginLogID] = loginLog
		loginLogDaoMock.LoginLogIDs = append(loginLogDaoMock.LoginLogIDs, loginLog.LoginLogID)
	}
	return loginLogDaoMock
}

func (m *LoginLogDaoMock) Create(loginLog *entity.LoginLogModel) {
	m.LoginLogMap[loginLog.LoginLogID] = loginLog
}

func (m *LoginLogDaoMock) Get(loginLogID primitive.ObjectID) (*entity.LoginLogModel, error) {
	loginLog, ok := m.LoginLogMap[loginLogID]
	if !ok {
		return nil, nil
	}
	return loginLog, nil
}

func (m *LoginLogDaoMock) GenerateLoginLogModel() *entity.LoginLogModel {
	userID := m.UserMock.RandomUserID()
	ipAddress, userAgent := GenerateLoginLog()
	loginLogID, err := m.LoginLogDao.InsertLoginLog(context.Background(), userID, ipAddress, userAgent)
	if err != nil {
		panic(err)
	}

	loginLog, err := m.LoginLogDao.GetLoginLogByID(context.Background(), loginLogID)
	if err != nil {
		panic(err)
	}

	return loginLog
}

func (m *LoginLogDaoMock) GenerateLoginLogWithUserID(userID primitive.ObjectID) *entity.LoginLogModel {
	ipAddress, userAgent := GenerateLoginLog()
	loginLogID, err := m.LoginLogDao.InsertLoginLog(context.Background(), userID, ipAddress, userAgent)
	if err != nil {
		panic(err)
	}

	loginLog, err := m.LoginLogDao.GetLoginLogByID(context.Background(), loginLogID)
	if err != nil {
		panic(err)
	}
	return loginLog
}

func (m *LoginLogDaoMock) GenerateLoginLogWithIpAddress(ipAddress string) *entity.LoginLogModel {
	userID := m.UserMock.RandomUserID()
	userAgent := RandomEnum([]string{"Chrome", "Firefox", "Safari", "Edge"})
	loginLogID, err := m.LoginLogDao.InsertLoginLog(context.Background(), userID, ipAddress, userAgent)
	if err != nil {
		panic(err)
	}

	loginLog, err := m.LoginLogDao.GetLoginLogByID(context.Background(), loginLogID)
	if err != nil {
		panic(err)
	}
	return loginLog
}

func (m *LoginLogDaoMock) GenerateLoginLogWithUserAgent(userAgent string) *entity.LoginLogModel {
	userID := m.UserMock.RandomUserID()
	ipAddress := RandomIp()
	loginLogID, err := m.LoginLogDao.InsertLoginLog(context.Background(), userID, ipAddress, userAgent)
	if err != nil {
		panic(err)
	}

	loginLog, err := m.LoginLogDao.GetLoginLogByID(context.Background(), loginLogID)
	if err != nil {
		panic(err)
	}
	return loginLog
}

func (m *LoginLogDaoMock) RandomLoginLogID() primitive.ObjectID {
	return m.LoginLogIDs[rand.Intn(len(m.LoginLogIDs))]
}

func (m *LoginLogDaoMock) Delete() {
	for _, loginLogID := range m.LoginLogIDs {
		_ = m.LoginLogDao.DeleteLoginLog(context.Background(), loginLogID)
	}
}

func GenerateLoginLog() (ipAddress, userAgent string) {
	ipAddress = RandomIp()
	userAgent = RandomEnum([]string{"Chrome", "Firefox", "Safari", "Edge"})
	return ipAddress, userAgent
}
