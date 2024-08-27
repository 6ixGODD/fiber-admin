package mods

import (
	"fiber-admin/internal/pkg/domain/vo"
	commonservice "fiber-admin/internal/pkg/service/common/mods"
	"fiber-admin/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type IdempotencyApi struct {
	IdempotencyService commonservice.IdempotencyService
}

// GenerateIdempotencyToken generates an idempotency token.
//
//	@description	Generate an idempotency token.
//	@id				common-generate-idempotency-token
//	@summary		generate idempotency token
//	@tags			Common API
//	@accept			json
//	@produce		json
//	@security		Bearer
//	@success		200							{object}	vo.Response{data=string}	"Success"
//	@failure		400							{object}	vo.Response{data=nil}		"Invalid request"
//	@failure		401							{object}	vo.Response{data=nil}		"Unauthorized"
//	@failure		500							{object}	vo.Response{data=nil}		"Internal server error"
//	@router			/common/idempotency-token	[get]
func (i *IdempotencyApi) GenerateIdempotencyToken(c *fiber.Ctx) error {
	resp, err := i.IdempotencyService.GenerateIdempotencyToken(c.UserContext())
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
