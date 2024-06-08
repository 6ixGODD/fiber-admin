package mods

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/internal/pkg/domain/entity"
	"fiber-admin/pkg/utils/common"
	"github.com/goccy/go-json"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	opt "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UserDao interface {
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*entity.UserModel, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.UserModel, error)
	GetUserList(
		ctx context.Context,
		offset, limit int64, desc bool, organization, role *string,
		createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
		query *string,
	) ([]entity.UserModel, *int64, error)
	CountUser(
		ctx context.Context, organization, role *string,
		createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
	) (*int64, error)
	InsertUser(ctx context.Context, username, email, password, role, organization string) (primitive.ObjectID, error)
	UpdateUser(
		ctx context.Context, userID primitive.ObjectID, username, email, password, role, organization *string,
	) error
	UpdateUserLastLogin(ctx context.Context, userID primitive.ObjectID) error
	SoftDeleteUser(ctx context.Context, userID primitive.ObjectID) error
	SoftDeleteUserList(
		ctx context.Context, organization, role *string,
		createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
	) (*int64, error)
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	DeleteUserList(
		ctx context.Context, organization, role *string,
		createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
	) (*int64, error)
}

type UserDaoImpl struct {
	Core  *dao.Core
	Cache *dao.Cache
}

func NewUserDao(ctx context.Context, core *dao.Core, cache *dao.Cache) (UserDao, error) {
	var _ UserDao = (*UserDaoImpl)(nil)
	coll := core.Mongo.MongoClient.Database(core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.CreateIndexes(
		ctx, []options.IndexModel{
			{
				Key:          []string{"username"},
				IndexOptions: opt.Index().SetUnique(true),
			},
			{
				Key:          []string{"email"},
				IndexOptions: opt.Index().SetUnique(true),
			},
			{Key: []string{"created_at"}}, {Key: []string{"updated_at"}},
		},
	); err != nil {
		core.Logger.Error(
			fmt.Sprintf("Failed to create indexes for %s", config.UserCollectionName),
			zap.Error(err),
		)
		return nil, err
	}
	return &UserDaoImpl{core, cache}, nil
}

func (u *UserDaoImpl) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*entity.UserModel, error) {
	var user entity.UserModel
	key := fmt.Sprintf("%s:userID:%s", config.UserCachePrefix, userID.Hex())
	cache, err := u.Cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		u.Core.Logger.Info("UserDaoImpl.GetUserByID: cache miss", zap.String("key", key))
	} else if err != nil {
		u.Core.Logger.Error("UserDaoImpl.GetUserByID: cache get failed", zap.Error(err), zap.String("key", key))
	} else {
		err := json.Unmarshal([]byte(*cache), &user)
		if err != nil {
			u.Core.Logger.Error(
				"UserDaoImpl.GetUserByID: failed to unmarshal cache", zap.Error(err), zap.String("key", key),
			)
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByID: cache hit", zap.String("key", key))
			return &user, nil
		}
	}
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.Find(
		ctx, bson.M{"_id": userID, "deleted": false},
	).One(&user); err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserByID: failed to find user", zap.Error(err), zap.String("userID", userID.Hex()),
		)
		return nil, err
	} else {
		u.Core.Logger.Info("UserDaoImpl.GetUserByID: success", zap.String("userID", userID.Hex()))
		userJSON, _ := json.Marshal(user)
		if err := u.Cache.Set(ctx, key, string(userJSON), &u.Core.Config.CacheConfig.UserCacheTTL); err != nil {
			u.Core.Logger.Error("UserDaoImpl.GetUserByID: cache set failed", zap.Error(err), zap.String("key", key))
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByID: cache set success", zap.String("key", key))
		}
		return &user, nil
	}
}

func (u *UserDaoImpl) GetUserByEmail(ctx context.Context, email string) (*entity.UserModel, error) {
	var user entity.UserModel
	key := fmt.Sprintf("%s:email:%s", config.UserCachePrefix, email)
	cache, err := u.Cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		u.Core.Logger.Info("UserDaoImpl.GetUserByEmail: cache miss", zap.String("key", key))
	} else if err != nil {
		u.Core.Logger.Error("UserDaoImpl.GetUserByEmail: cache get failed", zap.Error(err), zap.String("key", key))
	} else {
		if err := json.Unmarshal([]byte(*cache), &user); err != nil {
			u.Core.Logger.Error(
				"UserDaoImpl.GetUserByEmail: failed to unmarshal cache", zap.Error(err), zap.String("key", key),
			)
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByEmail: cache hit", zap.String("key", key))
			return &user, nil
		}
	}
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.Find(
		ctx, bson.M{"email": email, "deleted": false},
	).One(&user); err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserByEmail: failed to find user", zap.Error(err), zap.String("email", email),
		)
		return nil, err
	} else {
		u.Core.Logger.Info(
			"UserDaoImpl.GetUserByEmail: success",
			zap.String("email", email),
		)
		userJSON, _ := json.Marshal(user)
		if err := u.Cache.Set(ctx, key, string(userJSON), &u.Core.Config.CacheConfig.UserCacheTTL); err != nil {
			u.Core.Logger.Error("UserDaoImpl.GetUserByEmail: cache set failed", zap.Error(err), zap.String("key", key))
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByEmail: cache set success", zap.String("key", key))
		}
		return &user, nil
	}
}

