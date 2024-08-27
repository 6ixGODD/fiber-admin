package mods

import (
	"fmt"

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

type NoticeApi struct {
	NoticeService adminservice.NoticeService
	LogsService   sysservice.LogsService
	Validator     *validator.Validate
}

// InsertNotice inserts a new notice.
//
//	@description	Insert a new notice.
//	@id				admin-insert-notice
//	@summary		insert notice
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.InsertNoticeRequest	body	admin.InsertNoticeRequest	true	"Insert notice request"
//	@security		Bearer
//	@success		200				{object}	vo.Response{data=nil}	"Success"
//	@failure		400				{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401				{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403				{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500				{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/notice	[post]
func (n *NoticeApi) InsertNotice(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.InsertNoticeRequest)

	if err := c.BodyParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := n.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	noticeIDHex, err := n.NoticeService.InsertNotice(ctx, req.Title, req.Content, req.NoticeType)
	var (
		userID, _  = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr     = c.IP()
		userAgent  = c.Get(fiber.HeaderUserAgent)
		operation  = config.OperationTypeCreate
		entityType = config.EntityTypeNotice
	)
	if err != nil {
		var (
			description = fmt.Sprintf("Insert notice failed: %s", err.Error())
			status      = config.OperationStatusFailure
		)
		_ = n.LogsService.CacheOperationLog(
			ctx, &userID, nil, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		noticeID, _ = primitive.ObjectIDFromHex(noticeIDHex)
		description = fmt.Sprintf("Insert notice: %s", noticeIDHex)
		status      = config.OperationStatusSuccess
	)
	_ = n.LogsService.CacheOperationLog(
		ctx, &userID, &noticeID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)
	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}

// UpdateNotice updates the notice.
//
//	@description	Update the notice.
//	@id				admin-update-notice
//	@summary		update notice
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.UpdateNoticeRequest	body	admin.UpdateNoticeRequest	true	"Update notice request"
//	@security		Bearer
//	@success		200				{object}	vo.Response{data=nil}	"Success"
//	@failure		400				{object}	vo.Response{data=nil}	"Invalid request"
//	@failure		401				{object}	vo.Response{data=nil}	"Unauthorized"
//	@failure		403				{object}	vo.Response{data=nil}	"Forbidden"
//	@failure		500				{object}	vo.Response{data=nil}	"Internal server error"
//	@router			/admin/notice																																																																																																																																																																																		[put]
func (n *NoticeApi) UpdateNotice(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.UpdateNoticeRequest)

	if err := c.BodyParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := n.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	noticeID, err := primitive.ObjectIDFromHex(*req.NoticeID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid notice id"))
	}
	err = n.NoticeService.UpdateNotice(ctx, &noticeID, req.Title, req.Content, req.NoticeType)
	var (
		userID, _  = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr     = c.IP()
		userAgent  = c.Get(fiber.HeaderUserAgent)
		operation  = config.OperationTypeUpdate
		entityType = config.EntityTypeNotice
	)
	if err != nil {
		var (
			description = fmt.Sprintf("Update notice failed: %s", err.Error())
			status      = config.OperationStatusFailure
		)
		_ = n.LogsService.CacheOperationLog(
			ctx, &userID, &noticeID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		description = fmt.Sprintf("Update notice: %s", *req.NoticeID)
		status      = config.OperationStatusSuccess
	)
	_ = n.LogsService.CacheOperationLog(
		ctx, &userID, &noticeID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)

	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}

// DeleteNotice deletes the notice.
//
//	@description	Delete the notice.
//	@id				admin-delete-notice
//	@summary		delete notice
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.DeleteNoticeRequest	query	admin.DeleteNoticeRequest	true	"Delete notice request"
//	@security		Bearer
//	@success		200	{object}	vo.Response{data=nil}	"Success"
//	@failure		400	{object}	vo.Response{data=nil}	"Invalid request"
func (n *NoticeApi) DeleteNotice(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(admin.DeleteNoticeRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := n.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	noticeID, err := primitive.ObjectIDFromHex(*req.NoticeID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid notice id"))
	}

	var (
		userID, _  = primitive.ObjectIDFromHex(ctx.Value(config.UserIDKey).(string))
		ipAddr     = c.IP()
		userAgent  = c.Get(fiber.HeaderUserAgent)
		operation  = config.OperationTypeDelete
		entityType = config.EntityTypeNotice
	)
	err = n.NoticeService.DeleteNotice(ctx, &noticeID)
	if err != nil {
		var (
			description = fmt.Sprintf("Delete notice failed: %s", err.Error())
			status      = config.OperationStatusFailure
		)
		_ = n.LogsService.CacheOperationLog(
			ctx, &userID, &noticeID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
		)
		return err
	}

	var (
		description = fmt.Sprintf("Delete notice: %s", *req.NoticeID)
		status      = config.OperationStatusSuccess
	)
	_ = n.LogsService.CacheOperationLog(
		ctx, &userID, &noticeID, &ipAddr, &userAgent, &operation, &entityType, &description, &status,
	)

	return c.JSON(
		vo.Response{
			Code:    errors.CodeSuccess,
			Message: errors.MessageSuccess,
			Data:    nil,
		},
	)
}
