package mods

import (
	"context"
	"fmt"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"fiber-admin/pkg/utils/common"
	"go.uber.org/zap"
)

type IdempotencyService interface {
	GenerateIdempotencyToken(ctx context.Context) (string, error)
	CheckIdempotencyToken(ctx context.Context, token string) error
}

type idempotencyServiceImpl struct {
	core  *service.Core
	cache *dao.Cache
}

func NewIdempotencyService(core *service.Core, cache *dao.Cache) IdempotencyService {
	return &idempotencyServiceImpl{
		core:  core,
		cache: cache,
	}
}

func (s *idempotencyServiceImpl) GenerateIdempotencyToken(ctx context.Context) (string, error) {
	token, err := common.GenerateUUID4()
	if err != nil {
		s.core.Logger.Error("failed to generate idempotency token", zap.Error(err))
		return "", errors.ServiceError(fmt.Errorf("failed to generate idempotency token"))
	}
	key := fmt.Sprintf("%s:%s", config.IdempotencyCachePrefix, token)
	if err = s.cache.Set(ctx, key, config.CacheTrue, &s.core.Config.IdempotencyConfig.TTL); err != nil {
		s.core.Logger.Error("failed to set idempotency token", zap.Error(err))
		return "", errors.ServiceError(fmt.Errorf("failed to set idempotency token"))
	}
	return token, nil
}

func (s *idempotencyServiceImpl) CheckIdempotencyToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s:%s", config.IdempotencyCachePrefix, token)
	result, err := s.cache.Get(ctx, key)
	if err != nil {
		s.core.Logger.Error("failed to get idempotency token", zap.Error(err))
		return errors.Idempotency(fmt.Errorf("failed to get idempotency token"))
	}
	if err = s.cache.Delete(ctx, key); err != nil {
		s.core.Logger.Error("failed to delete idempotency token", zap.Error(err))
		return errors.Idempotency(fmt.Errorf("failed to delete idempotency token"))
	}
	if result == nil {
		return errors.Idempotency(fmt.Errorf("idempotency token not found"))
	}
	return nil
}
