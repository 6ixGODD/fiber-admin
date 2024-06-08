package mods

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/internal/pkg/domain/entity"
	"fiber-admin/pkg/utils/common"
	"github.com/goccy/go-json"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type LoginLogDao interface {
	GetLoginLogByID(ctx context.Context, loginLogID primitive.ObjectID) (*entity.LoginLogModel, error)
	GetLoginLogList(
		ctx context.Context,
		offset, limit int64, desc bool, startTime, endTime *time.Time, userID *primitive.ObjectID,
		ipAddress, userAgent, query *string,
	) ([]entity.LoginLogModel, *int64, error)
	InsertLoginLog(
		ctx context.Context,
		UserID primitive.ObjectID, IPAddress, UserAgent string,
	) (primitive.ObjectID, error)
	CacheLoginLog(
		ctx context.Context, userID primitive.ObjectID, IPAddress, UserAgent string,
	) error
	SyncLoginLog(ctx context.Context)
	DeleteLoginLog(ctx context.Context, LoginLogID primitive.ObjectID) error
	DeleteLoginLogList(
		ctx context.Context, startTime, endTime *time.Time, userID *primitive.ObjectID,
		ipAddress, userAgent *string,
	) (*int64, error)
}

type LoginLogDaoImpl struct {
	core    *dao.Core
	cache   *dao.Cache
	userDao UserDao
}

func NewLoginLogDao(ctx context.Context, core *dao.Core, cache *dao.Cache, userDao UserDao) (LoginLogDao, error) {
	var _ LoginLogDao = (*LoginLogDaoImpl)(nil) // Ensure that the interface is implemented
	coll := core.Mongo.MongoClient.Database(core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	err := coll.CreateIndexes(
		ctx, []options.IndexModel{{Key: []string{"created_at"}}, {Key: []string{"user_id"}}},
	)
	if err != nil {
		core.Logger.Error(
			fmt.Sprintf("Failed to create index for %s", config.LoginLogCollectionName),
			zap.Error(err),
		)
		return nil, err
	}
	return &LoginLogDaoImpl{
		core:    core,
		userDao: userDao,
		cache:   cache,
	}, nil
}

func (l *LoginLogDaoImpl) GetLoginLogByID(
	ctx context.Context, loginLogID primitive.ObjectID,
) (*entity.LoginLogModel, error) {
	coll := l.core.Mongo.MongoClient.Database(l.core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	var loginLog entity.LoginLogModel
	err := coll.Find(ctx, bson.M{"_id": loginLogID}).One(&loginLog)
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.GetLoginLogByID: failed to find login log",
			zap.Error(err), zap.String("loginLogID", loginLogID.Hex()),
		)
		return nil, err
	} else {
		l.core.Logger.Info("LoginLogDaoImpl.GetLoginLogByID: success", zap.String("loginLogID", loginLogID.Hex()))
		return &loginLog, nil
	}
}

func (l *LoginLogDaoImpl) GetLoginLogList(
	ctx context.Context,
	offset, limit int64, desc bool, startTime, endTime *time.Time, userID *primitive.ObjectID,
	ipAddress, userAgent, query *string,
) ([]entity.LoginLogModel, *int64, error) {
	coll := l.core.Mongo.MongoClient.Database(l.core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	var loginLogList []entity.LoginLogModel
	var err error
	doc := bson.M{}
	if startTime != nil && endTime != nil {
		doc["created_at"] = bson.M{"$gte": startTime, "$lte": endTime}
	}
	if userID != nil {
		doc["user_id"] = *userID
	}
	if ipAddress != nil {
		doc["ip_address"] = *ipAddress
	}
	if userAgent != nil {
		doc["user_agent"] = *userAgent
	}
	if query != nil {
		safetyQuery := common.EscapeSpecialChars(*query)
		pattern := fmt.Sprintf(".*%s.*", safetyQuery)
		doc["$or"] = []bson.M{
			{"username": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"email": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"ip_address": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"user_agent": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
		}
	}
	docJSON, _ := json.Marshal(doc)
	if desc {
		err = coll.Find(ctx, doc).Sort("-created_at").Skip(offset).Limit(limit).All(&loginLogList)
	} else {
		err = coll.Find(ctx, doc).Skip(offset).Limit(limit).All(&loginLogList)
	}
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.GetLoginLogList: failed to find login logs",
			zap.Error(err), zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
		return nil, nil, err
	}
	count, err := coll.Find(ctx, doc).Count()
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.GetLoginLogList: failed to count login logs",
			zap.Error(err), zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
		return nil, nil, err
	}
	l.core.Logger.Info(
		"LoginLogDaoImpl.GetLoginLogList: success",
		zap.Int64("count", count), zap.ByteString(config.LoginLogCollectionName, docJSON),
	)
	return loginLogList, &count, nil
}

