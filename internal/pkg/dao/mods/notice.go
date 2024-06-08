package mods

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/internal/pkg/domain/entity"
	"github.com/goccy/go-json"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	opt "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type NoticeDao interface {
	GetNoticeByID(ctx context.Context, noticeID primitive.ObjectID) (*entity.NoticeModel, error)
	GetNoticeList(
		ctx context.Context,
		offset, limit int64, desc bool, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
		noticeType *string,
	) ([]entity.NoticeModel, *int64, error)
	InsertNotice(ctx context.Context, title, content, noticeType string) (primitive.ObjectID, error)
	UpdateNotice(ctx context.Context, noticeID primitive.ObjectID, title, content, noticeType *string) error
	DeleteNotice(ctx context.Context, noticeID primitive.ObjectID) error
	DeleteNoticeList(
		ctx context.Context,
		createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
		noticeType *string,
	) (*int64, error)
}

type NoticeDaoImpl struct {
	core  *dao.Core
	cache *dao.Cache
}

func NewNoticeDao(ctx context.Context, core *dao.Core, cache *dao.Cache) (NoticeDao, error) {
	var _ NoticeDao = (*NoticeDaoImpl)(nil)
	collection := core.Mongo.MongoClient.Database(core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	err := collection.CreateIndexes(
		ctx, []options.IndexModel{
			{
				Key:          []string{"title"},
				IndexOptions: opt.Index().SetUnique(true),
			},
			{Key: []string{"created_at"}}, {Key: []string{"updated_at"}},
		},
	)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("Failed to create indexes for %s", config.NoticeCollectionName), zap.Error(err))
		return nil, err
	}
	return &NoticeDaoImpl{core, cache}, nil
}

func (n *NoticeDaoImpl) GetNoticeByID(ctx context.Context, noticeID primitive.ObjectID) (*entity.NoticeModel, error) {
	var notice entity.NoticeModel
	key := fmt.Sprintf("%s:noticeID:%s", config.NoticeCachePrefix, noticeID.Hex())
	cache, err := n.cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		n.core.Logger.Info("NoticeDaoImpl.GetNoticeByID: cache miss", zap.String("noticeID", noticeID.Hex()))
	} else if err != nil {
		n.core.Logger.Error("NoticeDaoImpl.GetNoticeByID: failed to get cache", zap.Error(err), zap.String("key", key))
	} else {
		n.core.Logger.Info("NoticeDaoImpl.GetNoticeByID: cache hit", zap.String("noticeID", noticeID.Hex()))
		err = json.Unmarshal([]byte(*cache), &notice)
		if err != nil {
			n.core.Logger.Error(
				"NoticeDaoImpl.GetNoticeByID: failed to unmarshal cache", zap.Error(err),
				zap.String("noticeID", noticeID.Hex()), zap.String("cache", *cache),
			)
			return nil, err
		}
		return &notice, nil
	}
	coll := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	if err = coll.Find(ctx, bson.M{"_id": noticeID}).One(&notice); err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.GetNoticeByID: failed to find notice", zap.Error(err),
			zap.String("noticeID", noticeID.Hex()),
		)
		return nil, err
	} else {
		docJSON, _ := json.Marshal(notice)
		if err := n.cache.Set(ctx, key, string(docJSON), &n.core.Config.CacheConfig.NoticeCacheTTL); err != nil {
			n.core.Logger.Error(
				"NoticeDaoImpl.GetNoticeByID: failed to set cache", zap.Error(err),
				zap.String("key", key), zap.ByteString(config.NoticeCollectionName, docJSON),
			)
		} else {
			n.core.Logger.Info(
				"NoticeDaoImpl.GetNoticeByID: cache set", zap.String("key", key),
				zap.ByteString(config.NoticeCollectionName, docJSON),
			)
		}
		n.core.Logger.Info("NoticeDaoImpl.GetNoticeByID: success", zap.String("noticeID", noticeID.Hex()))
		return &notice, nil
	}
}