func (u *UserDaoImpl) GetUserByUsername(ctx context.Context, username string) (*entity.UserModel, error) {
	var user entity.UserModel
	key := fmt.Sprintf("%s:username:%s", config.UserCachePrefix, username)
	cache, err := u.Cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		u.Core.Logger.Info("UserDaoImpl.GetUserByUsername: cache miss", zap.String("key", key))
	} else if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserByUsername: cache get failed", zap.Error(err), zap.String("key", key),
		)
	} else {
		if err := json.Unmarshal([]byte(*cache), &user); err != nil {
			u.Core.Logger.Error(
				"UserDaoImpl.GetUserByUsername: failed to unmarshal cache", zap.Error(err), zap.String("key", key),
			)
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByUsername: cache hit", zap.String("key", key))
			return &user, nil
		}
	}
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.Find(
		ctx, bson.M{"username": username, "deleted": false},
	).One(&user); err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserByUsername: failed to find user", zap.Error(err), zap.String("username", username),
		)
		return nil, err
	} else {
		u.Core.Logger.Info("UserDaoImpl.GetUserByUsername: success", zap.String("username", username))
		userJSON, _ := json.Marshal(user)
		if err := u.Cache.Set(ctx, key, string(userJSON), &u.Core.Config.CacheConfig.UserCacheTTL); err != nil {
			u.Core.Logger.Error(
				"UserDaoImpl.GetUserByUsername: cache set failed", zap.Error(err), zap.String("key", key),
			)
		} else {
			u.Core.Logger.Info("UserDaoImpl.GetUserByUsername: cache set success", zap.String("key", key))
		}
		return &user, nil
	}
}

