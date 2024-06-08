package mods

import (
	"context"
	"fmt"
	"time"

	"fiber-admin/internal/pkg/config"
	dao "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/vo/common"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (*common.GetProfileResponse, error)
}

type profileServiceImpl struct {
	core    *service.Core
	userDao dao.UserDao
}

func NewProfileService(core *service.Core, userDao dao.UserDao) ProfileService {
	return &profileServiceImpl{
		core:    core,
		userDao: userDao,
	}
}

func (p profileServiceImpl) GetProfile(ctx context.Context) (*common.GetProfileResponse, error) {
	userIDHex, ok := ctx.Value(config.UserIDKey).(string)
	if !ok {
		return nil, errors.NotAuthorized(fmt.Errorf("user is not authorized"))
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, errors.NotAuthorized(fmt.Errorf("user is not authorized"))
	}
	user, err := p.userDao.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.NotAuthorized(fmt.Errorf("user not exist"))
	}
	return &common.GetProfileResponse{
		UserID:       user.UserID.Hex(),
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		Organization: user.Organization,
		LastLogin:    user.LastLogin.Format(time.RFC3339),
	}, nil
}
