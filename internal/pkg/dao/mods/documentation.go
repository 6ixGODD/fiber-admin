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

// DocumentationDao defines the crud methods that the infrastructure layer should implement
type DocumentationDao interface {
	GetDocumentationByID(ctx context.Context, documentationID primitive.ObjectID) (*entity.DocumentationModel, error)
	GetDocumentationList(
		ctx context.Context,
		offset, limit int64, desc bool, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
	) ([]entity.DocumentationModel, *int64, error)
	InsertDocumentation(ctx context.Context, title, content string) (primitive.ObjectID, error)
	UpdateDocumentation(ctx context.Context, documentationID primitive.ObjectID, title, content *string) error
	DeleteDocumentation(ctx context.Context, documentationID primitive.ObjectID) error
	DeleteDocumentationList(
		ctx context.Context, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
	) (*int64, error)
}

// DocumentationDaoImpl implements the DocumentationDao interface and contains a qmgo.Collection instance
type DocumentationDaoImpl struct {
	Dao   *dao.Core
	Cache *dao.Cache
}

// NewDocumentationDao creates a new instance of DocumentationDaoImpl with the qmgo.Collection instance
func NewDocumentationDao(ctx context.Context, core *dao.Core, cache *dao.Cache) (DocumentationDao, error) {
	var _ DocumentationDao = (*DocumentationDaoImpl)(nil) // Ensure that the interface is implemented
	coll := core.Mongo.MongoClient.Database(core.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	if err := coll.CreateIndexes(
		ctx, []options.IndexModel{
			{
				Key:          []string{"title"},
				IndexOptions: opt.Index().SetUnique(true),
			},
			{Key: []string{"created_at"}}, {Key: []string{"updated_at"}},
		},
	); err != nil {
		core.Logger.Error(
			fmt.Sprintf("Failed to create indexes for %s", config.DocumentationCollectionName), zap.Error(err),
		)
		return nil, err
	}
	return &DocumentationDaoImpl{core, cache}, nil
}

func (d *DocumentationDaoImpl) GetDocumentationByID(
	ctx context.Context, documentationID primitive.ObjectID,
) (*entity.DocumentationModel, error) {
	var documentation entity.DocumentationModel
	key := fmt.Sprintf("%s:documentationID:%s", config.DocumentationCachePrefix, documentationID.Hex())
	cache, err := d.Cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		d.Dao.Logger.Info("DocumentationDaoImpl.GetDocumentationByID: cache miss", zap.String("key", key))
	} else if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationByID: failed to get cache",
			zap.Error(err), zap.String("key", key),
		)
	} else {
		err = json.Unmarshal([]byte(*cache), &documentation)
		if err != nil {
			d.Dao.Logger.Error(
				"DocumentationDaoImpl.GetDocumentationByID: failed to unmarshal cache",
				zap.Error(err), zap.String("key", key),
			)
		} else {
			d.Dao.Logger.Info(
				"DocumentationDaoImpl.GetDocumentationByID: cache hit",
				zap.String("documentationID", documentationID.Hex()),
			)
			return &documentation, nil
		}
	}
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	if err = coll.Find(ctx, bson.M{"_id": documentationID}).One(&documentation); err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationByID: failed to find documentation",
			zap.Error(err), zap.String("documentationID", documentationID.Hex()),
		)
		return nil, err
	} else {
		docJSON, _ := json.Marshal(documentation)
		if err = d.Cache.Set(ctx, key, string(docJSON), &d.Dao.Config.CacheConfig.DocumentationCacheTTL); err != nil {
			d.Dao.Logger.Error(
				"DocumentationDaoImpl.GetDocumentationByID: failed to set cache",
				zap.Error(err), zap.String("key", key),
			)
		} else {
			d.Dao.Logger.Info(
				"DocumentationDaoImpl.GetDocumentationByID: cache set",
				zap.String("key", key), zap.ByteString(config.DocumentationCollectionName, docJSON),
			)
		}
		d.Dao.Logger.Info(
			"DocumentationDaoImpl.GetDocumentationByID: success",
			zap.String("documentationID", documentationID.Hex()),
		)
		return &documentation, nil
	}
}

