package service_test

import (
	"testing"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIdempotencyToken(t *testing.T) {
	var (
		injector           = wire.GetInjector()
		ctx                = injector.Ctx
		idempotencyService = injector.CommonIdempotencyService
	)
	resp, err := idempotencyService.GenerateIdempotencyToken(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	t.Logf("Response Data: %+v", resp)
}

func TestCheckIdempotencyToken(t *testing.T) {
	var (
		injector           = wire.GetInjector()
		ctx                = injector.Ctx
		idempotencyService = injector.CommonIdempotencyService
	)
	token, err := idempotencyService.GenerateIdempotencyToken(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	err = idempotencyService.CheckIdempotencyToken(ctx, token)
	assert.NoError(t, err)
}
