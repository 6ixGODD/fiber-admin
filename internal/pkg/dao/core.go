package dao

import (
	"context"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/pkg/mongo"
	logging "fiber-admin/pkg/zap"
	"go.uber.org/zap"
)

type Core struct {
	Mongo  *mongo.Mongo
	Logger *zap.Logger
	Config *config.Config
}

func NewCore(ctx context.Context, mongo *mongo.Mongo, zap *logging.Zap, config *config.Config) (*Core, error) {
	c := zap.SetTagInContext(ctx, logging.MongoTag)
	logger, err := zap.GetLogger(c)
	if err != nil {
		return nil, err
	}
	return &Core{
		Mongo:  mongo,
		Logger: logger,
		Config: config,
	}, nil
}