func (d *DocumentationDaoImpl) GetDocumentationList(
	ctx context.Context,
	offset, limit int64, desc bool, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
) ([]entity.DocumentationModel, *int64, error) {
	var documentationList []entity.DocumentationModel
	var err error
	doc := bson.M{}
	key := fmt.Sprintf("%s:offset:%d:limit:%d", config.DocumentationCachePrefix, offset, limit)
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
		key += fmt.Sprintf(":createStartTime:%s:createEndTime:%s", createStartTime, createEndTime)
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
		key += fmt.Sprintf(":updateStartTime:%s:updateEndTime:%s", updateStartTime, updateEndTime)
	}
	docJSON, _ := json.Marshal(doc)

	if desc {
		key += ":desc"
	}
	var cache entity.DocumentationCacheList
	err = d.Cache.GetList(ctx, key, &cache)
	if errors.Is(err, dao.CacheNil{}) {
		d.Dao.Logger.Info("DocumentationDaoImpl.GetDocumentationList: cache miss", zap.String("key", key))
	} else if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationList: failed to get cache", zap.Error(err), zap.String("key", key),
		)
	} else {
		d.Dao.Logger.Info(
			"DocumentationDaoImpl.GetDocumentationList: cache hit", zap.String("key", key),
			zap.Int64("count", cache.Total),
		)
		for _, documentation := range cache.List {
			documentationID, err := primitive.ObjectIDFromHex(documentation.DocumentID)
			if err != nil {
				d.Dao.Logger.Error(
					"DocumentationDaoImpl.GetDocumentationList: failed to convert ObjectID",
					zap.String("documentationID", documentation.DocumentID), zap.Error(err),
				)
				break
			}
			documentationList = append(
				documentationList, entity.DocumentationModel{
					DocumentID: documentationID, Title: documentation.Title, Content: documentation.Content,
					CreatedAt: documentation.CreatedAt, UpdatedAt: documentation.UpdatedAt,
				},
			)
		}
		return documentationList, &cache.Total, nil
	}
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	if desc {
		err = coll.Find(ctx, doc).Sort("-created_at").Skip(offset).Limit(limit).All(&documentationList)
	} else {
		err = coll.Find(ctx, doc).Skip(offset).Limit(limit).All(&documentationList)
	}
	if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationList: failed to find documents",
			zap.ByteString(config.DocumentationCollectionName, docJSON), zap.Error(err),
		)
		return nil, nil, err
	}
	count, err := coll.Find(ctx, doc).Count()
	if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationList: failed to count documents",
			zap.ByteString(config.DocumentationCollectionName, docJSON), zap.Error(err),
		)
		return nil, nil, err
	}

	var documentationCacheList []entity.DocumentationCache
	for _, documentation := range documentationList {
		documentationCacheList = append(
			documentationCacheList, entity.DocumentationCache{
				DocumentID: documentation.DocumentID.Hex(), Title: documentation.Title, Content: documentation.Content,
				CreatedAt: documentation.CreatedAt, UpdatedAt: documentation.UpdatedAt,
			},
		)
	}
	if err := d.Cache.SetList(
		ctx, key, &entity.CacheList{List: documentationCacheList, Total: count},
		&d.Dao.Config.CacheConfig.DocumentationCacheTTL,
	); err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.GetDocumentationList: failed to set cache", zap.Error(err), zap.String("key", key),
		)
	} else {
		d.Dao.Logger.Info("DocumentationDaoImpl.GetDocumentationList: cache set", zap.String("key", key))
	}
	d.Dao.Logger.Info(
		"DocumentationDaoImpl.GetDocumentationList: success", zap.Int64("count", count),
		zap.ByteString(config.DocumentationCollectionName, docJSON),
	)
	return documentationList, &count, nil
}

