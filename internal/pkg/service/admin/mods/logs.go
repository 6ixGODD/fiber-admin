package mods

import (
	"context"
	"fmt"
	"time"

	dao "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/vo/admin"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
)

type LogsService interface {
	GetLoginLogList(
		ctx context.Context, page, pageSize *int64, desc *bool, query *string,
		createStartTime, createEndTime *time.Time,
	) (*admin.GetLoginLogListResponse, error)
	GetOperationLogList(
		ctx context.Context, page, pageSize *int64, desc *bool, query, operation, entityType, status *string,
		createStartTime, createEndTime *time.Time,
	) (*admin.GetOperationLogListResponse, error)
}

type LogsServiceImpl struct {
	core            *service.Core
	loginLogDao     dao.LoginLogDao
	operationLogDao dao.OperationLogDao
}

func NewLogsService(
	core *service.Core, loginLogDao dao.LoginLogDao, operationLogDao dao.OperationLogDao,
) LogsService {
	return &LogsServiceImpl{
		core:            core,
		loginLogDao:     loginLogDao,
		operationLogDao: operationLogDao,
	}
}

func (l LogsServiceImpl) GetLoginLogList(
	ctx context.Context, page, pageSize *int64, desc *bool, query *string, createStartTime, createEndTime *time.Time,
) (*admin.GetLoginLogListResponse, error) {
	l.loginLogDao.SyncLoginLog(ctx) // Sync log before querying
	offset := (*page - 1) * *pageSize
	loginLogs, total, err := l.loginLogDao.GetLoginLogList(
		ctx, offset, *pageSize, *desc, createStartTime, createEndTime, nil, nil, nil, query,
	)
	if err != nil {
		return nil, errors.OperationFailed(fmt.Errorf("failed to get login log list"))
	}
	loginLogList := make([]*admin.GetLoginLogResponse, 0, len(loginLogs))
	for _, loginLog := range loginLogs {
		loginLogList = append(
			loginLogList, &admin.GetLoginLogResponse{
				LoginLogID: loginLog.LoginLogID.Hex(),
				UserID:     loginLog.UserID.Hex(),
				Username:   loginLog.Username,
				Email:      loginLog.Email,
				IPAddress:  loginLog.IPAddress,
				UserAgent:  loginLog.UserAgent,
				CreatedAt:  loginLog.CreatedAt.Format(time.RFC3339),
			},
		)
	}
	return &admin.GetLoginLogListResponse{
		Total:        *total,
		LoginLogList: loginLogList,
	}, nil
}

func (l LogsServiceImpl) GetOperationLogList(
	ctx context.Context, page, pageSize *int64, desc *bool, query, operation, entityType, status *string,
	createStartTime, createEndTime *time.Time,
) (*admin.GetOperationLogListResponse, error) {
	l.operationLogDao.SyncOperationLog(ctx) // Sync log before querying
	offset := (*page - 1) * *pageSize
	operationLogs, total, err := l.operationLogDao.GetOperationLogList(
		ctx, offset, *pageSize, *desc, createStartTime, createEndTime, nil, nil,
		nil, operation, entityType, status, query,
	)
	if err != nil {
		return nil, errors.OperationFailed(fmt.Errorf("failed to get operation log list"))
	}
	operationLogList := make([]*admin.GetOperationLogResponse, 0, len(operationLogs))
	for _, operationLog := range operationLogs {
		operationLogList = append(
			operationLogList, &admin.GetOperationLogResponse{
				OperationLogID: operationLog.OperationLogID.Hex(),
				UserID:         operationLog.UserID.Hex(),
				Username:       operationLog.Username,
				Email:          operationLog.Email,
				IPAddress:      operationLog.IPAddress,
				UserAgent:      operationLog.UserAgent,
				Operation:      operationLog.Operation,
				EntityID:       operationLog.EntityID.Hex(),
				EntityType:     operationLog.EntityType,
				Description:    operationLog.Description,
				Status:         operationLog.Status,
				CreatedAt:      operationLog.CreatedAt.Format(time.RFC3339),
			},
		)
	}
	return &admin.GetOperationLogListResponse{
		Total:            *total,
		OperationLogList: operationLogList,
	}, nil
}
