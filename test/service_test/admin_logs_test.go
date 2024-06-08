package service_test

import (
	"testing"
	"time"

	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginLogList(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		logsService     = injector.AdminLogsService
		page            = int64(1)
		pageSize        = int64(10)
		desc            = false
		createStartTime = time.Now().AddDate(0, 0, -1)
		createEndTime   = time.Now()
		query           = "a"
	)
	resp, err := logsService.GetLoginLogList(ctx, &page, &pageSize, &desc, &query, &createStartTime, &createEndTime)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = logsService.GetLoginLogList(ctx, &page, &pageSize, &desc, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.LoginLogList)
	assert.Equal(t, pageSize, int64(len(resp.LoginLogList)))

	t.Logf("Response Data: %+v", resp)
}

func TestGetOperationLogList(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		logsService     = injector.AdminLogsService
		page            = int64(1)
		pageSize        = int64(10)
		desc            = false
		createStartTime = time.Now().AddDate(0, 0, -1)
		createEndTime   = time.Now()
		query           = "a"
		operation       = mock.RandomEnum([]string{"CREATE", "UPDATE", "DELETE"})
		entityType      = mock.RandomEnum([]string{"USER", "DOCUMENTATION"})
		status          = mock.RandomEnum([]string{"SUCCESS", "FAILED"})
	)
	resp, err := logsService.GetOperationLogList(
		ctx, &page, &pageSize, &desc, &query, &operation, &entityType, &status, &createStartTime, &createEndTime,
	)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = logsService.GetOperationLogList(ctx, &page, &pageSize, &desc, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.OperationLogList)
	assert.Equal(t, pageSize, int64(len(resp.OperationLogList)))

	t.Logf("Response Data: %+v", resp)
}
