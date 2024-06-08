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

type DocumentationService interface {
	InsertDocumentation(ctx context.Context, title, content *string) (string, error)
	UpdateDocumentation(ctx context.Context, documentationID *primitive.ObjectID, title, content *string) error
	DeleteDocumentation(ctx context.Context, documentationID *primitive.ObjectID) error
}

type DocumentationServiceImpl struct {
	core             *service.Core
	documentationDao dao.DocumentationDao
}

func NewDocumentationService(core *service.Core, documentationDao dao.DocumentationDao) DocumentationService {
	return &DocumentationServiceImpl{
		core:             core,
		documentationDao: documentationDao,
	}
}

func (d DocumentationServiceImpl) InsertDocumentation(ctx context.Context, title, content *string) (string, error) {
	documentationID, err := d.documentationDao.InsertDocumentation(ctx, *title, *content)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errors.DuplicateKeyError(fmt.Errorf("documentation with title %s already exists", *title))
		} else {
			return "", errors.OperationFailed(fmt.Errorf("failed to insert documentation")) // TODO: Consider index error, duplicate key error, etc.
		}
	}
	return documentationID.Hex(), nil
}

func (d DocumentationServiceImpl) UpdateDocumentation(
	ctx context.Context, documentationID *primitive.ObjectID, title, content *string,
) error {
	err := d.documentationDao.UpdateDocumentation(ctx, *documentationID, title, content)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.DuplicateKeyError(fmt.Errorf("documentation with title %s already exists", *title))
		} else if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("documentation (id: %s) not found", documentationID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to update documentation (id: %s)", documentationID.Hex()))
		}
	}
	return nil
}

func (d DocumentationServiceImpl) DeleteDocumentation(ctx context.Context, documentationID *primitive.ObjectID) error {
	err := d.documentationDao.DeleteDocumentation(ctx, *documentationID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("documentation (id: %s) not found", documentationID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to delete documentation (id: %s)", documentationID.Hex()))
		}
	}
	return nil
}