func (u *UserDaoImpl) GetUserList(
	ctx context.Context,
	offset, limit int64, desc bool, organization, role *string,
	createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
	query *string,
) ([]entity.UserModel, *int64, error) {
	var (
		userList []entity.UserModel
		err      error
	)
	doc := bson.M{"deleted": false}
	key := fmt.Sprintf("%s:offset:%d:limit:%d", config.UserCachePrefix, offset, limit) // for caching
	if organization != nil {
		doc["organization"] = *organization
		key += fmt.Sprintf(":organization:%s", *organization)
	}
	if role != nil {
		doc["role"] = *role
		key += fmt.Sprintf(":role:%s", *role)
	}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
		key += fmt.Sprintf(":createStartTime:%s:createEndTime:%s", createStartTime, createEndTime)
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
		key += fmt.Sprintf(":updateStartTime:%s:updateEndTime:%s", updateStartTime, updateEndTime)
	}
	if lastLoginStartTime != nil && lastLoginEndTime != nil {
		doc["last_login"] = bson.M{"$gte": lastLoginStartTime, "$lte": lastLoginEndTime}
		key += fmt.Sprintf(":lastLoginStartTime:%s:lastLoginEndTime:%s", lastLoginStartTime, lastLoginEndTime)
	}
	if query != nil {
		safetyQuery := common.EscapeSpecialChars(*query)
		pattern := fmt.Sprintf(".*%s.*", safetyQuery)
		doc["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"email": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"user_name": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
			{"organization": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}},
		}
	}
	docJSON, _ := json.Marshal(doc)

	if desc {
		key += ":desc"
	}
	var cache entity.UserCacheList
	err = u.Cache.GetList(ctx, key, &cache)
	if errors.Is(err, dao.CacheNil{}) {
		u.Core.Logger.Info("UserDaoImpl.GetUserList: cache miss", zap.String("key", key))
	} else if err != nil {
		u.Core.Logger.Error("UserDaoImpl.GetUserList: failed to get cache", zap.String("key", key), zap.Error(err))
	} else {
		u.Core.Logger.Info("UserDaoImpl.GetUserList: cache hit", zap.String("key", key))
		// convert cache to userList. cannot use cache directly because of primitive.ObjectID
		for _, user := range cache.List {
			userID, err := primitive.ObjectIDFromHex(user.UserID)
			if err != nil {
				u.Core.Logger.Error(
					"UserDaoImpl.GetUserList: failed to parse userID",
					zap.String("userID", user.UserID), zap.Error(err),
				)
				break
			}
			userList = append(
				userList, entity.UserModel{
					UserID:       userID,
					Username:     user.Username,
					Email:        user.Email,
					Password:     user.Password,
					Role:         user.Role,
					Organization: user.Organization,
					LastLogin:    user.LastLogin,
					Deleted:      user.Deleted,
					CreatedAt:    user.CreatedAt,
					UpdatedAt:    user.UpdatedAt,
					DeletedAt:    user.DeletedAt,
				},
			)
		}
		return userList, &cache.Total, nil
	}

	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	cursor := coll.Find(ctx, doc)
	count, err := cursor.Count()
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserList: failed to count userList",
			zap.ByteString(config.UserCollectionName, docJSON), zap.Error(err),
		)
		return nil, nil, err
	}
	if desc {
		err = cursor.Sort("-created_at").Skip(offset).Limit(limit).All(&userList)
	} else {
		err = cursor.Skip(offset).Limit(limit).All(&userList)
	}
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.GetUserList: failed to find userList",
			zap.ByteString(config.UserCollectionName, docJSON), zap.Error(err),
		)
		return nil, nil, err
	}
	u.Core.Logger.Info(
		"UserDaoImpl.GetUserList: success",
		zap.ByteString(config.UserCollectionName, docJSON), zap.Int64("count", count),
	)

	var userCacheList entity.UserCacheList
	userCacheList.Total = count
	for _, user := range userList {
		userCacheList.List = append(
			userCacheList.List, entity.UserCache{
				UserID:       user.UserID.Hex(),
				Username:     user.Username,
				Email:        user.Email,
				Password:     user.Password,
				Role:         user.Role,
				Organization: user.Organization,
				LastLogin:    user.LastLogin,
				Deleted:      user.Deleted,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
				DeletedAt:    user.DeletedAt,
			},
		)
	}
	if err := u.Cache.SetList(
		ctx, key, &userCacheList, &u.Core.Config.CacheConfig.UserCacheTTL,
	); err != nil {
		u.Core.Logger.Error(
			"NoticeDaoImpl.GetUserList: failed to set cache",
			zap.Error(err), zap.String("key", key), zap.ByteString(config.UserCollectionName, docJSON),
		)
	} else {
		u.Core.Logger.Info(
			"NoticeDaoImpl.GetUserList: cache set",
			zap.String("key", key), zap.ByteString(config.UserCollectionName, docJSON),
		)
	}
	return userList, &count, nil
}

func (u *UserDaoImpl) CountUser(
	ctx context.Context,
	organization, role *string,
	createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
) (*int64, error) {
	doc := bson.M{"deleted": false}
	key := fmt.Sprintf("%s:count", config.UserCachePrefix)
	if organization != nil {
		doc["organization"] = *organization
		key += fmt.Sprintf(":organization:%s", *organization)
	}
	if role != nil {
		doc["role"] = *role
		key += fmt.Sprintf(":role:%s", *role)
	}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
		key += fmt.Sprintf(":createStartTime:%s:createEndTime:%s", createStartTime, createEndTime)
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
		key += fmt.Sprintf(":updateStartTime:%s:updateEndTime:%s", updateStartTime, updateEndTime)
	}
	if lastLoginStartTime != nil && lastLoginEndTime != nil {
		doc["last_login"] = bson.M{"$gte": lastLoginStartTime, "$lte": lastLoginEndTime}
		key += fmt.Sprintf(":lastLoginStartTime:%s:lastLoginEndTime:%s", lastLoginStartTime, lastLoginEndTime)
	}
	docJSON, _ := json.Marshal(doc)
	cache, err := u.Cache.Get(ctx, key)
	if errors.Is(err, dao.CacheNil{}) {
		u.Core.Logger.Info("UserDaoImpl.CountUser: cache miss", zap.String("key", key))
	} else if err != nil {
		u.Core.Logger.Error("UserDaoImpl.CountUser: cache get failed", zap.Error(err), zap.String("key", key))
	} else {
		u.Core.Logger.Info("UserDaoImpl.CountUser: cache hit", zap.String("key", key))
		if count, err := strconv.ParseInt(*cache, 10, 64); err != nil {
			u.Core.Logger.Error(
				"UserDaoImpl.CountUser: failed to parse cache",
				zap.Error(err), zap.String("key", key), zap.String("value", *cache),
			)
		} else {
			return &count, nil
		}
	}
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	count, err := coll.Find(ctx, doc).Count()
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.CountUser: failed to count user", zap.Error(err),
			zap.ByteString(config.UserCollectionName, docJSON),
		)
	} else {
		u.Core.Logger.Info(
			"UserDaoImpl.CountUser: success", zap.Int64("count", count),
			zap.ByteString(config.UserCollectionName, docJSON),
		)
	}
	return &count, err
}

