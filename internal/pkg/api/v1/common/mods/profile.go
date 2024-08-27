package mods

import (
	"fiber-admin/internal/pkg/domain/vo"
	commonservice "fiber-admin/internal/pkg/service/common/mods"
	"fiber-admin/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type ProfileApi struct {
	ProfileService commonservice.ProfileService
}

// GetProfile returns the profile.
//
//	@description	Get the profile.
//	@id				common-get-profile
//	@summary		get profile
//	@tags			Common API
//	@accept			json
//	@produce		json
//	@security		Bearer
//	@success		200			{object}	vo.Response{data=common.GetProfileResponse}	"Success"
//	@failure		400			{object}	vo.Response{data=nil}						"Invalid request"
//	@failure		401			{object}	vo.Response{data=nil}						"Unauthorized"
//	@failure		500			{object}	vo.Response{data=nil}						"Internal server error"
//	@router			/profile	[get]
func (api *ProfileApi) GetProfile(c *fiber.Ctx) error {
	resp, err := api.ProfileService.GetProfile(c.UserContext())
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
