package mods

import (
	"fmt"
	"time"

	"fiber-admin/internal/pkg/domain/vo"
	"fiber-admin/internal/pkg/domain/vo/common"
	commonservice "fiber-admin/internal/pkg/service/common/mods"
	"fiber-admin/pkg/errors"
	utils "fiber-admin/pkg/utils/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoticeApi struct {
	NoticeService commonservice.NoticeService
	Validator     *validator.Validate
}

// GetNotice returns the notice by ID.
//
//	@description	Get the notice by ID.
//	@id				common-get-notice
//	@summary		get notice by ID
//	@tags			Notice API
//	@accept			json
//	@produce		json
//	@param			common.GetNoticeRequest	query	common.GetNoticeRequest	true	"Get notice request"
//	@security		Bearer
//	@success		200		{object}	vo.Response{data=common.GetNoticeResponse}	"Success"
//	@failure		400		{object}	vo.Response{data=nil}						"Invalid request"
//	@failure		401		{object}	vo.Response{data=nil}						"Unauthorized"
//	@failure		404		{object}	vo.Response{data=nil}						"Notice not found"
//	@failure		500		{object}	vo.Response{data=nil}						"Internal server error"
//	@router			/notice	[get]
func (n *NoticeApi) GetNotice(c *fiber.Ctx) error {
	req := new(common.GetNoticeRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := n.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(utils.FormatValidateError(errs))
	}

	noticeID, err := primitive.ObjectIDFromHex(*req.NoticeID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("invalid notice id"))
	}
	resp, err := n.NoticeService.GetNotice(c.UserContext(), &noticeID)
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

// GetNoticeList returns the notice list.
//
//	@description	Get the notice list.
//	@id				common-get-notice-list
//	@summary		get notice list
//	@tags			Notice API
//	@accept			json
//	@produce		json
//	@param			common.GetNoticeListRequest	query	common.GetNoticeListRequest	true	"Get notice list request"
//	@security		Bearer
//	@success		200				{object}	vo.Response{data=common.GetNoticeListResponse}	"Success"
//	@failure		400				{object}	vo.Response{data=nil}							"Invalid request"
//	@failure		401				{object}	vo.Response{data=nil}							"Unauthorized"
//	@failure		500				{object}	vo.Response{data=nil}							"Internal server error"
//	@router			/notice/list	[get]
func (n *NoticeApi) GetNoticeList(c *fiber.Ctx) error {
	req := new(common.GetNoticeListRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := n.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(utils.FormatValidateError(errs))
	}

	var (
		updateStartTime, updateEndTime       time.Time
		updateStartTimePtr, updateEndTimePtr *time.Time
		err                                  error
	)

	if req.UpdateStartTime != nil {
		updateStartTime, err = time.Parse(time.RFC3339, *req.UpdateStartTime)
		if err != nil {
			return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
		}
		updateStartTimePtr = &updateStartTime
	}
	if req.UpdateEndTime != nil {
		updateEndTime, err = time.Parse(time.RFC3339, *req.UpdateEndTime)
		if err != nil {
			return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
		}
		updateEndTimePtr = &updateEndTime
	}

	resp, err := n.NoticeService.GetNoticeList(
		c.UserContext(), req.Page, req.PageSize, req.NoticeType, updateStartTimePtr, updateEndTimePtr,
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