func (u *UserDaoImpl) InsertUser(
	ctx context.Context,
	username, email, password, role, organization string,
) (primitive.ObjectID, error) {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	doc := bson.M{
		"username":     username,
		"email":        email,
		"password":     password,
		"role":         role,
		"organization": organization,
		"last_login":   time.Time{},
		"deleted":      false,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
		"deleted_at":   nil,
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.InsertUser", zap.Error(err), zap.ByteString(config.UserCollectionName, docJSON),
		)
		return primitive.NilObjectID, err
	} else {
		u.Core.Logger.Info(
			"UserDaoImpl.InsertUser",
			zap.String("userID", result.InsertedID.(primitive.ObjectID).Hex()),
			zap.ByteString(config.UserCollectionName, docJSON),
		)
		prefix := config.UserCachePrefix
		if err = u.Cache.Flush(ctx, &prefix); err != nil {
			u.Core.Logger.Error("UserDaoImpl.InsertUser: failed to flush cache", zap.Error(err))
		} else {
			u.Core.Logger.Info("UserDaoImpl.InsertUser: cache flushed")
		}
		return result.InsertedID.(primitive.ObjectID), nil
	}
}

func (u *UserDaoImpl) UpdateUser(
	ctx context.Context, userID primitive.ObjectID, username, email, password, role, organization *string,
) error {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	doc := bson.M{"updated_at": time.Now()}
	if username != nil {
		doc["username"] = *username
	}
	if email != nil {
		doc["email"] = *email
	}
	if password != nil {
		doc["password"] = *password
	}
	if role != nil {
		doc["role"] = *role
	}
	if organization != nil {
		doc["organization"] = *organization
	}
	docJSON, _ := json.Marshal(doc)
	if err := coll.UpdateId(ctx, userID, bson.M{"$set": doc}); err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.UpdateUser: failed",
			zap.Error(err), zap.String("userID", userID.Hex()), zap.ByteString(config.UserCollectionName, docJSON),
		)
		return err
	}
	u.Core.Logger.Info(
		"UserDaoImpl.UpdateUser: success",
		zap.String("userID", userID.Hex()), zap.ByteString(config.UserCollectionName, docJSON),
	)
	prefix := config.UserCachePrefix
	if err := u.Cache.Flush(ctx, &prefix); err != nil {
		u.Core.Logger.Error("UserDaoImpl.UpdateUser: failed to flush cache", zap.Error(err))
	} else {
		u.Core.Logger.Info("UserDaoImpl.UpdateUser: cache flushed")
	}
	return nil
}

func (u *UserDaoImpl) UpdateUserLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	doc := bson.M{"last_login": time.Now()}
	docJSON, _ := json.Marshal(doc)
	if err := coll.UpdateId(ctx, userID, bson.M{"$set": doc}); err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.UpdateUserLastLogin: failed",
			zap.Error(err), zap.String("userID", userID.Hex()), zap.ByteString(config.UserCollectionName, docJSON),
		)
		return err
	}
	u.Core.Logger.Info(
		"UserDaoImpl.UpdateUserLastLogin: success",
		zap.String("userID", userID.Hex()), zap.ByteString(config.UserCollectionName, docJSON),
	)
	prefix := config.UserCachePrefix
	if err := u.Cache.Flush(ctx, &prefix); err != nil {
		u.Core.Logger.Error("UserDaoImpl.UpdateUserLastLogin: failed to flush cache", zap.Error(err))
	} else {
		u.Core.Logger.Info("UserDaoImpl.UpdateUserLastLogin: cache flushed")
	}
	return nil
}

