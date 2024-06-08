package mods

import (
	"context"
	"fmt"

	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogsService interface {
	InsertLoginLog(
		ctx context.Context, userID *primitive.ObjectID, ipAddress, userAgent *string,
	) error
	CacheLoginLog(ctx context.Context, userID *primitive.ObjectID, ipAddress, userAgent *string) error
	InsertOperationLog(
		ctx context.Context, userID, entityID *primitive.ObjectID,
		ipAddress, userAgent, operation, entityType, description, status *string,
	) error
	CacheOperationLog(
		ctx context.Context, userID, entityID *primitive.ObjectID, ipAddress, userAgent *string,
		operation, entityType, description, status *string,
	) error
}

type logsServiceImpl struct {
	core            *service.Core
	loginLogDao     mods.LoginLogDao
	operationLogDao mods.OperationLogDao
	userDao         mods.UserDao
}

func NewLogsService(
	core *service.Core, loginLogDao mods.LoginLogDao, operationLogDao mods.OperationLogDao,
) LogsService {
	return &logsServiceImpl{
		core:            core,
		loginLogDao:     loginLogDao,
		operationLogDao: operationLogDao,
	}
}

func (l logsServiceImpl) InsertLoginLog(
	ctx context.Context, userID *primitive.ObjectID, ipAddress, userAgent *string,
) error {
	_, err := l.loginLogDao.InsertLoginLog(ctx, *userID, *ipAddress, *userAgent)
	if err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to insert login log"))
	}
	return nil
}

func (l logsServiceImpl) CacheLoginLog(
	ctx context.Context, userID *primitive.ObjectID, ipAddress, userAgent *string,
) error {
	err := l.loginLogDao.CacheLoginLog(ctx, *userID, *ipAddress, *userAgent)
	if err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to cache login log"))
	}
	return nil
}

func (l logsServiceImpl) InsertOperationLog(
	ctx context.Context, userID, entityID *primitive.ObjectID,
	ipAddress, userAgent, operation, entityType, description, status *string,
) error {
	_, err := l.operationLogDao.InsertOperationLog(
		ctx, *userID, *entityID, *ipAddress, *userAgent, *operation, *entityType, *description,
		*status,
	)
	if err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to insert operation log"))
	}
	return nil
}

func (l logsServiceImpl) CacheOperationLog(
	ctx context.Context, userID, entityID *primitive.ObjectID, ipAddress, userAgent *string,
	operation, entityType, description, status *string,
) error {
	var _entityID primitive.ObjectID
	if entityID == nil {
		_entityID = primitive.NilObjectID
	}
	err := l.operationLogDao.CacheOperationLog(
		ctx, *userID, _entityID, *ipAddress, *userAgent, *operation, *entityType, *description, *status,
	)
	if err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to cache operation log"))
	}
	return nil
}
