package mods

import (
	"context"
	e "errors"
	"fmt"
	"time"

	dao "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/vo/common"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentationService interface {
	GetDocumentation(ctx context.Context, documentationID *primitive.ObjectID) (*common.GetDocumentationResponse, error)
	GetDocumentationList(
		ctx context.Context, page, pageSize *int64, updateStartTime, updateEndTime *time.Time,
	) (*common.GetDocumentationListResponse, error)
}

type documentationServiceImpl struct {
	core             *service.Core
	documentationDao dao.DocumentationDao
}

func NewDocumentationService(core *service.Core, documentationDao dao.DocumentationDao) DocumentationService {
	return &documentationServiceImpl{
		core:             core,
		documentationDao: documentationDao,
	}
}

func (d documentationServiceImpl) GetDocumentation(
	ctx context.Context, documentationID *primitive.ObjectID,
) (*common.GetDocumentationResponse, error) {
	documentation, err := d.documentationDao.GetDocumentationByID(ctx, *documentationID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.NotFound(fmt.Errorf("documentation (id: %s) not found", documentationID.Hex()))
		} else {
			return nil, errors.OperationFailed(
				fmt.Errorf(
					"failed to get documentation (id: %s)", documentationID.Hex(),
				),
			)
		}
	}
	return &common.GetDocumentationResponse{
		DocumentID: documentation.DocumentID.Hex(),
		Title:      documentation.Title,
		Content:    documentation.Content,
		CreatedAt:  documentation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  documentation.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (d documentationServiceImpl) GetDocumentationList(
	ctx context.Context, page, pageSize *int64, updateBefore, updateAfter *time.Time,
) (*common.GetDocumentationListResponse, error) {
	offset := (*page - 1) * *pageSize
	documentations, count, err := d.documentationDao.GetDocumentationList(
		ctx, offset, *pageSize, false, nil, nil, updateBefore, updateAfter,
	)
	if err != nil {
		return nil, errors.OperationFailed(fmt.Errorf("failed to get documentation list"))
	}
	resp := make([]*common.DocumentationSummary, 0, len(documentations))
	for _, documentation := range documentations {
		resp = append(
			resp, &common.DocumentationSummary{
				DocumentID: documentation.DocumentID.Hex(),
				Title:      documentation.Title,
				CreatedAt:  documentation.CreatedAt.Format(time.RFC3339),
			},
		)
	}
	return &common.GetDocumentationListResponse{
		Total:                    *count,
		DocumentationSummaryList: resp,
	}, nil
}
