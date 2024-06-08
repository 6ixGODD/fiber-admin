package mock

import (
	"context"
	"math/rand"

	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoticeDaoMock struct {
	NoticeMap map[primitive.ObjectID]*entity.NoticeModel
	NoticeIDs []primitive.ObjectID
	NoticeDao mods.NoticeDao
}

func NewNoticeDaoMock(noticeDao mods.NoticeDao) *NoticeDaoMock {
	return &NoticeDaoMock{
		NoticeMap: make(map[primitive.ObjectID]*entity.NoticeModel),
		NoticeDao: noticeDao,
	}
}

func NewNoticeDaoMockWithRandomData(n int, noticeDao mods.NoticeDao) *NoticeDaoMock {
	noticeDaoMock := NewNoticeDaoMock(noticeDao)
	for i := 0; i < n; i++ {
		notice := noticeDaoMock.GenerateNoticeModel()
		noticeDaoMock.NoticeMap[notice.NoticeID] = notice
		noticeDaoMock.NoticeIDs = append(noticeDaoMock.NoticeIDs, notice.NoticeID)
	}
	return noticeDaoMock
}

func (m *NoticeDaoMock) Create(notice *entity.NoticeModel) error {
	m.NoticeMap[notice.NoticeID] = notice
	return nil
}

func (m *NoticeDaoMock) Get(noticeID primitive.ObjectID) (*entity.NoticeModel, error) {
	notice, ok := m.NoticeMap[noticeID]
	if !ok {
		return nil, nil
	}
	return notice, nil
}

func (m *NoticeDaoMock) RandomNoticeID() primitive.ObjectID {
	return m.NoticeIDs[rand.Intn(len(m.NoticeIDs))]
}

func (m *NoticeDaoMock) GenerateNoticeModel() *entity.NoticeModel {
	title, content, noticeType := GenerateNotice()
	noticeID, err := m.NoticeDao.InsertNotice(context.Background(), title, content, noticeType)
	if err != nil {
		panic(err)
	}

	notice, err := m.NoticeDao.GetNoticeByID(context.Background(), noticeID)
	if err != nil {
		panic(err)
	}

	return notice
}

func (m *NoticeDaoMock) Delete() {
	for _, noticeID := range m.NoticeIDs {
		_ = m.NoticeDao.DeleteNotice(context.Background(), noticeID)
	}
}

// GenerateNotice generates random notice data, and returns the title, content, and notice type
func GenerateNotice() (string, string, string) {
	return RandomString(10), RandomString(100), RandomEnum([]string{"URGENT", "NORMAL"})
}
