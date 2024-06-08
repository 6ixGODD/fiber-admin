package dao_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var operationLogID primitive.ObjectID

func TestInsertOperationLog(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
		userID          = injector.UserDaoMock.RandomUserID()
		entityID        = injector.UserDaoMock.RandomUserID()
		entityType      = "USER"
		ipAddress       = "123.456.789.100"
		operation       = "CREATE"
		userAgent       = "Mozilla/5.0"
		description     = "Create user"
		status          = "SUCCESS"
		err             error
	)
	operationLogID, err = operationLogDao.InsertOperationLog(
		ctx, userID, entityID, ipAddress, userAgent, operation, entityType, description, status,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, operationLogID)

	operationLog, err := operationLogDao.GetOperationLogByID(ctx, operationLogID)
	assert.NoError(t, err)
	assert.NotNil(t, operationLog)
	assert.Equal(t, userID, operationLog.UserID)
	assert.Equal(t, entityID, operationLog.EntityID)
	assert.Equal(t, ipAddress, operationLog.IPAddress)
	assert.Equal(t, userAgent, operationLog.UserAgent)
	assert.Equal(t, operation, operationLog.Operation)
	assert.Equal(t, entityType, operationLog.EntityType)
	assert.Equal(t, description, operationLog.Description)
	assert.Equal(t, status, operationLog.Status)
}

func TestCacheOperationLog(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
		userID          = injector.UserDaoMock.RandomUserID()
		entityID        = injector.UserDaoMock.RandomUserID()
		entityType      = "cache USER"
		ipAddress       = "cache 123.456.789.100"
		operation       = "cache CREATE"
		userAgent       = "cache Mozilla/5.0"
		description     = "cache Create user"
		status          = "cache SUCCESS"
	)
	err := operationLogDao.CacheOperationLog(
		ctx, userID, entityID, ipAddress, userAgent, operation, entityType, description, status,
	)
	assert.NoError(t, err)
	pop, err := wire.GetInjector().Cache.LeftPop(ctx, "log:operation")
	assert.NoError(t, err)
	assert.NotNil(t, pop)
	assert.NotEmpty(t, *pop)
	t.Logf("Operation log: %s", *pop)
	t.Logf("=====================================")
	err = operationLogDao.CacheOperationLog(
		ctx, userID, entityID, ipAddress, userAgent, operation, entityType, description, status,
	)
	assert.NoError(t, err)
}

func TestSyncOperationLog(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
		entityType      = "cache USER"
		ipAddress       = "cache 123.456.789.100"
		operation       = "cache CREATE"
		userAgent       = "cache Mozilla/5.0"
		description     = "cache Create user"
		status          = "cache SUCCESS"
	)
	operationLogDao.SyncOperationLog(ctx)
	operationLogList, count, err := operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, &ipAddress, &operation, &entityType, &status, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotEmpty(t, *count)
	assert.NotEmpty(t, operationLogList)
	assert.Equal(t, 1, len(operationLogList))
	t.Logf("User ID: %s", userID)
	t.Logf("Entity type: %s", entityType)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Operation: %s", operation)
	t.Logf("User agent: %s", userAgent)
	t.Logf("Description: %s", description)
	t.Logf("Status: %s", status)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")
}

func TestGetOperationLogList(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
		userID          = injector.UserDaoMock.RandomUserID()
		startTime       = time.Now().AddDate(0, 0, -1)
		endTime         = time.Now()
		entityID        = wire.GetInjector().UserDaoMock.RandomUserID()
		entityType      = "USER"
		ipAddress       = "123.456.789.100"
		operation       = "CREATE"
		status          = "SUCCESS"
		query           = "a"
	)
	for i := 0; i < 100; i++ {
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithUserID(userID)
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithIpAddress(ipAddress)
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithEntityID(entityID)
	}
	operationLogList, count, err := operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.NotEmpty(t, *count)
	assert.NotEmpty(t, operationLogList)
	assert.Equal(t, 10, len(operationLogList))
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, &startTime, &endTime, nil, nil, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Start time: %s", startTime)
	t.Logf("End time: %s", endTime)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, &userID, nil, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("User ID: %s", userID)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, &entityID, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Entity ID: %s", entityID)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, &entityType, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Entity type: %s", entityType)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, &ipAddress, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, &operation, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Operation: %s", operation)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, &entityType, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Entity type: %s", entityType)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, nil, &status, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Status: %s", status)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, nil, nil, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Query: %s", query)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")

	operationLogList, count, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, &startTime, &endTime, &userID, &entityID, &ipAddress, &operation, &entityType,
		&status, &query,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Start time: %s", startTime)
	t.Logf("End time: %s", endTime)
	t.Logf("User ID: %s", userID)
	t.Logf("Entity ID: %s", entityID)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Operation: %s", operation)
	t.Logf("Entity type: %s", entityType)
	t.Logf("Status: %s", status)
	t.Logf("Query: %s", query)
	t.Logf("Operation log count: %d", *count)
	t.Logf("Operation log list: %v", operationLogList)
	t.Logf("=====================================")
}

func TestDeleteOperationLog(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
	)
	err := operationLogDao.DeleteOperationLog(ctx, operationLogID)
	assert.NoError(t, err)

	operationLog, err := operationLogDao.GetOperationLogByID(ctx, operationLogID)
	assert.Error(t, err)
	assert.Nil(t, operationLog)
}

func TestDeleteOperationLogList(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		operationLogDao = injector.OperationLogDao
		userID          = injector.UserDaoMock.RandomUserID()
		entityID        = injector.UserDaoMock.RandomUserID()
		entityType      = "USER"
		ipAddress       = "123.456.789.100"
		operation       = "CREATE"
		status          = "SUCCESS"
	)
	for i := 0; i < 10; i++ {
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithUserID(userID)
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithIpAddress(ipAddress)
		_ = injector.OperationLogDaoMock.GenerateOperationLogWithEntityID(entityID)
	}
	count, err := operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, &userID, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err := operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, &userID, nil, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("User ID: %s", userID)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	count, err = operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, nil, &entityID, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, &entityID, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("Entity ID: %s", entityID)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	count, err = operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, nil, nil, &ipAddress, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, &ipAddress, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("IP address: %s", ipAddress)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	count, err = operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, nil, nil, nil, &operation, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, &operation, nil, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("Operation: %s", operation)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	count, err = operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, nil, nil, nil, nil, &entityType, nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, nil, &entityType, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("Entity type: %s", entityType)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")

	count, err = operationLogDao.DeleteOperationLogList(
		ctx, nil, nil, nil, nil, nil, nil, nil, &status,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	operationLogList, _, err = operationLogDao.GetOperationLogList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil, nil, nil, &status, nil,
	)
	assert.NoError(t, err)
	assert.Empty(t, operationLogList)
	t.Logf("Status: %s", status)
	t.Logf("Delete count: %d", *count)
	t.Logf("=====================================")
}
