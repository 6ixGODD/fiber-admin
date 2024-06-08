package tasks

import (
	"context"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/pkg/cron"
	"fiber-admin/pkg/jwt"
	logging "fiber-admin/pkg/zap"
	"go.uber.org/zap"
)

type Tasks struct {
	cron            *cron.Cron
	config          *config.Config
	loginLogDao     mods.LoginLogDao
	operationLogDao mods.OperationLogDao
	jwt             *jwt.Jwt
	logger          *zap.Logger
}

func New(
	ctx context.Context, config *config.Config, loginLogDao mods.LoginLogDao, operationLogDao mods.OperationLogDao,
	jwt *jwt.Jwt, zap *logging.Zap,
) (*Tasks, error) {
	ctx = zap.SetTagInContext(ctx, logging.CronTag)
	logger, err := zap.GetLogger(ctx)
	if err != nil {
		return nil, err
	}
	return &Tasks{
		cron:            cron.New(ctx),
		config:          config,
		loginLogDao:     loginLogDao,
		operationLogDao: operationLogDao,
		jwt:             jwt,
		logger:          logger,
	}, nil
}

func (t *Tasks) syncLogs() {
	t.logger.Info("Syncing logs from cache to database")
	t.loginLogDao.SyncLoginLog(t.cron.Context())
	t.operationLogDao.SyncOperationLog(t.cron.Context())
}

func (t *Tasks) updateKey() {
	t.logger.Info("Updating JWT key")
	if err := t.jwt.UpdateKey(); err != nil {
		t.logger.Error("Failed to update JWT key", zap.Error(err))
	}
}

func (t *Tasks) Start() error {
	syncLogsID, err := t.cron.AddFunc(t.config.TasksConfig.SyncLogsSpec, t.syncLogs)
	if err != nil {
		t.logger.Error("Failed to add sync logs task", zap.Error(err))
		return err
	}
	t.logger.Info("Added sync logs task", zap.Int("id", int(syncLogsID)))
	updateKeyID, err := t.cron.AddFunc(t.config.TasksConfig.UpdateKeySpec, t.updateKey)
	if err != nil {
		return err
	}
	t.logger.Info("Added update key task", zap.Int("id", int(updateKeyID)))
	t.logger.Info("Starting tasks")
	t.cron.Start()
	return nil
}

func (t *Tasks) Stop() {
	t.logger.Info("Stopping tasks")
	t.cron.Stop()
}
