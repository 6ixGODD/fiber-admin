package service_test

import (
	"context"
	"testing"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	var (
		injector       = wire.GetInjector()
		ctx            = injector.Ctx
		profileService = injector.CommonProfileService
		userID         = injector.UserDaoMock.RandomUserID()
	)
	ctx = context.WithValue(ctx, config.UserIDKey, userID.Hex())
	resp, err := profileService.GetProfile(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	t.Logf("Response Data: %+v", resp)
}