func (l *LoginLogDaoImpl) InsertLoginLog(
	ctx context.Context, userID primitive.ObjectID, ipAddress, userAgent string,
) (primitive.ObjectID, error) {
	coll := l.core.Mongo.MongoClient.Database(l.core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	user, err := l.userDao.GetUserByID(ctx, userID)
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.InsertLoginLog: failed to get user",
			zap.Error(err), zap.String("userID", userID.Hex()),
		)
		return primitive.NilObjectID, err
	}
	doc := bson.M{
		"user_id":    userID,
		"username":   user.Username,
		"email":      user.Email,
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"created_at": time.Now(),
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.InsertLoginLog: failed to insert login log",
			zap.Error(err), zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
	} else {
		l.core.Logger.Info(
			"LoginLogDaoImpl.InsertLoginLog: success",
			zap.String("loginLogID", result.InsertedID.(primitive.ObjectID).Hex()),
			zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
	}
	return result.InsertedID.(primitive.ObjectID), err
}

// CacheLoginLog caches login logs in cache
func (l *LoginLogDaoImpl) CacheLoginLog(
	ctx context.Context, userID primitive.ObjectID, IPAddress, UserAgent string,
) error {
	loginLog := entity.LoginLogCache{
		UserIDHex: userID.Hex(),
		IPAddress: IPAddress,
		UserAgent: UserAgent,
		CreatedAt: time.Now(),
	}
	loginLogJSON, err := json.Marshal(loginLog)
	if err != nil {
		l.core.Logger.Error("LoginLogDaoImpl.CacheLoginLog: failed to marshal login log", zap.Error(err))
		return err
	}
	return l.cache.RightPush(ctx, config.LoginLogCacheKey, string(loginLogJSON))
}

// SyncLoginLog syncs login logs from cache to database
func (l *LoginLogDaoImpl) SyncLoginLog(ctx context.Context) {
	for {
		loginLogJSON, err := l.cache.LeftPop(ctx, config.LoginLogCacheKey)
		if err != nil {
			if errors.Is(err, dao.CacheNil{}) {
				break
			}
			l.core.Logger.Error(
				"LoginLogDaoImpl.SyncLoginLog: failed to pop login log from cache", zap.Error(err),
			)
			return
		}
		var loginLog entity.LoginLogCache
		if err := json.Unmarshal([]byte(*loginLogJSON), &loginLog); err != nil {
			l.core.Logger.Error(
				"LoginLogDaoImpl.SyncLoginLog: failed to unmarshal login log",
				zap.Error(err), zap.String("loginLogJSON", *loginLogJSON),
			)
			continue
		}
		userID, err := primitive.ObjectIDFromHex(loginLog.UserIDHex)
		if err != nil {
			l.core.Logger.Error(
				"LoginLogDaoImpl.SyncLoginLog: failed to convert user ID",
				zap.Error(err), zap.String("loginLogJSON", *loginLogJSON),
			)
			continue
		}
		if _, err := l.InsertLoginLog(
			ctx, userID, loginLog.IPAddress, loginLog.UserAgent,
		); err != nil {
			l.core.Logger.Error(
				"LoginLogDaoImpl.SyncLoginLog: failed to insert login log",
				zap.Error(err), zap.String("loginLogJSON", *loginLogJSON),
			)
		}
	}
}

func (l *LoginLogDaoImpl) DeleteLoginLog(ctx context.Context, loginLogID primitive.ObjectID) error {
	coll := l.core.Mongo.MongoClient.Database(l.core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	err := coll.RemoveId(ctx, loginLogID)
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.DeleteLoginLog: failed to delete login log",
			zap.Error(err), zap.String("loginLogID", loginLogID.Hex()),
		)
	} else {
		l.core.Logger.Info("LoginLogDaoImpl.DeleteLoginLog: success", zap.String("loginLogID", loginLogID.Hex()))
	}
	return err
}

func (l *LoginLogDaoImpl) DeleteLoginLogList(
	ctx context.Context, startTime, endTime *time.Time, userID *primitive.ObjectID,
	ipAddress, userAgent *string,
) (*int64, error) {
	coll := l.core.Mongo.MongoClient.Database(l.core.Mongo.DatabaseName).Collection(config.LoginLogCollectionName)
	doc := bson.M{}
	if startTime != nil && endTime != nil {
		doc["created_at"] = bson.M{"$gte": startTime, "$lte": endTime}
	}
	if userID != nil {
		doc["user_id"] = *userID
	}
	if ipAddress != nil {
		doc["ip_address"] = *ipAddress
	}
	if userAgent != nil {
		doc["user_agent"] = *userAgent
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.RemoveAll(ctx, doc)
	if err != nil {
		l.core.Logger.Error(
			"LoginLogDaoImpl.DeleteLoginLogList: failed to delete login logs",
			zap.Error(err), zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
	} else {
		l.core.Logger.Info(
			"LoginLogDaoImpl.DeleteLoginLogList: success",
			zap.Int64("count", result.DeletedCount), zap.ByteString(config.LoginLogCollectionName, docJSON),
		)
	}
	return &result.DeletedCount, err
}
