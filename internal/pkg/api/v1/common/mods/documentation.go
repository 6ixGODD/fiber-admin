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

type DocumentationApi struct {
	DocumentationService commonservice.DocumentationService
	Validator            *validator.Validate
}

//	@description	Get the documentation by ID.
//	@id				common-get-documentation
//	@summary		get documentation by ID
//	@tags			Documentation API
//	@accept			json
//	@produce		json
//	@param			common.GetDocumentationRequest	query	common.GetDocumentationRequest	true	"Get documentation request"
//	@security		Bearer
//	@success		200				{object}	vo.Response{data=common.GetDocumentationResponse}	"Success"
//	@failure		400				{object}	vo.Response{data=nil}								"Invalid request"
//	@failure		401				{object}	vo.Response{data=nil}								"Unauthorized"
//	@failure		404				{object}	vo.Response{data=nil}								"Documentation not found"
//	@failure		500				{object}	vo.Response{data=nil}								"Internal server error"
//	@router			/documentation	[get]

// GetDocumentation returns the documentation by ID.
func (d DocumentationApi) GetDocumentation(c *fiber.Ctx) error {
	req := new(common.GetDocumentationRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := d.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(utils.FormatValidateError(errs))
	}

	documentationID, err := primitive.ObjectIDFromHex(*req.DocumentationID)
	if err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	resp, err := d.DocumentationService.GetDocumentation(c.UserContext(), &documentationID)
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

//	@description	Get a list of documentation.
//	@id				common-get-documentation-list
//	@summary		get documentation list
//	@tags			Documentation API
//	@accept			json
//	@produce		json
//	@param			common.GetDocumentationListRequest	query	common.GetDocumentationListRequest	true	"Get documentation list request"
//	@security		Bearer
//	@success		200					{object}	vo.Response{data=common.GetDocumentationListResponse}	"Success"
//	@failure		400					{object}	vo.Response{data=nil}									"Invalid request"
//	@failure		401					{object}	vo.Response{data=nil}									"Unauthorized"
//	@failure		500					{object}	vo.Response{data=nil}									"Internal server error"
//	@router			/documentation/list	[get]

// GetDocumentationList returns a list of documentation.
func (d DocumentationApi) GetDocumentationList(c *fiber.Ctx) error {
	req := new(common.GetDocumentationListRequest)

	if err := c.QueryParser(req); err != nil {
		return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
	}
	if errs := d.Validator.Struct(req); errs != nil {
		return errors.InvalidRequest(utils.FormatValidateError(errs))
	}

	var (
		updateBefore, updateAfter       time.Time
		updateBeforePtr, updateAfterPtr *time.Time
		err                             error
	)

	if req.UpdateStartTime != nil && req.UpdateEndTime != nil {
		updateBefore, err = time.Parse(time.RFC3339, *req.UpdateStartTime)
		if err != nil {
			return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
		}
		updateAfter, err = time.Parse(time.RFC3339, *req.UpdateEndTime)
		if err != nil {
			return errors.InvalidRequest(fmt.Errorf("failed to parse request"))
		}
		updateBeforePtr = &updateBefore
		updateAfterPtr = &updateAfter
	}

	resp, err := d.DocumentationService.GetDocumentationList(
		c.UserContext(), req.Page, req.PageSize, updateBeforePtr, updateAfterPtr,
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