func (n *NoticeDaoImpl) GetNoticeList(
	ctx context.Context,
	offset, limit int64, desc bool, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
	noticeType *string,
) ([]entity.NoticeModel, *int64, error) {
	var noticeList []entity.NoticeModel
	var err error
	doc := bson.M{}
	key := fmt.Sprintf("%s:offset:%d:limit:%d", config.NoticeCachePrefix, offset, limit)
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
		key += fmt.Sprintf(
			":createStartTime:%s:createEndTime:%s", createStartTime.String(), createEndTime.String(),
		)
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
		key += fmt.Sprintf(
			":updateStartTime:%s:updateEndTime:%s", updateStartTime.String(), updateEndTime.String(),
		)
	}
	if noticeType != nil {
		doc["notice_type"] = *noticeType
		key += fmt.Sprintf(":noticeType:%s", *noticeType)
	}
	docJSON, _ := json.Marshal(doc)

	if desc {
		key += ":desc"
	}
	// cache, err := n.cache.GetList(ctx, key)
	var cache entity.NoticeCacheList
	err = n.cache.GetList(ctx, key, &cache)
	if errors.Is(err, dao.CacheNil{}) {
		n.core.Logger.Info("NoticeDaoImpl.GetNoticeList: cache miss", zap.String("key", key))
	} else if err != nil {
		n.core.Logger.Error("NoticeDaoImpl.GetNoticeList: failed to get cache", zap.Error(err), zap.String("key", key))
	} else {
		n.core.Logger.Info("NoticeDaoImpl.GetNoticeList: cache hit", zap.String("key", key))
		for _, noticeCache := range cache.List {
			noticeID, err := primitive.ObjectIDFromHex(noticeCache.NoticeID)
			if err != nil {
				n.core.Logger.Error(
					"NoticeDaoImpl.GetNoticeList: failed to parse ObjectID", zap.Error(err),
					zap.String("noticeID", noticeCache.NoticeID),
				)
				break
			}
			noticeList = append(
				noticeList, entity.NoticeModel{
					NoticeID:   noticeID,
					Title:      noticeCache.Title,
					Content:    noticeCache.Content,
					NoticeType: noticeCache.NoticeType,
					CreatedAt:  noticeCache.CreatedAt,
					UpdatedAt:  noticeCache.UpdatedAt,
				},
			)
		}
		return noticeList, &cache.Total, nil
	}

	collection := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	if desc {
		err = collection.Find(ctx, doc).Sort("-created_at").Skip(offset).Limit(limit).All(&noticeList)
	} else {
		err = collection.Find(ctx, doc).Skip(offset).Limit(limit).All(&noticeList)
	}
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.GetNoticeList: failed to find notices",
			zap.Error(err), zap.ByteString(config.NoticeCollectionName, docJSON),
		)
		return nil, nil, err
	}
	count, err := collection.Find(ctx, doc).Count()
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.GetNoticeList: count failed",
			zap.Error(err), zap.ByteString(config.NoticeCollectionName, docJSON),
		)
		return nil, nil, err
	}
	n.core.Logger.Info(
		"NoticeDaoImpl.GetNoticeList: success",
		zap.Int64("count", count), zap.ByteString(config.NoticeCollectionName, docJSON),
	)

	var noticeCacheList []entity.NoticeCache
	for _, notice := range noticeList {
		noticeCacheList = append(
			noticeCacheList, entity.NoticeCache{
				NoticeID:   notice.NoticeID.Hex(),
				Title:      notice.Title,
				Content:    notice.Content,
				NoticeType: notice.NoticeType,
				CreatedAt:  notice.CreatedAt,
				UpdatedAt:  notice.UpdatedAt,
			},
		)
	}
	if err := n.cache.SetList(
		ctx, key, &entity.CacheList{Total: count, List: noticeCacheList}, &n.core.Config.CacheConfig.NoticeCacheTTL,
	); err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.GetNoticeList: failed to set cache",
			zap.Error(err), zap.String("key", key), zap.ByteString(config.NoticeCollectionName, docJSON),
		)
	} else {
		n.core.Logger.Info(
			"NoticeDaoImpl.GetNoticeList: cache set", zap.String("key", key),
			zap.ByteString(config.NoticeCollectionName, docJSON),
		)
	}
	return noticeList, &count, nil
}