func (u *UserDaoImpl) SoftDeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.UpdateId(
		ctx, userID, bson.M{"$set": bson.M{"deleted": true, "deleted_at": time.Now()}},
	); err != nil {
		u.Core.Logger.Error("UserDaoImpl.DeleteUser", zap.Error(err), zap.String("userID", userID.Hex()))
		return err
	}
	u.Core.Logger.Info("UserDaoImpl.DeleteUser", zap.String("userID", userID.Hex()))
	prefix := config.UserCachePrefix
	if err := u.Cache.Flush(ctx, &prefix); err != nil {
		u.Core.Logger.Error("UserDaoImpl.SoftDeleteUser: failed to flush cache", zap.Error(err))
	} else {
		u.Core.Logger.Info("UserDaoImpl.SoftDeleteUser: cache flushed")
	}
	return nil
}

func (u *UserDaoImpl) SoftDeleteUserList(
	ctx context.Context, organization, role *string,
	createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
) (*int64, error) {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	doc := bson.M{"deleted": false}
	if organization != nil {
		doc["organization"] = *organization
	}
	if role != nil {
		doc["role"] = *role
	}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
	}
	if lastLoginStartTime != nil && lastLoginEndTime != nil {
		doc["last_login"] = bson.M{"$gte": lastLoginStartTime, "$lte": lastLoginEndTime}
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.UpdateAll(ctx, doc, bson.M{"$set": bson.M{"deleted": true, "deleted_at": time.Now()}})
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.DeleteUserList: failed", zap.Error(err), zap.ByteString(config.UserCollectionName, docJSON),
		)
	} else {
		u.Core.Logger.Info(
			"UserDaoImpl.DeleteUserList: success",
			zap.Int64("count", result.ModifiedCount), zap.ByteString(config.UserCollectionName, docJSON),
		)
		prefix := config.UserCachePrefix
		if err = u.Cache.Flush(ctx, &prefix); err != nil {
			u.Core.Logger.Error("UserDaoImpl.SoftDeleteUserList: failed to flush cache", zap.Error(err))
		} else {
			u.Core.Logger.Info("UserDaoImpl.SoftDeleteUserList: cache flushed")
		}
	}
	return &result.ModifiedCount, err
}

func (u *UserDaoImpl) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	if err := coll.RemoveId(ctx, userID); err != nil {
		u.Core.Logger.Error("UserDaoImpl.DeleteUser: failed", zap.Error(err), zap.String("userID", userID.Hex()))
		return err
	}
	u.Core.Logger.Info("UserDaoImpl.DeleteUser: success", zap.String("userID", userID.Hex()))
	prefix := config.UserCachePrefix
	if err := u.Cache.Flush(ctx, &prefix); err != nil {
		u.Core.Logger.Error("UserDaoImpl.DeleteUser: failed to flush cache", zap.Error(err))
	} else {
		u.Core.Logger.Info("UserDaoImpl.DeleteUser: cache flushed")
	}
	return nil
}

func (u *UserDaoImpl) DeleteUserList(
	ctx context.Context, organization, role *string,
	createStartTime, createEndTime, updateStartTime, updateEndTime, lastLoginStartTime, lastLoginEndTime *time.Time,
) (*int64, error) {
	coll := u.Core.Mongo.MongoClient.Database(u.Core.Mongo.DatabaseName).Collection(config.UserCollectionName)
	doc := bson.M{}
	if organization != nil {
		doc["organization"] = *organization
	}
	if role != nil {
		doc["role"] = *role
	}
	if createStartTime != nil && createEndTime != nil {
		doc["created_at"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
	}
	if updateStartTime != nil && updateEndTime != nil {
		doc["updated_at"] = bson.M{"$gte": updateStartTime, "$lte": updateEndTime}
	}
	if lastLoginStartTime != nil && lastLoginEndTime != nil {
		doc["last_login"] = bson.M{"$gte": lastLoginStartTime, "$lte": lastLoginEndTime}
	}
	docJSON, _ := json.Marshal(doc)
	result, err := coll.RemoveAll(ctx, doc)
	if err != nil {
		u.Core.Logger.Error(
			"UserDaoImpl.DeleteUserList: failed", zap.Error(err), zap.ByteString(config.UserCollectionName, docJSON),
		)
	} else {
		u.Core.Logger.Info(
			"UserDaoImpl.DeleteUserList: success",
			zap.Int64("count", result.DeletedCount), zap.ByteString(config.UserCollectionName, docJSON),
		)
		prefix := config.UserCachePrefix
		if err = u.Cache.Flush(ctx, &prefix); err != nil {
			u.Core.Logger.Error("UserDaoImpl.DeleteUserList: failed to flush cache", zap.Error(err))
		} else {
			u.Core.Logger.Info("UserDaoImpl.DeleteUserList: cache flushed")
		}
	}
	return &result.DeletedCount, err
}
