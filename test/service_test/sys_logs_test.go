package service_test

import (
	"testing"

	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestInsertLoginLog(t *testing.T) {
	var (
		injector       = wire.GetInjector()
		ctx            = injector.Ctx
		sysLogsService = injector.SysLogsService
		userID         = injector.UserDaoMock.RandomUserID()
		ip             = mock.RandomIp()
		userAgent      = mock.RandomEnum([]string{"Chrome", "Firefox", "Safari"})
	)
	err := sysLogsService.InsertLoginLog(ctx, &userID, &ip, &userAgent)
	assert.NoError(t, err)
}

func TestCacheLoginLog(t *testing.T) {}

func TestInsertOperationLog(t *testing.T) {

}

func TestCacheOperationLog(t *testing.T) {}