func (n *NoticeDaoImpl) InsertNotice(
	ctx context.Context, title, content, noticeType string,
) (primitive.ObjectID, error) {
	collection := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	doc := bson.M{
		"title":       title,
		"content":     content,
		"notice_type": noticeType,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.InsertNotice: failed to marshal notice", zap.Error(err),
			zap.String("title", title), zap.String("content", content), zap.String("noticeType", noticeType),
		)
		return primitive.NilObjectID, err
	}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.InsertNotice: failed to insert notice", zap.Error(err),
			zap.ByteString(config.NoticeCollectionName, docJSON),
		)
		return primitive.NilObjectID, err
	}
	n.core.Logger.Info(
		"NoticeDaoImpl.InsertNotice: success",
		zap.String("noticeID", result.InsertedID.(primitive.ObjectID).Hex()),
		zap.ByteString(config.NoticeCollectionName, docJSON),
	)
	prefix := config.NoticeCachePrefix
	if err = n.cache.Flush(ctx, &prefix); err != nil {
		n.core.Logger.Error("NoticeDaoImpl.InsertNotice: failed to flush cache", zap.Error(err))
	} else {
		n.core.Logger.Info("NoticeDaoImpl.InsertNotice: cache flush success")
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (n *NoticeDaoImpl) UpdateNotice(
	ctx context.Context, noticeID primitive.ObjectID, title, content, noticeType *string,
) error {
	collection := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	doc := bson.M{"updated_at": time.Now()}
	if title != nil {
		doc["title"] = *title
	}
	if content != nil {
		doc["content"] = *content
	}
	if noticeType != nil {
		doc["notice_type"] = *noticeType
	}
	docJSON, _ := json.Marshal(doc)
	err := collection.UpdateId(ctx, noticeID, bson.M{"$set": doc})

	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.UpdateNotice: failed to update notice",
			zap.Error(err), zap.String("noticeID", noticeID.Hex()),
			zap.ByteString(config.NoticeCollectionName, docJSON),
		)
	} else {
		n.core.Logger.Info(
			"NoticeDaoImpl.UpdateNotice: success",
			zap.String("noticeID", noticeID.Hex()), zap.ByteString(config.NoticeCollectionName, docJSON),
		)
		prefix := config.NoticeCachePrefix
		if err = n.cache.Flush(ctx, &prefix); err != nil {
			n.core.Logger.Error("NoticeDaoImpl.UpdateNotice: failed to flush cache", zap.Error(err))
		} else {
			n.core.Logger.Info("NoticeDaoImpl.UpdateNotice: cache flush success")
		}
	}
	return err
}

func (n *NoticeDaoImpl) DeleteNotice(ctx context.Context, noticeID primitive.ObjectID) error {
	collection := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	err := collection.RemoveId(ctx, noticeID)
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.DeleteNotice: failed to delete notice",
			zap.Error(err), zap.String("noticeID", noticeID.Hex()),
		)
	} else {
		prefix := config.NoticeCachePrefix
		n.core.Logger.Info("NoticeDaoImpl.DeleteNotice, success", zap.String("noticeID", noticeID.Hex()))
		if err = n.cache.Flush(ctx, &prefix); err != nil {
			n.core.Logger.Error("NoticeDaoImpl.DeleteNotice: failed to flush cache", zap.Error(err))
		} else {
			n.core.Logger.Info("NoticeDaoImpl.DeleteNotice: cache flush success")
		}
	}
	return err
}

func (n *NoticeDaoImpl) DeleteNoticeList(
	ctx context.Context,
	createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
	noticeType *string,
) (*int64, error) {
	collection := n.core.Mongo.MongoClient.Database(n.core.Mongo.DatabaseName).Collection(config.NoticeCollectionName)
	doc := bson.M{}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
	}
	if noticeType != nil {
		doc["notice_type"] = *noticeType
	}
	docJSON, _ := json.Marshal(doc)
	result, err := collection.RemoveAll(ctx, doc)
	if err != nil {
		n.core.Logger.Error(
			"NoticeDaoImpl.DeleteNoticeList: failed to delete notices",
			zap.Error(err), zap.ByteString(config.NoticeCollectionName, docJSON),
		)
	} else {
		n.core.Logger.Info(
			"NoticeDaoImpl.DeleteNoticeList: success", zap.ByteString(config.NoticeCollectionName, docJSON),
		)
		prefix := config.NoticeCachePrefix
		if err = n.cache.Flush(ctx, &prefix); err != nil {
			n.core.Logger.Error("NoticeDaoImpl.DeleteNoticeList: failed to flush cache", zap.Error(err))
		} else {
			n.core.Logger.Info("NoticeDaoImpl.DeleteNoticeList: cache flush success")
		}
	}
	return &result.DeletedCount, err
}
