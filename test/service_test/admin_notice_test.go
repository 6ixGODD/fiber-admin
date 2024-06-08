package service_test

import (
	"testing"

	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertNotice(t *testing.T) {
	var (
		injector      = wire.GetInjector()
		ctx           = injector.Ctx
		noticeService = injector.AdminNoticeService
		title         = mock.RandomString(10)
		content       = mock.RandomString(10)
		noticeType    = mock.RandomEnum([]string{"NORMAL", "URGENT"})
	)
	noticeIDHex, err := noticeService.InsertNotice(ctx, &title, &content, &noticeType)
	assert.NoError(t, err)
	assert.NotEmpty(t, noticeIDHex)

	t.Logf("Notice ID: %s", noticeIDHex)

	noticeID, err := primitive.ObjectIDFromHex(noticeIDHex)
	assert.NoError(t, err)
	notice, err := injector.NoticeDao.GetNoticeByID(ctx, noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, title, notice.Title)
	assert.Equal(t, content, notice.Content)
	assert.Equal(t, noticeType, notice.NoticeType)
}

func TestUpdateNotice(t *testing.T) {
	var (
		injector      = wire.GetInjector()
		ctx           = injector.Ctx
		noticeService = injector.AdminNoticeService
		noticeID      = injector.NoticeDaoMock.RandomNoticeID()
		title         = mock.RandomString(10)
		content       = mock.RandomString(10)
		noticeType    = mock.RandomEnum([]string{"NORMAL", "URGENT"})
	)
	err := noticeService.UpdateNotice(ctx, &noticeID, &title, &content, &noticeType)
	assert.NoError(t, err)

	notice, err := injector.NoticeDao.GetNoticeByID(ctx, noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, title, notice.Title)
	assert.Equal(t, content, notice.Content)
	assert.Equal(t, noticeType, notice.NoticeType)

	t.Logf("Notice Data: %+v", notice)
}

func TestDeleteNotice(t *testing.T) {
	var (
		injector      = wire.GetInjector()
		ctx           = injector.Ctx
		noticeService = injector.AdminNoticeService
		noticeID      = injector.NoticeDaoMock.RandomNoticeID()
	)
	err := noticeService.DeleteNotice(ctx, &noticeID)
	assert.NoError(t, err)

	notice, err := injector.NoticeDao.GetNoticeByID(ctx, noticeID)
	assert.Error(t, err)
	assert.Nil(t, notice)
}
