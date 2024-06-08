package errors

import (
	"errors"

	"fiber-admin/internal/pkg/domain/vo"
	e "fiber-admin/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *e.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Status()).JSON(
			vo.Response{
				Code:    appErr.Code(),
				Message: appErr.Error(),
				Data:    nil,
			},
		)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(
		vo.Response{
			Code:    e.CodeUnknownError,
			Message: err.Error(),
			Data:    nil,
		},
	)
}
