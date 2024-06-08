package mongo

import (
	"context"
	"sync"

	"github.com/qiniu/qmgo"
)

type Mongo struct {
	MongoClient  *qmgo.Client
	DatabaseName string
	qmgoConfig   *qmgo.Config
	pingTimeout  int64
}

var (
	mongoInstance *Mongo
	once          sync.Once
)

func New(ctx context.Context, config *qmgo.Config, PingTimeout int64, databaseName string) (*Mongo, error) {
	var err error
	once.Do(
		func() {
			m := &Mongo{
				qmgoConfig:   config,
				pingTimeout:  PingTimeout,
				DatabaseName: databaseName,
			}
			if err = m.init(ctx); err != nil {
				return
			}
			mongoInstance = m
		},
	)
	return mongoInstance, err
}

func Update(ctx context.Context, config *qmgo.Config, PingTimeout int64, databaseName string) error {
	if err := mongoInstance.Close(ctx); err != nil {
		return err
	}
	*mongoInstance = Mongo{
		qmgoConfig:   config,
		pingTimeout:  PingTimeout,
		DatabaseName: databaseName,
	}
	return mongoInstance.init(ctx)
}

func (m *Mongo) init(ctx context.Context) error {
	client, err := qmgo.NewClient(ctx, m.qmgoConfig)
	if err != nil {
		return err
	}
	if err = client.Ping(m.pingTimeout); err != nil {
		return err
	}
	m.MongoClient = client
	return nil
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.MongoClient.Close(ctx)
}
