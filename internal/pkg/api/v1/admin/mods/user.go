package mods

import (
	"fmt"
	"time"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/domain/vo"
	"fiber-admin/internal/pkg/domain/vo/admin"
	adminservice "fiber-admin/internal/pkg/service/admin/mods"
	sysservice "fiber-admin/internal/pkg/service/sys/mods"
	"fiber-admin/pkg/errors"
	"fiber-admin/pkg/utils/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserApi struct {
	UserService adminservice.UserService
	LogsService sysservice.LogsService
	Validator   *validator.Validate
}

//	@description	Insert a new user.
//	@id				admin-insert-user
//	@summary		insert user
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.InsertUserRequest	body	admin.InsertUserRequest	true	"Insert user request"
//	@security		Bearer
//	@success		200			{object}	vo.Response{data=nil}	"Success"
//	@failure		400			{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401			{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403			{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500			{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/user	[post]

// InsertUser inserts a new user.
func (u *UserApi) InsertUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.InsertUserRequest)

	if err := c.BodyParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	userIDHex, err := u.UserService.InsertUser(
		ctx, req.Username, req.Email, req.Password, req.Organization,
	)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)
	var (
		operatorID, _ = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr        = c.IP()
		userAgent     = c.Get(fiber.HeaderUserAgent)
		operation     = config.OperationTypeCreate
		entityType    = config.EntityTypeUser
	)
	if err != nil {
		var (
			description = fmt.Sprintf("Failed to insert user %s", *req.Username)
			status      = config.OperationStatusFailure
		)
		_ = u.LogsService.CacheOperationLog(
			ctx, &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		description = fmt.Sprintf("Insert user %s", *req.Username)
		status      = config.OperationStatusSuccess
	)
	_ = u.LogsService.CacheOperationLog(
		c.UserContext(), &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)
	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}

//	@description	Get the user by ID.
//	@id				admin-get-user
//	@summary		get user by ID
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.GetUserRequest	query	admin.GetUserRequest	true	"Get user request"
//	@security		Bearer
//	@success		200			{object}	vo.Response{data=admin.GetUserResponse}	"Success"
//	@failure		400			{object}	vo.Response{data=nil}					"Invalid request"
//	@failure		401			{object}	vo.Response{data=nil}					"Unauthorized"
//	@failure		403			{object}	vo.Response{data=nil}					"Forbidden"
//	@failure		404			{object}	vo.Response{data=nil}					"User not found"
//	@failure		500			{object}	vo.Response{data=nil}					"Internal server error"
//	@router			/admin/user	[get]

// GetUser returns the user by ID.
func (u *UserApi) GetUser(c *fiber.Ctx) error {
	req := new(admin.GetUserRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	userID, err := primitive.ObjectIDFromHex(*req.UserID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid user id"))
	}

	resp, err := u.UserService.GetUser(c.UserContext(), &userID)
	if err != nil {
		return err
	}

	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    resp,
		},
	)
}

// GetUserList returns the list of users based on the query parameters.
//
//	@description	Get the list of users based on the query parameters.
//	@id				admin-get-user-list
//	@summary		get user list
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.GetUserListRequest	query	admin.GetUserListRequest	true	"Get user list request"
//	@security		Bearer
//	@success		200					{object}	vo.Response{data=admin.GetUserListResponse}	"Success"
//	@failure		400					{object}	vo.Response{data=nil}						"Invalid request"
//	@failure		401					{object}	vo.Response{data=nil}						"Unauthorized"
//	@failure		403					{object}	vo.Response{data=nil}						"Forbidden"
//	@failure		500					{object}	vo.Response{data=nil}						"Internal server error"
//	@router			/admin/user/list																																																																																																																																																																																																																																																																																																				[get]
func (u *UserApi) GetUserList(c *fiber.Ctx) error {
	req := new(admin.GetUserListRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	var (
		lastLoginStartTime, lastLoginEndTime,
		createdStartTime, createdEndTime time.Time
		lastLoginStartTimePtr, lastLoginEndTimePtr,
		createdStartTimePtr, createdEndTimePtr *time.Time
		err error
	)

	if req.LastLoginStartTime != nil && req.LastLoginEndTime != nil {
		lastLoginStartTime, err = time.Parse(time.RFC3339, *req.LastLoginStartTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid last login time start %s (should be in RFC3339 format)", *req.LastLoginStartTime,
				),
			)
		}
		lastLoginEndTime, err = time.Parse(time.RFC3339, *req.LastLoginEndTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid last login time end %s (should be in RFC3339 format)", *req.LastLoginEndTime,
				),
			)
		}
		lastLoginStartTimePtr = &lastLoginStartTime
		lastLoginEndTimePtr = &lastLoginEndTime
	}
	if req.CreateStartTime != nil && req.CreateEndTime != nil {
		createdStartTime, err = time.Parse(time.RFC3339, *req.CreateStartTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create time start %s (should be in RFC3339 format)", *req.CreateStartTime,
				),
			)
		}
		createdEndTime, err = time.Parse(time.RFC3339, *req.CreateEndTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create time end %s (should be in RFC3339 format)", *req.CreateEndTime,
				),
			)
		}
		createdStartTimePtr = &createdStartTime
		createdEndTimePtr = &createdEndTime
	}

	resp, err := u.UserService.GetUserList(
		c.UserContext(), req.Page, req.PageSize, req.Desc, req.Role, lastLoginStartTimePtr, lastLoginEndTimePtr,
		createdStartTimePtr, createdEndTimePtr, req.Query,
	)
	if err != nil {
		return err
	}

	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    resp,
		},
	)
}

