package mods

import (
	"context"
	e "errors"
	"fmt"
	"time"

	"fiber-admin/internal/pkg/config"
	dao "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/vo/admin"
	"fiber-admin/internal/pkg/service"
	"fiber-admin/pkg/errors"
	"fiber-admin/pkg/utils/crypt"
	"github.com/casbin/casbin/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserService interface {
	InsertUser(ctx context.Context, username, email, password, organization *string) (string, error)
	GetUser(ctx context.Context, userID *primitive.ObjectID) (*admin.GetUserResponse, error)
	GetUserList(
		ctx context.Context, page, pageSize *int64, desc *bool, role *string,
		lastLoginBefore, lastLoginAfter, createdBefore, createdAfter *time.Time, query *string,
	) (*admin.GetUserListResponse, error)
	UpdateUser(ctx context.Context, userID *primitive.ObjectID, username, email, organization *string) error
	DeleteUser(ctx context.Context, userID *primitive.ObjectID) error
	ChangeUserPassword(ctx context.Context, userID *primitive.ObjectID, newPassword *string) error
}

// UserServiceImpl implements the UserService.
type UserServiceImpl struct {
	core     *service.Core
	userDao  dao.UserDao
	enforcer *casbin.Enforcer
}

// NewUserService is a wire provider function that returns a UserServiceImpl.
func NewUserService(core *service.Core, userDao dao.UserDao, enforcer *casbin.Enforcer) UserService {
	return &UserServiceImpl{
		core:     core,
		userDao:  userDao,
		enforcer: enforcer,
	}
}

// InsertUser inserts a new user into mongodb and creates a new role for the user in casbin.
// Returns the user ID if successful.
func (u UserServiceImpl) InsertUser(
	ctx context.Context, username, email, password, organization *string,
) (string, error) {
	passwordHash, err := crypt.Hash(*password)
	if err != nil {
		u.core.Logger.Error("failed to hash password", zap.Error(err))
		return "", errors.ServiceError(fmt.Errorf("failed to hash password"))
	}
	userID, err := u.userDao.InsertUser(ctx, *username, *email, passwordHash, config.UserRoleUser, *organization)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errors.DuplicateKeyError(fmt.Errorf("user with email %s already exists", *email))
		} else {
			return "", errors.OperationFailed(fmt.Errorf("failed to insert user"))
		}
	}
	_, err = u.enforcer.AddRoleForUser(userID.Hex(), config.UserRoleUser)
	if err != nil {
		u.core.Logger.Error("failed to create role for user", zap.Error(err))
		return "", errors.ServiceError(fmt.Errorf("failed to create role for user"))
	}
	return userID.Hex(), nil
}

// GetUser retrieves a user by user ID.
// Returns the user if successful.
func (u UserServiceImpl) GetUser(ctx context.Context, userID *primitive.ObjectID) (*admin.GetUserResponse, error) {
	user, err := u.userDao.GetUserByID(ctx, *userID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return nil, errors.OperationFailed(fmt.Errorf("failed to get user (id: %s)", userID.Hex()))
		}
	}
	return &admin.GetUserResponse{
		UserID:       user.UserID.Hex(),
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		Organization: user.Organization,
		LastLogin:    user.LastLogin.Format(time.RFC3339),
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetUserList retrieves a list of users based on the query parameters.
// Returns the list of users if successful.
func (u UserServiceImpl) GetUserList(
	ctx context.Context, page, pageSize *int64, desc *bool, role *string,
	lastLoginBefore, lastLoginAfter, createdBefore, createdAfter *time.Time, query *string,
) (*admin.GetUserListResponse, error) {
	offset := (*page - 1) * *pageSize
	users, count, err := u.userDao.GetUserList(
		ctx, offset, *pageSize, *desc, nil, role, createdBefore, createdAfter,
		nil, nil, lastLoginBefore, lastLoginAfter, query,
	)
	if err != nil {
		return nil, errors.OperationFailed(fmt.Errorf("failed to get user list"))
	}
	resp := make([]*admin.GetUserResponse, 0, len(users))
	for _, user := range users {
		resp = append(
			resp, &admin.GetUserResponse{
				UserID:       user.UserID.Hex(),
				Username:     user.Username,
				Email:        user.Email,
				Role:         user.Role,
				Organization: user.Organization,
				LastLogin:    user.LastLogin.Format(time.RFC3339),
				CreatedAt:    user.CreatedAt.Format(time.RFC3339),
				UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
			},
		)
	}
	return &admin.GetUserListResponse{
		Total:    *count,
		UserList: resp,
	}, nil
}

// UpdateUser updates a user's information.
// Returns nil if successful.
func (u UserServiceImpl) UpdateUser(
	ctx context.Context, userID *primitive.ObjectID, username, email, organization *string,
) error {
	err := u.userDao.UpdateUser(ctx, *userID, username, email, nil, nil, organization)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.DuplicateKeyError(fmt.Errorf("user with email %s already exists", *email))
		} else if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to update user (id: %s)", userID.Hex()))
		}
	}
	return nil
}

// DeleteUser deletes a user by user ID.
// Returns nil if successful.
func (u UserServiceImpl) DeleteUser(ctx context.Context, userID *primitive.ObjectID) error {
	err := u.userDao.DeleteUser(ctx, *userID)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to delete user (id: %s)", userID.Hex()))
		}
	}
	return nil
}

// ChangeUserPassword changes a user's password.
// Returns nil if successful.
func (u UserServiceImpl) ChangeUserPassword(
	ctx context.Context, userID *primitive.ObjectID, newPassword *string,
) error {
	newPasswordHash, err := crypt.Hash(*newPassword)
	if err != nil {
		u.core.Logger.Error("failed to hash password", zap.Error(err))
		return errors.ServiceError(fmt.Errorf("failed to hash password"))
	}
	err = u.userDao.UpdateUser(ctx, *userID, nil, nil, &newPasswordHash, nil, nil)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFound(fmt.Errorf("user (id: %s) not found", userID.Hex()))
		} else {
			return errors.OperationFailed(fmt.Errorf("failed to update user (id: %s)", userID.Hex()))
		}
	}
	return nil
}
