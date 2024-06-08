package mods

import (
	"context"
	e "errors"
	"fmt"

	dao "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoticeService interface {
	InsertNotice(ctx context.Context, title, content, noticeType *string) (string, error)
	UpdateNotice(ctx context.Context, noticeID *primitive.ObjectID, title, content, noticeType *string) error
	DeleteNotice(ctx context.Context, noticeID *primitive.ObjectID) error
}

type NoticeServiceImpl struct {
	core      *service.Core
	noticeDao dao.NoticeDao
}

func NewNoticeService(core *service.Core, noticeDao dao.NoticeDao) NoticeService {
	return &NoticeServiceImpl{
		core:      core,
		noticeDao: noticeDao,
	}
}

func (n NoticeServiceImpl) InsertNotice(ctx context.Context, title, content, noticeType *string) (string, error) {
	noticeID, err := n.noticeDao.InsertNotice(ctx, *title, *content, *noticeType)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errors.DuplicateKeyError(fmt.Errorf("notice with title %s already exists", *title))
		} else {
			return "", errors.OperationFailed(fmt.Errorf("failed to insert notice"))
		}
	}
	return noticeID.Hex(), nil
}

func (n NoticeServiceImpl) UpdateNotice(
	ctx context.Context, noticeID *primitive.ObjectID, title, content, noticeType *string,
) error {
	err := n.noticeDao.UpdateNotice(ctx, *noticeID, title, content, noticeType)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.DuplicateKeyError(fmt.Errorf("notice with title %s already exists", *title))
		} else if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("notice (id: %s) not found", noticeID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to update notice (id: %s)", noticeID.Hex()))
		}
	}
	return nil
}

func (n NoticeServiceImpl) DeleteNotice(ctx context.Context, noticeID *primitive.ObjectID) error {
	err := n.noticeDao.DeleteNotice(ctx, *noticeID)
	if err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to delete notice (id: %s)", noticeID.Hex()))
	}
	return nil
}
