package dao_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var loginLogID primitive.ObjectID

func TestInsertLoginLog(t *testing.T) {
	// t.Skip("Skip TestInsertLoginLog")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		userID      = injector.UserDaoMock.RandomUserID()
		ipAddress   = "123.456.789.012"
		userAgent   = "User Agent"
		err         error
	)

	loginLogID, err = loginLogDao.InsertLoginLog(ctx, userID, ipAddress, userAgent)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginLogID)

	loginLog, err := loginLogDao.GetLoginLogByID(ctx, loginLogID)
	assert.NoError(t, err)
	assert.NotNil(t, loginLog)
	assert.NotEmpty(t, loginLog.LoginLogID)
	assert.NotEmpty(t, loginLog.UserID)
	assert.NotEmpty(t, loginLog.IPAddress)
	assert.NotEmpty(t, loginLog.UserAgent)
	assert.NotEmpty(t, loginLog.CreatedAt)
}

func TestCacheLoginLog(t *testing.T) {
	// t.Skip("Skip TestCacheLoginLog")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		userID      = injector.UserDaoMock.RandomUserID()
		ipAddress   = "cache 123.456.789.012"
		userAgent   = "cache User Agent"
		err         error
	)

	err = loginLogDao.CacheLoginLog(ctx, userID, ipAddress, userAgent)
	assert.NoError(t, err)
	pop, err := wire.GetInjector().Cache.LeftPop(ctx, "log:login")
	assert.NoError(t, err)
	assert.NotNil(t, pop)
	assert.NotEmpty(t, *pop)
	t.Logf("Login log: %s", *pop)
	t.Logf("=====================================")
	err = loginLogDao.CacheLoginLog(ctx, userID, ipAddress, userAgent)
	assert.NoError(t, err)
}

func TestSyncLoginLog(t *testing.T) {
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		ipAddress   = "cache 123.456.789.012"
		userAgent   = "cache User Agent"
		err         error
	)
	// t.Skip("Skip TestSyncLoginLog")
	loginLogDao.SyncLoginLog(ctx)
	loginLogList, count, err := loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, &ipAddress, &userAgent, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotEmpty(t, *count)
	assert.NotEmpty(t, loginLogList)
	// assert.Equal(t, 1, len(loginLogList))
	t.Logf("IP address: %s", ipAddress)
	t.Logf("User agent: %s", userAgent)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")
}

func TestGetLoginLogList(t *testing.T) {
	// t.Skip("Skip TestGetLoginLogList")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		userID      = injector.UserDaoMock.RandomUserID()
		startTime   = time.Now().AddDate(0, 0, -1)
		endTime     = time.Now()
		ipAddress   = "123.456.789.012"
		userAgent   = "User Agent"
		query       = "a"
		err         error
	)
	for i := 0; i < 100; i++ {
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithUserID(userID)
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithIpAddress(ipAddress)
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithUserAgent(userAgent)
	}

	loginLogList, count, err := loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotEmpty(t, *count)
	assert.NotEmpty(t, loginLogList)
	assert.Equal(t, 10, len(loginLogList))
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, &startTime, &endTime, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Start time: %s", startTime)
	t.Logf("End time: %s", endTime)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, &userID, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("User ID: %s", userID)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, &ipAddress, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, &userAgent, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("User agent: %s", userAgent)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Query: %s", query)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")

	loginLogList, count, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, &startTime, &endTime, &userID, &ipAddress, &userAgent, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Start time: %s", startTime)
	t.Logf("End time: %s", endTime)
	t.Logf("User ID: %s", userID)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("User agent: %s", userAgent)
	t.Logf("Query: %s", query)
	t.Logf("Login log count: %d", *count)
	t.Logf("Login log list: %v", loginLogList)
	t.Logf("=====================================")
}

func TestDeleteLoginLog(t *testing.T) {
	// t.Skip("Skip TestDeleteLoginLog")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		err         error
	)
	err = loginLogDao.DeleteLoginLog(ctx, loginLogID)
	assert.NoError(t, err)

	loginLog, err := loginLogDao.GetLoginLogByID(ctx, loginLogID)
	assert.Error(t, err)
	assert.Nil(t, loginLog)
}

func TestDeleteLoginLogList(t *testing.T) {
	// t.Skip("Skip TestDeleteLoginLogList")
	var (
		injector    = wire.GetInjector()
		ctx         = injector.Ctx
		loginLogDao = injector.LoginLogDao
		userID      = injector.UserDaoMock.RandomUserID()
		ipAddress   = "123.456.789.012"
		userAgent   = "User Agent"
		err         error
	)
	for i := 0; i < 100; i++ {
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithUserID(userID)
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithIpAddress(ipAddress)
		_ = wire.GetInjector().LoginLogDaoMock.GenerateLoginLogWithUserAgent(userAgent)
	}

	count, err := wire.GetInjector().LoginLogDao.DeleteLoginLogList(ctx, nil, nil, &userID, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	loginLogList, _, err := loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, &userID, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, loginLogList)
	t.Logf("User ID: %s", userID)
	t.Logf("Delete count: %d", count)
	t.Logf("=====================================")

	count, err = loginLogDao.DeleteLoginLogList(ctx, nil, nil, nil, &ipAddress, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	loginLogList, _, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, &ipAddress, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, loginLogList)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Delete count: %d", count)
	t.Logf("=====================================")

	count, err = loginLogDao.DeleteLoginLogList(ctx, nil, nil, nil, nil, &userAgent)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	loginLogList, _, err = loginLogDao.GetLoginLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, &userAgent, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, loginLogList)
	t.Logf("User agent: %s", userAgent)
	t.Logf("Delete count: %d", count)
	t.Logf("=====================================")
}
