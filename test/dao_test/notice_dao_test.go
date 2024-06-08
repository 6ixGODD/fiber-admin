package dao_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var noticeID primitive.ObjectID

func TestInsertNotice(t *testing.T) {
	// t.Skip("Skip TestInsertNotice")
	var (
		injector   = wire.GetInjector()
		noticeDao  = injector.NoticeDao
		ctx        = injector.Ctx
		title      = "Title"
		content    = "Content"
		noticeType = "NORMAL"
		err        error
	)

	noticeID, err = noticeDao.InsertNotice(ctx, title, content, noticeType)
	assert.NoError(t, err)
	assert.NotEmpty(t, noticeID)

	notice, err := noticeDao.GetNoticeByID(ctx, noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, title, notice.Title)
	assert.Equal(t, content, notice.Content)
	assert.Equal(t, noticeType, notice.NoticeType)
}

func TestGetNotice(t *testing.T) {
	// t.Skip("Skip TestGetNotice")
	var (
		injector  = wire.GetInjector()
		ctx       = injector.Ctx
		noticeDao = injector.NoticeDao
		err       error
	)
	notice, err := noticeDao.GetNoticeByID(ctx, noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.NotEmpty(t, notice.NoticeID)
	assert.NotEmpty(t, notice.Title)
	assert.NotEmpty(t, notice.Content)
	assert.NotEmpty(t, notice.NoticeType)
	assert.NotEmpty(t, notice.CreatedAt)
	assert.NotEmpty(t, notice.UpdatedAt)
}

func TestGetNoticeList(t *testing.T) {
	// t.Skip("Skip TestGetNoticeList")
	var (
		injector        = wire.GetInjector()
		ctx             = injector.Ctx
		noticeDao       = injector.NoticeDao
		createStartTime = time.Now().Add(-time.Hour)
		createEndTime   = time.Now().Add(time.Hour)
		updateStartTime = time.Now().Add(-time.Hour)
		updateEndTime   = time.Now().Add(time.Hour)
		noticeType      = "NORMAL"
	)

	noticeList, count, err := noticeDao.GetNoticeList(
		ctx, 0, 10, false, nil, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	assert.NotEmpty(t, noticeList)
	t.Logf("Notice Count: %d", *count)
	t.Logf("Notice List: %v", noticeList)
	t.Logf("=====================================")

	noticeList, count, err = noticeDao.GetNoticeList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Create Start Time: %s", createStartTime)
	t.Logf("Create End Time: %s", createEndTime)
	t.Logf("Notice Count: %d", *count)
	t.Logf("Notice List: %v", noticeList)
	t.Logf("=====================================")

	noticeList, count, err = noticeDao.GetNoticeList(
		ctx, 0, 10, false, nil, nil, &updateStartTime, &updateEndTime, nil,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Update Start Time: %s", updateStartTime)
	t.Logf("Update End Time: %s", updateEndTime)
	t.Logf("Notice Count: %d", *count)
	t.Logf("Notice List: %v", noticeList)
	t.Logf("=====================================")

	noticeList, count, err = noticeDao.GetNoticeList(
		ctx, 0, 10, false, nil, nil, nil, nil, &noticeType,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Notice Type: %s", noticeType)
	t.Logf("Notice Count: %d", *count)
	t.Logf("Notice List: %v", noticeList)
	t.Logf("=====================================")

	noticeList, count, err = noticeDao.GetNoticeList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, &updateStartTime, &updateEndTime, &noticeType,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Create Start Time: %s", createStartTime)
	t.Logf("Create End Time: %s", createEndTime)
	t.Logf("Update Start Time: %s", updateStartTime)
	t.Logf("Update End Time: %s", updateEndTime)
	t.Logf("Notice Type: %s", noticeType)
	t.Logf("Notice Count: %d", *count)
	t.Logf("Notice List: %v", noticeList)
	t.Logf("=====================================")
}

func TestUpdateNotice(t *testing.T) {
	// t.Skip("Skip TestUpdateNotice")
	var (
		injector   = wire.GetInjector()
		ctx        = injector.Ctx
		noticeDao  = injector.NoticeDao
		title      = "New Title"
		content    = "New Content"
		noticeType = "NORMAL"
	)

	err := noticeDao.UpdateNotice(ctx, noticeID, &title, &content, &noticeType)
	assert.NoError(t, err)

	notice, err := noticeDao.GetNoticeByID(ctx, noticeID)
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, title, notice.Title)
	assert.Equal(t, content, notice.Content)
	assert.Equal(t, noticeType, notice.NoticeType)
}

func TestDeleteNotice(t *testing.T) {
	// t.Skip("Skip TestDeleteNotice")
	var (
		injector  = wire.GetInjector()
		ctx       = injector.Ctx
		noticeDao = injector.NoticeDao
		err       error
	)
	err = noticeDao.DeleteNotice(ctx, noticeID)
	assert.NoError(t, err)

	notice, err := noticeDao.GetNoticeByID(ctx, noticeID)
	assert.Error(t, err)
	assert.Nil(t, notice)
}

func TestDeleteNoticeList(t *testing.T) {
	// t.Skip("Skip TestDeleteNoticeList")
	var (
		injector  = wire.GetInjector()
		ctx       = injector.Ctx
		noticeDao = injector.NoticeDao
		err       error
	)
	noticeType := "NORMAL"
	count, err := noticeDao.DeleteNoticeList(ctx, nil, nil, nil, nil, &noticeType)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Notice Type: %s", noticeType)
	t.Logf("Delete Count: %d", *count)
	t.Logf("=====================================")

	noticeList, count, err := noticeDao.GetNoticeList(
		ctx, 0, 10, false, nil, nil, nil, nil, &noticeType,
	)
	assert.Empty(t, noticeList)
}