// UpdateUser updates the user.
//
//	@description	Update the user.
//	@id				admin-update-user
//	@summary		update user
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.UpdateUserRequest	body	admin.UpdateUserRequest	true	"Update user request"
//	@security		Bearer
//	@success		200			{object}	vo.Response{data=nil}	"Success"
//	@failure		400			{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401			{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403			{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500			{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/user																																																																																																																																																																																		[put]
func (u *UserApi) UpdateUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.UpdateUserRequest)

	if err := c.BodyParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	userID, err := primitive.ObjectIDFromHex(*req.UserID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid user id"))
	}

	err = u.UserService.UpdateUser(ctx, &userID, req.Username, req.Email, req.Organization)

	var (
		operatorID, _ = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr        = c.IP()
		userAgent     = c.Get(fiber.HeaderUserAgent)
		operation     = config.OperationTypeUpdate
		entityType    = config.EntityTypeUser
	)

	if err != nil {
		var (
			description = fmt.Sprintf("Failed to update user %s", *req.Username)
			status      = config.OperationStatusFailure
		)
		_ = u.LogsService.CacheOperationLog(
			ctx, &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		description = fmt.Sprintf("Update user %s", *req.Username)
		status      = config.OperationStatusSuccess
	)
	_ = u.LogsService.CacheOperationLog(
		ctx, &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)
	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}

// DeleteUser deletes a user by user ID.
//
//	@description	Delete the user by ID.
//	@id				admin-delete-user
//	@summary		delete user by ID
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.DeleteUserRequest	query	admin.DeleteUserRequest	true	"Delete user request"
//	@security		Bearer
//	@success		200			{object}	vo.Response{data=nil}	"Success"
//	@failure		400			{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401			{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403			{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500			{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/user																																																																																																																																																																																					[delete]
func (u *UserApi) DeleteUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.DeleteUserRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	userID, err := primitive.ObjectIDFromHex(*req.UserID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid user id"))
	}
	err = u.UserService.DeleteUser(ctx, &userID)

	var (
		operatorID, _ = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr        = c.IP()
		userAgent     = c.Get(fiber.HeaderUserAgent)
		operation     = config.OperationTypeDelete
		entityType    = config.EntityTypeUser
	)

	if err != nil {
		var (
			description = fmt.Sprintf("Failed to delete user %s", *req.UserID)
			status      = config.OperationStatusFailure
		)
		_ = u.LogsService.CacheOperationLog(
			ctx, &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		description = fmt.Sprintf("Delete user %s", *req.UserID)
		status      = config.OperationStatusSuccess
	)
	_ = u.LogsService.CacheOperationLog(
		ctx, &operatorID, &userID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)
	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}

// ChangeUserPassword changes a user's password.
//
//	@description	Change the user's password.
//	@id				admin-change-user-password
//	@summary		change user password
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.ChangeUserPasswordRequest	body	admin.ChangeUserPasswordRequest	true	"Change user password request"
//	@security		Bearer
//	@success		200						{object}	vo.Response{data=nil}	"Success"
//	@failure		400						{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401						{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403						{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500						{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/user/password	[put]
func (u *UserApi) ChangeUserPassword(c *fiber.Ctx) error {
	req := new(admin.ChangeUserPasswordRequest)

	if err := c.BodyParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := u.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	userID, err := primitive.ObjectIDFromHex(*req.UserID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid user id"))
	}
	err = u.UserService.ChangeUserPassword(c.UserContext(), &userID, req.NewPassword)
	if err != nil {
		return err
	}

	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}
