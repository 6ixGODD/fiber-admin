package service_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestGetDocumentation(t *testing.T) {
	var (
		injector             = wire.GetInjector()
		ctx                  = injector.Ctx
		documentationService = injector.CommonDocumentationService
		documentationID      = injector.DocumentationDaoMock.RandomDocumentationID()
	)

	resp, err := documentationService.GetDocumentation(ctx, &documentationID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	t.Logf("Response Data: %+v", resp)
}

func TestGetDocumentationList(t *testing.T) {
	var (
		injector             = wire.GetInjector()
		ctx                  = injector.Ctx
		documentationService = injector.CommonDocumentationService
		page                 = int64(1)
		pageSize             = int64(10)
		updateStartTime      = time.Now().AddDate(0, 0, -1)
		updateEndTime        = time.Now()
	)

	resp, err := documentationService.GetDocumentationList(ctx, &page, &pageSize, &updateStartTime, &updateEndTime)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = documentationService.GetDocumentationList(ctx, &page, &pageSize, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.DocumentationSummaryList)
	assert.Equal(t, pageSize, int64(len(resp.DocumentationSummaryList)))

	t.Logf("Response Data: %+v", resp)
}
