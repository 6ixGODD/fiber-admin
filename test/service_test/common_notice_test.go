package service_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

func TestGetNotice(t *testing.T) {
	var (
		injector      = wire.GetInjector()
		ctx           = injector.Ctx
		noticeService = injector.CommonNoticeService
		noticeID      = injector.NoticeDaoMock.RandomNoticeID()
	)

	resp, err := noticeService.GetNotice(ctx, &noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	t.Logf("Response Data: %+v", resp)
}

func TestGetNoticeList(t *testing.T) {
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		noticeService   = injector.CommonNoticeService
		page            = int64(1)
		pageSize        = int64(10)
		noticeType      = "URGENT"
		updateStartTime = time.Now().AddDate(0, 0, -1)
		updateEndTime   = time.Now()
	)

	resp, err := noticeService.GetNoticeList(ctx, &page, &pageSize, &noticeType, &updateStartTime, &updateEndTime)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = noticeService.GetNoticeList(ctx, &page, &pageSize, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.NoticeSummaryList)
	assert.Equal(t, pageSize, int64(len(resp.NoticeSummaryList)))

	t.Logf("Response Data: %+v", resp)
}
