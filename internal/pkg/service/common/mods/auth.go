package mods

import (
	"context"
	e "errors"
	"fmt"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	daos "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/vo/common"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"fiber-admin/pkg/jwt"
	"fiber-admin/pkg/utils/crypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, email, password *string) (*common.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken *string) (*common.RefreshTokenResponse, error)
	Logout(ctx context.Context, accessToken *string) error
	ChangePassword(ctx context.Context, oldPassword, newPassword *string) error
}

type authServiceImpl struct {
	core    *service.Core
	cache   *dao.Cache
	userDao daos.UserDao
	jwt     *jwt.Jwt
}

func NewAuthService(core *service.Core, userDao daos.UserDao, cache *dao.Cache, jwt *jwt.Jwt) AuthService {
	return &authServiceImpl{
		core:    core,
		cache:   cache,
		userDao: userDao,
		jwt:     jwt,
	}
}

func (a authServiceImpl) Login(ctx context.Context, email, password *string) (*common.LoginResponse, error) {
	user, err := a.userDao.GetUserByEmail(ctx, *email)
	if err != nil {
		return nil, errors.AuthFailed(fmt.Errorf("user not exist or password wrong"))
	}
	if !crypt.Compare(*password, user.Password) {
		return nil, errors.AuthFailed(fmt.Errorf("user not exist or password wrong"))
	}
	accessToken, err := a.jwt.GenerateAccessToken(user.UserID.Hex())
	if err != nil {
		a.core.Logger.Error("failed to generate access token", zap.Error(err))
		return nil, errors.ServiceError(fmt.Errorf("failed to generate access token"))
	}
	refreshToken, err := a.jwt.GenerateRefreshToken(user.UserID.Hex())
	if err != nil {
		a.core.Logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, errors.ServiceError(fmt.Errorf("failed to generate refresh token"))
	}
	err = a.userDao.UpdateUserLastLogin(ctx, user.UserID)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.DuplicateKeyError(fmt.Errorf("user last login already updated"))
		} else if e.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.NotFound(fmt.Errorf("user (id: %s) not found", user.UserID.Hex()))
		} else {
			return nil, errors.OperationFailed(
				fmt.Errorf(
					"failed to update user last login (id: %s)", user.UserID.Hex(),
				),
			)
		}
	}
	return &common.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    a.core.Config.JWTConfig.TokenDuration.Seconds(),
		Meta: struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}(struct {
			UserID   string
			Username string
			Email    string
			Role     string
		}{
			UserID:   user.UserID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}),
	}, nil
}

func (a authServiceImpl) RefreshToken(ctx context.Context, refreshToken *string) (*common.RefreshTokenResponse, error) {
	userIDHex, err := a.jwt.VerifyRefreshToken(*refreshToken)
	if err != nil {
		return nil, errors.TokenInvalid(fmt.Errorf("refresh token invalid"))
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, errors.NotAuthorized(fmt.Errorf("user id invalid"))
	}
	user, err := a.userDao.GetUserByID(ctx, userID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return nil, errors.OperationFailed(fmt.Errorf("failed to get user (id: %s)", userID.Hex()))
		}
	}
	accessToken, err := a.jwt.GenerateAccessToken(userID.Hex())
	if err != nil {
		a.core.Logger.Error("failed to generate access token", zap.Error(err))
		return nil, errors.ServiceError(fmt.Errorf("failed to generate access token"))
	}
	newRefreshToken, err := a.jwt.GenerateRefreshToken(userID.Hex())
	if err != nil {
		a.core.Logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, errors.ServiceError(fmt.Errorf("failed to generate refresh token"))
	}
	return &common.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    a.core.Config.JWTConfig.TokenDuration.Seconds(),
		Meta: struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}(struct {
			UserID   string
			Username string
			Email    string
			Role     string
		}{
			UserID:   user.UserID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}),
	}, nil
}

func (a authServiceImpl) Logout(ctx context.Context, accessToken *string) error {
	if err := a.cache.Set(
		ctx, fmt.Sprintf("%s:%s", config.TokenBlacklistCachePrefix, crypt.MD5(*accessToken)), config.CacheTrue,
		&a.core.Config.CacheConfig.TokenBlacklistTTL,
	); err != nil {
		return errors.OperationFailed(fmt.Errorf("failed to blacklist token"))
	}
	return nil
}

func (a authServiceImpl) ChangePassword(ctx context.Context, oldPassword, newPassword *string) error {
	var (
		userIDHex string
		ok        bool
	)
	if userIDHex, ok = ctx.Value(config.UserIDKey).(string); !ok {
		return errors.NotAuthorized(fmt.Errorf("user id not found in context"))
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return errors.NotAuthorized(fmt.Errorf("user id invalid"))
	}
	user, err := a.userDao.GetUserByID(ctx, userID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to get user (id: %s)", userID.Hex()))
		}
	}
	if !crypt.Compare(*oldPassword, user.Password) {
		return errors.AuthFailed(fmt.Errorf("old password wrong"))
	}
	hashedPassword, err := crypt.Hash(*newPassword)
	if err != nil {
		a.core.Logger.Error("failed to hash password", zap.Error(err))
		return errors.ServiceError(fmt.Errorf("failed to hash password"))
	}
	if err = a.userDao.UpdateUser(ctx, userID, nil, nil, &hashedPassword, nil, nil); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.DuplicateKeyError(
				fmt.Errorf(
					"user with username %s or email %s already exists", user.Username, user.Email,
				),
			)
		} else if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to update user (id: %s)", userID.Hex()))
		}
	}
	return nil
}
