package mods

import (
	"fmt"
	"time"

	"fiber-admin/internal/pkg/domain/vo"
	"fiber-admin/internal/pkg/domain/vo/admin"
	adminservice "fiber-admin/internal/pkg/service/admin/mods"
	"fiber-admin/pkg/errors"
	"fiber-admin/pkg/utils/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LogsApi struct {
	LogsService adminservice.LogsService
	Validator   *validator.Validate
}

// GetLoginLogList returns the login log list.
//
//	@description	Get the login log list.
//	@id				admin-get-login-log-list
//	@summary		get login log list
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.GetLoginLogListRequest	query	admin.GetLoginLogListRequest	true	"Get login log list request"
//	@security		Bearer
//	@success		200						{object}	vo.Response{data=admin.GetLoginLogListResponse}	"Success"
//	@failure		400						{object}	vo.Response{data=nil}							"Invalid request"
//	@failure		401						{object}	vo.Response{data=nil}							"Unauthorized"
//	@failure		403						{object}	vo.Response{data=nil}							"Forbidden"
//	@failure		500						{object}	vo.Response{data=nil}							"Internal server error"
//	@router			/admin/login-log/list	[get]
func (l *LogsApi) GetLoginLogList(c *fiber.Ctx) error {
	req := new(admin.GetLoginLogListRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := l.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	var (
		createdBefore, createdAfter *time.Time
		err                         error
	)

	if req.CreateStartTime != nil {
		*createdBefore, err = time.Parse(time.RFC3339, *req.CreateStartTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create start time %s (should be in RFC3339 format)", *req.CreateStartTime,
				),
			)
		}
	}
	if req.CreateEndTime != nil {
		*createdAfter, err = time.Parse(time.RFC3339, *req.CreateEndTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create end time %s (should be in RFC3339 format)", *req.CreateEndTime,
				),
			)
		}
	}

	resp, err := l.LogsService.GetLoginLogList(
		c.UserContext(), req.Page, req.PageSize, req.Desc, req.Query, createdBefore, createdAfter,
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

// GetOperationLogList returns the operation log list.
//
//	@description	Get the operation log list.
//	@id				admin-get-operation-log-list
//	@summary		get operation log list
//	@tags			Admin API
//	@accept			json
//	@produce		json
//	@param			admin.GetOperationLogListRequest	query	admin.GetOperationLogListRequest	true	"Get operation log list request"
//	@security		Bearer
//	@success		200							{object}	vo.Response{data=admin.GetOperationLogListResponse}	"Success"
//	@failure		400							{object}	vo.Response{data=nil}								"Invalid request"
//	@failure		401							{object}	vo.Response{data=nil}								"Unauthorized"
//	@failure		403							{object}	vo.Response{data=nil}								"Forbidden"
//	@failure		500							{object}	vo.Response{data=nil}								"Internal server error"
//	@router			/admin/operation-log/list	[get]
func (l *LogsApi) GetOperationLogList(c *fiber.Ctx) error {
	req := new(admin.GetOperationLogListRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := l.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(common.FormatValidateError(errs))
	}

	var (
		createdBefore, createdAfter *time.Time
		err                         error
	)

	if req.CreateStartTime != nil {
		*createdBefore, err = time.Parse(time.RFC3339, *req.CreateStartTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create start time %s (should be in RFC3339 format)", *req.CreateStartTime,
				),
			)
		}
	}
	if req.CreateEndTime != nil {
		*createdAfter, err = time.Parse(time.RFC3339, *req.CreateEndTime)
		if err != nil {
			return errors.InvalidRequest(
				fmt.Errorf(
					"invalid create end time %s (should be in RFC3339 format)", *req.CreateEndTime,
				),
			)
		}
	}

	resp, err := l.LogsService.GetOperationLogList(
		c.UserContext(), req.Page, req.PageSize, req.Desc, req.Query, req.Operation, req.EntityType, req.Status,
		createdBefore, createdAfter,
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