func (d *DocumentationDaoImpl) InsertDocumentation(
	ctx context.Context, title, content string,
) (primitive.ObjectID, error) {
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	doc := bson.M{
		"title": title, "content": content, "created_at": time.Now(), "updated_at": time.Now(),
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.InsertDocumentation: failed to insert documentation",
			zap.Error(err), zap.ByteString(config.DocumentationCollectionName, docJSON),
		)
		return primitive.NilObjectID, err
	}
	d.Dao.Logger.Info(
		"DocumentationDaoImpl.InsertDocumentation: success",
		zap.String("documentation_id", result.InsertedID.(primitive.ObjectID).Hex()),
		zap.ByteString(config.DocumentationCollectionName, docJSON),
	)
	prefix := config.DocumentationCachePrefix
	err = d.Cache.Flush(ctx, &prefix)
	if err != nil {
		d.Dao.Logger.Error("DocumentationDaoImpl.InsertDocumentation: failed to flush cache", zap.Error(err))
	} else {
		d.Dao.Logger.Info("DocumentationDaoImpl.InsertDocumentation: cache flushed")
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (d *DocumentationDaoImpl) UpdateDocumentation(
	ctx context.Context, documentationID primitive.ObjectID, title, content *string,
) error {
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	doc := bson.M{"updated_at": time.Now()}
	if title != nil {
		doc["title"] = *title
	}
	if content != nil {
		doc["content"] = *content
	}
	docJSON, _ := json.Marshal(doc)
	err := coll.UpdateId(ctx, documentationID, bson.M{"$set": doc})
	if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.UpdateDocumentation: failed to update documentation",
			zap.Error(err), zap.String("documentationID", documentationID.Hex()),
			zap.ByteString(config.DocumentationCollectionName, docJSON),
		)
		return err
	}
	d.Dao.Logger.Info(
		"DocumentationDaoImpl.UpdateDocumentation: success",
		zap.String("documentationID", documentationID.Hex()),
		zap.ByteString(config.DocumentationCollectionName, docJSON),
	)
	prefix := config.DocumentationCachePrefix
	err = d.Cache.Flush(ctx, &prefix)
	if err != nil {
		d.Dao.Logger.Error("DocumentationDaoImpl.UpdateDocumentation: failed to flush cache", zap.Error(err))
	} else {
		d.Dao.Logger.Info("DocumentationDaoImpl.UpdateDocumentation: cache flushed")
	}
	return nil
}

func (d *DocumentationDaoImpl) DeleteDocumentation(
	ctx context.Context, documentationID primitive.ObjectID,
) error {
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	if err := coll.RemoveId(ctx, documentationID); err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.DeleteDocumentation: failed to delete documentation",
			zap.Error(err), zap.String("documentationID", documentationID.Hex()),
		)
		return err
	}
	d.Dao.Logger.Info(
		"DocumentationDaoImpl.DeleteDocumentation: success",
		zap.String("documentationID", documentationID.Hex()),
	)
	prefix := config.DocumentationCachePrefix
	if err := d.Cache.Flush(ctx, &prefix); err != nil {
		d.Dao.Logger.Error("DocumentationDaoImpl.DeleteDocumentation: failed to flush cache", zap.Error(err))
	} else {
		d.Dao.Logger.Info("DocumentationDaoImpl.DeleteDocumentation: cache flushed")
	}
	return nil
}

func (d *DocumentationDaoImpl) DeleteDocumentationList(
	ctx context.Context, createStartTime, createEndTime, updateStartTime, updateEndTime *time.Time,
) (*int64, error) {
	coll := d.Dao.Mongo.MongoClient.Database(d.Dao.Mongo.DatabaseName).Collection(config.DocumentationCollectionName)
	doc := bson.M{}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.RemoveAll(ctx, doc)
	if err != nil {
		d.Dao.Logger.Error(
			"DocumentationDaoImpl.DeleteDocumentationList: failed to delete documents",
			zap.Error(err), zap.ByteString(config.DocumentationCollectionName, docJSON),
		)
		return nil, err
	}
	d.Dao.Logger.Info(
		"DocumentationDaoImpl.DeleteDocumentationList: success",
		zap.Int64("count", result.DeletedCount), zap.ByteString(config.DocumentationCollectionName, docJSON),
	)
	prefix := config.DocumentationCachePrefix
	if err = d.Cache.Flush(ctx, &prefix); err != nil {
		d.Dao.Logger.Error("DocumentationDaoImpl.DeleteDocumentationList: failed to flush cache", zap.Error(err))
	} else {
		d.Dao.Logger.Info("DocumentationDaoImpl.DeleteDocumentationList: cache flushed")
	}
	return &result.DeletedCount, nil
}
